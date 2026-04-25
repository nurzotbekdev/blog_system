package main

import (
	"blog_system/config"
	"blog_system/logging"
	"blog_system/routes"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	config.EnvConfig()
	config.DatabaseConfig()
	config.MigrateConfig()
	config.SetupGoogleAuth()

	logging.Init()

	router := gin.New()

	router.Use(logging.GinLogger())
	router.Use(logging.GinErrorLogger())
	router.Use(gin.Recovery())

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	routes.UserRoutes(router)
	routes.ChannelRoutes(router)
	routes.SubscriptionRoutes(router)

	logging.Log.Info("Server started on :8080")

	router.Run()
}
