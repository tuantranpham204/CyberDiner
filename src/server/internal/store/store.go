package store

import (
	"sync"
	"time"

	"github.com/tuantranpham204/CyberDiner.git/src/server/internal/model"
	"gorm.io/gorm"
)

// --- UserStore ---

type UserStore interface {
	CreateUser(user *model.User) error
	FindByEmail(email string) (*model.User, error)
	FindByUsername(username string) (*model.User, error)
}

type userStore struct {
	db *gorm.DB
}

func NewUserStore(db *gorm.DB) UserStore {
	return &userStore{db: db}
}

func (s *userStore) CreateUser(user *model.User) error {
	return s.db.Create(user).Error
}

func (s *userStore) FindByEmail(email string) (*model.User, error) {
	var user model.User
	err := s.db.Where("email = ?", email).First(&user).Error
	return &user, err
}

func (s *userStore) FindByUsername(username string) (*model.User, error) {
	var user model.User
	err := s.db.Where("username = ?", username).First(&user).Error
	return &user, err
}

// --- DenylistStore ---

// DenylistStore tracks revoked JWT IDs so sign-out takes effect before natural expiry.
type DenylistStore interface {
	Add(jti string, expiry time.Time)
	Contains(jti string) bool
}

type inMemoryDenylist struct {
	mu      sync.RWMutex
	entries map[string]time.Time // jti -> expiry
}

// NewInMemoryDenylist creates a denylist that prunes expired entries every 10 minutes.
func NewInMemoryDenylist() DenylistStore {
	d := &inMemoryDenylist{entries: make(map[string]time.Time)}
	go d.cleanup()
	return d
}

func (d *inMemoryDenylist) Add(jti string, expiry time.Time) {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.entries[jti] = expiry
}

func (d *inMemoryDenylist) Contains(jti string) bool {
	d.mu.RLock()
	expiry, ok := d.entries[jti]
	d.mu.RUnlock()
	if !ok {
		return false
	}
	// Treat naturally-expired entries as absent; cleanup will remove them.
	return time.Now().Before(expiry)
}

func (d *inMemoryDenylist) cleanup() {
	ticker := time.NewTicker(10 * time.Minute)
	defer ticker.Stop()
	for range ticker.C {
		now := time.Now()
		d.mu.Lock()
		for jti, expiry := range d.entries {
			if now.After(expiry) {
				delete(d.entries, jti)
			}
		}
		d.mu.Unlock()
	}
}
