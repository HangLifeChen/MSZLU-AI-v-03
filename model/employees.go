package model

import (
	"time"
)

type Employee struct {
	// Id: 设置为主键，类型为 bigint，让数据库自动生成
	Id int64 `json:"id" gorm:"primaryKey;autoIncrement"`
	// EmployeeNo: 员工编号，唯一，非空
	EmployeeNo string `json:"employeeNo" gorm:"type:varchar(50);uniqueIndex;not null"`
	// Name: 姓名，非空
	Name string `json:"name" gorm:"type:varchar(100);not null"`
	// Gender: 性别
	Gender *int16 `json:"gender" gorm:"type:smallint"`
	// Phone: 电话
	Phone string `json:"phone" gorm:"type:varchar(20)"`
	// Email: 邮箱
	Email string `json:"email" gorm:"type:varchar(100)"`
	// DepartmentId: 部门ID
	DepartmentId *int64 `json:"departmentId" gorm:"type:bigint"`
	// Position: 职位
	Position string `json:"position" gorm:"type:varchar(100)"`
	// HireDate: 入职日期
	HireDate *time.Time `json:"hireDate" gorm:"type:date"`
	// Status: 状态
	Status EmployeeStatusEnum `json:"status" gorm:"type:smallint;default:1"`
	// Username: 用户名，唯一
	Username string `json:"username" gorm:"type:varchar(100);uniqueIndex"`
	// Password: 密码
	Password string `json:"password" gorm:"type:varchar(255)"`
	// Remark: 备注
	Remark string `json:"remark" gorm:"type:text"`
	// CreatedAt: 创建时间
	CreatedAt time.Time `json:"createdAt" gorm:"type:timestamp;default:CURRENT_TIMESTAMP"`
	// UpdatedAt: 更新时间
	UpdatedAt time.Time `json:"updatedAt" gorm:"type:timestamp;default:CURRENT_TIMESTAMP"`
}

// TableName 指定表名
func (Employee) TableName() string {
	return "employees"
}

type EmployeeStatusEnum int16

var (
	EmployeeStatusActive   EmployeeStatusEnum = 1 // 在职
	EmployeeStatusResigned EmployeeStatusEnum = 2 // 离职
	EmployeeStatusOnLeave  EmployeeStatusEnum = 3 // 休假
)

type EmployeeDTO struct {
	Id           int64              `json:"id"`
	EmployeeNo   string             `json:"employeeNo"`
	Name         string             `json:"name"`
	Gender       *int16             `json:"gender"`
	Phone        string             `json:"phone"`
	Email        string             `json:"email"`
	DepartmentId *int64             `json:"departmentId"`
	Position     string             `json:"position"`
	HireDate     *time.Time         `json:"hireDate"`
	Status       EmployeeStatusEnum `json:"status"`
	Username     string             `json:"username"`
	Remark       string             `json:"remark"`
	CreatedAt    time.Time          `json:"createdAt"`
	UpdatedAt    time.Time          `json:"updatedAt"`
}
