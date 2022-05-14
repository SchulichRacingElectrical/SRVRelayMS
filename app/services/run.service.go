package services

import (
	"bytes"
	"context"
	"database-ms/app/models"
	"database-ms/config"
	"database-ms/databases"
	"database-ms/utils"
	"errors"
	"fmt"
	"io/ioutil"
	"mime/multipart"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/gridfs"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type RunServiceI interface {
	CreateRun(context.Context, *models.Run) error
	FindById(context.Context, string) (*models.Run, error)
	GetRunsByThingId(context.Context, string) ([]*models.Run, error)
	UpdateRun(context.Context, *models.RunUpdate) error
	DeleteRun(context.Context, string) error
	GetRunFileMetaData(context.Context, string) (*models.RunFileMetaData, error)
	UploadFile(context.Context, *models.RunFileUpload, *multipart.FileHeader) error
	DownloadFile(context.Context, string) ([]byte, error)
}

type RunService struct {
	config *config.Configuration
}

func NewRunService(c *config.Configuration) RunServiceI {
	return &RunService{config: c}
}

func (service *RunService) CreateRun(ctx context.Context, run *models.Run) error {
	database, err := databases.GetDatabase(service.config.AtlasUri, service.config.MongoDbName, ctx)
	if err != nil {
		panic(err)
	}
	defer database.Client().Disconnect(ctx)

	run.ID = primitive.NewObjectID()
	// Check if Thing exists
	res := database.Collection("Thing").FindOne(ctx, bson.M{"_id": run.ThingID})
	if res.Err() == mongo.ErrNoDocuments {
		return errors.New("thing does not exist")
	}

	// Check of Session exists
	res = database.Collection("Session").FindOne(ctx, bson.M{"_id": run.SessionId})
	if res.Err() == mongo.ErrNoDocuments {
		return errors.New("session does not exist")
	}

	_, err = database.Collection("Run").InsertOne(ctx, run)
	return err
}

func (service *RunService) FindById(ctx context.Context, runId string) (*models.Run, error) {
	database, err := databases.GetDatabase(service.config.AtlasUri, service.config.MongoDbName, ctx)
	if err != nil {
		panic(err)
	}
	defer database.Client().Disconnect(ctx)

	var run models.Run
	bsonRunId, err := primitive.ObjectIDFromHex(runId)
	if err != nil {
		return nil, err
	}

	err = database.Collection("Run").FindOne(ctx, bson.M{"_id": bsonRunId}).Decode(&run)
	if err != nil {
		return nil, err
	}

	return &run, err
}

func (service *RunService) GetRunsByThingId(ctx context.Context, thingId string) ([]*models.Run, error) {
	database, err := databases.GetDatabase(service.config.AtlasUri, service.config.MongoDbName, ctx)
	if err != nil {
		panic(err)
	}
	defer database.Client().Disconnect(ctx)

	bsonThingId, err := primitive.ObjectIDFromHex(thingId)
	if err != nil {
		return nil, err
	}

	var runs []*models.Run
	cursor, err := database.Collection("Run").Find(ctx, bson.M{"thingId": bsonThingId})
	if err != nil {
		return nil, err
	}

	if err = cursor.All(ctx, &runs); err != nil {
		return nil, err
	}

	return runs, nil
}

func (service *RunService) UpdateRun(ctx context.Context, updatedRun *models.RunUpdate) error {
	database, err := databases.GetDatabase(service.config.AtlasUri, service.config.MongoDbName, ctx)
	if err != nil {
		panic(err)
	}
	defer database.Client().Disconnect(ctx)

	// Check if Thing exists
	res := database.Collection("Thing").FindOne(ctx, bson.M{"_id": updatedRun.ThingID})
	if res.Err() == mongo.ErrNoDocuments {
		return errors.New("thing does not exist")
	}

	// Check of Session exists
	res = database.Collection("Session").FindOne(ctx, bson.M{"_id": updatedRun.SessionId})
	if res.Err() == mongo.ErrNoDocuments {
		return errors.New("session does not exist")
	}

	_, err = database.Collection("Run").UpdateOne(ctx, bson.M{"_id": updatedRun.ID}, bson.M{"$set": updatedRun})

	return err
}

