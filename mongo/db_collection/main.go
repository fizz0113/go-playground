package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"mongo/models"

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

	// 列出集合名稱
	collectionFilter := bson.M{
		// "name": "exampleCollectionName",
	}
	collectionNames, err := db.ListCollectionNames(ctx, collectionFilter)
	if err != nil {
		panic(err)
	}
	log.Println("collectionNames:", collectionNames)

	// 列出集合
	collections, err := db.ListCollections(ctx, collectionFilter)
	if err != nil {
		panic(err)
	}
	log.Println("collections:", collections)

	// 建立集合 (如果用 CreateCollection 會直接進 mongodb, 直接建立實例的話也可以, 之後寫入時會自動建)
	collectionName := "users"
	// err = db.CreateCollection(ctx, collectionName) // 直接在 mongodb 上建立
	// if err != nil {
	// 	panic(err)
	// }
	collection := db.Collection(collectionName)

	// 刪除集合
	// collection.Drop(ctx)

	// 插入一條 data
	data := models.User{Name: "ian", Email: "ian@example.com", Role: "admin"}
	insertDataResult, err := collection.InsertOne(ctx, data)
	if err != nil {
		panic(err)
	}
	log.Println("insertDataResult:", insertDataResult)

	// 插入很多 data
	datas := bson.A{
		models.User{Name: "ian001", Email: "ian001@example.com", Role: "user"},
		models.User{Name: "ian002", Email: "ian002@example.com", Role: "user"},
		models.User{Name: "ian003", Email: "ian003@example.com", Role: "user"},
	}
	insertDatasResult, err := collection.InsertMany(ctx, datas)
	if err != nil {
		panic(err)
	}
	log.Println("insertDatasResult:", insertDatasResult)

	// 查詢一條
	dataFilter := bson.M{
		"name": "ian002",
	}
	expData := models.User{}
	err = collection.FindOne(ctx, dataFilter).Decode(&expData)
	if err != nil {
		panic(err)
	}
	log.Println("expData:", expData)

	// 查詢多條
	datasFilter := bson.M{
		"role": "user",
	}
	expDatas := []models.User{}
	findCursor, err := collection.Find(ctx, datasFilter)
	if err != nil {
		panic(err)
	}
	// defer findCursor.Close(ctx)
	// for findCursor.Next(ctx) {
	// 	data := models.User{}
	// 	if err := findCursor.Decode(&data); err != nil {
	// 		panic(err)
	// 	}
	// 	expDatas = append(expDatas, data)
	// }
	err = findCursor.All(ctx, &expDatas)
	if err != nil {
		panic(err)
	}
	log.Println("expDatas:", expDatas)
}
