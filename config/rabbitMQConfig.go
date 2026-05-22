package config

import (
	"blog_system/logging"
	"fmt"
	"os"
	"time"

	"github.com/streadway/amqp"
	"go.uber.org/zap"
)

var RabbitConn *amqp.Connection
var RabbitChannel *amqp.Channel

func ConnectingRabbitMQ() {
	user := os.Getenv("RABBITMQ_USER")
	password := os.Getenv("RABBITMQ_PASSWORD")
	host := os.Getenv("RABBITMQ_HOST")
	port := os.Getenv("RABBITMQ_PORT")

	logging.Log.Info("Initializing RabbitMQ connection", zap.String("host", host), zap.String("port", port), zap.String("user", user))
	dsn := fmt.Sprintf("amqp://%s:%s@%s:%s/", user, password, host, port)
	var err error

	for i := 1; i <= 5; i++ {
		RabbitConn, err = amqp.Dial(dsn)
		if err == nil {
			logging.Log.Info("RabbitMQ connection established", zap.Int("attempt", i))
			break
		}

		logging.Log.Warn("RabbitMQ connection failed, retrying", zap.Int("attempt", i), zap.Int("max_attempts", 5), zap.Error(err))
		time.Sleep(3 * time.Second)
	}

	if err != nil {
		logging.Log.Fatal("Failed to connect to RabbitMQ", zap.String("host", host), zap.String("port", port), zap.Error(err))
	}

	RabbitChannel, err = RabbitConn.Channel()
	if err != nil {
		logging.Log.Fatal("Failed to open RabbitMQ channel", zap.Error(err))
	}

	logging.Log.Info("RabbitMQ connected successfully")
}

func DeclareNotificationQueue() error {
	logging.Log.Info("Declaring RabbitMQ queue", zap.String("queue", "notification_queue"))
	_, err := RabbitChannel.QueueDeclare("notification_queue", true, false, false, false, nil)
	if err != nil {
		logging.Log.Error("Failed to declare RabbitMQ queue", zap.String("queue", "notification_queue"), zap.Error(err))
		return err
	}

	logging.Log.Info("RabbitMQ queue declared successfully", zap.String("queue", "notification_queue"))
	return nil
}

func CloseRabbitMQ() {
	logging.Log.Info("Closing RabbitMQ connections")

	if RabbitChannel != nil {
		if err := RabbitChannel.Close(); err != nil {
			logging.Log.Error("Failed to close RabbitMQ channel", zap.Error(err))
		} else {
			logging.Log.Info("RabbitMQ channel closed successfully")
		}
	}

	if RabbitConn != nil {
		if err := RabbitConn.Close(); err != nil {
			logging.Log.Error("Failed to close RabbitMQ connection", zap.Error(err))
		} else {
			logging.Log.Info("RabbitMQ connection closed successfully")
		}
	}
}
