package service

import (
	"context"
	"database/sql"
	"golang-rest-api/internal/model"
	"golang-rest-api/internal/model/request"
	"golang-rest-api/internal/model/response"
	"golang-rest-api/internal/repository"
	"golang-rest-api/pkg"
	"log"
	"time"
)

type TaskService interface {
	Create(ctx context.Context, request request.CreateTaskRequest) (*response.TaskResponse, error)
	GetById(ctx context.Context, id int) (*response.TaskResponse, error)
}

type TaskServiceImpl struct {
	TaskRepository repository.TaskRepository
	DB             *sql.DB
}

func NewTaskService(taskRepository repository.TaskRepository, db *sql.DB) TaskService {
	return &TaskServiceImpl{
		TaskRepository: taskRepository,
		DB:             db,
	}
}

// Method
func (s *TaskServiceImpl) Create(ctx context.Context, request request.CreateTaskRequest) (*response.TaskResponse, error) {
	// Begin transaction from DB
	tx, err := s.DB.Begin()
	if err != nil {
		log.Printf("Error starting transaction: %v", err)
		return nil, err
	}

	// Defer rollback jika error
	defer func() {
		if err != nil {
			log.Printf("Rolling back transaction due to error: %v", err)
			tx.Rollback()
		}
	}()

	// Convert request ke model
	task := s.constructCreateTask(request)
	log.Printf("Constructed Task: %+v", task)

	// Call repository create
	id, err := s.TaskRepository.Create(ctx, tx, task)
	if err != nil {
		log.Printf("Error creating task: %v", err)
		return nil, err
	}

	// Commit transaction
	err = tx.Commit()
	if err != nil {
		log.Printf("Error committing transaction: %v", err)
		tx.Rollback()
		return nil, pkg.ErrInternalServerError
	}

	// Get by id
	response, err := s.GetById(ctx, *id)
	if err != nil {
		log.Printf("Error getting task by ID %d: %v", *id, err)
		return nil, pkg.ErrNotFound
	}

	return response, nil
}

func (s *TaskServiceImpl) GetById(ctx context.Context, id int) (*response.TaskResponse, error) {
	data, err := s.TaskRepository.GetById(ctx, s.DB, id)
	if err != nil {
		return nil, err
	}

	task := s.constructTaskResponse(*data)
	return &task, nil
}

func (s *TaskServiceImpl) constructTaskResponse(task model.Task) response.TaskResponse {
	return response.TaskResponse{
		Id:          task.Id,
		Title:       task.Title,
		Description: task.Description,
		Status:      task.Status,
		IsActive:    task.IsActive,
		DueDate:     task.DueDate,
		CreatedAt:   task.CreatedAt,
		UpdatedAt:   task.UpdatedAt,
	}
}

func (s *TaskServiceImpl) constructCreateTask(request request.CreateTaskRequest) model.Task {
	return model.Task{
		Id:          0,
		Title:       request.Title,
		Description: request.Description,
		DueDate:     request.DueDate,
		Status:      0,
		IsActive:    1,
		CreatedAt:   time.Now(),
		UpdatedAt:   nil,
	}
}
