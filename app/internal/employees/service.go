package employees

import (
	"common/biz"
	"common/utils"
	"context"
	"model"
	"time"

	"github.com/google/uuid"
	"github.com/mszlu521/thunder/database"
	"github.com/mszlu521/thunder/errs"
	"github.com/mszlu521/thunder/logs"
	"gorm.io/gorm/logger"
)

type service struct {
	repo repository
}

func (s *service) createEmployee(ctx context.Context, operatorID uuid.UUID, req CreateEmployeeReq) (*EmployeeResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	// 检查员工编号是否已存在
	existingEmployee, err := s.repo.getEmployeeByEmployeeNo(ctx, req.EmployeeNo)
	if err != nil {
		logs.Errorf("查询员工失败: %v", err)
		return nil, errs.DBError
	}
	if existingEmployee != nil {
		return nil, biz.ErrEmployeeNoExisted
	}

	// 检查用户名是否已存在
	if req.Username != "" {
		existingUsername, err := s.repo.getEmployeeByUsername(ctx, req.Username)
		if err != nil {
			logs.Errorf("查询员工用户名失败: %v", err)
			return nil, errs.DBError
		}
		if existingUsername != nil {
			return nil, biz.ErrEmployeeUsernameExisted
		}
	}

	employee := &model.Employee{
		EmployeeNo:   req.EmployeeNo,
		Name:         req.Name,
		Gender:       req.Gender,
		Phone:        req.Phone,
		Email:        req.Email,
		DepartmentId: req.DepartmentId,
		Position:     req.Position,
		HireDate:     convertCustomTimeToTime(req.HireDate),
		Status:       req.Status,
		Username:     req.Username,
		Password:     req.Password,
		Remark:       req.Remark,
	}

	err = s.repo.createEmployee(ctx, employee)
	if err != nil {
		logs.Errorf("创建员工失败: %v", err)
		return nil, errs.DBError
	}

	return toEmployeeResponse(employee), nil
}

func (s *service) getEmployee(ctx context.Context, id int64) (*EmployeeResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	// 检查员工是否存在
	employee, err := s.repo.getEmployee(ctx, id)
	if err != nil {
		logs.Errorf("查询员工失败: %v", err)
		return nil, errs.DBError
	}
	if employee == nil {
		return nil, biz.ErrEmployeeNotFound
	}

	return toEmployeeResponse(employee), nil
}

func (s *service) updateEmployee(ctx context.Context, operatorID uuid.UUID, req UpdateEmployeeReq) (*EmployeeResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	// 检查员工是否存在
	employee, err := s.repo.getEmployee(ctx, req.Id)
	if err != nil {
		logs.Errorf("查询员工失败: %v", err)
		return nil, errs.DBError
	}
	if employee == nil {
		return nil, biz.ErrEmployeeNotFound
	}

	// 更新员工信息
	if req.EmployeeNo != "" {
		employee.EmployeeNo = req.EmployeeNo
	}
	if req.Name != "" {
		employee.Name = req.Name
	}
	if req.Gender != nil {
		employee.Gender = req.Gender
	}
	if req.Phone != "" {
		employee.Phone = req.Phone
	}
	if req.Email != "" {
		employee.Email = req.Email
	}
	if req.DepartmentId != nil {
		employee.DepartmentId = req.DepartmentId
	}
	if req.Position != "" {
		employee.Position = req.Position
	}
	if req.HireDate != nil {
		employee.HireDate = convertCustomTimeToTime(req.HireDate)
	}

	employee.Status = req.Status

	if req.Username != "" {
		employee.Username = req.Username
	}
	if req.Password != "" {
		employee.Password = req.Password
	}
	if req.Remark != "" {
		employee.Remark = req.Remark
	}

	err = s.repo.updateEmployee(ctx, employee)
	if err != nil {
		logs.Errorf("更新员工失败: %v", err)
		return nil, errs.DBError
	}

	return toEmployeeResponse(employee), nil
}

func (s *service) deleteEmployee(ctx context.Context, operatorID uuid.UUID, id int64) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	// 检查员工是否存在
	employee, err := s.repo.getEmployee(ctx, id)
	if err != nil {
		logs.Errorf("查询员工失败: %v", err)
		return errs.DBError
	}
	if employee == nil {
		return biz.ErrEmployeeNotFound
	}

	err = s.repo.deleteEmployee(ctx, id)
	if err != nil {
		logs.Errorf("删除员工失败: %v", err)
		return errs.DBError
	}

	return nil
}

func (s *service) listEmployees(ctx context.Context, operatorID uuid.UUID, req ListEmployeesReq) (*ListEmployeeResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	filter := EmployeeFilter{
		EmployeeNo:   req.EmployeeNo,
		Name:         req.Name,
		DepartmentId: req.DepartmentId,
		Status:       req.Status,
		Page:         req.Page,
		PageSize:     req.PageSize,
	}

	if filter.Page <= 0 {
		filter.Page = 1
	}
	if filter.PageSize <= 0 {
		filter.PageSize = 10
	}

	employees, count, err := s.repo.listEmployees(ctx, filter)
	if err != nil {
		logs.Errorf("查询员工列表失败: %v", err)
		return nil, errs.DBError
	}

	var list []*EmployeeResponse
	for _, employee := range employees {
		list = append(list, toEmployeeResponse(employee))
	}

	return &ListEmployeeResponse{Employees: list, Total: count}, nil
}

func newService() *service {
	return &service{
		repo: newModels(database.GetPostgresDB().GormDB.Logger.LogMode(logger.Info)),
	}
}

// convertCustomTimeToTime 将 *utils.CustomTime 转换为 *time.Time
func convertCustomTimeToTime(ct *utils.CustomTime) *time.Time {
	if ct == nil {
		return nil
	}
	t := ct.ToTime()
	return &t
}
