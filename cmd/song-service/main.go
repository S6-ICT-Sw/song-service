package main

import (
	"context"
	"encoding/json"

	//"fmt"
	"log"
	"net/http"

	//"flag"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/TonyJ3/song-service/api"
	"github.com/TonyJ3/song-service/messaging"
	"github.com/TonyJ3/song-service/repository"
	"github.com/TonyJ3/song-service/services"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"

	"github.com/TonyJ3/song-service/models"
	"github.com/gorilla/handlers"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	//"github.com/TonyJ3/song-service/services"
)

/*func LambdaHandler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var newSong models.Song

	// Decode the incoming request body
	err := json.Unmarshal([]byte(req.Body), &newSong)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest,
			Body:       `{"message": "Invalid JSON format"}`,
		}, nil
	}

	// Validate required fields
	if newSong.Title == "" || newSong.Artist == "" || newSong.Genre == "" {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest,
			Body:       `{"message": "Missing required fields: title, artist, or genre"}`,
		}, nil
	}

	// Call the service to create the song
	createdSong := services.CreateSong(newSong)

	// Return the created song as the response
	responseBody, err := json.Marshal(createdSong)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       `{"message": "Failed to process response"}`,
		}, nil
	}

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusCreated,
		Body:       string(responseBody),
	}, nil
}*/

var mongoClient *mongo.Client

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	// Load MongoDB URI from environment variable
	dbURI := os.Getenv("MONGO_URI")
	if dbURI == "" {
		log.Fatal("MONGO_URI environment variable is not set")
	}

	mongoClient, err = mongo.Connect(context.Background(), options.Client().ApplyURI(dbURI))
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}

	if err := messaging.InitRabbitMQ(); err != nil {
		log.Fatalf("Failed to initialize RabbitMQ: %v, running without RabbitMQ", err)
	} else {
		log.Println("RabbitMQ successfully initialized.")
	}
}

func LambdaHandler(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// Parse the incoming request body
	var newSong models.Song
	err := json.Unmarshal([]byte(req.Body), &newSong)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 400,
			Body:       "Invalid JSON format",
		}, nil
	}

	// Validate required fields
	if newSong.Title == "" || newSong.Artist == "" || newSong.Genre == "" {
		return events.APIGatewayProxyResponse{
			StatusCode: 400,
			Body:       "Missing required fields: title, artist, or genre",
		}, nil
	}

	// Set up MongoDB repository
	repo := repository.NewMongoSongRepository(mongoClient.Database("songDB").Collection("songs"))
	services.SetRepository(repo)

	// Create the song
	createdSong, err := services.CreateSong(newSong)
	if err != nil {
		log.Printf("Failed to create song: %v", err)
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       "Failed to create song",
		}, nil
	}

	// Try to publish the event to RabbitMQ
	if messaging.IsRabbitMQEnabled() {
		err = messaging.PublishMessage(
			messaging.GetChannel(),
			"created",
			createdSong.ID.Hex(),
			createdSong.Title,
			createdSong.Artist,
		)
		if err != nil {
			log.Printf("Failed to publish song creation event: %v", err)
		} else {
			log.Println("Successfully published song creation event to RabbitMQ")
		}
	} else {
		log.Printf("Mock Publish: Created event for song ID=%s, Title=%s, Artist=%s",
			createdSong.ID.Hex(), createdSong.Title, createdSong.Artist)
	}

	// Return the created song as a response
	responseBody, err := json.Marshal(createdSong)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       "Failed to encode response",
		}, nil
	}

	return events.APIGatewayProxyResponse{
		StatusCode: 201,
		Body:       string(responseBody),
	}, nil
}

