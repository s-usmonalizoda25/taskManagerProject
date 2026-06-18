package service

import (
	"context"
	"errors"

	"github.com/s-usmonalizoda25/taskManagerProject/internal/models"
	"github.com/s-usmonalizoda25/taskManagerProject/internal/repository"
)

var ErrEmptyTitle = errors.New("title cannot be empty")

type ITaskService interface {
	CreateTask(ctx context.Context, task *models.Task) error
	GetTask(ctx context.Context, status string, page, limit int) ([]models.Task, error)
	GetTaskByID(ctx context.Context, id uint) (*models.Task, error)
	UpdateTask(ctx context.Context, id uint, input *models.Task) error
	DeleteTask(ctx context.Context, id uint) error
}

type taskService struct {
	repo repository.ITaskRepository
}

func NewTaskService(repo repository.ITaskRepository) ITaskService {
	return &taskService{repo: repo}
}

func (s *taskService) CreateTask(ctx context.Context, task *models.Task) error {
	if task.Title == "" {
		return ErrEmptyTitle
	}

	if task.Status == "" {
		task.Status = "new"
	}

	return s.repo.Create(ctx, task)
}

func (s *taskService) GetTask(ctx context.Context, status string, page, limit int) ([]models.Task, error) {
	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 10
	}

	offset := (page - 1) * limit

	return s.repo.GetAll(ctx, status, limit, offset)
}

func (s *taskService) GetTaskByID(ctx context.Context, id uint) (*models.Task, error) {
	return s.repo.GetById(ctx, id)
}

func (s *taskService) UpdateTask(ctx context.Context, id uint, input *models.Task) error {
	if input.Title == "" {
		return ErrEmptyTitle
	}

	task, err := s.repo.GetById(ctx, id)
	if err != nil {
		return err
	}

	task.Title = input.Title
	task.Description = input.Description
	task.Status = input.Status

	return s.repo.Update(ctx, task)
}

func (s *taskService) DeleteTask(ctx context.Context, id uint) error {
	return s.repo.Delete(ctx, id)
}
