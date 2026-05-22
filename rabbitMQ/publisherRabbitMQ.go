package rabbitmq

import (
	"blog_system/config"
	"encoding/json"

	"github.com/streadway/amqp"
)

func PublishNotification(data interface{}) error {
	body, err := json.Marshal(data)
	if err != nil {
		return err
	}

	err = config.RabbitChannel.Publish("", "notification_queue", false, false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		})

	return err
}
