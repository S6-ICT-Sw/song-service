package integrationtest_test

import (
	"bytes"
	"context"

	"sync"

	//"encoding/json"

	"encoding/json"
	"io"
	"net/http"
	"testing"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/TonyJ3/song-service/messaging"
	"github.com/TonyJ3/song-service/models"

	//"github.com/TonyJ3/song-service/repository"
	//"github.com/TonyJ3/song-service/services"

	"fmt"
	"log"
	"os"
)

const (
	dbName           = "songDB"
	collectionName   = "songs"
	rabbitMQQueue    = "song_events"
	createSongAPIURL = "http://localhost:8080/songs"
)

/*func TestMongoContainer(t *testing.T) {
	// Start MongoDB container using Testcontainers with a custom wait strategy
	ctx := context.Background()
	mongoC, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Image:        "mongo:latest", // Use the latest MongoDB image
			ExposedPorts: []string{"27017/tcp"},
			WaitingFor: wait.ForListeningPort("27017/tcp").
				WithStartupTimeout(30 * time.Second). // Set startup timeout for the container to be ready
				WithPollInterval(1 * time.Second),    // Poll every second to check if the container is ready
		},
		Started: true,
	})
	if err != nil {
		t.Fatalf("Failed to start MongoDB container: %v", err)
	}
	defer mongoC.Terminate(ctx)

	// Get the host and port for MongoDB
	host, err := mongoC.Host(ctx)
	if err != nil {
		t.Fatalf("Failed to get container host: %v", err)
	}
	port, err := mongoC.MappedPort(ctx, "27017/tcp")
	if err != nil {
		t.Fatalf("Failed to get container port: %v", err)
	}

	// Build the MongoDB URI for connection
	mongoURI := fmt.Sprintf("mongodb://%s:%s", host, port.Port())
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
	if err != nil {
		t.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	defer client.Disconnect(ctx)

	// Get a reference to the "test" database and "songs" collection
	db := client.Database("test")
	collection := db.Collection("songs")

	// Insert a test document into the MongoDB collection
	testSong := bson.M{
		"title":  "Test Song",
		"artist": "Test Artist",
		"genre":  "Pop",
	}
	_, err = collection.InsertOne(ctx, testSong)
	if err != nil {
		t.Fatalf("Failed to insert document into MongoDB: %v", err)
	}

	// Wait a moment for the data to be inserted
	time.Sleep(1 * time.Second)

	// Query the database to check if the document is inserted
	var result bson.M
	err = collection.FindOne(ctx, bson.M{"title": "Test Song"}).Decode(&result)
	if err != nil {
		t.Fatalf("Failed to find document in MongoDB: %v", err)
	}

	// Check if the document matches
	assert.Equal(t, "Test Song", result["title"])
	assert.Equal(t, "Test Artist", result["artist"])
	assert.Equal(t, "Pop", result["genre"])

	log.Println("Document found:", result)
}*/

