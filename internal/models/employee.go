package models

import "time"

type Department struct {
	ID        int64  `json:"id" db:"id" gorm:"primaryKey;autoIncrement"`
	Name      string `json:"name" db:"name" gorm:"column:name;type:varchar(100);not null"`
	WorkFloor int    `json:"work_floor" db:"work_floor" gorm:"column:work_floor;not null"`

	Employees []Employee `json:"employees,omitempty" gorm:"foreignKey:DepartmentID"`
}

type Employee struct {
	ID           int64     `json:"id" db:"id" gorm:"primaryKey;autoIncrement"`
	FirstName    string    `json:"firstname" db:"firstname" gorm:"column:firstname;type:varchar(100);not null"`
	LastName     string    `json:"lastname" db:"lastname" gorm:"column:lastname;type:varchar(100);not null"`
	DepartmentID int64     `json:"department_id" db:"department_id" gorm:"column:department_id"`
	CreatedAt    time.Time `json:"created_at" db:"created_at" gorm:"autoCreateTime"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at" gorm:"autoUpdateTime"`

	Department *Department `json:"department,omitempty" gorm:"foreignKey:DepartmentID"`
}
type CreateEmpParams struct {
	FirstName    string
	LastName     string
	DepartmentID int64
}

type UpdateEmpParams struct {
	ID           int64
	FirstName    *string
	LastName     *string
	DepartmentID *int64
}
type SearchEmpParams struct {
	FirstName      *string
	LastName       *string
	DepartmentName *string
}

type DepartmentEmpCount struct {
	DepartmentName string `json:"department_name"`
	EmpCount       int64  `json:"emp_count"`
}
