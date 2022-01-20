package databases

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type MongoDB struct {
	Context context.Context
	Client  *mongo.Client
	Db      *mongo.Database
}

func (db *MongoDB) Init(dbUri string, dbName string) {
	// client, err := mongo.NewClient(options.Client().ApplyURI(dbUri))
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(dbUri))
	if err != nil {
		panic(err)
	}
	db.Client = client

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second) // Not sure about the timeout
	if err != nil {
		panic(err)
	}
	db.Context = ctx

	// Ping the primary
	if err := client.Ping(context.TODO(), readpref.Primary()); err != nil {
		panic(err)
	}

	fmt.Println("Successfully connected and pinged")

	database := client.Database(dbName)
	db.Db = database
}

func (db *MongoDB) Close() {
	db.Client.Disconnect(db.Context)
}
