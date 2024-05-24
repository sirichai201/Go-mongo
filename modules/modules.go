package modules

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func ConnectToMongoDB(uri string) (*mongo.Client, context.Context, context.CancelFunc, error) {
	clientOptions := options.Client().ApplyURI(uri)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, nil, nil, err
	}
	if err := client.Ping(ctx, nil); err != nil {
		return nil, nil, nil, err
	}
	fmt.Println("Connected to MongoDB!")
	return client, ctx, cancel, nil
}
