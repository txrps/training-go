package repository

import (
	"context"
	"errors"
	"strings"
	"time"
	"training-go/internal/models"

	"gorm.io/gorm"
)

// 1. สร้าง service สำหรับเพิ่ม/แก้ไข employee
// 2. สร้าง service สำหรับค้นหาพนักงานด้วยชื่อนามสกุลหรือแผนก
// 3. สร้าง service สำหรับหาจำนวนพนักงานในแต่ละแผนก โดย filter แผนกได้

func CreateEmpRepo(db *gorm.DB, params models.CreateEmpParams) (*models.Employee, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	emp := &models.Employee{
		FirstName:    params.FirstName,
		LastName:     params.LastName,
		DepartmentID: params.DepartmentID,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	if err := db.WithContext(ctx).Create(emp).Error; err != nil {
		return nil, err
	}

	return emp, nil
}

func UpdateEmpRepo(db *gorm.DB, params models.UpdateEmpParams) (*models.Employee, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	updates := map[string]interface{}{}

	if params.FirstName != nil {
		updates["firstname"] = *params.FirstName
	}

	if params.LastName != nil {
		updates["lastname"] = *params.LastName
	}

	if params.DepartmentID != nil {
		updates["department_id"] = *params.DepartmentID
	}

	if len(updates) == 0 {
		return nil, errors.New("no fields to update")
	}

	updates["updated_at"] = time.Now()

	var emp models.Employee

	result := db.WithContext(ctx).
		Model(&emp).
		Where("id = ?", params.ID).
		Updates(updates)

	if result.Error != nil {
		return nil, result.Error
	}

	if result.RowsAffected == 0 {
		return nil, ErrTodoNotFound
	}

	if err := db.WithContext(ctx).First(&emp, params.ID).Error; err != nil {
		return nil, err
	}

	return &emp, nil
}

func SearchEmpRepo(db *gorm.DB, params models.SearchEmpParams) ([]models.Employee, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var employees []models.Employee
	query := db.WithContext(ctx).Model(&models.Employee{}).Preload("Department")

	if params.FirstName != nil && strings.TrimSpace(*params.FirstName) != "" {
		search := "%" + strings.ToLower(strings.TrimSpace(*params.FirstName)) + "%"
		query = query.Where("LOWER(firstname) LIKE ?", search)
	}

	if params.LastName != nil && strings.TrimSpace(*params.LastName) != "" {
		search := "%" + strings.ToLower(strings.TrimSpace(*params.LastName)) + "%"
		query = query.Where("LOWER(lastname) LIKE ?", search)
	}

	if params.DepartmentName != nil && strings.TrimSpace(*params.DepartmentName) != "" {
		search := "%" + strings.ToLower(strings.TrimSpace(*params.DepartmentName)) + "%"
		query = query.Joins("LEFT JOIN departments ON employees.department_id = departments.id").
			Where("LOWER(departments.name) LIKE ?", search)
	}

	result := query.Find(&employees)
	err := result.Error
	if err != nil {
		return nil, err
	}

	return employees, nil
}

func CountEmpInDepartmentRepo(db *gorm.DB, department_name *string) ([]models.DepartmentEmpCount, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var results []models.DepartmentEmpCount

	query := db.WithContext(ctx).
		Model(&models.Department{}).
		Select("departments.name as department_name, COUNT(employees.id) as emp_count").
		Joins("LEFT JOIN employees ON employees.department_id = departments.id").
		Group("departments.id, departments.name")

	if department_name != nil && strings.TrimSpace(*department_name) != "" {
		search := "%" + strings.ToLower(strings.TrimSpace(*department_name)) + "%"
		query = query.Where("LOWER(departments.name) LIKE ?", search)
	}

	if err := query.Scan(&results).Error; err != nil {
		return nil, err
	}

	return results, nil
}
