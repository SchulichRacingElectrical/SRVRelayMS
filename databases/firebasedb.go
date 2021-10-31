package databases

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

func (db *FirebaseDB) Init() error {
	db.Context = context.Background()
	gopath := os.Getenv("GOPATH")
	if gopath == "" {
		gopath = build.Default.GOPATH
	}

	sa := option.WithCredentialsFile(gopath + "/src/github.com/SchulichRacingElectrical/srv-database-ms/config/firebase_config.json")
	app, err := firebase.NewApp(db.Context, nil, sa)
	if err != nil {
		return err
	}

	db.App = app

	client, err := app.Firestore(db.Context)
	if err != nil {
		return err
	}

	db.Client = client
	return nil
}

func (db *FirebaseDB) Close() {
	db.Client.Close()
}
