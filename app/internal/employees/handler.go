package employees

import (
	"github.com/gin-gonic/gin"
	"github.com/mszlu521/thunder/req"
	"github.com/mszlu521/thunder/res"
)

type Handler struct {
	service *service
}

func (h *Handler) CreateEmployee(c *gin.Context) {
	var createReq CreateEmployeeReq
	if err := req.JsonParam(c, &createReq); err != nil {
		return
	}
	userID, ok := req.GetUserIdUUID(c)
	if !ok {
		return
	}
	resp, err := h.service.createEmployee(c.Request.Context(), userID, createReq)
	if err != nil {
		res.Error(c, err)
		return
	}
	res.Success(c, resp)
}

func (h *Handler) GetEmployee(c *gin.Context) {
	var id int64
	if err := req.Path(c, "id", &id); err != nil {
		return
	}
	resp, err := h.service.getEmployee(c.Request.Context(), id)
	if err != nil {
		res.Error(c, err)
		return
	}
	res.Success(c, resp)
}

func (h *Handler) UpdateEmployee(c *gin.Context) {
	var updateReq UpdateEmployeeReq
	if err := req.JsonParam(c, &updateReq); err != nil {
		return
	}
	userID, ok := req.GetUserIdUUID(c)
	if !ok {
		return
	}
	resp, err := h.service.updateEmployee(c.Request.Context(), userID, updateReq)
	if err != nil {
		res.Error(c, err)
		return
	}
	res.Success(c, resp)
}

func (h *Handler) DeleteEmployee(c *gin.Context) {
	var id int64
	if err := req.Path(c, "id", &id); err != nil {
		return
	}
	userID, ok := req.GetUserIdUUID(c)
	if !ok {
		return
	}
	err := h.service.deleteEmployee(c.Request.Context(), userID, id)
	if err != nil {
		res.Error(c, err)
		return
	}
	res.Success(c, nil)
}

func (h *Handler) ListEmployees(c *gin.Context) {
	var listReq ListEmployeesReq
	if err := req.QueryParam(c, &listReq); err != nil {
		return
	}
	userID, ok := req.GetUserIdUUID(c)
	if !ok {
		return
	}
	resp, err := h.service.listEmployees(c.Request.Context(), userID, listReq)
	if err != nil {
		res.Error(c, err)
		return
	}
	res.Success(c, resp)
}

func NewHandler() *Handler {
	return &Handler{
		service: newService(),
	}
}
