package handlers

import (
	"errors"
	"net/http"
	"strconv"
	"training-go/internal/repository"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type CreateTodoInput struct {
	Title     string `json:"title" binding:"required"`
	Completed bool   `json:"completed"`
}

type UpdateTodoInput struct {
	Title     *string `json:"title"`
	Completed *bool   `json:"completed"`
}

const ErrInvalidID = "invalid id"
const ErrNotFound = "todo not found"
const ErrInternalServerError = "internal server error"

// CreateTodo godoc
// @Summary      Create todo
// @Description  Create a new todo
// @Tags         todos
// @Accept       json
// @Produce      json
// @Param        todo  body     handlers.CreateTodoInput  true  "Todo payload"
// @Success      201   {object} map[string]interface{}
// @Failure      400   {object} map[string]string
// @Failure      500   {object} map[string]string
// @Router       /todos [post]
func CreateTodoHandlerGorm(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var input CreateTodoInput

		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		todo, err := repository.CreateTodoGorm(db, input.Title, input.Completed)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(http.StatusCreated, todo)
	}
}

// GetAllTodos godoc
// @Summary      Get all todos
// @Description  Retrieve all todo items
// @Tags         todos
// @Produce      json
// @Success      200   {object} map[string]interface{}
// @Failure      500   {object} map[string]string
// @Router       /todos [get]
func GetAllTodosHandlerGorm(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {

		todos, err := repository.GetAllTodosGorm(db)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": ErrInternalServerError,
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"status": "success",
			"data":   todos,
		})
	}
}

// GetTodoByID godoc
// @Summary      Get todo by ID
// @Description  Retrieve a single todo by ID
// @Tags         todos
// @Produce      json
// @Param        id   path      int  true  "Todo ID"
// @Success      200  {object}  map[string]interface{}
// @Failure      400  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /todos/{id} [get]
func GetTodoByIDHandlerGorm(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {

		idStr := c.Param("id")
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": ErrInvalidID,
			})
			return
		}

		todo, err := repository.GetTodoByIDGorm(db, id)
		if err != nil {
			if errors.Is(err, repository.ErrTodoNotFound) {
				c.JSON(http.StatusNotFound, gin.H{
					"error": ErrNotFound,
				})
				return
			}

			c.JSON(http.StatusInternalServerError, gin.H{
				"error": ErrInternalServerError,
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"status": "success",
			"data":   todo,
		})
	}
}

// UpdateTodo godoc
// @Summary      Update todo
// @Description  Update a todo by ID
// @Tags         todos
// @Accept       json
// @Produce      json
// @Param        id    path      int                        true  "Todo ID"
// @Param        todo  body      handlers.UpdateTodoInput  true  "Update payload"
// @Success      200   {object}  map[string]interface{}
// @Failure      400   {object}  map[string]string
// @Failure      404   {object}  map[string]string
// @Failure      500   {object}  map[string]string
// @Router       /todos/{id} [put]
func UpdateTodoHandlerGorm(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {

		idStr := c.Param("id")
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": ErrInvalidID,
			})
			return
		}

		var input UpdateTodoInput
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		if input.Title == nil && input.Completed == nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "at least one field must be provided",
			})
			return
		}

		todo, err := repository.UpdateTodoGorm(
			db,
			id,
			input.Title,
			input.Completed,
		)

		if err != nil {
			if errors.Is(err, repository.ErrTodoNotFound) {
				c.JSON(http.StatusNotFound, gin.H{
					"error": ErrNotFound,
				})
				return
			}

			c.JSON(http.StatusInternalServerError, gin.H{
				"error": ErrInternalServerError,
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"status": "success",
			"data":   todo,
		})
	}
}

// DeleteTodo godoc
// @Summary      Delete todo
// @Description  Delete a todo by ID
// @Tags         todos
// @Produce      json
// @Param        id   path      int  true  "Todo ID"
// @Success      200  {object}  map[string]string
// @Failure      400  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /todos/{id} [delete]
func DeleteTodoHandlerGorm(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {

		idStr := c.Param("id")
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": ErrInvalidID,
			})
			return
		}

		err = repository.DeleteTodoGorm(db, id)
		if err != nil {

			if errors.Is(err, repository.ErrTodoNotFound) {
				c.JSON(http.StatusNotFound, gin.H{
					"error": ErrNotFound,
				})
				return
			}

			c.JSON(http.StatusInternalServerError, gin.H{
				"error": ErrInternalServerError,
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"status":  "success",
			"message": "todo deleted successfully",
		})
	}
}
