package router

import (
	"app/internal/users"

	"github.com/gin-gonic/gin"
)

type UserRouter struct {
}

func (u *UserRouter) Register(r *gin.Engine) {
	userGroup := r.Group("/api/user")
	usersHandler := users.NewHandler()
	userGroup.POST("/", usersHandler.CreateUser)
	userGroup.GET("/", usersHandler.GetUser)
	userGroup.PUT("/", usersHandler.UpdateUser)
	userGroup.DELETE("/", usersHandler.DeleteUser)
}
