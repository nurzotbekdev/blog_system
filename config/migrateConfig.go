package config

import "blog_system/models"

func MigrateConfig() {
	DB.AutoMigrate(
		&models.User{},
		&models.Channel{},
		&models.Subscription{},
		&models.Category{},
		&models.Language{},
		&models.Video{},
		&models.VideoShare{},
		&models.VideoDownload{},
		&models.VideoQuality{},
		&models.Comment{},
		&models.Like{},
		&models.History{},
	)
}
