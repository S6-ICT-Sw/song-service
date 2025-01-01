package main

import (
	"context"
	//"encoding/json"
	"fmt"
	"log"
	"net/http"

	//"os"
	//"time"

	"github.com/TonyJ3/song-service/api"
	"github.com/TonyJ3/song-service/repository"

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

func StartLocalServer() {
	// MongoDB connection
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI("mongodb+srv://song-snippets-admin:DQv4P9LXNBQ2xsdb@songsnippets.ci2mt.mongodb.net/?retryWrites=true&w=majority&appName=SongSnippets"))
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	defer client.Disconnect(context.Background())

	// Initialize repository
	repository.InitRepository(client, "songDB", "songs")

	// Setup the router (mux)
	router := api.SetupRouter()

	// Notifaction that the system is running
	fmt.Println("Localhost:8080 is running")

	// Start the server
	log.Fatal(http.ListenAndServe(":8080", router))
}

func main() {
	// Check if we are running in AWS Lambda environment
	/*if os.Getenv("AWS_LAMBDA_FUNCTION_NAME") != "" {
		// Run as AWS Lambda function
		lambda.Start(LambdaHandler)
	} else {
		StartLocalServer()
	}*/

	StartLocalServer()
}
