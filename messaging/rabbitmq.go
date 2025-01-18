package messaging

import (
	"log"
	"os"

	//"sync"

	"github.com/rabbitmq/amqp091-go"
)

var connection *amqp091.Connection
var channel *amqp091.Channel
var rabbitMQEnabled bool

func InitRabbitMQ() error {
	rabbitURI := os.Getenv("RABBITMQ_URI") //"amqp://user:password@rabbitmq:5672/" amqp://guest:guest@localhost:5672/
	if rabbitURI == "" {
		log.Println("RABBITMQ_URI not set. RabbitMQ is disabled.")
		rabbitMQEnabled = false
		return nil
	}

	var err error
	connection, err = amqp091.Dial(rabbitURI)
	if err != nil {
		log.Printf("Failed to connect to RabbitMQ: %v", err)
		return err
	}

	channel, err = connection.Channel()
	if err != nil {
		log.Printf("Failed to open a channel: %v", err)
		rabbitMQEnabled = false
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
		rabbitMQEnabled = false
		return err
	}

	// Set the global channel variable
	//channelLock.Lock()
	//channelLock.Unlock()
	rabbitMQEnabled = true
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

func GetChannel() *amqp091.Channel {
	//channelLock.Lock()
	//defer channelLock.Unlock()
	return channel
}

func IsRabbitMQEnabled() bool {
	//return channel != nil
	return rabbitMQEnabled
}

// This is use for the integration test
func ConsumeQueue(queueName string) (<-chan amqp091.Delivery, error) {
	msgs, err := channel.Consume(
		queueName, // queue
		"",        // consumer
		true,      // auto-ack
		false,     // exclusive
		false,     // no-local
		false,     // no-wait
		nil,       // args
	)
	if err != nil {
		log.Printf("Failed to consume RabbitMQ queue: %v", err)
		return nil, err
	}
	return msgs, nil
}
