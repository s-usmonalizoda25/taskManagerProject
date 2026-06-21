package repository

import (
	"context"
	"errors"
	"time"

	"github.com/s-usmonalizoda25/taskManagerProject/internal/models"
	"github.com/s-usmonalizoda25/taskManagerProject/pkg/errs"
	"gorm.io/gorm"
)

type IUserRepository interface {
	CreateUser(ctx context.Context, user *models.User) error
	GetActiveUser(ctx context.Context, limit, offset int) ([]models.User, error)
	GetAllUsers(ctx context.Context, limit, offset int) ([]models.User, error)
	GetByUserID(ctx context.Context, id uint) (*models.User, error)
	UpdateUser(ctx context.Context, user *models.User) error
	DeactivateUser(ctx context.Context, id uint) error
	DeleteUser(ctx context.Context, id uint) error
	GetByUsernameOrEmail(ctx context.Context, username, email string) (*models.User, error)
}

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) IUserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) CreateUser(ctx context.Context, user *models.User) error {
	user.UpdatedAt = nil
	user.DeletedAt = nil
	return r.db.WithContext(ctx).Create(user).Error
}

func (r *UserRepository) GetActiveUser(ctx context.Context, limit, offset int) ([]models.User, error) {
	var users []models.User

	err := r.db.WithContext(ctx).Where("deleted_at IS NULL").Limit(limit).Offset(offset).Find(&users).Error
	return users, err
}

func (r *UserRepository) GetAllUsers(ctx context.Context, limit, offset int) ([]models.User, error) {
	var users []models.User

	err := r.db.WithContext(ctx).Limit(limit).Offset(offset).Find(&users).Error
	return users, err
}

func (r *UserRepository) GetByUserID(ctx context.Context, id uint) (*models.User, error) {
	var user models.User
	err := r.db.WithContext(ctx).Where("deleted_at IS NULL").First(&user, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errs.ErrUserNotFound
	}
	return &user, err
}

func (r *UserRepository) UpdateUser(ctx context.Context, user *models.User) error {
	now := time.Now()
	user.UpdatedAt = &now

	return r.db.WithContext(ctx).Model(user).
		Select("Username", "Email", "UpdatedAt").
		Updates(user).Error
}

func (r *UserRepository) DeactivateUser(ctx context.Context, id uint) error {
	now := time.Now()
	result := r.db.WithContext(ctx).
		Model(&models.User{}).
		Where("id = ? AND deleted_at IS NULL", id).
		Update("deleted_at", now)

	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errs.ErrUserNotFound
	}
	return nil
}

func (r *UserRepository) DeleteUser(ctx context.Context, id uint) error {
	result := r.db.WithContext(ctx).Unscoped().Delete(&models.User{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errs.ErrUserNotFound
	}
	return nil
}

func (r *UserRepository) GetByUsernameOrEmail(ctx context.Context, username, email string) (*models.User, error) {
	var user models.User

	err := r.db.WithContext(ctx).Where("(username = ? OR email = ?) AND deleted_at IS NULL", username, email).First(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &user, err
}
