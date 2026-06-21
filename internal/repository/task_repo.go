package repository

import (
	"context"
	"errors"
	"time"

	"github.com/s-usmonalizoda25/taskManagerProject/internal/models"
	"github.com/s-usmonalizoda25/taskManagerProject/pkg/errs"
	"gorm.io/gorm"
)

type ITaskRepository interface {
	Create(ctx context.Context, task *models.Task) error
	GetActive(ctx context.Context, status string, limit, offset int) ([]models.Task, error)
	GetAll(ctx context.Context, status string, limit, offset int) ([]models.Task, error)
	GetByID(ctx context.Context, id uint) (*models.Task, error)
	Update(ctx context.Context, task *models.Task) error
	Deactivate(ctx context.Context, id uint) error
	Delete(ctx context.Context, id uint) error
	UpdateStatus(ctx context.Context, id uint, status models.TaskStatus) (*models.Task, error)
}

type taskRepository struct {
	db *gorm.DB
}

func NewTaskRepository(db *gorm.DB) ITaskRepository {
	return &taskRepository{db: db}
}

func (r *taskRepository) Create(ctx context.Context, task *models.Task) error {
	task.UpdatedAt = nil
	task.DeletedAt = nil
	return r.db.WithContext(ctx).Create(task).Error
}

func (r *taskRepository) GetActive(ctx context.Context, status string, limit, offset int) ([]models.Task, error) {
	var tasks []models.Task
	query := r.db.WithContext(ctx).Where("deleted_at IS NULL").Limit(limit).Offset(offset)
	if status != "" {
		query = query.Where("status = ?", status)
	}
	err := query.Find(&tasks).Error
	return tasks, err
}

func (r *taskRepository) GetAll(ctx context.Context, status string, limit, offset int) ([]models.Task, error) {
	var tasks []models.Task
	query := r.db.WithContext(ctx).Limit(limit).Offset(offset)
	if status != "" {
		query = query.Where("status = ?", status)
	}
	err := query.Find(&tasks).Error
	return tasks, err
}

func (r *taskRepository) GetByID(ctx context.Context, id uint) (*models.Task, error) {
	var task models.Task
	err := r.db.WithContext(ctx).Where("deleted_at IS NULL").First(&task, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errs.ErrTaskNotFound
	}
	return &task, err
}

func (r *taskRepository) Update(ctx context.Context, task *models.Task) error {
	now := time.Now()
	task.UpdatedAt = &now

	return r.db.WithContext(ctx).Model(task).
		Select("Title", "Description", "Status", "UpdatedAt").
		Updates(task).Error
}

func (r *taskRepository) Deactivate(ctx context.Context, id uint) error {
	now := time.Now()
	result := r.db.WithContext(ctx).
		Model(&models.Task{}).
		Where("id = ? AND deleted_at IS NULL", id).
		Update("deleted_at", now)

	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errs.ErrTaskNotFound
	}
	return nil
}

func (r *taskRepository) Delete(ctx context.Context, id uint) error {
	result := r.db.WithContext(ctx).Unscoped().Delete(&models.Task{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errs.ErrTaskNotFound
	}
	return nil
}

func (r *taskRepository) UpdateStatus(ctx context.Context, id uint, status models.TaskStatus) (*models.Task, error) {
	now := time.Now()
	var task models.Task

	err := r.db.WithContext(ctx).Where("id = ? AND deleted_at IS NULL", id).First(&task).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errs.ErrTaskNotFound
	}
	if err != nil {
		return nil, err
	}

	task.Status = status
	task.UpdatedAt = &now

	err = r.db.WithContext(ctx).Model(&task).Select("Status", "UpdatedAt").Updates(&task).Error
	if err != nil {
		return nil, err
	}

	return &task, nil
}
