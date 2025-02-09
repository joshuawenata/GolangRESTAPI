package model

import "time"

type Task struct {
	Id          int
	Title       string
	Description string
	Status      int
	IsActive    int
	DueDate     time.Time
	CreatedAt   time.Time
	UpdatedAt   *time.Time
}
