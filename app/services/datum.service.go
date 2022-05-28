package services

import (
	"context"
	"database-ms/app/model"
	"database-ms/config"

	"gorm.io/gorm"
)

type DatumServiceInterface interface {
	Create(context.Context, *model.Datum) error
	CreateMany(context.Context, []*model.Datum) error
	FindBySessionIdAndSensorId(context.Context, string, string) ([]model.Datum, error)
}

type DatumService struct {
	db     *gorm.DB
	config *config.Configuration
}

func NewDatumService(db *gorm.DB, c *config.Configuration) DatumServiceInterface {
	return &DatumService{config: c, db: db}
}

func (service *DatumService) Create(ctx context.Context, datum *model.Datum) error {
	// result, err := service.DatumCollection(ctx).InsertOne(ctx, datum)
	// if err == nil {
	// 	datum.ID = (result.InsertedID).(primitive.ObjectID)
	// }
	// return err
	return nil
}

func (service *DatumService) CreateMany(ctx context.Context, datumArray []*model.Datum) error {
	// docs := make([]interface{}, len(datumArray))
	// for i, datum := range datumArray {
	// 	docs[i] = datum
	// }
	// _, err := service.DatumCollection(ctx).InsertMany(ctx, docs)
	// return err
	return nil
}

func (service *DatumService) FindBySessionIdAndSensorId(ctx context.Context, sessionId string, sensorId string) ([]model.Datum, error) {
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
