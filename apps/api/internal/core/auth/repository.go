package auth

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

var ErrUserNotFound = errors.New("user not found")

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) CreateUser(ctx context.Context, user *User) error {
	return r.db.WithContext(ctx).Create(user).Error
}

func (r *Repository) FindUserByEmail(ctx context.Context, email string) (*User, error) {
	var user User
	if err := r.db.WithContext(ctx).Where("email = ?", email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	return &user, nil
}

func (r *Repository) FindUserByID(ctx context.Context, id uuid.UUID) (*User, error) {
	var user User
	if err := r.db.WithContext(ctx).First(&user, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	return &user, nil
}

func (r *Repository) UpdateLastLogin(ctx context.Context, userID uuid.UUID) error {
	now := time.Now().UTC()
	return r.db.WithContext(ctx).Model(&User{}).Where("id = ?", userID).Update("last_login_at", now).Error
}

func (r *Repository) CreateRefreshToken(ctx context.Context, token *RefreshToken) error {
	return r.db.WithContext(ctx).Create(token).Error
}

func (r *Repository) FindRefreshTokenByHash(ctx context.Context, tokenHash string) (*RefreshToken, error) {
	var token RefreshToken
	if err := r.db.WithContext(ctx).Where("token_hash = ?", tokenHash).First(&token).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("refresh token not found")
		}
		return nil, err
	}
	return &token, nil
}

func (r *Repository) RevokeRefreshToken(ctx context.Context, tokenID uuid.UUID) error {
	now := time.Now().UTC()
	return r.db.WithContext(ctx).Model(&RefreshToken{}).Where("id = ?", tokenID).Update("revoked_at", now).Error
}

func (r *Repository) RevokeAllUserTokens(ctx context.Context, userID uuid.UUID) error {
	now := time.Now().UTC()
	return r.db.WithContext(ctx).Model(&RefreshToken{}).
		Where("user_id = ? AND revoked_at IS NULL", userID).
		Update("revoked_at", now).Error
}