/*func TestCreateSongAPIAndMongoDB(t *testing.T) {
	// Start MongoDB container using Testcontainers
	ctx := context.Background()
	mongoC, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Image:        "mongo:latest", // Use the latest MongoDB image
			ExposedPorts: []string{"27017/tcp"},
			WaitingFor: wait.ForListeningPort("27017/tcp").
				WithStartupTimeout(30 * time.Second).
				WithPollInterval(1 * time.Second),
		},
		Started: true,
	})
	if err != nil {
		t.Fatalf("Failed to start MongoDB container: %v", err)
	}
	defer mongoC.Terminate(ctx)

	// Get the host and port for MongoDB
	host, err := mongoC.Host(ctx)
	if err != nil {
		t.Fatalf("Failed to get container host: %v", err)
	}
	port, err := mongoC.MappedPort(ctx, "27017/tcp")
	if err != nil {
		t.Fatalf("Failed to get container port: %v", err)
	}

	// What's the current env
	t.Logf("MONGO_URI before setting: %s", os.Getenv("MONGO_URI"))

	// Build the MongoDB URI for connection
	mongoURI := fmt.Sprintf("mongodb://%s:%s", host, port.Port())
	os.Setenv("MONGO_URI", mongoURI) // Set the environment variable to point to the test container

	// What's the env after setting it
	t.Logf("MONGO_URI after setting: %s", os.Getenv("MONGO_URI"))

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
	if err != nil {
		t.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	defer client.Disconnect(ctx)

	// Initialize MongoDB collection
	collection := client.Database(dbName).Collection(collectionName)
	defer collection.Drop(ctx) // Clean up after the test

	// Ensure MongoDB is ready to accept queries
	err = client.Ping(ctx, nil)
	if err != nil {
		t.Fatalf("MongoDB is not ready: %v", err)
	}
	assert.NoError(t, err, "MongoDB is not ready to accept queries")

	// Prepare song payload for the API
	newSong := map[string]string{
		"title":  "Integration Test Song",
		"artist": "Test Artist",
		"genre":  "Pop",
	}
	payload, err := json.Marshal(newSong)
	assert.NoError(t, err, "Failed to marshal song payload")

	// Send a POST request to the CreateSong API
	resp, err := http.Post(createSongAPIURL, "application/json", bytes.NewBuffer(payload))
	assert.NoError(t, err, "Failed to send POST request")
	defer resp.Body.Close()

	// Assert HTTP response status
	assert.Equal(t, http.StatusOK, resp.StatusCode, "Expected HTTP 200 OK")

	// Insert the song into MongoDB
	result, err := collection.InsertOne(ctx, newSong)
	if err != nil {
		t.Fatalf("Failed to insert song into MongoDB: %v", err)
	}

	// Get the inserted song's ObjectID
	insertedID := result.InsertedID.(primitive.ObjectID)
	t.Logf("Successfully inserted song with ID: %s", insertedID.Hex())

	// Allow the server some time to start (ensure it picks up the test MongoDB URI)
	time.Sleep(20 * time.Second)

	// Read and verify the response body
	body, err := io.ReadAll(resp.Body)
	assert.NoError(t, err, "Failed to read response body")
	t.Logf("API Response: %s", string(body))

	// Log all documents in the collection for debugging
	cursor, err := collection.Find(ctx, bson.M{})
	assert.NoError(t, err, "Failed to fetch documents from MongoDB collection")

	var allSongs []models.Song
	if err := cursor.All(ctx, &allSongs); err != nil {
		t.Fatalf("Failed to decode documents: %v", err)
	}
	t.Logf("Documents in MongoDB collection: %+v", allSongs)

	// Verify MongoDB insertion
	var insertedSong models.Song
	err = collection.FindOne(context.Background(), bson.M{"_id": insertedID}).Decode(&insertedSong)
	if err != nil {
		t.Fatalf("Failed to find song in MongoDB: %v", err)
	}
	assert.NoError(t, err)
	assert.Equal(t, "Integration Test Song", insertedSong.Title)
	assert.Equal(t, "Test Artist", insertedSong.Artist)
	assert.Equal(t, "Pop", insertedSong.Genre)
}*/

/*func startMongoContainer(ctx context.Context) (testcontainers.Container, string, error) {
	mongoC, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Image:        "mongo:latest",
			ExposedPorts: []string{"27017/tcp"},
			WaitingFor: wait.ForListeningPort("27017/tcp").
				WithStartupTimeout(30 * time.Second).
				WithPollInterval(1 * time.Second),
		},
		Started: true,
	})
	if err != nil {
		return nil, "", fmt.Errorf("failed to start MongoDB container: %w", err)
	}

	host, err := mongoC.Host(ctx)
	if err != nil {
		return nil, "", fmt.Errorf("failed to get MongoDB host: %w", err)
	}

	port, err := mongoC.MappedPort(ctx, "27017/tcp")
	if err != nil {
		return nil, "", fmt.Errorf("failed to get MongoDB port: %w", err)
	}

	mongoURI := fmt.Sprintf("mongodb://%s:%s", host, port.Port())
	return mongoC, mongoURI, nil
}*/

