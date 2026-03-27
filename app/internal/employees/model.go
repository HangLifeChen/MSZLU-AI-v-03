package employees

import (
	"context"
	"model"

	"github.com/mszlu521/thunder/gorms"
	"gorm.io/gorm"
)

type models struct {
	db *gorm.DB
}

func (m *models) createEmployee(ctx context.Context, employee *model.Employee) error {
	return m.db.WithContext(ctx).Create(employee).Error
}

func (m *models) getEmployee(ctx context.Context, id int64) (*model.Employee, error) {
	var employee model.Employee
	err := m.db.WithContext(ctx).Where("id = ?", id).First(&employee).Error
	if gorms.IsRecordNotFoundError(err) {
		return nil, nil
	}
	return &employee, nil
}

func (m *models) updateEmployee(ctx context.Context, employee *model.Employee) error {
	return m.db.WithContext(ctx).Select("status").Updates(employee).Error
}

func (m *models) deleteEmployee(ctx context.Context, id int64) error {
	return m.db.WithContext(ctx).Where("id = ?", id).Unscoped().Delete(&model.Employee{}).Error
}

func (m *models) listEmployees(ctx context.Context, filter EmployeeFilter) ([]*model.Employee, int64, error) {
	var employees []*model.Employee
	var count int64
	query := m.db.WithContext(ctx).Model(&model.Employee{})
	if filter.EmployeeNo != "" {
		query = query.Where("employee_no like ?", "%"+filter.EmployeeNo+"%")
	}
	if filter.Name != "" {
		query = query.Where("name like ?", "%"+filter.Name+"%")
	}
	if filter.DepartmentId != nil {
		query = query.Where("department_id = ?", *filter.DepartmentId)
	}
	if filter.Status != 0 {
		query = query.Where("status = ?", filter.Status)
	}
	// 获取总数
	countQuery := m.db.WithContext(ctx).Model(&model.Employee{})
	if filter.EmployeeNo != "" {
		countQuery = countQuery.Where("employee_no like ?", "%"+filter.EmployeeNo+"%")
	}
	if filter.Name != "" {
		countQuery = countQuery.Where("name like ?", "%"+filter.Name+"%")
	}
	if filter.DepartmentId != nil {
		countQuery = countQuery.Where("department_id = ?", *filter.DepartmentId)
	}
	if filter.Status != 0 {
		countQuery = countQuery.Where("status = ?", filter.Status)
	}
	if err := countQuery.Count(&count).Error; err != nil {
		return nil, 0, err
	}
	if filter.Page > 0 && filter.PageSize > 0 {
		query = query.Limit(filter.PageSize).Offset((filter.Page - 1) * filter.PageSize)
	}
	if err := query.Find(&employees).Error; err != nil {
		return nil, 0, err
	}
	return employees, count, nil
}

func (m *models) getEmployeeByEmployeeNo(ctx context.Context, employeeNo string) (*model.Employee, error) {
	var employee model.Employee
	err := m.db.WithContext(ctx).Where("employee_no = ?", employeeNo).First(&employee).Error
	if gorms.IsRecordNotFoundError(err) {
		return nil, nil
	}
	return &employee, nil
}

func (m *models) getEmployeeByUsername(ctx context.Context, username string) (*model.Employee, error) {
	var employee model.Employee
	err := m.db.WithContext(ctx).Where("username = ?", username).First(&employee).Error
	if gorms.IsRecordNotFoundError(err) {
		return nil, nil
	}
	return &employee, nil
}

func newModels(db *gorm.DB) *models {
	return &models{
		db: db,
	}
}
