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
	r.HandleFunc("/api/login", collections.Login(client.Database("Go-mongo").Collection("go-mongo"))).Methods("POST")
	r.HandleFunc("/api/register", collections.Register(client.Database("Go-mongo").Collection("go-mongo"))).Methods("POST")

	return r
}
