package employees

import (
	"model"
	"time"
)

type CreateEmployeeReq struct {
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
	Password     string                   `json:"password"`
	Remark       string                   `json:"remark"`
}

type UpdateEmployeeReq struct {
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
	Password     string                   `json:"password"`
	Remark       string                   `json:"remark"`
}

type ListEmployeesReq struct {
	EmployeeNo   string                   `json:"employeeNo" form:"employeeNo"`
	Name         string                   `json:"name" form:"name"`
	DepartmentId *int64                   `json:"departmentId" form:"departmentId"`
	Status       model.EmployeeStatusEnum `json:"status" form:"status"`
	Page         int                      `json:"page" form:"page"`
	PageSize     int                      `json:"pageSize" form:"pageSize"`
}