/*func StartLocalServer() {
	// Initialize RabbitMQ
	if err := messaging.InitRabbitMQ(); err != nil {
		log.Fatalf("Failed to initialize RabbitMQ: %v, running without RabbitMQ", err)
	} else {
		log.Println("RabbitMQ successfully initialized.")
	}
	defer messaging.CloseRabbitMQ() // Ensure RabbitMQ connection is closed on shutdown

	// Load MongoDB URI from environment variable
	dbURI := os.Getenv("MONGO_URI")
	if dbURI == "" {
		log.Fatal("MONGO_URI environment variable is not set")
	}

	// MongoDB connection
	// "mongodb+srv://song-snippets-admin:DQv4P9LXNBQ2xsdb@songsnippets.ci2mt.mongodb.net/?retryWrites=true&w=majority&appName=SongSnippets"
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(dbURI))
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	//defer client.Disconnect(context.Background())

	defer func() {
		if err := client.Disconnect(context.Background()); err != nil {
			log.Printf("Error disconnecting from MongoDB: %v", err)
		}
	}()

	// Initialize repository
	repo := repository.NewMongoSongRepository(client.Database("songDB").Collection("songs"))
	services.SetRepository(repo)

	// Setup the router (mux)
	router := api.SetupRouter()

	// Add CORS support
	cors := handlers.CORS(
		handlers.AllowedOrigins([]string{"*"}), // Allow all origins (or specify your frontend's origin)
		handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}),
		handlers.AllowedHeaders([]string{"Content-Type", "Authorization"}),
	)

	// Notifaction that the system is running
	fmt.Println("Localhost:8080 is running")

	// Start the server
	//log.Fatal(http.ListenAndServe(":8080", router))

	// Channel to handle shutdown signals
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// Create a custom HTTP server
	server := &http.Server{
		Addr:    ":8080",
		Handler: cors(router),
	}

	// Graceful shutdown handling
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()
	log.Println("Server is ready to handle requests.")

	// Wait for interrupt signal to gracefully shut down
	<-quit
	log.Println("Shutting down server...")

	// Gracefully shutdown the server with a timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server shutdown failed: %v", err)
	}
	log.Println("Server exited.")
}*/

func StartLocalServer() {
	defer func() {
		if err := mongoClient.Disconnect(context.Background()); err != nil {
			log.Printf("Error disconnecting from MongoDB: %v", err)
		}
	}()
	defer messaging.CloseRabbitMQ()

	// Initialize repository
	repo := repository.NewMongoSongRepository(mongoClient.Database("songDB").Collection("songs"))
	services.SetRepository(repo)

	// Setup the router (mux)
	router := api.SetupRouter()

	// Add CORS support
	cors := handlers.CORS(
		handlers.AllowedOrigins([]string{"*"}), // Allow all origins (or specify your frontend's origin)
		handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}),
		handlers.AllowedHeaders([]string{"Content-Type", "Authorization"}),
	)

	// Channel to handle shutdown signals
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// Create a custom HTTP server
	server := &http.Server{
		Addr:    ":8080",
		Handler: cors(router),
	}

	// Graceful shutdown handling
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()
	log.Println("Server is ready to handle requests on http://localhost:8080.")

	// Wait for interrupt signal to gracefully shut down
	<-quit
	log.Println("Shutting down server...")

	// Gracefully shutdown the server with a timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server shutdown failed: %v", err)
	}
	log.Println("Server exited.")
}

func main() {
	//StartLocalServer()
	defer func() {
		if err := mongoClient.Disconnect(context.Background()); err != nil {
			log.Printf("Error disconnecting from MongoDB: %v", err)
		}
	}()
	defer messaging.CloseRabbitMQ()

	// Use a flag to determine the mode (local or Lambda)
	/*runLocal := flag.Bool("local", false, "Run the application in local server mode")
	flag.Parse()

	if *runLocal {
		StartLocalServer()
	} else {
		lambda.Start(LambdaHandler)
	}*/

	//lambda.Start(LambdaHandler)

	if os.Getenv("LOCAL") == "true" {
		StartLocalServer()
	} else {
		lambda.Start(LambdaHandler)
	}
}
