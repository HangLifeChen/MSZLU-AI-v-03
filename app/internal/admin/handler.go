package admin

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/mszlu521/thunder/res"
)

type Handler struct {
	service *service
}

func NewHandler() *Handler {
	return &Handler{
		service: newService(),
	}
}

func (h *Handler) GetUserGrowthTrend(c *gin.Context) {
	year := getYearFromQuery(c)
	resp, err := h.service.getUserGrowthTrend(c.Request.Context(), year)
	if err != nil {
		res.Error(c, err)
		return
	}
	res.Success(c, resp)
}

func (h *Handler) GetKnowledgeBaseGrowthTrend(c *gin.Context) {
	year := getYearFromQuery(c)
	resp, err := h.service.getKnowledgeBaseGrowthTrend(c.Request.Context(), year)
	if err != nil {
		res.Error(c, err)
		return
	}
	res.Success(c, resp)
}

func (h *Handler) GetAgentGrowthTrend(c *gin.Context) {
	year := getYearFromQuery(c)
	resp, err := h.service.getAgentGrowthTrend(c.Request.Context(), year)
	if err != nil {
		res.Error(c, err)
		return
	}
	res.Success(c, resp)
}

func getYearFromQuery(c *gin.Context) int {
	yearStr := c.Query("year")
	if yearStr == "" {
		return 0
	}
	year, err := strconv.Atoi(yearStr)
	if err != nil {
		return 0
	}
	return year
}
