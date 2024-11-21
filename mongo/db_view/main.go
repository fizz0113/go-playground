package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"mongo/models"
	"slices"

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

	// 建立 view
	viewName := "find_role_user"
	viewOn := "users"
	opts := options.CreateView()
	// opts.SetCollation(&options.Collation{})
	pipeline := mongo.Pipeline{
		{
			{
				Key: "$match",
				Value: bson.D{
					{
						Key: "role",
						Value: bson.D{
							{Key: "$eq", Value: "user"},
						},
					},
				},
			},
		},
		{
			{
				Key: "$project",
				Value: bson.D{
					{Key: "name", Value: 1},
					{Key: "email", Value: 1},
					{Key: "_id", Value: 0},
				},
			},
		},
	}

	// 檢查 view 存不存在
	collectionNames, err := db.ListCollectionNames(ctx, bson.M{})
	if err != nil {
		panic(err)
	}
	if !slices.Contains(collectionNames, viewName) {
		// 建立 view
		err = db.CreateView(ctx, viewName, viewOn, pipeline, opts)
		if err != nil {
			panic(err)
		}
	}

	// 直接用建立好的 view 查詢
	findRoleUserView := db.Collection(viewName)
	expDatas := []models.User{}
	findCursor, err := findRoleUserView.Find(ctx, bson.D{})
	if err != nil {
		panic(err)
	}
	err = findCursor.All(ctx, &expDatas)
	if err != nil {
		panic(err)
	}
	log.Println("expDatas:", expDatas)
}
