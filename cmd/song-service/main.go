package main

import (
	"context"
	//"encoding/json"
	"fmt"
	"log"
	"net/http"

	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/TonyJ3/song-service/api"
	"github.com/TonyJ3/song-service/messaging"
	"github.com/TonyJ3/song-service/repository"
	"github.com/TonyJ3/song-service/services"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	//"github.com/TonyJ3/song-service/models"
	//"github.com/TonyJ3/song-service/services"
	//"github.com/aws/aws-lambda-go/events"
	//"github.com/aws/aws-lambda-go/lambda"
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

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, using system environment variables")
	}
}

func StartLocalServer() {
	// Initialize RabbitMQ
	if err := messaging.InitRabbitMQ(); err != nil {
		log.Fatalf("Failed to initialize RabbitMQ: %v", err)
	}
	defer messaging.CloseRabbitMQ() // Ensure RabbitMQ connection is closed on shutdown

	// Load MongoDB URI from environment variable
	dbURI := os.Getenv("MONGO_URI")
	if dbURI == "" {
		log.Fatal("MONGO_URI environment variable is not set")
	}

	//Test
	log.Printf("Using MongoDB URI: %s", os.Getenv("MONGO_URI"))

	// MongoDB connection
	// "mongodb+srv://song-snippets-admin:DQv4P9LXNBQ2xsdb@songsnippets.ci2mt.mongodb.net/?retryWrites=true&w=majority&appName=SongSnippets"
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(dbURI))
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	defer client.Disconnect(context.Background())

	// Initialize repository
	//repository.InitRepository(client, "songDB", "songs")
	//repository.InitRepository(client.Database("songDB").Collection("songs"))
	repo := repository.NewMongoSongRepository(client.Database("songDB").Collection("songs"))
	services.SetRepository(repo)

	//svc := services.NewSongService(repo)
	//h := song.NewSongHandler(svc)
	//r := api.NewRouter(h)

	// Setup the router (mux)
	router := api.SetupRouter()

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
		Handler: router,
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
}

func main() {

	StartLocalServer()
}
