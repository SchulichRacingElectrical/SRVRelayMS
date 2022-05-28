package services

import (
	"context"
	"database-ms/app/model"
	"database-ms/config"
	"mime/multipart"

	"github.com/google/uuid"
)

type SessionServiceInterface interface {
	CreateSession(context.Context, *model.Session) error
	FindById(context.Context, uuid.UUID) (*model.Session, error)
	GetSessionsByThingId(context.Context, uuid.UUID) ([]*model.Session, error)
	UpdateSession(context.Context, *model.Session) error
	DeleteSession(context.Context, uuid.UUID) error
	GetSessionFileMetaData(context.Context, uuid.UUID) (*model.Session, error)
	UploadFile(context.Context, *model.Session, *multipart.FileHeader) error
	DownloadFile(context.Context, uuid.UUID) ([]byte, error)
}

type SessionService struct {
	config *config.Configuration
}

func NewSessionService(c *config.Configuration) SessionServiceInterface {
	return &SessionService{config: c}
}

func (service *SessionService) CreateSession(ctx context.Context, session *model.Session) error {
	// database, err := databases.GetDatabase(service.config.AtlasUri, service.config.MongoDbName, ctx)
	// if err != nil {
	// 	panic(err)
	// }
	// defer database.Client().Disconnect(ctx)

	// run.ID = primitive.NewObjectID()
	// // Check if Thing exists
	// res := database.Collection("Thing").FindOne(ctx, bson.M{"_id": run.ThingID})
	// if res.Err() == mongo.ErrNoDocuments {
	// 	return errors.New("thing does not exist")
	// }

	// // Check of Session exists
	// res = database.Collection("Session").FindOne(ctx, bson.M{"_id": run.SessionId})
	// if res.Err() == mongo.ErrNoDocuments {
	// 	return errors.New("session does not exist")
	// }

	// _, err = database.Collection("Run").InsertOne(ctx, run)
	// return err
	return nil
}

func (service *SessionService) FindById(ctx context.Context, sessionId uuid.UUID) (*model.Session, error) {
	// database, err := databases.GetDatabase(service.config.AtlasUri, service.config.MongoDbName, ctx)
	// if err != nil {
	// 	panic(err)
	// }
	// defer database.Client().Disconnect(ctx)

	// var run models.Run
	// bsonRunId, err := primitive.ObjectIDFromHex(runId)
	// if err != nil {
	// 	return nil, err
	// }

	// err = database.Collection("Run").FindOne(ctx, bson.M{"_id": bsonRunId}).Decode(&run)
	// if err != nil {
	// 	return nil, err
	// }

	// return &run, err
	return nil, nil
}

func (service *SessionService) GetSessionsByThingId(ctx context.Context, thingId uuid.UUID) ([]*model.Session, error) {
	// database, err := databases.GetDatabase(service.config.AtlasUri, service.config.MongoDbName, ctx)
	// if err != nil {
	// 	panic(err)
	// }
	// defer database.Client().Disconnect(ctx)

	// bsonThingId, err := primitive.ObjectIDFromHex(thingId)
	// if err != nil {
	// 	return nil, err
	// }

	// var runs []*models.Run
	// cursor, err := database.Collection("Run").Find(ctx, bson.M{"thingId": bsonThingId})
	// if err != nil {
	// 	return nil, err
	// }

	// if err = cursor.All(ctx, &runs); err != nil {
	// 	return nil, err
	// }

	// return runs, nil
	return nil, nil
}

func (service *SessionService) UpdateSession(ctx context.Context, updatedSession *model.Session) error {
	// database, err := databases.GetDatabase(service.config.AtlasUri, service.config.MongoDbName, ctx)
	// if err != nil {
	// 	panic(err)
	// }
	// defer database.Client().Disconnect(ctx)

	// // Check if Thing exists
	// res := database.Collection("Thing").FindOne(ctx, bson.M{"_id": updatedRun.ThingID})
	// if res.Err() == mongo.ErrNoDocuments {
	// 	return errors.New("thing does not exist")
	// }

	// // Check of Session exists
	// res = database.Collection("Session").FindOne(ctx, bson.M{"_id": updatedRun.SessionId})
	// if res.Err() == mongo.ErrNoDocuments {
	// 	return errors.New("session does not exist")
	// }

	// _, err = database.Collection("Run").UpdateOne(ctx, bson.M{"_id": updatedRun.ID}, bson.M{"$set": updatedRun})

	// return err
	return nil
}

