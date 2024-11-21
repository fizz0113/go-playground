package main

import (
	"context"
	"flag"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	mongodbServerHost = flag.String("mongodbServerHost", "mongodb", "MongoDB server host")
	mongodbServerPort = flag.Int("mongodbServerPort", 27017, "MongoDB server port")
	mongodbUsername   = flag.String("mongodbUsername", "admin", "MongoDB username")
	mongodbPassword   = flag.String("mongodbPassword", "admin", "MongoDB password")
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// 設置 MongoDB 連接選項
	clientOpts := options.Client()
	clientOpts.SetMaxPoolSize(20)

	// 設置連接 URI
	uri := fmt.Sprintf("mongodb://%s:%d", *mongodbServerHost, *mongodbServerPort)
	clientOpts.ApplyURI(uri)

	// 設置認證
	cred := options.Credential{
		Username: *mongodbUsername,
		Password: *mongodbPassword,
		// AuthSource:    "admin",
		// AuthMechanism: "SCRAM-SHA-256",
	}
	clientOpts.SetAuth(cred)

	// 連接到 mongodb
	client, err := mongo.Connect(ctx, clientOpts)
	if err != nil {
		panic(err)
	}
	defer client.Disconnect(ctx)

	// 檢查連線
	err = client.Ping(ctx, nil)
	if err != nil {
		panic(err)
	}
	log.Println("Connected to MongoDB!")
}
