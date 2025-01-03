package messaging

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/rabbitmq/amqp091-go"
)

type Message struct {
	Event   string `json:"event"`
	Song_ID string `json:"song_id"`
	Title   string `json:"title,omitempty"`
	Artist  string `json:"artist,omitempty"`
}

func PublishMessage(eventType, song_ID, title, artist string) error {
	if channel == nil {
		log.Println("RabbitMQ channel is not initialized")
		return fmt.Errorf("RabbitMQ channel is nil")
	}

	// Construct the message with a struct to ensure field order
	message := Message{
		Event:   eventType,
		Song_ID: song_ID,
	}

	// Add title and artist only for "created" event
	if eventType == "created" {
		message.Title = title
		message.Artist = artist
	}

	body, err := json.Marshal(message)
	if err != nil {
		log.Printf("Failed to marshal message: %v", err)
		return err
	}

	err = channel.Publish(
		"",            // exchange
		"song_events", // routing key
		false,         // mandatory
		false,         // immediate
		amqp091.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)

	if err != nil {
		log.Printf("Failed to publish message: %v", err)
		return err
	}

	log.Println("Published message to RabbitMQ:", string(body))
	return nil
}
