package repository

import (
	"context"
	"log"
	"strings"
	"time"
	"training-go/internal/models"

	"gorm.io/gorm"
)

// 1. สร้าง service สำหรับเพิ่ม/แก้ไข employee
// 2. สร้าง service สำหรับค้นหาพนักงานด้วยชื่อนามสกุลหรือแผนก
// 3. สร้าง service สำหรับหาจำนวนพนักงานในแต่ละแผนก โดย filter แผนกได้

// CREATE TABLE departments (
//     id SERIAL PRIMARY KEY,
//     name VARCHAR(100) NOT NULL,
//     work_floor INT NOT NULL
// );

// CREATE TABLE employees (
//     id SERIAL PRIMARY KEY,
//     firstname VARCHAR(100) NOT NULL,
//     lastname VARCHAR(100) NOT NULL,
// 	created_at TIMESTAMPTZ DEFAULT NOW(),
//     updated_at TIMESTAMPTZ DEFAULT NOW(),
//     department_id INT,
//     CONSTRAINT fk_departments
//         FOREIGN KEY(department_id)
//         REFERENCES departments(id)
//         ON DELETE SET NULL
// );

// INSERT INTO departments (name, work_floor) VALUES
// ('Developer', 1),
// ('CO', 2),
// ('Admin', 3),
// ('Tester', 4),
// ('SA', 5),
// ('UX/UI', 2);

func SaveEmployee(db *gorm.DB, empID *int64, firstName string, lastName string, departmentID int64) (*models.Employee, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var employee models.Employee

	if empID != nil {
		if err := db.WithContext(ctx).First(&employee, *empID).Error; err != nil {
			return nil, err
		}

		employee.FirstName = firstName
		employee.LastName = lastName
		employee.DepartmentID = departmentID

		if err := db.WithContext(ctx).Save(&employee).Error; err != nil {
			return nil, err
		}

		return &employee, nil
	}

	employee = models.Employee{
		FirstName:    firstName,
		LastName:     lastName,
		DepartmentID: departmentID,
	}

	if err := db.WithContext(ctx).Create(&employee).Error; err != nil {
		return nil, err
	}

	return &employee, nil
}

func GetAllEmployeeByFilter(db *gorm.DB, filter string) ([]models.Employee, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	log.Printf("🚀 filter is %s", filter)

	var employees []models.Employee

	query := db.WithContext(ctx).
		Model(&models.Employee{}).
		Joins("LEFT JOIN departments ON departments.id = employees.department_id")

	if strings.TrimSpace(filter) != "" {
		like := "%" + filter + "%"

		query = query.Where(`
			employees.firstname ILIKE ? OR
			employees.lastname ILIKE ? OR
			departments.name ILIKE ?
		`, like, like, like)
	}

	if err := query.
		Preload("Department").
		Order("employees.id DESC").
		Find(&employees).Error; err != nil {
		return nil, err
	}

	return employees, nil
}

type EmployeeCountByDepartment struct {
	DepartmentID   int64  `json:"department_id"`
	DepartmentName string `json:"department_name"`
	TotalEmployees int64  `json:"total_employees"`
}

func GetEmployeeCountByDepartment(db *gorm.DB, filter string) ([]EmployeeCountByDepartment, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var results []EmployeeCountByDepartment

	query := db.WithContext(ctx).
		Table("departments").
		Select(`
			departments.id as department_id,
			departments.name as department_name,
			COUNT(employees.id) as total_employees
		`).
		Joins("LEFT JOIN employees ON employees.department_id = departments.id").
		Group("departments.id, departments.name")

	if strings.TrimSpace(filter) != "" {
		like := "%" + filter + "%"

		query = query.Where(`
			departments.name ILIKE ?
		`, like)
	}

	if err := query.
		Order("departments.id ASC").
		Scan(&results).Error; err != nil {
		return nil, err
	}

	return results, nil
}
