package employees

import (
	"model"
	"time"
)

type EmployeeResponse struct {
	Id           int64                    `json:"id"`
	EmployeeNo   string                   `json:"employeeNo"`
	Name         string                   `json:"name"`
	Gender       *int16                   `json:"gender"`
	Phone        string                   `json:"phone"`
	Email        string                   `json:"email"`
	DepartmentId *int64                   `json:"departmentId"`
	Position     string                   `json:"position"`
	HireDate     *time.Time               `json:"hireDate"`
	Status       model.EmployeeStatusEnum `json:"status"`
	Username     string                   `json:"username"`
	Remark       string                   `json:"remark"`
	CreatedAt    time.Time                `json:"createdAt"`
	UpdatedAt    time.Time                `json:"updatedAt"`
}

type ListEmployeeResponse struct {
	Employees []*EmployeeResponse `json:"employees"`
	Total     int64               `json:"total"`
}

func toEmployeeResponse(employee *model.Employee) *EmployeeResponse {
	return &EmployeeResponse{
		Id:           employee.Id,
		EmployeeNo:   employee.EmployeeNo,
		Name:         employee.Name,
		Gender:       employee.Gender,
		Phone:        employee.Phone,
		Email:        employee.Email,
		DepartmentId: employee.DepartmentId,
		Position:     employee.Position,
		HireDate:     employee.HireDate,
		Status:       employee.Status,
		Username:     employee.Username,
		Remark:       employee.Remark,
		CreatedAt:    employee.CreatedAt,
		UpdatedAt:    employee.UpdatedAt,
	}
}