func (service *SessionService) DeleteSession(ctx context.Context, sessionId uuid.UUID) error {
	// bsonRunId, err := primitive.ObjectIDFromHex(runId)
	// if err != nil {
	// 	return err
	// }

	// client, err := databases.GetDBClient(service.config.AtlasUri, ctx)
	// if err != nil {
	// 	return err
	// }
	// defer client.Disconnect(ctx)

	// callback := func(sessCtx mongo.SessionContext) (interface{}, error) {
	// 	db := client.Database(service.config.MongoDbName)

	// 	// Delete related comments
	// 	commentFilter := bson.M{"associatedId": bsonRunId, "type": utils.Run}
	// 	if _, err := db.Collection("Comment").DeleteMany(ctx, commentFilter); err != nil {
	// 		return nil, err
	// 	}

	// 	// Delete run
	// 	_, err := db.Collection("Run").DeleteOne(ctx, bson.M{"_id": bsonRunId})
	// 	return nil, err
	// }

	// _, err = databases.WithTransaction(client, ctx, callback)
	// return err
	return nil
}

func (service *SessionService) GetSessionFileMetaData(ctx context.Context, sessionId uuid.UUID) (*model.Session, error) {
	// database, err := databases.GetDatabase(service.config.AtlasUri, service.config.MongoDbName, ctx)
	// if err != nil {
	// 	panic(err)
	// }
	// defer database.Client().Disconnect(ctx)

	// bsonRunId, err := primitive.ObjectIDFromHex(runId)
	// if err != nil {
	// 	return nil, err
	// }

	// var runFileMetaData models.RunFileMetaData
	// err = database.Collection("fs.files").FindOne(ctx, bson.M{"_id": bsonRunId}).Decode(&runFileMetaData)

	// // TODO for now the id of the file is the same as the run id
	// // err = database.Collection("fs.files").FindOne(
	// // 	ctx,
	// // 	bson.D{
	// // 		{"metadata", bson.D{
	// // 			{"runId", bsonRunId},
	// // 		}},
	// // 	}).Decode(&runFileMetaData)
	// if err != nil {
	// 	fmt.Println(err.Error())
	// 	return nil, err
	// }

	// return &runFileMetaData, nil
	return nil, nil
}

func (service *SessionService) UploadFile(ctx context.Context, metadata *model.Session, file *multipart.FileHeader) error {
	// fileContent, err := file.Open()
	// if err != nil {
	// 	return err
	// }

	// byteContainer, err := ioutil.ReadAll(fileContent)
	// if err != nil {
	// 	return err
	// }

	// database, err := databases.GetDatabase(service.config.AtlasUri, service.config.MongoDbName, ctx)
	// if err != nil {
	// 	panic(err)
	// }
	// defer database.Client().Disconnect(ctx)

	// bucket, err := gridfs.NewBucket(
	// 	database,
	// )
	// if err != nil {
	// 	return err
	// }

	// opts := options.GridFSUpload()
	// opts.SetMetadata(metadata)
	// uploadStream, err := bucket.OpenUploadStreamWithID(
	// 	metadata.RunId,
	// 	file.Filename,
	// 	opts,
	// )
	// if err != nil {
	// 	return err
	// }
	// defer uploadStream.Close()

	// fileSize, err := uploadStream.Write(byteContainer)
	// if err != nil {
	// 	return err
	// }

	// fmt.Printf("Write file to DB was successful. File size: %d M\n", fileSize)
	return nil
}

func (service *SessionService) DownloadFile(ctx context.Context, sessionId uuid.UUID) ([]byte, error) {
	// database, err := databases.GetDatabase(service.config.AtlasUri, service.config.MongoDbName, ctx)
	// if err != nil {
	// 	panic(err)
	// }
	// defer database.Client().Disconnect(ctx)

	// bsonRunId, err := primitive.ObjectIDFromHex(runId)
	// if err != nil {
	// 	return nil, err
	// }

	// var result bson.M
	// err = database.Collection("fs.files").FindOne(ctx, bson.M{"_id": bsonRunId}).Decode(&result)
	// if err != nil {
	// 	return nil, err
	// }

	// bucket, _ := gridfs.NewBucket(database)

	// var buf bytes.Buffer
	// _, err = bucket.DownloadToStream(bsonRunId, &buf)
	// if err != nil {
	// 	return nil, err
	// }

	// return buf.Bytes(), nil
	return nil, nil
}
