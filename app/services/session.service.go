package services

import (
	"context"
	"database-ms/app/model"
	"database-ms/config"
	"database-ms/utils"
	"mime/multipart"

	"github.com/google/uuid"
	"github.com/jackc/pgconn"
	"gorm.io/gorm"
)

type SessionServiceInterface interface {
	CreateSession(context.Context, *model.Session) error
	FindById(context.Context, uuid.UUID) (*model.Session, *pgconn.PgError)
	GetSessionsByThingId(context.Context, uuid.UUID) ([]*model.Session, *pgconn.PgError)
	UpdateSession(context.Context, *model.Session) *pgconn.PgError
	DeleteSession(context.Context, uuid.UUID) *pgconn.PgError
	GetSessionFileMetaData(context.Context, uuid.UUID) (*model.Session, *pgconn.PgError)
	UploadFile(context.Context, *model.Session, *multipart.FileHeader) error
	DownloadFile(context.Context, uuid.UUID) ([]byte, error)
}

type SessionService struct {
	db     *gorm.DB
	config *config.Configuration
}

func NewSessionService(db *gorm.DB, c *config.Configuration) SessionServiceInterface {
	return &SessionService{config: c, db: db}
}

func (service *SessionService) CreateSession(ctx context.Context, session *model.Session) error {
	result := service.db.Create(&session)
	if result.Error != nil {
		return result.Error
		// var perr *pgconn.PgError
		// errors.As(result.Error, &perr)
		// return utils.GetPostgresError(result.Error)
	}
	return nil
}

func (service *SessionService) FindById(ctx context.Context, sessionId uuid.UUID) (*model.Session, *pgconn.PgError) {
	var session *model.Session
	result := service.db.Where("id = ?", sessionId).First(&session)
	if result.Error != nil {
		return nil, &pgconn.PgError{}
	}
	return session, nil
}

func (service *SessionService) GetSessionsByThingId(ctx context.Context, thingId uuid.UUID) ([]*model.Session, *pgconn.PgError) {
	var sessions []*model.Session
	result := service.db.Where("thing_id = ?", thingId).Find(&sessions)
	if result.Error != nil {
		return nil, utils.GetPostgresError(result.Error)
	}
	return sessions, nil
}

func (service *SessionService) UpdateSession(ctx context.Context, updatedSession *model.Session) *pgconn.PgError {
	var session model.Session
	if err := service.db.Where("id = ?", updatedSession.Id).First(&session).Error; err != nil {
		return &pgconn.PgError{}
	}

	result := service.db.Model(&session).Updates(&updatedSession)
	if result.Error != nil {
		return utils.GetPostgresError(result.Error)
	}
	return nil
}

func (service *SessionService) DeleteSession(ctx context.Context, sessionId uuid.UUID) *pgconn.PgError {
	session := model.Session{Base: model.Base{Id: sessionId}}
	result := service.db.Delete(&session)
	return utils.GetPostgresError(result.Error)
}

func (service *SessionService) GetSessionFileMetaData(ctx context.Context, sessionId uuid.UUID) (*model.Session, *pgconn.PgError) {
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