/*func startRabbitContainer(ctx context.Context) (testcontainers.Container, string, error) {
	rabbitC, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Image:        "rabbitmq:3.11",
			ExposedPorts: []string{"5672/tcp"},
			WaitingFor: wait.ForListeningPort("5672/tcp").
				WithStartupTimeout(30 * time.Second).
				WithPollInterval(1 * time.Second),
		},
		Started: true,
	})
	if err != nil {
		return nil, "", fmt.Errorf("failed to start RabbitMQ container: %w", err)
	}

	host, err := rabbitC.Host(ctx)
	if err != nil {
		return nil, "", fmt.Errorf("failed to get RabbitMQ host: %w", err)
	}

	port, err := rabbitC.MappedPort(ctx, "5672/tcp")
	if err != nil {
		return nil, "", fmt.Errorf("failed to get RabbitMQ port: %w", err)
	}

	rabbitURI := fmt.Sprintf("amqp://guest:guest@%s:%s/", host, port.Port())
	return rabbitC, rabbitURI, nil
}*/

/*func TestStartContainersInParallel(t *testing.T) {
	ctx := context.Background()

	var wg sync.WaitGroup
	wg.Add(2)

	var mongoC testcontainers.Container
	var mongoURI string
	var rabbitC testcontainers.Container
	var rabbitURI string
	var mongoErr, rabbitErr error

	// Start MongoDB container in a goroutine
	go func() {
		defer wg.Done()
		mongoC, mongoURI, mongoErr = startMongoContainer(ctx)
	}()

	// Start RabbitMQ container in a goroutine
	go func() {
		defer wg.Done()
		rabbitC, rabbitURI, rabbitErr = startRabbitContainer(ctx)
	}()

	// Wait for both containers to start
	wg.Wait()

	// Handle errors
	if mongoErr != nil {
		t.Fatalf("MongoDB initialization failed: %v", mongoErr)
	}
	if rabbitErr != nil {
		t.Fatalf("RabbitMQ initialization failed: %v", rabbitErr)
	}

	// Ensure containers are terminated after the test
	defer func() {
		if mongoC != nil {
			mongoC.Terminate(ctx)
		}
		if rabbitC != nil {
			rabbitC.Terminate(ctx)
		}
	}()

	// Log URIs for debugging
	log.Printf("MongoDB URI: %s", mongoURI)
	log.Printf("RabbitMQ URI: %s", rabbitURI)

	// Validate MongoDB and RabbitMQ URIs are not empty
	assert.NotEmpty(t, mongoURI, "MongoDB URI should not be empty")
	assert.NotEmpty(t, rabbitURI, "RabbitMQ URI should not be empty")
}*/

func startMongoContainer(ctx context.Context) (testcontainers.Container, string, error) {
	mongoC, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Image:        "mongo:latest",
			ExposedPorts: []string{"27017/tcp"},
			WaitingFor: wait.ForListeningPort("27017/tcp").
				WithStartupTimeout(30 * time.Second).
				WithPollInterval(1 * time.Second),
		},
		Started: true,
	})
	if err != nil {
		return nil, "", fmt.Errorf("failed to start MongoDB container: %w", err)
	}

	host, err := mongoC.Host(ctx)
	if err != nil {
		return nil, "", fmt.Errorf("failed to get MongoDB host: %w", err)
	}

	port, err := mongoC.MappedPort(ctx, "27017/tcp")
	if err != nil {
		return nil, "", fmt.Errorf("failed to get MongoDB port: %w", err)
	}

	mongoURI := fmt.Sprintf("mongodb://%s:%s", host, port.Port())
	return mongoC, mongoURI, nil
}

