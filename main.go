package main

import (
	"context"
	"log"
	"net/http"

	"github.com/rs/cors"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"Go-mongo/collections"
	"Go-mongo/middlewares"
	"Go-mongo/routers"
)

func main() {
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	// Use your middleware here
	authMiddleware := middlewares.BasicAuth("First", "112233")

	r := routers.InitializeRoutes(client)

	// Apply the middleware to authenticated routes
	authRoutes := r.PathPrefix("/api").Subrouter()
	authRoutes.Use(authMiddleware)
	authRoutes.HandleFunc("/people", collections.CreatePerson(client.Database("Go-mongo").Collection("go-mongo"))).Methods("POST")
	authRoutes.HandleFunc("/people/{id}", collections.UpdatePerson(client.Database("Go-mongo").Collection("go-mongo"))).Methods("PUT")
	authRoutes.HandleFunc("/people/{id}", collections.DeletePerson(client.Database("Go-mongo").Collection("go-mongo"))).Methods("DELETE")

	// CORS handling
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowCredentials: true,
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE"},
		AllowedHeaders:   []string{"Authorization", "Content-Type"},
	})

	handler := c.Handler(r)

	log.Println("Server is running on port 8000")
	log.Fatal(http.ListenAndServe(":8000", handler))
}
