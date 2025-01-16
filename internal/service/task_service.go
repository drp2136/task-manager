package service

import (
	"task-manager/internal/model"

	"github.com/jinzhu/gorm"
)

type TaskService struct {
	DB *gorm.DB
}

func (s *TaskService) CreateTask(task *model.Task) error {
	return s.DB.Create(task).Error
}

func (s *TaskService) GetAllTasks() ([]model.Task, error) {
	var tasks []model.Task
	err := s.DB.Find(&tasks).Error
	return tasks, err
}

func (s *TaskService) GetTaskByID(id uint) (*model.Task, error) {
	var task model.Task
	err := s.DB.First(&task, id).Error
	return &task, err
}

func (s *TaskService) UpdateTask(id uint, task *model.Task) error {
	return s.DB.Model(&model.Task{}).Where("id = ?", id).Updates(task).Error
}

func (s *TaskService) DeleteTask(id uint) error {
	return s.DB.Delete(&model.Task{}, id).Error
}
