package config

import "blog_system/models"

func MigrateConfig() {
	DB.AutoMigrate(
		&models.User{},
		&models.Channel{},
		&models.Subscription{},
		&models.Category{},
	)
}
