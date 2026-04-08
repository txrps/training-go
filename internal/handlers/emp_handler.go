package handlers

import (
	"errors"
	"net/http"
	"training-go/internal/models"
	"training-go/internal/repository"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type CreateEmployeeInput struct {
	FirstName    string `json:"firstname" binding:"required"`
	LastName     string `json:"lastname" binding:"required"`
	DepartmentID int64  `json:"department_id"`
}

type UpdateEmpInput struct {
	ID           int64   `json:"id"`
	FirstName    *string `json:"firstname"`
	LastName     *string `json:"lastname"`
	DepartmentID *int64  `json:"department_id"`
}

type SearchEmpInput struct {
	FirstName      *string `json:"firstname"`
	LastName       *string `json:"lastname"`
	DepartmentName *string `json:"name"`
}
type SearchEmpInDepartmentInput struct {
	DepartmentName *string `json:"name"`
}

// CreateEmp godoc
// @Summary      Create emp
// @Description  Create a new emp
// @Tags         employees
// @Accept       json
// @Produce      json
// @Param        create_emp  body     handlers.CreateEmployeeInput  true  "Employee payload"
// @Success      201   {object} map[string]interface{}
// @Failure      400   {object} map[string]string
// @Failure      500   {object} map[string]string
// @Router       /create_employee [post]
func CreateEmpHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var input CreateEmployeeInput

		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		params := models.CreateEmpParams{
			FirstName:    input.FirstName,
			LastName:     input.LastName,
			DepartmentID: input.DepartmentID,
		}

		todo, err := repository.CreateEmpRepo(db, params)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, todo)
	}
}

// UpdateEmp godoc
// @Summary      Update emp
// @Description  Update a emp by ID
// @Tags         employees
// @Accept       json
// @Produce      json
// @Param        update_emp  body      handlers.UpdateEmpInput  true  "Update payload"
// @Success      200   {object}  map[string]interface{}
// @Failure      400   {object}  map[string]string
// @Failure      404   {object}  map[string]string
// @Failure      500   {object}  map[string]string
// @Router       /update_employee [put]
func UpdateEmpHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var input UpdateEmpInput

		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if input.FirstName == nil && input.LastName == nil && input.DepartmentID == nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "at least one field must be provided"})
			return
		}

		if input.ID == 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "id is required"})
			return
		}

		params := models.UpdateEmpParams{
			ID:           input.ID,
			FirstName:    input.FirstName,
			LastName:     input.LastName,
			DepartmentID: input.DepartmentID,
		}

		emp, err := repository.UpdateEmpRepo(db, params)
		if err != nil {
			if errors.Is(err, repository.ErrTodoNotFound) {
				c.JSON(http.StatusNotFound, gin.H{"error": ErrNotFound})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": ErrInternalServerError})
			return
		}

		c.JSON(http.StatusOK, gin.H{"status": "success", "data": emp})
	}
}

// SearchEmp godoc
// @Summary      Search emp
// @Description  Search a emp
// @Tags         employees
// @Accept       json
// @Produce      json
// @Param        search_emp  body     handlers.SearchEmpInput  true  "Employee payload"
// @Success      200   {object} map[string]interface{}
// @Failure      400   {object} map[string]string
// @Failure      500   {object} map[string]string
// @Router       /search_employee [post]
func SearchEmpHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var input SearchEmpInput

		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		// map input → params
		params := models.SearchEmpParams{
			FirstName:      input.FirstName,
			LastName:       input.LastName,
			DepartmentName: input.DepartmentName,
		}

		empls, err := repository.SearchEmpRepo(db, params)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"status": "success",
			"data":   empls,
		})
	}
}

// SearchEmpInDepartment godoc
// @Summary      Count emp in department
// @Description  Get employee count per department with optional filter
// @Tags         employees
// @Accept       json
// @Produce      json
// @Param        search_emp  body     handlers.SearchEmpInDepartmentInput  false  "Department filter"
// @Success      200   {object} map[string]interface{}
// @Failure      404   {object} map[string]string
// @Failure      500   {object} map[string]string
// @Router       /employee_count [post]
func CountEmpInDepartmentHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var input SearchEmpInDepartmentInput

		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		results, err := repository.CountEmpInDepartmentRepo(db, input.DepartmentName)
		if err != nil {
			if errors.Is(err, repository.ErrTodoNotFound) {
				c.JSON(http.StatusNotFound, gin.H{
					"error": "department not found",
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
			"data":   results,
		})
	}
}
