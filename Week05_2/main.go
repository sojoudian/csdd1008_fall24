package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"sync"
)

// Product struct represents a basic product entity
type Product struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Price int    `json:"price"`
}

var (
	products    = make(map[int]Product) // In-memory product storage
	nextID      = 1                     // Auto-incrementing product ID
	productsMux sync.Mutex              // Mutex for concurrent access
)

// Handler for creating a new product (POST /products)
func createProductHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var newProduct Product
	err := json.NewDecoder(r.Body).Decode(&newProduct)
	if err != nil {
		http.Error(w, "Error parsing request body", http.StatusBadRequest)
		return
	}

	productsMux.Lock()
	defer productsMux.Unlock()

	newProduct.ID = nextID
	nextID++
	products[newProduct.ID] = newProduct

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newProduct)
}

// Handler for getting a single product by ID (GET /products/{id})
func getProductHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Path[len("/products/"):]
	id, err := strconv.Atoi(idStr)
	if err != nil || id < 1 {
		http.Error(w, "Invalid product ID", http.StatusBadRequest)
		return
	}

	productsMux.Lock()
	defer productsMux.Unlock()

	product, exists := products[id]
	if !exists {
		http.Error(w, "Product not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(product)
}

// Handler for updating a product by ID (PUT /products/{id})
func updateProductHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	idStr := r.URL.Path[len("/products/"):]
	id, err := strconv.Atoi(idStr)
	if err != nil || id < 1 {
		http.Error(w, "Invalid product ID", http.StatusBadRequest)
		return
	}

	var updatedProduct Product
	err = json.NewDecoder(r.Body).Decode(&updatedProduct)
	if err != nil {
		http.Error(w, "Error parsing request body", http.StatusBadRequest)
		return
	}

	productsMux.Lock()
	defer productsMux.Unlock()

	_, exists := products[id]
	if !exists {
		http.Error(w, "Product not found", http.StatusNotFound)
		return
	}

	updatedProduct.ID = id
	products[id] = updatedProduct

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updatedProduct)
}

// Handler for deleting a product by ID (DELETE /products/{id})
func deleteProductHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	idStr := r.URL.Path[len("/products/"):]
	id, err := strconv.Atoi(idStr)
	if err != nil || id < 1 {
		http.Error(w, "Invalid product ID", http.StatusBadRequest)
		return
	}

	productsMux.Lock()
	defer productsMux.Unlock()

	_, exists := products[id]
	if !exists {
		http.Error(w, "Product not found", http.StatusNotFound)
		return
	}

	delete(products, id)

	w.WriteHeader(http.StatusNoContent)
}

// Handler for listing all products (GET /products)
func listProductsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	productsMux.Lock()
	defer productsMux.Unlock()

	var productList []Product
	for _, product := range products {
		productList = append(productList, product)
	}

	json.NewEncoder(w).Encode(productList)
}

// Main function to set up the routes and start the HTTP server
func main() {
	http.HandleFunc("/products", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			listProductsHandler(w, r)
		case http.MethodPost:
			createProductHandler(w, r)
		default:
			http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		}
	})

	http.HandleFunc("/products/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			getProductHandler(w, r)
		case http.MethodPut:
			updateProductHandler(w, r)
		case http.MethodDelete:
			deleteProductHandler(w, r)
		default:
			http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		}
	})

	fmt.Println("Server is running on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

