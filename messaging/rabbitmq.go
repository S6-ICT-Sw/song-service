package messaging

import (
	//"fmt"
	"log"

	"github.com/rabbitmq/amqp091-go"
)

var connection *amqp091.Connection
var channel *amqp091.Channel

func InitRabbitMQ() error {
	var err error
	connection, err = amqp091.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		log.Printf("Failed to connect to RabbitMQ: %v", err)
		return err
	}

	channel, err = connection.Channel()
	if err != nil {
		log.Printf("Failed to open a channel: %v", err)
		return err
	}

	// Declare a queue
	_, err = channel.QueueDeclare(
		"song_events", // name
		true,          // durable
		false,         // delete when unused
		false,         // exclusive
		false,         // no-wait
		nil,           // arguments
	)
	if err != nil {
		log.Printf("Failed to declare a queue: %v", err)
		return err
	}

	log.Println("RabbitMQ connection and channel initialized")
	return nil
}

func CloseRabbitMQ() {
	if channel != nil {
		channel.Close()
	}

	if connection != nil {
		connection.Close()
	}
}
