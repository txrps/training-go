package handlers

type CreateEmployeeInput struct {
	FirstName    string `json:"firstname" binding:"required"`
	LastName     string `json:"lastname" binding:"required"`
	DepartmentID *int64 `json:"department_id"`
}
