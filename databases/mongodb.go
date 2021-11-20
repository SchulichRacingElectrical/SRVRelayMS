package databases

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoDB struct {
	Context context.Context
	Client mongo.Client
}

func (db* MongoDB) Init() error {
	client, err := mongo.NewClient(options.Client().ApplyURI("TODO"))
	if err != nil {
		return err
	}

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second) // Not sure about the timeout
	err = client.Connect(ctx)
	if err != nil {
		return err
	}

	defer client.Disconnect(ctx)

	return nil
}

func (db *MongoDB) Close() {
	// Close somehow
}