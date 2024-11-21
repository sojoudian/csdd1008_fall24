package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var mongoClient *mongo.Client

// POST API to receive the IP address of the client
func saveIPHandler(w http.ResponseWriter, r *http.Request) {
	// TODO Handle FW resp
	if r.Method != "POST" {
		http.Error(w, "Only POST method is allowed!", http.StatusMethodNotAllowed)
		return
	}
	var data map[string]interface{}
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		http.Error(w, "Error parsing JSON request body", http.StatusBadRequest)
		return
	}
	currentTime := time.Now().Format("2006-01-02 15:04:05")
	data["time"] = currentTime
	collection := mongoClient.Database("portfolio").Collection("ips")
	_, err = collection.InsertOne(context.TODO(), data)
	if err != nil {
		http.Error(w, "Error saving data to MongoDB", http.StatusInternalServerError)
		return
	}

}

func main() {
	var err error
	mongoClient, err = mongo.Connect(context.TODO(), options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Fatal(err)
	}
	// Define the port and create the server
	port := ":80"

	// Create a log file to capture all logs
	logFile, err := os.OpenFile("error.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("Failed to open log file: %v", err)
	}
	defer logFile.Close()

	// Set log output to the log file
	log.SetOutput(logFile)

	// Define a handler to serve the index.html file
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		htmlContent, err := os.ReadFile("templates/index.html")
		if err != nil {
			http.Error(w, "Could not read index.html", http.StatusInternalServerError)
			log.Printf("Failed to read index.html: %v", err)
			return
		}
		w.Write(htmlContent)
	})

	http.HandleFunc("/api/saveIP", saveIPHandler)

	// Create the server
	// server := &http.Server{
	// 	Addr:    ":" + port,
	// 	Handler: http.DefaultServeMux,
	// }

	log.Fatal(http.ListenAndServe(port, nil))

	// Start the server in a goroutine
	// go func() {
	// 	log.Printf("Starting server on port %s...\n", port)
	// 	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
	// 		log.Fatalf("Server failed: %v", err)
	// 	}
	// }()

	// Set up channel to wait for an interrupt or termination signal
	// quit := make(chan os.Signal, 1)
	// signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	// <-quit // Wait for the signal
	// log.Println("Shutting down server...")

	// // Create a deadline to wait for active requests to finish
	// ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	// defer cancel()

	// if err := server.Shutdown(ctx); err != nil {
	// 	log.Fatalf("Server forced to shutdown: %v", err)
	// }

	// log.Println("Server exiting")
}
