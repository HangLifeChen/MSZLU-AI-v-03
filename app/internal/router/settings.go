package router

import (
	"app/internal/settings"

	"github.com/gin-gonic/gin"
)

type SettingsRouter struct {
}

func (s *SettingsRouter) Register(engine *gin.Engine) {
	settingsGroup := engine.Group("/api/v1/settings")
	{
		settingsHandler := settings.NewHandler()
		settingsGroup.GET("", settingsHandler.GetSettings)
		settingsGroup.POST("", settingsHandler.SaveSettings)
		settingsGroup.PUT("", settingsHandler.UpdateSettings)
		settingsGroup.GET("/:module", settingsHandler.GetSettingsModule)
		settingsGroup.PUT("/:module", settingsHandler.UpdateSettingsModule)
	}
}
