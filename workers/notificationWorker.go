package workers

import (
	"blog_system/config"
	"blog_system/models"
	"context"
	"fmt"
)

func WarmUpUnread() {

	var results []struct {
		UserID uint
		Count  int64
	}

	config.DB.
		Model(&models.Notification{}).
		Select("user_id, count(*) as count").
		Where("is_read = false").
		Group("user_id").
		Scan(&results)

	for _, r := range results {
		key := fmt.Sprintf("notif:unread:%d", r.UserID)
		config.RedisClient.Set(context.Background(), key, r.Count, 0)
	}
}
