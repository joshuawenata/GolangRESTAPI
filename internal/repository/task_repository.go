package repository

import (
	"context"
	"database/sql"
	"golang-rest-api/internal/model"
	"log"
)

type TaskRepository interface {
	Create(ctx context.Context, tx *sql.Tx, task model.Task) (*int, error)
	// GetAll()
	GetById(ctx context.Context, db *sql.DB, id int) (*model.Task, error)
	// Update()
	// Delete()
}

type TaskRepositoryImpl struct {
}

func NewTaskRepository() TaskRepository {
	return &TaskRepositoryImpl{}
}

// Method
func (r *TaskRepositoryImpl) Create(ctx context.Context, tx *sql.Tx, task model.Task) (*int, error) {
	// Define query
	query := `INSERT INTO tasks(title, description, status, is_active, due_date) VALUES ($1, $2, $3, 1, $4) RETURNING id;`

	log.Printf("Executing Query: %s\n", query)
	log.Printf("With Values: %v, %v, %v, %v", task.Title, task.Description, task.Status, task.DueDate)

	// Execute query
	row := tx.QueryRowContext(ctx, query, task.Title, task.Description, task.Status, task.DueDate)

	// Get id
	var id int
	err := row.Scan(&id)
	if err != nil {
		log.Printf("Error scanning row: %v", err)
		return nil, err
	}

	log.Printf("Task created with ID: %d", id)
	return &id, nil
}

// func (r *TaskRepositoryImpl) GetAll() {

// }

func (r *TaskRepositoryImpl) GetById(ctx context.Context, db *sql.DB, id int) (*model.Task, error) {
	// Define Query
	query := `SELECT id, title, description, status, is_active, due_date FROM tasks WHERE id = $1`

	// Execute Query
	row := db.QueryRowContext(ctx, query, id)

	// Scan row to task model
	var task model.Task
	if err := row.Scan(
		&task.Id, &task.Title, &task.Description, &task.Status, &task.IsActive, &task.DueDate,
	); err != nil {
		return nil, err
	}

	// Return task and error
	return &task, nil
}

// func (r *TaskRepositoryImpl) Update() {

// }

// func (r *TaskRepositoryImpl) Delete() {

// }
