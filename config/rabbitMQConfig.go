package config

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/streadway/amqp"
)

var RabbitConn *amqp.Connection
var RabbitChannel *amqp.Channel

func ConnectingRabbitMQ() {
	user := os.Getenv("RABBITMQ_USER")
	password := os.Getenv("RABBITMQ_PASSWORD")
	host := os.Getenv("RABBITMQ_HOST")
	port := os.Getenv("RABBITMQ_PORT")

	log.Printf("Initializing RabbitMQ connection (host: %s, port: %s, user: %s)", host, port, user)
	dsn := fmt.Sprintf("amqp://%s:%s@%s:%s/", user, password, host, port)
	var err error

	for i := 1; i <= 5; i++ {
		RabbitConn, err = amqp.Dial(dsn)
		if err == nil {
			log.Printf("RabbitMQ connection established (attempt %d/5)", i)
			break
		}

		log.Printf("WARNING: RabbitMQ connection failed, retrying (attempt %d/5): %v", i, err)
		time.Sleep(3 * time.Second)
	}

	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ (host: %s, port: %s): %v", host, port, err)
	}

	RabbitChannel, err = RabbitConn.Channel()
	if err != nil {
		log.Fatalf("Failed to open RabbitMQ channel: %v", err)
	}

	log.Println("RabbitMQ connected successfully")
}

func DeclareNotificationQueue() error {
	log.Println("Declaring RabbitMQ queue: notification_queue")
	_, err := RabbitChannel.QueueDeclare("notification_queue", true, false, false, false, nil)
	if err != nil {
		log.Printf("ERROR: Failed to declare RabbitMQ queue (notification_queue): %v", err)
		return err
	}

	log.Println("RabbitMQ queue declared successfully: notification_queue")
	return nil
}

func CloseRabbitMQ() {
	log.Println("Closing RabbitMQ connections")

	if RabbitChannel != nil {
		if err := RabbitChannel.Close(); err != nil {
			log.Printf("ERROR: Failed to close RabbitMQ channel: %v", err)
		} else {
			log.Println("RabbitMQ channel closed successfully")
		}
	}

	if RabbitConn != nil {
		if err := RabbitConn.Close(); err != nil {
			log.Printf("ERROR: Failed to close RabbitMQ connection: %v", err)
		} else {
			log.Println("RabbitMQ connection closed successfully")
		}
	}
}
