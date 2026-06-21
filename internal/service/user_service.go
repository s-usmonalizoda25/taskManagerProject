package service

import (
	"context"

	"github.com/s-usmonalizoda25/taskManagerProject/internal/models"
	"github.com/s-usmonalizoda25/taskManagerProject/internal/repository"
	"github.com/s-usmonalizoda25/taskManagerProject/pkg/errs"
	"github.com/s-usmonalizoda25/taskManagerProject/pkg/logger"
	"go.uber.org/zap"
)

type IUserService interface {
	CreateUser(ctx context.Context, user *models.User) error
	UpdateUser(ctx context.Context, id uint, input *models.User) error
	DeactivateUser(ctx context.Context, id uint) error
	DeleteUser(ctx context.Context, id uint) error
	GetActiveUsers(ctx context.Context, page, limit int) ([]models.User, error)
	GetAllUsers(ctx context.Context, page, limit int) ([]models.User, error)
}

type userService struct {
	repo repository.IUserRepository
	log  *logger.Logger
}

func NewUserService(repo repository.IUserRepository, log *logger.Logger) IUserService {
	return &userService{repo: repo, log: log}
}

func (s *userService) CreateUser(ctx context.Context, user *models.User) error {
	s.log.Info("Attempting to create user", zap.String("username", user.Username))
	if user.Username == "" {
		return errs.ErrEmptyUsername
	}
	if user.Email == "" {
		return errs.ErrEmptyEmail
	}

	existingUser, err := s.repo.GetByUsernameOrEmail(ctx, user.Username, user.Email)
	if err != nil {
		return err
	}
	if existingUser != nil {
		if existingUser.Username == user.Username {
			return errs.ErrUsernameTaken
		}
		if existingUser.Email == user.Email {
			return errs.ErrEmailTaken
		}
	}

	return s.repo.CreateUser(ctx, user)
}

func (s *userService) UpdateUser(ctx context.Context, id uint, input *models.User) error {
	s.log.Info("Attempting to update user", zap.Uint("user_id", id))
	if id == 0 {
		return errs.ErrInvalidID
	}
	if input.Username == "" {
		return errs.ErrEmptyUsername
	}
	if input.Email == "" {
		return errs.ErrEmptyEmail
	}

	user, err := s.repo.GetByUserID(ctx, id)
	if err != nil {
		return err
	}

	user.Username = input.Username
	user.Email = input.Email

	return s.repo.UpdateUser(ctx, user)
}

func (s *userService) DeactivateUser(ctx context.Context, id uint) error {
	s.log.Info("Attempting to deactivate (soft-delete) user", zap.Uint("user_id", id))
	if id == 0 {
		return errs.ErrInvalidID
	}
	return s.repo.DeactivateUser(ctx, id)
}

func (s *userService) DeleteUser(ctx context.Context, id uint) error {
	s.log.Info("Attempting to hard delete user", zap.Uint("user_id", id))
	if id == 0 {
		return errs.ErrInvalidID
	}
	return s.repo.DeleteUser(ctx, id)
}

func (s *userService) GetActiveUsers(ctx context.Context, page, limit int) ([]models.User, error) {
	s.log.Info("Fetching active users list")
	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 10
	}
	offset := (page - 1) * limit
	return s.repo.GetActiveUser(ctx, limit, offset)
}

func (s *userService) GetAllUsers(ctx context.Context, page, limit int) ([]models.User, error) {
	s.log.Info("Fetching absolutely all users list")
	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 10
	}
	offset := (page - 1) * limit
	return s.repo.GetAllUsers(ctx, limit, offset)
}
