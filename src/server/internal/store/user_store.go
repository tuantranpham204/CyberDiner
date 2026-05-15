package store

import (
	"context"
	"errors"

	"github.com/tuantranpham204/CyberDiner.git/src/server/internal/model/entity"
	"gorm.io/gorm"
)

var ErrNotFound = errors.New("record not found")

type UserStore interface {
	ExistsByEmail(ctx context.Context, email string) (bool, error)
	ExistsByUsername(ctx context.Context, username string) (bool, error)
	CreateWithProfile(ctx context.Context, user *entity.User, profile *entity.Profile) error
	FindByEmail(ctx context.Context, email string) (*entity.User, error)
	FindByUsername(ctx context.Context, username string) (*entity.User, error)
	FindByID(ctx context.Context, id int64) (*entity.User, error)
	UpdateProfile(ctx context.Context, userID int64, fields map[string]any) error
}

type userStore struct {
	db *gorm.DB
}

func NewUserStore(db *gorm.DB) UserStore {
	return &userStore{db: db}
}

func (s *userStore) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	var count int64
	err := s.db.WithContext(ctx).Model(&entity.User{}).Where("email = ?", email).Count(&count).Error
	return count > 0, err
}

func (s *userStore) ExistsByUsername(ctx context.Context, username string) (bool, error) {
	var count int64
	err := s.db.WithContext(ctx).Model(&entity.User{}).Where("username = ?", username).Count(&count).Error
	return count > 0, err
}

func (s *userStore) CreateWithProfile(ctx context.Context, user *entity.User, profile *entity.Profile) error {
	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(user).Error; err != nil {
			return err
		}
		profile.UserID = user.ID
		if err := tx.Create(profile).Error; err != nil {
			return err
		}
		return nil
	})
}

func (s *userStore) FindByEmail(ctx context.Context, email string) (*entity.User, error) {
	var u entity.User
	err := s.db.WithContext(ctx).Preload("Profile").Where("email = ?", email).First(&u).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (s *userStore) FindByUsername(ctx context.Context, username string) (*entity.User, error) {
	var u entity.User
	err := s.db.WithContext(ctx).Preload("Profile").Where("username = ?", username).First(&u).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (s *userStore) UpdateProfile(ctx context.Context, userID int64, fields map[string]any) error {
	if len(fields) == 0 {
		return nil
	}
	res := s.db.WithContext(ctx).
		Model(&entity.Profile{}).
		Where("user_id = ?", userID).
		Updates(fields)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return ErrNotFound
	}
	return nil
}

func (s *userStore) FindByID(ctx context.Context, id int64) (*entity.User, error) {
	var u entity.User
	err := s.db.WithContext(ctx).Preload("Profile").Where("id = ?", id).First(&u).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return &u, nil
}
