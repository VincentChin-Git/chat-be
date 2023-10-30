package storage

import (
	"client-server-go/config"
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var ClientDatabase *mongo.Database

func ConnectDatabase() {
	configGet := config.GetConfig()
	// Use the SetServerAPIOptions() method to set the Stable API version to 1
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(configGet.DatabaseURL).SetServerAPIOptions(serverAPI)
	// Create a new client and connect to the server
	client, err := mongo.Connect(context.TODO(), opts)
	ClientDatabase = client.Database("chat")

	if err != nil {
		panic(err)
	}

}
