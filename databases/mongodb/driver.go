package mongo_driver

import (
	"context"
	"fmt"
	"time"

	_config "go-articles/server/config"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func GetConnection() string {
	return fmt.Sprintf("mongodb://%v:%v",
		_config.GetConfiguration("mongo.host"),
		_config.GetConfiguration("mongo.port"),
	)
}

func Init() *mongo.Database {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(GetConnection()))
	if err != nil {
		panic(err)
	}

	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		panic(err)
	}

	db := client.Database(_config.GetConfiguration("mongo.name"))

	// _, err = db.Collection("logger").InsertOne(context.Background(), bson.M{"name": "Wick"})
	// fmt.Println("background:",context.Background())
	// if err != nil {
	// 	fmt.Println("error insert", err.Error())
	// }
	return db
}