func (service *RunService) DeleteRun(ctx context.Context, runId string) error {
	bsonRunId, err := primitive.ObjectIDFromHex(runId)
	if err != nil {
		return err
	}

	client, err := databases.GetDBClient(service.config.AtlasUri, ctx)
	if err != nil {
		return err
	}
	defer client.Disconnect(ctx)

	callback := func(sessCtx mongo.SessionContext) (interface{}, error) {
		db := client.Database(service.config.MongoDbName)

		// Delete related comments
		commentFilter := bson.M{"associatedId": bsonRunId, "type": utils.Run}
		if _, err := db.Collection("Comment").DeleteMany(ctx, commentFilter); err != nil {
			return nil, err
		}

		// Delete run
		_, err := db.Collection("Run").DeleteOne(ctx, bson.M{"_id": bsonRunId})
		return nil, err
	}

	_, err = databases.WithTransaction(client, ctx, callback)
	return err
}

func (service *RunService) GetRunFileMetaData(ctx context.Context, runId string) (*models.RunFileMetaData, error) {
	database, err := databases.GetDatabase(service.config.AtlasUri, service.config.MongoDbName, ctx)
	if err != nil {
		panic(err)
	}
	defer database.Client().Disconnect(ctx)

	bsonRunId, err := primitive.ObjectIDFromHex(runId)
	if err != nil {
		return nil, err
	}

	var runFileMetaData models.RunFileMetaData
	err = database.Collection("fs.files").FindOne(ctx, bson.M{"_id": bsonRunId}).Decode(&runFileMetaData)

	// TODO for now the id of the file is the same as the run id
	// err = database.Collection("fs.files").FindOne(
	// 	ctx,
	// 	bson.D{
	// 		{"metadata", bson.D{
	// 			{"runId", bsonRunId},
	// 		}},
	// 	}).Decode(&runFileMetaData)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}

	return &runFileMetaData, nil

}

func (service *RunService) UploadFile(ctx context.Context, metadata *models.RunFileUpload, file *multipart.FileHeader) error {
	fileContent, err := file.Open()
	if err != nil {
		return err
	}

	byteContainer, err := ioutil.ReadAll(fileContent)
	if err != nil {
		return err
	}

	database, err := databases.GetDatabase(service.config.AtlasUri, service.config.MongoDbName, ctx)
	if err != nil {
		panic(err)
	}
	defer database.Client().Disconnect(ctx)

	bucket, err := gridfs.NewBucket(
		database,
	)
	if err != nil {
		return err
	}

	opts := options.GridFSUpload()
	opts.SetMetadata(metadata)
	uploadStream, err := bucket.OpenUploadStreamWithID(
		metadata.RunId,
		file.Filename,
		opts,
	)
	if err != nil {
		return err
	}
	defer uploadStream.Close()

	fileSize, err := uploadStream.Write(byteContainer)
	if err != nil {
		return err
	}

	fmt.Printf("Write file to DB was successful. File size: %d M\n", fileSize)
	return nil
}

func (service *RunService) DownloadFile(ctx context.Context, runId string) ([]byte, error) {
	database, err := databases.GetDatabase(service.config.AtlasUri, service.config.MongoDbName, ctx)
	if err != nil {
		panic(err)
	}
	defer database.Client().Disconnect(ctx)

	bsonRunId, err := primitive.ObjectIDFromHex(runId)
	if err != nil {
		return nil, err
	}

	var result bson.M
	err = database.Collection("fs.files").FindOne(ctx, bson.M{"_id": bsonRunId}).Decode(&result)
	if err != nil {
		return nil, err
	}

	bucket, _ := gridfs.NewBucket(database)

	var buf bytes.Buffer
	_, err = bucket.DownloadToStream(bsonRunId, &buf)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
