package service

import (
	"context"
	"errors"
	"strings"

	"github.com/tuantranpham204/CyberDiner.git/src/server/internal/model/entity"
	"github.com/tuantranpham204/CyberDiner.git/src/server/internal/store"
)

var (
	ErrProfileNotFound = errors.New("profile not found")
	ErrNothingToUpdate = errors.New("no updatable fields supplied")
)

type UpdateProfileInput struct {
	Name        *string
	Surname     *string
	PhoneNumber *string
	Address     *string
}

type ProfileService interface {
	Get(ctx context.Context, userID int64) (*entity.User, error)
	Update(ctx context.Context, userID int64, in UpdateProfileInput) (*entity.User, error)
}

type profileService struct {
	users store.UserStore
}

func NewProfileService(users store.UserStore) ProfileService {
	return &profileService{users: users}
}

func (s *profileService) Get(ctx context.Context, userID int64) (*entity.User, error) {
	user, err := s.users.FindByID(ctx, userID)
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			return nil, ErrProfileNotFound
		}
		return nil, err
	}
	return user, nil
}

func (s *profileService) Update(ctx context.Context, userID int64, in UpdateProfileInput) (*entity.User, error) {
	fields := map[string]any{}
	if in.Name != nil {
		fields["name"] = strings.TrimSpace(*in.Name)
	}
	if in.Surname != nil {
		fields["surname"] = strings.TrimSpace(*in.Surname)
	}
	if in.PhoneNumber != nil {
		v := strings.TrimSpace(*in.PhoneNumber)
		if v == "" {
			fields["phone_number"] = nil
		} else {
			fields["phone_number"] = v
		}
	}
	if in.Address != nil {
		v := strings.TrimSpace(*in.Address)
		if v == "" {
			fields["address"] = nil
		} else {
			fields["address"] = v
		}
	}

	if len(fields) == 0 {
		return nil, ErrNothingToUpdate
	}

	if err := s.users.UpdateProfile(ctx, userID, fields); err != nil {
		if errors.Is(err, store.ErrNotFound) {
			return nil, ErrProfileNotFound
		}
		return nil, err
	}
	return s.Get(ctx, userID)
}
