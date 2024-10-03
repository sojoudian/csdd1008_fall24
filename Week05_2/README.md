
# Go CRUD API

This is a simple RESTful CRUD (Create, Read, Update, Delete) API built using Go version 1.23 without any third-party packages. It uses Go's built-in HTTP server, JSON encoding/decoding, and concurrency mechanisms like mutexes for handling in-memory product storage.

## Table of Contents
- [Features](#features)
- [Project Structure](#project-structure)
- [Routes](#routes)
- [Installation and Running](#installation-and-running)
- [API Usage](#api-usage)
- [Code Explanation](#code-explanation)

## Features
This API provides basic CRUD operations for managing products:
1. **Create** a new product.
2. **Read** a single product by ID or list all products.
3. **Update** an existing product by ID.
4. **Delete** a product by ID.

## Project Structure
```
go-crud-api/
|-- main.go
```

## Routes
The API exposes the following endpoints:
- `GET /products` - Retrieve all products.
- `POST /products` - Create a new product.
- `GET /products/{id}` - Retrieve a product by ID.
- `PUT /products/{id}` - Update a product by ID.
- `DELETE /products/{id}` - Delete a product by ID.

## Installation and Running

### Prerequisites:
- Go version 1.23 installed on your machine.

### Running the server:
1. Clone the repository or copy the `main.go` file.
2. Open a terminal and navigate to the project directory.
3. Run the following command to start the server:
   ```bash
   go run main.go
   ```

The server will start on `http://localhost:8080`.

## API Usage

### 1. Create a Product (POST /products)
- **Method**: POST
- **Endpoint**: `/products`
- **Request Body** (JSON):
  ```json
  {
    "name": "Product Name",
    "price": 100
  }
  ```
- **Response** (201 Created):
  ```json
  {
    "id": 1,
    "name": "Product Name",
    "price": 100
  }
  ```

### 2. Get All Products (GET /products)
- **Method**: GET
- **Endpoint**: `/products`
- **Response** (200 OK):
  ```json
  [
    {
      "id": 1,
      "name": "Product Name",
      "price": 100
    }
  ]
  ```

### 3. Get a Product by ID (GET /products/{id})
- **Method**: GET
- **Endpoint**: `/products/{id}`
- **Response** (200 OK):
  ```json
  {
    "id": 1,
    "name": "Product Name",
    "price": 100
  }
  ```

### 4. Update a Product by ID (PUT /products/{id})
- **Method**: PUT
- **Endpoint**: `/products/{id}`
- **Request Body** (JSON):
  ```json
  {
    "name": "Updated Product Name",
    "price": 120
  }
  ```
- **Response** (200 OK):
  ```json
  {
    "id": 1,
    "name": "Updated Product Name",
    "price": 120
  }
  ```

### 5. Delete a Product by ID (DELETE /products/{id})
- **Method**: DELETE
- **Endpoint**: `/products/{id}`
- **Response** (204 No Content)

## Code Explanation

### 1. Product Struct
The `Product` struct represents the structure of a product in our application. It contains the following fields:
- `ID`: A unique identifier for each product.
- `Name`: The name of the product.
- `Price`: The price of the product.

### 2. In-Memory Storage
We use an in-memory map (`products`) to store the products. The `nextID` variable is used to auto-increment product IDs.

```go
var (
    products    = make(map[int]Product) // In-memory product storage
    nextID      = 1                     // Auto-incrementing product ID
    productsMux sync.Mutex              // Mutex for concurrent access
)
```

The `productsMux` mutex is used to ensure thread-safe access to the products map, as multiple requests can modify it concurrently.

### 3. Handlers
We define separate handlers for each CRUD operation.

- **Create Product** (`POST /products`): Handles product creation by decoding the request body and adding the product to the in-memory storage.
- **Get Product** (`GET /products/{id}`): Fetches a product by its ID.
- **Update Product** (`PUT /products/{id}`): Updates a product's details.
- **Delete Product** (`DELETE /products/{id}`): Removes a product by its ID.
- **List Products** (`GET /products`): Lists all available products.

### 4. HTTP Server
We use Go's built-in `net/http` package to set up an HTTP server and define routes.

- `http.HandleFunc("/products", ...)`: This defines the route for listing and creating products.
- `http.HandleFunc("/products/{id}", ...)`: This defines routes for getting, updating, and deleting products by ID.

The server is started using:
```go
log.Fatal(http.ListenAndServe(":8080", nil))
```

This binds the server to port 8080 and listens for incoming requests.

## Conclusion

This CRUD API is a simple demonstration of using Go 1.23's built-in packages to create a basic API for managing products. The API supports all CRUD operations and ensures thread safety with mutexes for concurrent access to in-memory data.


