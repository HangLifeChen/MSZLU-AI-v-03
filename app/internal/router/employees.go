package router

import (
	"app/internal/employees"

	"github.com/gin-gonic/gin"
)

type EmployeeRouter struct {
}

func (e *EmployeeRouter) Register(engine *gin.Engine) {
	employeesGroup := engine.Group("/api/v1/employees")
	{
		employeesHandler := employees.NewHandler()
		employeesGroup.POST("/create", employeesHandler.CreateEmployee)
		employeesGroup.GET("/list", employeesHandler.ListEmployees)
		employeesGroup.GET("/:id", employeesHandler.GetEmployee)
		employeesGroup.PUT("/update", employeesHandler.UpdateEmployee)
		employeesGroup.DELETE("/:id", employeesHandler.DeleteEmployee)
	}
}
