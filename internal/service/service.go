package service

import (
	"context"

	"github.com/s-usmonalizoda25/taskManagerProject/internal/models"
	"github.com/s-usmonalizoda25/taskManagerProject/internal/repository"
	"github.com/s-usmonalizoda25/taskManagerProject/pkg/errs"
	"github.com/s-usmonalizoda25/taskManagerProject/pkg/logger"
	"go.uber.org/zap"
)

type ITaskService interface {
	CreateTask(ctx context.Context, task *models.Task) error
	GetActiveTasks(ctx context.Context, status string, page, limit int) ([]models.Task, error)
	GetAllTasks(ctx context.Context, status string, page, limit int) ([]models.Task, error)
	GetTaskByID(ctx context.Context, id uint) (*models.Task, error)
	UpdateTask(ctx context.Context, id uint, input *models.Task) error
	DeactivateTask(ctx context.Context, id uint) error
	DeleteTask(ctx context.Context, id uint) error
}

type taskService struct {
	repo repository.ITaskRepository
	log  *logger.Logger
}

func NewTaskService(repo repository.ITaskRepository, log *logger.Logger) ITaskService {
	return &taskService{repo: repo, log: log}
}

func (s *taskService) CreateTask(ctx context.Context, task *models.Task) error {
	s.log.Info("Attempting to create task", zap.String("title", task.Title))
	if task.Title == "" {
		return errs.ErrEmptyTitle
	}
	if task.Status == "" {
		task.Status = models.StatusNew
	}
	return s.repo.Create(ctx, task)
}

func (s *taskService) GetActiveTasks(ctx context.Context, status string, page, limit int) ([]models.Task, error) {
	s.log.Info("Fetching active tasks list", zap.String("status", status))
	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 10
	}
	offset := (page - 1) * limit
	return s.repo.GetActive(ctx, status, limit, offset)
}

func (s *taskService) GetAllTasks(ctx context.Context, status string, page, limit int) ([]models.Task, error) {
	s.log.Info("Fetching absolutely all tasks list", zap.String("status", status))
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
	s.log.Info("Fetching task by ID", zap.Uint("task_id", id))
	if id == 0 {
		return nil, errs.ErrInvalidID
	}
	return s.repo.GetByID(ctx, id)
}

func (s *taskService) UpdateTask(ctx context.Context, id uint, input *models.Task) error {
	s.log.Info("Attempting to update task", zap.Uint("task_id", id))
	if id == 0 {
		return errs.ErrInvalidID
	}
	if input.Title == "" {
		return errs.ErrEmptyTitle
	}

	task, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	task.Title = input.Title
	task.Description = input.Description
	task.Status = input.Status

	return s.repo.Update(ctx, task)
}

func (s *taskService) DeactivateTask(ctx context.Context, id uint) error {
	s.log.Info("Attempting to soft-delete (deactivate) task", zap.Uint("task_id", id))
	if id == 0 {
		return errs.ErrInvalidID
	}
	return s.repo.Deactivate(ctx, id)
}

func (s *taskService) DeleteTask(ctx context.Context, id uint) error {
	s.log.Info("Attempting hard delete task from DB", zap.Uint("task_id", id))
	if id == 0 {
		return errs.ErrInvalidID
	}
	return s.repo.Delete(ctx, id)
}
