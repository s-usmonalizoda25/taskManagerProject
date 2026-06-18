package repository

import (
	"context"
	"errors"

	"github.com/s-usmonalizoda25/taskManagerProject/internal/models"
	"gorm.io/gorm"
)

var ErrTaskNotFound = errors.New("task not found")

type ITaskRepository interface {
	Create(ctx context.Context, task *models.Task) error
	GetAll(ctx context.Context, status string, limit, offset int) ([]models.Task, error)
	GetById(ctx context.Context, id uint) (*models.Task, error)
	Update(ctx context.Context, task *models.Task) error
	Delete(ctx context.Context, id uint) error
}

type taskRepository struct {
	db *gorm.DB
}

func NewTaskRepository(db *gorm.DB) ITaskRepository {
	return &taskRepository{db: db}
}

func (r *taskRepository) Create(ctx context.Context, task *models.Task) error {
	return r.db.WithContext(ctx).Create(task).Error
}

func (r *taskRepository) GetAll(ctx context.Context, status string, limit, offset int) ([]models.Task, error) {
	var tasks []models.Task
	tx := r.db.WithContext(ctx).Order("id DESC").Limit(limit).Offset(offset)

	if status != "" {
		tx = tx.Where("status = ?", status)
	}

	err := tx.Find(&tasks).Error
	if err != nil {
		return nil, err
	}
	return tasks, nil
}
func (r *taskRepository) GetById(ctx context.Context, id uint) (*models.Task, error) {
	var task models.Task

	err := r.db.WithContext(ctx).First(&task, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrTaskNotFound
		}
		return nil, err
	}
	return &task, nil
}

func (r *taskRepository) Update(ctx context.Context, task *models.Task) error {
	return r.db.WithContext(ctx).Save(task).Error
}

func (r *taskRepository) Delete(ctx context.Context, id uint) error {
	res := r.db.WithContext(ctx).Delete(&models.Task{}, id)
	if res.Error != nil {
		return res.Error
	}

	if res.RowsAffected == 0 {
		return ErrTaskNotFound
	}

	return nil
}