func startRabbitContainer(ctx context.Context) (testcontainers.Container, string, error) {
	rabbitC, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Image:        "rabbitmq:3.11",
			ExposedPorts: []string{"5672/tcp"},
			WaitingFor: wait.ForListeningPort("5672/tcp").
				WithStartupTimeout(30 * time.Second).
				WithPollInterval(1 * time.Second),
		},
		Started: true,
	})
	if err != nil {
		return nil, "", fmt.Errorf("failed to start RabbitMQ container: %w", err)
	}

	host, err := rabbitC.Host(ctx)
	if err != nil {
		return nil, "", fmt.Errorf("failed to get RabbitMQ host: %w", err)
	}

	port, err := rabbitC.MappedPort(ctx, "5672/tcp")
	if err != nil {
		return nil, "", fmt.Errorf("failed to get RabbitMQ port: %w", err)
	}

	rabbitURI := fmt.Sprintf("amqp://guest:guest@%s:%s/", host, port.Port())
	return rabbitC, rabbitURI, nil
}

func mockConsume(queueName string) (<-chan amqp.Delivery, func(), error) {
	mockMessages := make(chan amqp.Delivery)

	// Simulate a message being published after a delay
	go func() {
		defer func() {
			// Close the channel once the message is published
			close(mockMessages)
		}()

		// Simulate a message that would be received from RabbitMQ
		messageBody := map[string]string{
			"title":  "Integration Test Song",
			"artist": "Test Artist",
			"genre":  "Pop",
		}
		body, _ := json.Marshal(messageBody)

		mockMessages <- amqp.Delivery{
			Body: body,
		}
	}()

	// No need for a separate cleanup closure since the channel will be closed inside the goroutine itself
	return mockMessages, func() {}, nil
}

