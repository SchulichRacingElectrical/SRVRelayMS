package databases

import (
	"context"
	"go/build"
	"os"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
	"firebase.google.com/go/auth"
	"google.golang.org/api/option"
)

type FirebaseDB struct {
	Context context.Context
	App     *firebase.App
	Auth 		*auth.Client
	Client  *firestore.Client
}

func (db *FirebaseDB) Init() error {
	// Configure the context
	db.Context = context.Background()
	gopath := os.Getenv("GOPATH")
	if gopath == "" {
		gopath = build.Default.GOPATH
	}

	// Configure the app
	sa := option.WithCredentialsFile(gopath + "/src/github.com/SchulichRacingElectrical/srv-database-ms/config/firebase_config.json")
	app, err := firebase.NewApp(db.Context, nil, sa)
	if err != nil {
		return err
	}
	db.App = app

	// Configure firestore
	client, err := app.Firestore(db.Context)
	if err != nil {
		return err
	}
	db.Client = client

	// Configure auth
	auth, err := db.App.Auth(db.Context)
	if err != nil {
		return err
	}
	db.Auth = auth

	return nil
}

func (db *FirebaseDB) Close() {
	db.Client.Close()
}
