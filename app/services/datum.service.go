package services

import (
	"context"
	"database-ms/app/model"
	"database-ms/config"

	"github.com/google/uuid"
	"github.com/jackc/pgconn"
	"gorm.io/gorm"
)

type DatumServiceInterface interface {
	// Public
	FindBySessionIdAndSensorId(context.Context, uuid.UUID, uuid.UUID) ([]model.Datum, *pgconn.PgError)
	CreateMany(context.Context, []*model.Datum) *pgconn.PgError
}

type DatumService struct {
	db     *gorm.DB
	config *config.Configuration
}

func NewDatumService(db *gorm.DB, c *config.Configuration) DatumServiceInterface {
	return &DatumService{config: c, db: db}
}

func (service *DatumService) FindBySessionIdAndSensorId(ctx context.Context, sessionId uuid.UUID, sensorId uuid.UUID) ([]model.Datum, *pgconn.PgError) {
	// database, err := databases.GetDatabase(service.config.AtlasUri, service.config.MongoDbName, ctx)
	// if err != nil {
	// 	panic(err)
	// }
	// defer database.Client().Disconnect(ctx)

	// var datumArray []*models.Datum
	// bsonSessionId, err := primitive.ObjectIDFromHex(sessionId)
	// if err != nil {
	// 	return nil, err
	// }
	// bsonSensorId, err := primitive.ObjectIDFromHex(sensorId)
	// if err != nil {
	// 	return nil, err
	// }

	// cursor, err := database.Collection("Datum").Find(ctx, bson.M{"sessionId": bsonSessionId, "sensorId": bsonSensorId})
	// if err != nil {
	// 	return nil, err
	// }

	// if err = cursor.All(ctx, &datumArray); err != nil {
	// 	return nil, err
	// }

	// var formattedDatumArray []models.FormattedDatum
	// for _, datum := range datumArray {
	// 	formattedDatumArray = append(formattedDatumArray, models.FormattedDatum{
	// 		X: datum.Value,
	// 		Y: datum.Timestamp,
	// 	})
	// }

	// return formattedDatumArray, nil
	return nil, nil
}

func (service *DatumService) CreateMany(ctx context.Context, datumArray []*model.Datum) *pgconn.PgError {
	// docs := make([]interface{}, len(datumArray))
	// for i, datum := range datumArray {
	// 	docs[i] = datum
	// }
	// _, err := service.DatumCollection(ctx).InsertMany(ctx, docs)
	// return err
	return nil
}
