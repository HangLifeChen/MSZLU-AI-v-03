package employees

import (
	"common/utils"
	"encoding/json"
	"model"
	"strconv"
)

type CreateEmployeeReq struct {
	EmployeeNo   string                   `json:"employeeNo"`
	Name         string                   `json:"name"`
	Gender       *int16                   `json:"gender"`
	Phone        string                   `json:"phone"`
	Email        string                   `json:"email"`
	DepartmentId *int64                   `json:"departmentId"`
	Position     string                   `json:"position"`
	HireDate     *utils.CustomTime        `json:"hireDate"`
	Status       model.EmployeeStatusEnum `json:"status"`
	Username     string                   `json:"username"`
	Password     string                   `json:"password"`
	Remark       string                   `json:"remark"`
}

// UnmarshalJSON 自定义 JSON 解析，支持 departmentId 为字符串或数字类型
func (r *CreateEmployeeReq) UnmarshalJSON(data []byte) error {
	// 定义临时结构体，用于解析原始 JSON
	type Alias CreateEmployeeReq
	tmp := &struct {
		*Alias
		DepartmentId interface{} `json:"departmentId"`
	}{
		Alias: (*Alias)(r),
	}

	if err := json.Unmarshal(data, &tmp); err != nil {
		return err
	}

	// 处理 departmentId 字段
	if tmp.DepartmentId != nil {
		switch v := tmp.DepartmentId.(type) {
		case float64:
			// JSON 数字默认解析为 float64
			val := int64(v)
			r.DepartmentId = &val
		case string:
			// 字符串类型，尝试转换为 int64
			if v != "" {
				val, err := strconv.ParseInt(v, 10, 64)
				if err != nil {
					return err
				}
				r.DepartmentId = &val
			}
			// nil 或其他类型保持不变
		}
	}

	return nil
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
	HireDate     *utils.CustomTime        `json:"hireDate"`
	Status       model.EmployeeStatusEnum `json:"status"`
	Username     string                   `json:"username"`
	Password     string                   `json:"password"`
	Remark       string                   `json:"remark"`
}

// UnmarshalJSON 自定义 JSON 解析，支持 departmentId 为字符串或数字类型
func (r *UpdateEmployeeReq) UnmarshalJSON(data []byte) error {
	// 定义临时结构体，用于解析原始 JSON
	type Alias UpdateEmployeeReq
	tmp := &struct {
		*Alias
		DepartmentId interface{} `json:"departmentId"`
	}{
		Alias: (*Alias)(r),
	}

	if err := json.Unmarshal(data, &tmp); err != nil {
		return err
	}

	// 处理 departmentId 字段
	if tmp.DepartmentId != nil {
		switch v := tmp.DepartmentId.(type) {
		case float64:
			// JSON 数字默认解析为 float64
			val := int64(v)
			r.DepartmentId = &val
		case string:
			// 字符串类型，尝试转换为 int64
			if v != "" {
				val, err := strconv.ParseInt(v, 10, 64)
				if err != nil {
					return err
				}
				r.DepartmentId = &val
			}
			// nil 或其他类型保持不变
		}
	}

	return nil
}

type ListEmployeesReq struct {
	EmployeeNo   string                   `json:"employeeNo" form:"employeeNo"`
	Name         string                   `json:"name" form:"name"`
	DepartmentId *int64                   `json:"departmentId" form:"departmentId"`
	Status       model.EmployeeStatusEnum `json:"status" form:"status"`
	Page         int                      `json:"page" form:"page"`
	PageSize     int                      `json:"pageSize" form:"pageSize"`
}
