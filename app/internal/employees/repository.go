package employees

import (
	"context"
	"model"
)

type repository interface {
	createEmployee(ctx context.Context, employee *model.Employee) error
	getEmployee(ctx context.Context, id int64) (*model.Employee, error)
	updateEmployee(ctx context.Context, employee *model.Employee) error
	deleteEmployee(ctx context.Context, id int64) error
	listEmployees(ctx context.Context, filter EmployeeFilter) ([]*model.Employee, int64, error)
	getEmployeeByEmployeeNo(ctx context.Context, employeeNo string) (*model.Employee, error)
	getEmployeeByUsername(ctx context.Context, username string) (*model.Employee, error)
}

type EmployeeFilter struct {
	EmployeeNo   string
	Name         string
	DepartmentId *int64
	Status       model.EmployeeStatusEnum
	Page         int
	PageSize     int
}
