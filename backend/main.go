package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

// Response structure
type Response struct {
	Message string `json:"message"`
	Visits  int    `json:"visits"`
	Error   string `json:"error,omitempty"`
}

// Get environment variable or default
func getenv(key, def string) string {
	val := os.Getenv(key)
	if val == "" {
		return def
	}
	return val
}

// Create a new DB connection with retries
func getDBConnection() (*sql.DB, error) {
	host := getenv("DB_HOST", "Pass me from docker-compose.yaml")
	user := getenv("DB_USER", "Pass me from docker-compose.yaml")
	pass := getenv("DB_PASSWORD", "Pass me from docker-compose.yaml")
	dbName := getenv("DB_NAME", "Pass me from docker-compose.yaml")

	dsn := fmt.Sprintf("%s:%s@tcp(%s:3306)/%s", user, pass, host, dbName)

	var db *sql.DB
	var err error
	retries := 10

	for i := 1; i <= retries; i++ {
		db, err = sql.Open("mysql", dsn)
		if err == nil {
			err = db.Ping()
		}
		if err == nil {
			log.Printf("Connected to MySQL (attempt %d)", i)
			return db, nil
		}
		log.Printf("MySQL connection failed (attempt %d/%d): %v", i, retries, err)
		time.Sleep(3 * time.Second)
	}

	return nil, fmt.Errorf("failed to connect to DB after %d attempts: %v", retries, err)
}

// HTTP handler
func handler(w http.ResponseWriter, r *http.Request) {
	// Allow CORS for all all
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "*")
	w.Header().Set("Access-Control-Allow-Headers", "*")

	// Handle preflight request
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	db, err := getDBConnection()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(Response{Error: err.Error()})
		return
	}
	defer db.Close()

	// Fetch message
	var message string
	err = db.QueryRow("SELECT text FROM messages LIMIT 1;").Scan(&message)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(Response{Error: err.Error()})
		return
	}

	// Increment visit count
	_, err = db.Exec("UPDATE visits SET count = count + 1 WHERE id = 1;")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(Response{Error: err.Error()})
		return
	}

	// Get current visit count
	var visits int
	err = db.QueryRow("SELECT count FROM visits WHERE id = 1;").Scan(&visits)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(Response{Error: err.Error()})
		return
	}

	// Return JSON
	resp := Response{Message: message, Visits: visits}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func main() {
	http.HandleFunc("/", handler)
	log.Println("Backend server running on port 9000")
	log.Fatal(http.ListenAndServe(":9000", nil))
}
