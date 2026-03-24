package users

import (
	"model"

	"github.com/google/uuid"
)

type CreateUserReq struct {
	Username        string           `json:"username"`
	Password        string           `json:"password"`
	Avatar          string           `json:"avatar"`
	Email           string           `json:"email"`
	Introduction    string           `json:"introduction"`
	TelephoneNumber string           `json:"telephone_number"`
	Status          model.StatusEnum `json:"status"`
}

type UpdateUserReq struct {
	ID              uuid.UUID        `json:"id"`
	Username        string           `json:"username"`
	Password        string           `json:"password"`
	Avatar          string           `json:"avatar"`
	Email           string           `json:"email"`
	Introduction    string           `json:"introduction"`
	TelephoneNumber string           `json:"telephone_number"`
	Status          model.StatusEnum `json:"status"`
}

type ListUsersReq struct {
	Username string           `json:"username" form:"username"`
	Email    string           `json:"email" form:"email"`
	Status   model.StatusEnum `json:"status" form:"status"`
	Page     int              `json:"page" form:"page"`
	PageSize int              `json:"pageSize" form:"pageSize"`
}
