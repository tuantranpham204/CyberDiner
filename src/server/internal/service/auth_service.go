package service

import (
	"context"
	"errors"
	"strconv"
	"strings"
	"time"

	"github.com/tuantranpham204/CyberDiner.git/src/server/internal/model/entity"
	"github.com/tuantranpham204/CyberDiner.git/src/server/internal/store"
	"github.com/tuantranpham204/CyberDiner.git/src/server/pkg/util"
)

var (
	ErrEmailAlreadyExists    = errors.New("email already registered")
	ErrUsernameAlreadyExists = errors.New("username already taken")
	ErrPasswordMismatch      = errors.New("password and confirm password do not match")
	ErrInvalidDateOfBirth    = errors.New("invalid date of birth")
	ErrInvalidGender         = errors.New("invalid gender")
	ErrInvalidCredentials    = errors.New("invalid credentials")
	ErrAccountBanned         = errors.New("account is banned")
	ErrAccountInactive       = errors.New("account is inactive")
)

type SignUpInput struct {
	Name            string
	Surname         string
	Username        string
	Email           string
	DateOfBirth     string
	Gender          string
	Password        string
	ConfirmPassword string
}

type SignInInput struct {
	Username string
	Password string
}

type AuthResult struct {
	User  *entity.User
	Token string
}

// SignUpResult kept as alias for backwards compatibility.
type SignUpResult = AuthResult

type AuthService interface {
	SignUp(ctx context.Context, in SignUpInput) (*AuthResult, error)
	SignIn(ctx context.Context, in SignInInput) (*AuthResult, error)
	SignOut(ctx context.Context, jti string, expiresAt time.Time) error
}

type authService struct {
	users     store.UserStore
	denylist  store.TokenDenylist
	jwt       *util.JWTManager
}

func NewAuthService(users store.UserStore, denylist store.TokenDenylist, jwt *util.JWTManager) AuthService {
	return &authService{users: users, denylist: denylist, jwt: jwt}
}

// SignOut adds the supplied jti to the denylist with a TTL equal to its
// remaining lifetime, so the entry expires together with the token itself.
// Tokens that are already expired are accepted as a no-op (UC-03 alt-flow).
func (s *authService) SignOut(ctx context.Context, jti string, expiresAt time.Time) error {
	if jti == "" {
		return nil
	}
	ttl := time.Until(expiresAt)
	if ttl <= 0 {
		return nil
	}
	return s.denylist.Add(ctx, jti, ttl)
}

func (s *authService) SignUp(ctx context.Context, in SignUpInput) (*AuthResult, error) {
	if in.Password != in.ConfirmPassword {
		return nil, ErrPasswordMismatch
	}

	email := strings.ToLower(strings.TrimSpace(in.Email))
	username := strings.TrimSpace(in.Username)

	dob, err := time.Parse("2006-01-02", in.DateOfBirth)
	if err != nil {
		return nil, ErrInvalidDateOfBirth
	}

	genderVal, err := util.ParseGender(in.Gender)
	if err != nil {
		return nil, ErrInvalidGender
	}

	emailExists, err := s.users.ExistsByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	if emailExists {
		return nil, ErrEmailAlreadyExists
	}

	usernameExists, err := s.users.ExistsByUsername(ctx, username)
	if err != nil {
		return nil, err
	}
	if usernameExists {
		return nil, ErrUsernameAlreadyExists
	}

	hash, err := util.HashPassword(in.Password)
	if err != nil {
		return nil, err
	}

	user := &entity.User{
		Username:     username,
		Email:        email,
		PasswordHash: hash,
		Role:         entity.RoleUser,
		IsActive:     true,
	}
	profile := &entity.Profile{
		Name:        strings.TrimSpace(in.Name),
		Surname:     strings.TrimSpace(in.Surname),
		DateOfBirth: dob,
		Gender:      genderVal,
	}

	if err := s.users.CreateWithProfile(ctx, user, profile); err != nil {
		return nil, err
	}
	user.Profile = profile

	token, err := s.jwt.Generate(strconv.FormatInt(user.ID, 10), user.Username, string(user.Role))
	if err != nil {
		return nil, err
	}

	return &AuthResult{User: user, Token: token}, nil
}

func (s *authService) SignIn(ctx context.Context, in SignInInput) (*AuthResult, error) {
	username := strings.TrimSpace(in.Username)
	if username == "" || in.Password == "" {
		return nil, ErrInvalidCredentials
	}

	user, err := s.users.FindByUsername(ctx, username)
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			return nil, ErrInvalidCredentials
		}
		return nil, err
	}

	if !util.CheckPassword(user.PasswordHash, in.Password) {
		return nil, ErrInvalidCredentials
	}

	if user.IsBanned {
		return nil, ErrAccountBanned
	}
	if !user.IsActive {
		return nil, ErrAccountInactive
	}

	token, err := s.jwt.Generate(strconv.FormatInt(user.ID, 10), user.Username, string(user.Role))
	if err != nil {
		return nil, err
	}
	return &AuthResult{User: user, Token: token}, nil
}
