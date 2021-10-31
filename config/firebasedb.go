package config

import (
	"context"
	"go/build"
	"os"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
	"google.golang.org/api/option"
)

type FirebaseDB struct {
	Context context.Context
	App     *firebase.App
	Client  *firestore.Client
}

func (f *FirebaseDB) FirestoreDBInit() error {
	f.Context = context.Background()
	gopath := os.Getenv("GOPATH")
	if gopath == "" {
		gopath = build.Default.GOPATH
	}

	sa := option.WithCredentialsFile(gopath + "/src/github.com/SchulichRacingElectrical/srv-database-ms/config/firebase_config.json")
	app, err := firebase.NewApp(f.Context, nil, sa)
	if err != nil {
		return err
	}

	f.App = app

	client, err := f.App.Firestore(f.Context)
	if err != nil {
		return err
	}

	f.Client = client
	// defer f.Client.Close() //TODO: close client later when server shutsdown

	return nil
}
