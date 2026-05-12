package service

import (
	"errors"
	"time"

	"github.com/tuantranpham204/CyberDiner.git/src/server/internal/model"
	"github.com/tuantranpham204/CyberDiner.git/src/server/internal/store"
	"github.com/tuantranpham204/CyberDiner.git/src/server/pkg/util"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// --- Request / Response types ---

type SignUpRequest struct {
	Name     string `json:"name"     binding:"required"`
	Surname  string `json:"surname"  binding:"required"`
	Username string `json:"username" binding:"required"`
	Email    string `json:"email"    binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
	Gender   string `json:"gender"`
	Dob      string `json:"dob"`
}

type SignInRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type UserProfile struct {
	ID       uint   `json:"id"`
	Name     string `json:"name"`
	Surname  string `json:"surname"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Role     int    `json:"role"`
}

type SignInResponse struct {
	AccessToken  string      `json:"access_token"`
	RefreshToken string      `json:"refresh_token"`
	ExpiresIn    int         `json:"expires_in"`
	User         UserProfile `json:"user"`
}

// --- Sentinel errors ---

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrAccountSuspended   = errors.New("account suspended")
)

// --- Interface ---

type AuthService interface {
	Register(req *SignUpRequest) error
	Login(req *SignInRequest) (*SignInResponse, error)
	// SignOut receives the JTI and expiry already extracted by JWTAuth middleware
	// and adds them to the denylist.
	SignOut(jti string, expiry time.Time)
}

// --- Implementation ---

type authService struct {
	userStore store.UserStore
	denylist  store.DenylistStore
}

func NewAuthService(userStore store.UserStore, denylist store.DenylistStore) AuthService {
	return &authService{userStore: userStore, denylist: denylist}
}

func (s *authService) Register(req *SignUpRequest) error {
	if _, err := s.userStore.FindByEmail(req.Email); err == nil {
		return errors.New("email already in use")
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	if _, err := s.userStore.FindByUsername(req.Username); err == nil {
		return errors.New("username already taken")
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user := &model.User{
		Name:     req.Name,
		Surname:  req.Surname,
		Username: req.Username,
		Email:    req.Email,
		Password: string(hashed),
		Gender:   req.Gender,
		Status:   model.UserStatusActive,
		Role:     model.RoleCustomer,
	}
	if req.Dob != "" {
		if dob, parseErr := time.Parse("2006-01-02", req.Dob); parseErr == nil {
			user.Dob = &dob
		}
	}

	return s.userStore.CreateUser(user)
}

func (s *authService) Login(req *SignInRequest) (*SignInResponse, error) {
	user, err := s.userStore.FindByUsername(req.Username)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrInvalidCredentials
	}
	if err != nil {
		return nil, err
	}

	if user.Status != model.UserStatusActive {
		return nil, ErrAccountSuspended
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, ErrInvalidCredentials
	}

	accessToken, err := util.GenerateAccessToken(user.ID, int(user.Role))
	if err != nil {
		return nil, err
	}
	refreshToken, err := util.GenerateRefreshToken(user.ID)
	if err != nil {
		return nil, err
	}

	return &SignInResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    900,
		User: UserProfile{
			ID:       user.ID,
			Name:     user.Name,
			Surname:  user.Surname,
			Username: user.Username,
			Email:    user.Email,
			Role:     int(user.Role),
		},
	}, nil
}

func (s *authService) SignOut(jti string, expiry time.Time) {
	if jti != "" {
		s.denylist.Add(jti, expiry)
	}
}
