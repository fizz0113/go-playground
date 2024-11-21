package main

import (
	"context"
	"flag"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/bson"
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

	// 列出資料庫名稱
	dbFilter := bson.M{
		// "name": []string{"admin", "local"},
	}
	dbNames, err := client.ListDatabaseNames(ctx, dbFilter)
	if err != nil {
		panic(err)
	}
	log.Println("dbNames:", dbNames)

	// 列出資料庫
	databases, err := client.ListDatabases(ctx, dbFilter)
	if err != nil {
		panic(err)
	}
	log.Println("databases:", databases)

	// 建立資料庫 (這操作不會直接進 mongodb, 之後操作才會)
	dbName := "test"
	db := client.Database(dbName)

	// 刪除資料庫
	// db.Drop(ctx)

	// 之後操作用這個實例
	_ = db
}
