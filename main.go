package main

import (
	"blog_system/config"
	"blog_system/logging"
	rabbitmq "blog_system/rabbitMQ"
	"blog_system/routes"
	"blog_system/workers"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	config.EnvConfig()
	config.DatabaseConfig()
	config.MigrateConfig()
	config.RedisConfig()
	config.SetupGoogleAuth()
	config.ConnectingRabbitMQ()

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
	routes.CategoryRoutes(router)
	routes.LanguageRoutes(router)
	routes.VideoRoutes(router)
	routes.VideoShareRoutes(router)
	routes.VideoDownloadRoutes(router)
	routes.CommentRoutes(router)
	routes.CommentReactionRoutes(router)
	routes.LikeRoutes(router)
	routes.HistoryRoutes(router)
	routes.PlaylistRoutes(router)
	routes.PlaylistVideoRoutes(router)
	routes.NotificationRoutes(router)

	go rabbitmq.ConsumeNotifications()
	go workers.StartVideoWorker()
	go workers.StartVideoViewWorker()
	go workers.WarmUpUnread()

	logging.Log.Info("Server started on :8080")

	router.Run()
}
