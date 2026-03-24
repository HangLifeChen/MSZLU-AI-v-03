package users

import (
	"common/biz"
	"context"
	"model"
	"time"

	"github.com/google/uuid"
	"github.com/mszlu521/thunder/database"
	"github.com/mszlu521/thunder/errs"
	"github.com/mszlu521/thunder/logs"
)

type service struct {
	repo repository
}

func (s *service) createUser(ctx context.Context, operatorID uuid.UUID, req CreateUserReq) (*UserResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	// 检查用户名是否已存在
	existingUser, err := s.repo.getUserByUsername(ctx, req.Username)
	if err != nil {
		logs.Errorf("查询用户失败: %v", err)
		return nil, errs.DBError
	}
	if existingUser != nil {
		return nil, biz.ErrUserNameExisted
	}

	// 检查邮箱是否已存在
	existingEmail, err := s.repo.getUserByEmail(ctx, req.Email)
	if err != nil {
		logs.Errorf("查询用户邮箱失败: %v", err)
		return nil, errs.DBError
	}
	if existingEmail != nil {
		return nil, biz.ErrEmailExisted
	}

	user := &model.User{
		Id:              uuid.New(),
		Username:        req.Username,
		Password:        req.Password,
		Avatar:          req.Avatar,
		Email:           req.Email,
		Introduction:    req.Introduction,
		TelephoneNumber: req.TelephoneNumber,
		Status:          req.Status,
		CurrentPlan:     model.FreePlan,
	}

	err = s.repo.createUser(ctx, user)
	if err != nil {
		logs.Errorf("创建用户失败: %v", err)
		return nil, errs.DBError
	}

	return toUserResponse(user), nil
}

func (s *service) getUser(ctx context.Context, id uuid.UUID) (*UserResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	// 检查用户是否存在
	user, err := s.repo.getUser(ctx, id)
	if err != nil {
		logs.Errorf("查询用户失败: %v", err)
		return nil, errs.DBError
	}
	if user == nil {
		return nil, biz.ErrUserNotFound
	}

	return toUserResponse(user), nil
}

func (s *service) updateUser(ctx context.Context, operatorID uuid.UUID, req UpdateUserReq) (*UserResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	// 检查用户是否存在
	user, err := s.repo.getUser(ctx, req.ID)
	if err != nil {
		logs.Errorf("查询用户失败: %v", err)
		return nil, errs.DBError
	}
	if user == nil {
		return nil, biz.ErrUserNotFound
	}

	// 更新用户信息
	if req.Username != "" {
		user.Username = req.Username
	}
	if req.Password != "" {
		user.Password = req.Password
	}
	if req.Avatar != "" {
		user.Avatar = req.Avatar
	}
	if req.Email != "" {
		user.Email = req.Email
	}
	if req.Introduction != "" {
		user.Introduction = req.Introduction
	}
	if req.TelephoneNumber != "" {
		user.TelephoneNumber = req.TelephoneNumber
	}
	if req.Status != 0 {
		user.Status = req.Status
	}

	err = s.repo.updateUser(ctx, user)
	if err != nil {
		logs.Errorf("更新用户失败: %v", err)
		return nil, errs.DBError
	}

	return toUserResponse(user), nil
}

func (s *service) deleteUser(ctx context.Context, operatorID uuid.UUID, id uuid.UUID) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	// 检查用户是否存在
	user, err := s.repo.getUser(ctx, id)
	if err != nil {
		logs.Errorf("查询用户失败: %v", err)
		return errs.DBError
	}
	if user == nil {
		return biz.ErrUserNotFound
	}

	err = s.repo.deleteUser(ctx, id)
	if err != nil {
		logs.Errorf("删除用户失败: %v", err)
		return errs.DBError
	}

	return nil
}

func (s *service) listUsers(ctx context.Context, operatorID uuid.UUID, req ListUsersReq) (*ListUserResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	filter := UserFilter{
		Username: req.Username,
		Email:    req.Email,
		Status:   req.Status,
		Page:     req.Page,
		PageSize: req.PageSize,
	}

	if filter.Page <= 0 {
		filter.Page = 1
	}
	if filter.PageSize <= 0 {
		filter.PageSize = 10
	}

	users, count, err := s.repo.listUsers(ctx, filter)
	if err != nil {
		logs.Errorf("查询用户列表失败: %v", err)
		return nil, errs.DBError
	}

	var list []*UserResponse
	for _, user := range users {
		list = append(list, toUserResponse(user))
	}

	return &ListUserResponse{Users: list, Total: count}, nil
}

func newService() *service {
	return &service{
		repo: newModels(database.GetPostgresDB().GormDB),
	}
}
