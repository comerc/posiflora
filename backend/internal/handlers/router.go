package handlers

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/posiflora/backend/config"
)

func SetupRouter(
	cfg *config.Config,
	telegramHandler *TelegramHandler,
	orderHandler *OrderHandler,
) *gin.Engine {
	r := gin.Default()

	// CORS middleware
	r.Use(cors.New(cors.Config{
		AllowOrigins:     cfg.Server.CORSOrigins,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	// API routes
	api := r.Group("/shops")
	{
		api.POST("/:shopId/telegram/connect", telegramHandler.Connect)
		api.GET("/:shopId/telegram/status", telegramHandler.GetStatus)
		api.POST("/:shopId/orders", orderHandler.CreateOrder)
	}

	return r
}
