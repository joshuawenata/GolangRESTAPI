package api

import (
	"database/sql"
	"golang-rest-api/internal/handler"
	"golang-rest-api/internal/repository"
	"golang-rest-api/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type Handlers struct {
	TaskHandler handler.TaskHandler
}

func InitRoutes(db *sql.DB) *gin.Engine {
	return setupRoutes(*initHandler(db))
}

func initHandler(db *sql.DB) *Handlers {
	// validator
	validator := validator.New()

	// task
	taskRepository := repository.NewTaskRepository()
	taskService := service.NewTaskService(taskRepository, db)
	taskHandler := handler.NewTaskHandler(taskService, validator)

	return &Handlers{
		TaskHandler: taskHandler,
	}
}

// kapital bisa diakses oleh semua package
func setupRoutes(handlers Handlers) *gin.Engine {
	route := gin.Default()

	api := route.Group("/api/v1")
	api.GET("/ping", index)

	// Task
	task := api.Group("/tasks")
	task.POST("", handlers.TaskHandler.Create)
	task.GET("/:id", handlers.TaskHandler.GetById)

	return route
}

func index(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "success",
	})
}
