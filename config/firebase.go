package config

import (
	"context"
	"go/build"
	"os"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
	"firebase.google.com/go/auth"
	"google.golang.org/api/option"
)

type Firebase struct {
	Context context.Context
	App     *firebase.App
	Auth    *auth.Client
	Client  *firestore.Client
}

func (f *Firebase) Init() error {
	// Configure the context
	f.Context = context.Background()
	gopath := os.Getenv("GOPATH")
	if gopath == "" {
		gopath = build.Default.GOPATH
	}

	// Configure the app
	sa := option.WithCredentialsFile(gopath + "/src/github.com/SchulichRacingElectrical/srv-database-ms/config/firebase_config.json")
	app, err := firebase.NewApp(f.Context, nil, sa)
	if err != nil {
		return err
	}
	f.App = app

	// Configure firestore
	client, err := f.App.Firestore(f.Context)
	if err != nil {
		return err
	}
	f.Client = client

	// Configure auth
	auth, err := f.App.Auth(f.Context)
	if err != nil {
		return err
	}
	f.Auth = auth

	// defer f.Client.Close() //TODO: close client later when server shutsdown

	return nil
}