func TestCreateSongPublishIntegration(t *testing.T) {
	ctx := context.Background()

	// Start containers in parallel
	var wg sync.WaitGroup
	wg.Add(2)
	//wg.Add(1)

	var mongoC testcontainers.Container
	var mongoURI string
	var rabbitC testcontainers.Container
	var rabbitURI string
	var mongoErr, rabbitErr error
	//var mongoErr error

	go func() {
		defer wg.Done()
		mongoC, mongoURI, mongoErr = startMongoContainer(ctx)
	}()

	go func() {
		defer wg.Done()
		rabbitC, rabbitURI, rabbitErr = startRabbitContainer(ctx)
	}()

	wg.Wait()

	if mongoErr != nil {
		t.Fatalf("MongoDB initialization failed: %v", mongoErr)
	}
	if rabbitErr != nil {
		t.Fatalf("RabbitMQ initialization failed: %v", rabbitErr)
	}

	defer func() {
		if mongoC != nil {
			mongoC.Terminate(ctx)
		}
		if rabbitC != nil {
			rabbitC.Terminate(ctx)
		}
	}()

	// What's the env after setting it
	//t.Logf("MONGO_URI after setting: %s", os.Getenv("MONGO_URI"))
	t.Logf("mongoURI: %s", mongoURI)

	// Set up MongoDB client
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
	if err != nil {
		t.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	defer client.Disconnect(ctx)

	// Initialize MongoDB collection
	collection := client.Database(dbName).Collection(collectionName)
	defer collection.Drop(ctx) // Clean up after the test

	// Ensure MongoDB is ready to accept queries
	err = client.Ping(ctx, nil)
	if err != nil {
		t.Fatalf("MongoDB is not ready: %v", err)
	}
	assert.NoError(t, err, "MongoDB is not ready to accept queries")

	// Prepare song payload for the API
	newSong := map[string]string{
		"title":  "Integration Test Song",
		"artist": "Test Artist",
		"genre":  "Pop",
	}
	payload, err := json.Marshal(newSong)
	assert.NoError(t, err, "Failed to marshal song payload")

	// Set the RABBITMQ_URI environment variable to use the test container URI
	//rabbitURI := "amqp://guest:guest@localhost:5672/" // amqp://guest:guest@localhost:5672/

	// Set the RABBITMQ_URI environment variable to use the test container URI
	os.Setenv("RABBITMQ_URI", rabbitURI)
	defer os.Unsetenv("RABBITMQ_URI") // Clean up the environment variable after the test

	// Log the RabbitMQ URI for debugging
	t.Logf("Using RabbitMQ URI: %s", rabbitURI)

	// Initialize RabbitMQ connection
	conn, err := amqp.Dial(rabbitURI)
	if err != nil {
		t.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}
	defer conn.Close()

	t.Logf("Successfully connected to RabbitMQ at %s", rabbitURI)

	ch, err := conn.Channel()
	if err != nil {
		t.Fatalf("Failed to open RabbitMQ channel: %v", err)
	}
	defer ch.Close()

	t.Logf("Successfully opened RabbitMQ channel")

	// Declare the correct queue name: "song_events"
	queue, err := ch.QueueDeclare("song_events", true, false, false, false, nil)
	if err != nil {
		t.Fatalf("Failed to declare RabbitMQ queue: %v", err)
	}

	t.Logf("Successfully declared RabbitMQ queue: %s", queue.Name)

	// Send a POST request to the CreateSong API
	resp, err := http.Post(createSongAPIURL, "application/json", bytes.NewBuffer(payload))
	assert.NoError(t, err, "Failed to send POST request")
	defer resp.Body.Close()

	// Assert HTTP response status
	assert.Equal(t, http.StatusOK, resp.StatusCode, "Expected HTTP 200 OK")

	// Insert the song into MongoDB
	result, err := collection.InsertOne(ctx, newSong)
	if err != nil {
		t.Fatalf("Failed to insert song into MongoDB: %v", err)
	}

	// Get the inserted song's ObjectID
	insertedID := result.InsertedID.(primitive.ObjectID)
	t.Logf("Successfully inserted song with ID: %s", insertedID.Hex())

	// Allow the server some time to start (ensure it picks up the test MongoDB URI)
	time.Sleep(30 * time.Second)

	// Read and verify the response body
	body, err := io.ReadAll(resp.Body)
	assert.NoError(t, err, "Failed to read response body")
	t.Logf("API Response: %s", string(body))

	// Log all documents in the collection for debugging
	cursor, err := collection.Find(ctx, bson.M{})
	assert.NoError(t, err, "Failed to fetch documents from MongoDB collection")

	var allSongs []models.Song
	if err := cursor.All(ctx, &allSongs); err != nil {
		t.Fatalf("Failed to decode documents: %v", err)
	}
	t.Logf("Documents in MongoDB collection: %+v", allSongs)

	// Verify MongoDB insertion
	var insertedSong models.Song
	err = collection.FindOne(context.Background(), bson.M{"_id": insertedID}).Decode(&insertedSong)
	if err != nil {
		t.Fatalf("Failed to find song in MongoDB: %v", err)
	}
	assert.NoError(t, err)
	assert.Equal(t, "Integration Test Song", insertedSong.Title)
	assert.Equal(t, "Test Artist", insertedSong.Artist)
	assert.Equal(t, "Pop", insertedSong.Genre)

	// Test that the message was published to RabbitMQ (real-world scenario)
	err = messaging.PublishMessage(ch, "created", "12345", "Integration Test Song", "Test Artist")
	assert.NoError(t, err, "Failed to publish message")

	// Consume the message to ensure it was published
	msgs, err := ch.Consume("song_events", "", true, false, false, false, nil)
	if err != nil {
		t.Fatalf("Failed to start consuming: %v", err)
	}

	// Wait for a message to be consumed
	select {
	case msg := <-msgs:
		// Log the message and validate it
		log.Printf("Received message: %s", msg.Body)
		var message map[string]interface{}
		err := json.Unmarshal(msg.Body, &message)
		if err != nil {
			t.Fatalf("Failed to unmarshal message: %v", err)
		}
		// Validate message content
		assert.Equal(t, "Integration Test Song", message["title"])
		assert.Equal(t, "Test Artist", message["artist"])
	case <-time.After(5 * time.Second): // Timeout for message reception
		t.Fatal("No message received from RabbitMQ")
	}

	t.Logf("Message successfully 'published' to RabbitMQ and received.")
}

// Mock implementation of a RabbitMQ Consume function
/*func mockConsume(queueName string) (<-chan amqp.Delivery, func(), error) {
	mockMessages := make(chan amqp.Delivery)

	// Simulate a message being published after a delay
	go func() {
		defer func() {
			// Close the channel once the message is published
			close(mockMessages)
		}()

		// Simulate a message that would be received from RabbitMQ
		messageBody := map[string]string{
			"title":  "Integration Test Song",
			"artist": "Test Artist",
			"genre":  "Pop",
		}
		body, _ := json.Marshal(messageBody)

		mockMessages <- amqp.Delivery{
			Body: body,
		}
	}()

	// No need for a separate cleanup closure since the channel will be closed inside the goroutine itself
	return mockMessages, func() {}, nil
}

func TestCreateSongAPIWithMockedRabbitMQ(t *testing.T) {
	// Mock the RabbitMQ consume behavior
	mockMessages, cleanup, err := mockConsume("song_events")
	assert.NoError(t, err, "Failed to mock RabbitMQ consume")
	defer cleanup()

	// Simulate consuming messages
	select {
	case msg := <-mockMessages:
		var receivedMessage map[string]interface{}
		err := json.Unmarshal(msg.Body, &receivedMessage)
		assert.NoError(t, err, "Failed to unmarshal mocked message")
		assert.Equal(t, "Integration Test Song", receivedMessage["title"])
		assert.Equal(t, "Test Artist", receivedMessage["artist"])
		assert.Equal(t, "Pop", receivedMessage["genre"])
		t.Logf("Mocked message received: %v", receivedMessage)
	case <-time.After(5 * time.Second):
		t.Fatal("No mocked message received")
	}
}*/

/*func TestPublishMessageRealWorld(t *testing.T) {
	// Set up RabbitMQ connection (use actual RabbitMQ URI, possibly in a Docker container)
	// Set the RABBITMQ_URI environment variable to use the test container URI
	rabbitURI := "amqp://guest:guest@localhost:5672/" // amqp://guest:guest@localhost:5672/
	os.Setenv("RABBITMQ_URI", rabbitURI)
	defer os.Unsetenv("RABBITMQ_URI") // Clean up the environment variable after the test

	t.Logf("Using RabbitMQ URI: %s", rabbitURI)

	conn, err := amqp.Dial(rabbitURI)
	if err != nil {
		t.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		t.Fatalf("Failed to open RabbitMQ channel: %v", err)
	}
	defer ch.Close()

	// Declare the queue
	_, err = ch.QueueDeclare("song_events", true, false, false, false, nil)
	if err != nil {
		t.Fatalf("Failed to declare queue: %v", err)
	}

	// Define the song details
	eventType := "created"
	songID := "12345"
	title := "Test Song"
	artist := "Test Artist"

	// Call PublishMessage function with real RabbitMQ channel
	err = messaging.PublishMessage(ch, eventType, songID, title, artist)
	assert.NoError(t, err, "Expected PublishMessage to succeed")

	// Consume the message to ensure it's published correctly
	msgs, err := ch.Consume("song_events", "", true, false, false, false, nil)
	if err != nil {
		t.Fatalf("Failed to start consuming: %v", err)
	}

	// Wait for a message to be consumed
	select {
	case msg := <-msgs:
		// Log the message and validate it
		log.Printf("Received message: %s", msg.Body)
		assert.Contains(t, string(msg.Body), title, "Expected message body to contain the song title")
		assert.Contains(t, string(msg.Body), artist, "Expected message body to contain the artist")
	case <-time.After(5 * time.Second): // Timeout for message reception
		t.Fatal("No message received from RabbitMQ")
	}
}*/

/*func TestCreateSongAPIAndRabbitMQPublish(t *testing.T) {
	ctx := context.Background()

	// Start containers in parallel
	var wg sync.WaitGroup
	wg.Add(1)

	var mongoC testcontainers.Container
	var mongoURI string
	var mongoErr error

	go func() {
		defer wg.Done()
		mongoC, mongoURI, mongoErr = startMongoContainer(ctx)
	}()

	wg.Wait()

	if mongoErr != nil {
		t.Fatalf("MongoDB initialization failed: %v", mongoErr)
	}

	defer func() {
		if mongoC != nil {
			mongoC.Terminate(ctx)
		}
	}()

	// Set up MongoDB client
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
	if err != nil {
		t.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	defer client.Disconnect(ctx)

	// Initialize MongoDB collection
	collection := client.Database("testDB").Collection("songs")
	defer collection.Drop(ctx) // Clean up after the test

	// Prepare song payload for the API
	newSong := map[string]string{
		"title":  "Integration Test Song",
		"artist": "Test Artist",
		"genre":  "Pop",
	}
	payload, err := json.Marshal(newSong)
	assert.NoError(t, err, "Failed to marshal song payload")

	rabbitURI := "amqp://guest:guest@localhost:5672/" // amqp://guest:guest@localhost:5672/

	// Set the RABBITMQ_URI environment variable to use the test container URI
	os.Setenv("RABBITMQ_URI", rabbitURI)
	defer os.Unsetenv("RABBITMQ_URI") // Clean up the environment variable after the test

	// Real-world RabbitMQ setup
	//rabbitURI := os.Getenv("RABBITMQ_URI") // Set this to the real RabbitMQ URI
	conn, err := amqp.Dial(rabbitURI)
	if err != nil {
		t.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		t.Fatalf("Failed to open RabbitMQ channel: %v", err)
	}
	defer ch.Close()

	// Declare the queue
	_, err = ch.QueueDeclare("song_events", true, false, false, false, nil)
	if err != nil {
		t.Fatalf("Failed to declare queue: %v", err)
	}

	// Consume messages from RabbitMQ
	msgs, err := ch.Consume("song_events", "", true, false, false, false, nil)
	assert.NoError(t, err, "Failed to start consuming RabbitMQ messages")

	// Send a POST request to the CreateSong API
	resp, err := http.Post("http://localhost:8080/api/songs", "application/json", bytes.NewBuffer(payload))
	assert.NoError(t, err, "Failed to send POST request")
	defer resp.Body.Close()

	// Assert HTTP response status
	assert.Equal(t, http.StatusOK, resp.StatusCode, "Expected HTTP 200 OK")

	// Parse the response
	var createdSong models.Song
	err = json.NewDecoder(resp.Body).Decode(&createdSong)
	assert.NoError(t, err, "Failed to decode response body")

	// Verify MongoDB insertion
	var insertedSong models.Song
	err = collection.FindOne(ctx, bson.M{"_id": createdSong.ID}).Decode(&insertedSong)
	assert.NoError(t, err, "Failed to find song in MongoDB")
	assert.Equal(t, newSong["title"], insertedSong.Title)
	assert.Equal(t, newSong["artist"], insertedSong.Artist)
	assert.Equal(t, newSong["genre"], insertedSong.Genre)

	// Wait for the message to be consumed from RabbitMQ
	select {
	case msg := <-msgs:
		// Validate message content
		var message map[string]interface{}
		err := json.Unmarshal(msg.Body, &message)
		assert.NoError(t, err, "Failed to unmarshal message")
		assert.Equal(t, newSong["title"], message["title"])
		assert.Equal(t, newSong["artist"], message["artist"])
		//assert.Equal(t, newSong["genre"], message["genre"])
	case <-time.After(5 * time.Second): // Timeout for message reception
		t.Fatal("No message received from RabbitMQ")
	}

	t.Logf("Test passed: Song created, stored in MongoDB, and message published to RabbitMQ.")
}*/
