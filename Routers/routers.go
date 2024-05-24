package routers

import (
	"Go-mongo/collections"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/mongo"
)

func InitializeRoutes(client *mongo.Client) *mux.Router {
	r := mux.NewRouter()
	peopleCollection := client.Database("Go-mongo").Collection("go-mongo")

	// Public route, no authentication required
	r.HandleFunc("/api/people", collections.GetPeople(peopleCollection)).Methods("GET")

	// Authenticated routes
	r.HandleFunc("/api/people", collections.CreatePerson(peopleCollection)).Methods("POST")
	r.HandleFunc("/api/people/{id}", collections.UpdatePerson(peopleCollection)).Methods("PUT")
	r.HandleFunc("/api/people/{id}", collections.DeletePerson(peopleCollection)).Methods("DELETE")

	return r
}
