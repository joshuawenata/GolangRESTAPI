package handler

import (
	"context"
	"errors"
	"golang-rest-api/internal/model/request"
	"golang-rest-api/internal/service"
	"golang-rest-api/pkg"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type TaskHandler interface {
	Create(c *gin.Context)
	GetById(c *gin.Context)
}

type TaskHandlerImpl struct {
	TaskService service.TaskService
	Validator   *validator.Validate
}

func NewTaskHandler(taskService service.TaskService, validator *validator.Validate) TaskHandler {
	return &TaskHandlerImpl{
		TaskService: taskService,
		Validator:   validator,
	}
}

// Method
func (h *TaskHandlerImpl) Create(c *gin.Context) {
	var request request.CreateTaskRequest
	err := c.ShouldBindJSON(&request)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "bad request",
		})
	}

	// Create Context
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	response, err := h.TaskService.Create(ctx, request)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": "internal server error",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "success",
		"data":    response,
	})
}

func (h *TaskHandlerImpl) GetById(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": pkg.ErrBadRequest.Error(),
		})
		return
	}

	data, err := h.TaskService.GetById(ctx, id)
	if err != nil {
		switch {
		case errors.Is(err, pkg.ErrInternalServerError):
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":    http.StatusInternalServerError,
				"message": pkg.ErrInternalServerError.Error(),
			})
		case errors.Is(err, pkg.ErrNotFound):
			c.JSON(http.StatusNotFound, gin.H{
				"code":    http.StatusNotFound,
				"message": pkg.ErrNotFound.Error(),
			})
		}
	} else {
		c.JSON(http.StatusOK, gin.H{
			"code":    http.StatusOK,
			"message": "success get data by id",
			"data":    data,
		})
	}
}
