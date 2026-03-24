package users

import (
	"model"

	"github.com/google/uuid"
)

type UserResponse struct {
	Id              uuid.UUID        `json:"id"`
	Username        string           `json:"username"`
	Avatar          string           `json:"avatar"`
	Status          model.StatusEnum `json:"status"`
	LastLoginTime   int64            `json:"lastLoginTime"`
	CurrentPlan     string           `json:"currentPlan"`
	Email           string           `json:"email"`
	EmailVerified   bool             `json:"emailVerified"`
	Introduction    string           `json:"introduction"`
	TelephoneNumber string           `json:"telephoneNumber"`
}

type ListUserResponse struct {
	Users []*UserResponse `json:"users"`
	Total int64           `json:"total"`
}

type UploadAvatarResponse struct {
	URL string `json:"url"`
}

func toUserResponse(user *model.User) *UserResponse {
	return &UserResponse{
		Id:              user.Id,
		Username:        user.Username,
		Avatar:          user.Avatar,
		Status:          user.Status,
		LastLoginTime:   user.LastLoginTime.UnixMilli(),
		CurrentPlan:     string(user.CurrentPlan),
		Email:           user.Email,
		EmailVerified:   user.EmailVerified,
		Introduction:    user.Introduction,
		TelephoneNumber: user.TelephoneNumber,
	}
}
