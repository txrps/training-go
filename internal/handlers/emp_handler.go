package handlers

import (
	"errors"
	"net/http"
	"training-go/internal/repository"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type CreateEmployeeInput struct {
	EmpID        *int64 `json:"emp_id,omitempty"`
	FirstName    string `json:"firstname" binding:"required"`
	LastName     string `json:"lastname" binding:"required"`
	DepartmentID int64  `json:"department_id"`
}

// CreateEmployeeHandler godoc
// @Summary      Create employee
// @Description  Create a new employee
// @Tags         employees
// @Accept       json
// @Produce      json
// @Param        employee  body     handlers.CreateEmployeeInput  true  "Employee payload"
// @Success      201   {object} map[string]interface{}
// @Failure      400   {object} map[string]string
// @Failure      500   {object} map[string]string
// @Router       /employees [post]
func SaveEmployeeHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var input CreateEmployeeInput

		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		employee, err := repository.SaveEmployee(db, input.EmpID, input.FirstName, input.LastName, input.DepartmentID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(http.StatusCreated, employee)
	}
}

// GetEmployeeByFilter godoc
// @Summary      Get employee by filter
// @Description  Retrieve a single employee by filter
// @Tags         employees
// @Produce      json
// @Param filter query string false "filter keyword"
// @Success      200  {object}  map[string]interface{}
// @Failure      400  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /employees [get]
func GetEmployeeByFilterHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {

		filterStr := c.Query("filter")

		employee, err := repository.GetAllEmployeeByFilter(db, filterStr)
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
			"data":   employee,
		})
	}
}

// GetEmployeeCountByDepartment godoc
// @Summary      Get employee count by department
// @Description  Retrieve employee count by department with optional filter
// @Tags         employees
// @Produce      json
// @Param filter query string false "filter keyword"
// @Success      200  {object}  map[string]interface{}
// @Failure      400  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /employees/department [get]
func GetEmployeeCountByDepartmentHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {

		filterStr := c.Query("filter")

		departmentCounts, err := repository.GetEmployeeCountByDepartment(db, filterStr)
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
			"data":   departmentCounts,
		})
	}
}
