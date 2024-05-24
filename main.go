package main

import (
	"fmt"
	"log"
	"net/http"

	"Go-mongo/middlewares"
	"Go-mongo/modules"
	"Go-mongo/routers"

	"github.com/gorilla/mux"
)

func main() {
	client, ctx, cancel, err := modules.ConnectToMongoDB("mongodb://localhost:27017")
	if err != nil {
		log.Fatal(err)
	}
	defer cancel()
	defer client.Disconnect(ctx)

	r := mux.NewRouter()
	apiRouter := routers.InitializeRoutes(client)

	// Apply the BasicAuth middleware to all routes except the GetPeople route
	apiRouter.Use(middlewares.BasicAuth("First", "112233"))
	r.PathPrefix("/api").Handler(apiRouter)

	fmt.Println("Server is running on port 8000")
	log.Fatal(http.ListenAndServe(":8000", r))
}
