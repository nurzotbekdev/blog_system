package rabbitmq

import (
	"blog_system/config"
	"blog_system/models"
	"blog_system/schemas"
	"context"
	"encoding/json"
	"fmt"
	"log"
)

func ConsumeNotifications() {
	msgs, err := config.RabbitChannel.Consume("notification_queue", "", true, false, false, false, nil)
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		for msg := range msgs {

			var payload schemas.NotificationPayload
			if err := json.Unmarshal(msg.Body, &payload); err != nil {
				continue
			}

			notification := models.Notification{
				UserID:     payload.UserID,
				FromUserID: payload.FromUserID,
				Type:       payload.Type,
				VideoID:    payload.VideoID,
				CommentID:  payload.CommentID,
				ChannelID:  payload.ChannelID,
			}

			err := config.DB.Create(&notification).Error
			if err != nil {
				continue
			}

			key := fmt.Sprintf("notif:unread:%d", payload.UserID)

			err = config.RedisClient.Incr(context.Background(), key).Err()
			if err != nil {
				log.Println("Redis INCR error:", err)
			}
		}
	}()
}
