package users

import (
	"common/biz"
	"context"

	"core/upload"
	"fmt"
	"mime/multipart"
	"model"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/mszlu521/thunder/database"
	"github.com/mszlu521/thunder/errs"
	"github.com/mszlu521/thunder/logs"
	// "github.com/mszlu521/thunder/upload"
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

func (s *service) uploadAvatar(ctx context.Context, userID uuid.UUID, file *multipart.FileHeader) (*UploadAvatarResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	// 检查文件大小（限制为5MB）
	if file.Size > 5*1024*1024 {
		return nil, fmt.Errorf("文件大小不能超过5MB")
	}

	// 检查文件类型
	ext := strings.ToLower(filepath.Ext(file.Filename))
	allowedExts := map[string]bool{
		".jpg":  true,
		".jpeg": true,
		".png":  true,
		".gif":  true,
		".webp": true,
	}
	if !allowedExts[ext] {
		return nil, fmt.Errorf("不支持的文件类型，仅支持 jpg, jpeg, png, gif, webp 格式")
	}

	// 打开文件
	src, err := file.Open()
	if err != nil {
		logs.Errorf("打开文件失败: %v", err)
		return nil, fmt.Errorf("打开文件失败")
	}
	defer src.Close()

	// 生成文件名
	filename := fmt.Sprintf("avatar/%s/%s%s", userID.String(), uuid.New().String(), ext)

	// 上传到阿里云OSS
	err = upload.AliyunOSSUpload.Upload(ctx, src, filename)
	if err != nil {
		logs.Errorf("上传文件失败: %v", err)
		return nil, fmt.Errorf("上传文件失败")
	}

	// 获取公开访问URL
	url := upload.AliyunOSSUpload.GetPublicUrl(filename)

	// 更新用户头像
	user, err := s.repo.getUser(ctx, userID)
	if err != nil {
		logs.Errorf("查询用户失败: %v", err)
		return nil, errs.DBError
	}
	if user == nil {
		return nil, biz.ErrUserNotFound
	}

	user.Avatar = url
	err = s.repo.updateUser(ctx, user)
	if err != nil {
		logs.Errorf("更新用户头像失败: %v", err)
		return nil, errs.DBError
	}

	return &UploadAvatarResponse{
		URL: url,
	}, nil
}
