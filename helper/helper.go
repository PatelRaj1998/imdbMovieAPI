package helper

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func ConnectToDB() *mongo.Collection {
	config := GetConfiguration()
	// Set client options
	clientOptions := options.Client().ApplyURI(config.ConnectionString)

	// Connect to MongoDB
	client, err := mongo.Connect(context.TODO(), clientOptions)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to MongoDB!")

	collection := client.Database("moviesDB").Collection("movies")

	return collection
}

/*
type ErrorResponse struct {
	StatusCode   int    `json:"status"`
	ErrorMessage string `json:"message"`
}

func GetError(err error, w http.ResponseWriter) {

	log.Fatal(err.Error())
	var response = ErrorResponse{
		ErrorMessage: err.Error(),
		StatusCode:   http.StatusBadRequest,
	}

	message, _ := json.Marshal(response)

	w.WriteHeader(response.StatusCode)
	w.Write(message)
}
*/
type Configuration struct {
	Port             string
	ConnectionString string
}

func GetConfiguration() Configuration {
	err := godotenv.Load("./.env")

	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	configuration := Configuration{
		os.Getenv("PORT"),
		os.Getenv("CONNECTION_STRING"),
	}

	return configuration
}
