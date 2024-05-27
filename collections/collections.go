package collections

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/dgrijalva/jwt-go"
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

func Login(collection *mongo.Collection) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		var credentials map[string]string
		if err := json.NewDecoder(r.Body).Decode(&credentials); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		username, ok := credentials["username"]
		if !ok {
			http.Error(w, "Username is required", http.StatusBadRequest)
			return
		}

		password, ok := credentials["password"]
		if !ok {
			http.Error(w, "Password is required", http.StatusBadRequest)
			return
		}

		filter := bson.M{"username": username, "password": password}
		var user bson.M
		if err := collection.FindOne(context.Background(), filter).Decode(&user); err != nil {
			http.Error(w, "Invalid username or password", http.StatusUnauthorized)
			return
		}

		// สร้าง JWT token สำหรับการยืนยันตัวตน
		token, err := generateJWTToken(username)
		if err != nil {
			http.Error(w, "Failed to generate token", http.StatusInternalServerError)
			return
		}

		// ส่ง token กลับไปยังผู้ใช้
		response := map[string]string{"token": token}
		json.NewEncoder(w).Encode(response)
	}
}

// generateJWTToken เป็นฟังก์ชั่นสำหรับสร้าง JWT token สำหรับการยืนยันตัวตน
func generateJWTToken(username string) (string, error) {
	// สร้าง claim ด้วย username เพื่อใช้ในการสร้าง token
	claims := jwt.MapClaims{
		"username": username,
		"exp":      time.Now().Add(time.Hour * 24).Unix(), // กำหนดเวลาหมดอายุของ token เป็น 24 ชั่วโมง
	}

	// สร้าง token ด้วย claims และใช้คีย์ลับสำหรับลงนาม
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte("Users")) // your_secret_key คือคีย์ลับที่คุณใช้สำหรับลงนาม token
	if err != nil {
		return "", err
	}

	return signedToken, nil
}

func Register(collection *mongo.Collection) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		var user map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			log.Println("Error decoding request body:", err)
			return
		}

		username := user["username"].(string)
		existingUser := bson.M{"username": username}
		if err := collection.FindOne(context.Background(), existingUser).Err(); err == nil {
			http.Error(w, "Username already exists", http.StatusBadRequest)
			log.Println("Username already exists:", username)
			return
		}

		insertResult, err := collection.InsertOne(context.Background(), user)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			log.Println("Error inserting user into database:", err)
			return
		}

		var newUser bson.M
		filter := bson.M{"_id": insertResult.InsertedID}
		if err := collection.FindOne(context.Background(), filter).Decode(&newUser); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			log.Println("Error finding new user in database:", err)
			return
		}

		json.NewEncoder(w).Encode(newUser)
	}
}
