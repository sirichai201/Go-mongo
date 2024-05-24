package collections

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func CreatePerson(collection *mongo.Collection) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		var person map[string]interface{}
		_ = json.NewDecoder(r.Body).Decode(&person)

		insertResult, err := collection.InsertOne(context.Background(), person)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		var insertedPerson bson.M
		filter := bson.M{"_id": insertResult.InsertedID}
		err = collection.FindOne(context.Background(), filter).Decode(&insertedPerson)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(insertedPerson)
	}
}

func GetPeople(collection *mongo.Collection) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		cursor, err := collection.Find(context.Background(), bson.M{})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer cursor.Close(context.Background())

		var people []bson.M
		if err = cursor.All(context.Background(), &people); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(people)
	}
}

func UpdatePerson(collection *mongo.Collection) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		params := mux.Vars(r)
		idParam := params["id"]
		id, err := primitive.ObjectIDFromHex(idParam)
		if err != nil {
			http.Error(w, "Invalid ID", http.StatusBadRequest)
			return
		}

		var updateData map[string]interface{}
		_ = json.NewDecoder(r.Body).Decode(&updateData)

		update := bson.M{"$set": updateData}
		filter := bson.M{"_id": id}
		_, err = collection.UpdateOne(context.Background(), filter, update)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		var updatedPerson bson.M
		err = collection.FindOne(context.Background(), filter).Decode(&updatedPerson)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(updatedPerson)
	}
}

func DeletePerson(collection *mongo.Collection) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		params := mux.Vars(r)
		idParam := params["id"]
		id, err := primitive.ObjectIDFromHex(idParam)
		if err != nil {
			http.Error(w, "Invalid ID", http.StatusBadRequest)
			return
		}

		filter := bson.M{"_id": id}
		result, err := collection.DeleteOne(context.Background(), filter)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(result.DeletedCount)
	}
}
