package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

const (
	rootUser     = "root"
	rootPassword = "password"
	dbName       = "time_api"
)

var db *sql.DB

// TimeResponse represents the JSON structure for the API response
type TimeResponse struct {
	CurrentTime string `json:"current_time"`
}

// LoggedTimesResponse represents the JSON structure for logged times
type LoggedTimesResponse struct {
	Timestamps []string `json:"timestamps"`
}

func main() {
	// Initialize MySQL root connection to check for database existence
	rootDSN := fmt.Sprintf("%s:%s@tcp(127.0.0.1:3306)/", rootUser, rootPassword)
	rootDB, err := sql.Open("mysql", rootDSN)
	if err != nil {
		log.Fatalf("Failed to connect to MySQL: %v", err)
	}
	defer rootDB.Close()

	// Check if the database exists
	exists, err := databaseExists(rootDB, dbName)
	if err != nil {
		log.Fatalf("Failed to check database existence: %v", err)
	}

	if !exists {
		// Create the database if it does not exist
		if _, err := rootDB.Exec(fmt.Sprintf("CREATE DATABASE %s", dbName)); err != nil {
			log.Fatalf("Failed to create database: %v", err)
		}
		log.Printf("Database %s created successfully", dbName)
	}

	// Connect to the specific database
	dbDSN := fmt.Sprintf("%s:%s@tcp(127.0.0.1:3306)/%s", rootUser, rootPassword, dbName)
	db, err = sql.Open("mysql", dbDSN)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Create the time_log table if it doesn't exist
	if _, err := db.Exec(`CREATE TABLE IF NOT EXISTS time_log (
		id INT AUTO_INCREMENT PRIMARY KEY,
		timestamp DATETIME NOT NULL
	)`); err != nil {
		log.Fatalf("Failed to create time_log table: %v", err)
	}

	// Setup routes
	http.HandleFunc("/current-time", currentTimeHandler)
	http.HandleFunc("/logged-times", loggedTimesHandler)

	// Start server
	log.Println("Server running on port 8080...")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}

func databaseExists(db *sql.DB, name string) (bool, error) {
	query := "SELECT SCHEMA_NAME FROM INFORMATION_SCHEMA.SCHEMATA WHERE SCHEMA_NAME = ?"
	var schemaName string
	err := db.QueryRow(query, name).Scan(&schemaName)
	if err == sql.ErrNoRows {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}

func currentTimeHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Received request at /current-time endpoint")

	// Get current time in Toronto timezone
	location, err := time.LoadLocation("America/Toronto")
	if err != nil {
		http.Error(w, "Failed to load timezone", http.StatusInternalServerError)
		log.Printf("Error loading timezone: %v", err)
		return
	}
	torontoTime := time.Now().In(location)

	// Insert time into the database
	if err := logTimeToDB(torontoTime); err != nil {
		http.Error(w, "Failed to log time to database", http.StatusInternalServerError)
		log.Printf("Error logging time to database: %v", err)
		return
	}

	// Respond with the current time in JSON format
	response := TimeResponse{CurrentTime: torontoTime.Format(time.RFC3339)}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		log.Printf("Error encoding response: %v", err)
	}
}

func loggedTimesHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Received request at /logged-times endpoint")

	// Retrieve all logged times from the database
	rows, err := db.Query("SELECT timestamp FROM time_log")
	if err != nil {
		http.Error(w, "Failed to retrieve logged times", http.StatusInternalServerError)
		log.Printf("Error retrieving logged times: %v", err)
		return
	}
	defer rows.Close()

	var times []string
	for rows.Next() {
		var timestamp string // Use string to handle DATETIME from MySQL
		if err := rows.Scan(&timestamp); err != nil {
			http.Error(w, "Error scanning database rows", http.StatusInternalServerError)
			log.Printf("Error scanning database rows: %v", err)
			return
		}
		times = append(times, timestamp)
	}

	// Respond with the logged times in JSON format
	response := LoggedTimesResponse{Timestamps: times}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		log.Printf("Error encoding response: %v", err)
	}
}

func logTimeToDB(t time.Time) error {
	query := "INSERT INTO time_log (timestamp) VALUES (?)"
	_, err := db.Exec(query, t)
	return err
}
