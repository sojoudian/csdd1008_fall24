package main

import (
	"fmt"
	"net/http"
)

type Message struct {
	Name    string `json:"name"`
	Content string `json:"content"`
}

func getHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}
	fmt.Fprintf(w, "Welcome to Go API!")
}

func main() {
	// p()
	http.HandleFunc("/get", getHandler)
	http.ListenAndServe(":8090", nil)
}
