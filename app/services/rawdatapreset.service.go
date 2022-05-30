package services

import (
	"context"
	"database-ms/app/model"
	"database-ms/config"
	"database-ms/utils"

	"github.com/google/uuid"
	"github.com/jackc/pgconn"
	"gorm.io/gorm"
)

type RawDataPresetServiceInterface interface {
	// Public
	FindByThingId(context.Context, uuid.UUID) ([]*model.RawDataPreset, *pgconn.PgError)
	Create(context.Context, *model.RawDataPreset) *pgconn.PgError
	Update(context.Context, *model.RawDataPreset) *pgconn.PgError
	Delete(context.Context, uuid.UUID) *pgconn.PgError

	// Private
	FindById(context.Context, uuid.UUID) (*model.RawDataPreset, *pgconn.PgError)
}

type RawDataPresetService struct {
	db     *gorm.DB
	config *config.Configuration
}

func NewRawDataPresetService(db *gorm.DB, c *config.Configuration) RawDataPresetServiceInterface {
	return &RawDataPresetService{db: db, config: c}
}

func (service *RawDataPresetService) FindByThingId(ctx context.Context, thingId uuid.UUID) ([]*model.RawDataPreset, *pgconn.PgError) {
	// bsonThingId, err := primitive.ObjectIDFromHex(thingId)
	// if err != nil {
	// 	return nil, err
	// }
	// var rawDataPresets []*models.RawDataPreset
	// cursor, err := service.RawDataPresetCollection(ctx).Find(ctx, bson.D{{"thingId", bsonThingId}})
	// if err = cursor.All(ctx, &rawDataPresets); err != nil {
	// 	return nil, err
	// }
	// if rawDataPresets == nil {
	// 	rawDataPresets = []*models.RawDataPreset{}
	// }
	// return rawDataPresets, nil
	return nil, nil
}

func (service *RawDataPresetService) Create(ctx context.Context, rawDataPreset *model.RawDataPreset) *pgconn.PgError {
	err := service.db.Transaction(func(db *gorm.DB) error {
		// Create the preset
		result := db.Create(&rawDataPreset)
		if result.Error != nil {
			return result.Error
		}

		// Generate the list of preset sensors
		var rawDataPresetSensors []model.RawDataPresetSensor
		for _, sensorId := range rawDataPreset.SensorIds {
			rawDataPresetSensors = append(rawDataPresetSensors, model.RawDataPresetSensor{
				RawDataPresetId: rawDataPreset.Id,
				SensorId:        sensorId,
			})
		}

		// Insert empty sensorIds
		if len(rawDataPreset.SensorIds) == 0 {
			rawDataPreset.SensorIds = []uuid.UUID{}
		}

		// Batch insert preset sensor
		result = db.Table(model.TableNameRawdataPresetSensor).CreateInBatches(rawDataPresetSensors, 100)
		return result.Error
	})
	return utils.GetPostgresError(err)
}

func (service *RawDataPresetService) Update(ctx context.Context, updatedRawDataPreset *model.RawDataPreset) *pgconn.PgError {
	// Remove duplicate sensor Ids from the preset
	// sensorIdMap := make(map[primitive.ObjectID]int)
	// for _, sensorId := range updatedRawDataPreset.SensorIds {
	// 	sensorIdMap[sensorId] = 0
	// }
	// updatedRawDataPreset.SensorIds = []primitive.ObjectID{}
	// for id, _ := range sensorIdMap {
	// 	updatedRawDataPreset.SensorIds = append(updatedRawDataPreset.SensorIds, id)
	// }
	// _, err := service.RawDataPresetCollection(ctx).UpdateOne(ctx, bson.M{"_id": updatedRawDataPreset.ID}, bson.M{"$set": updatedRawDataPreset})
	// return err
	return nil
}

func (service *RawDataPresetService) Delete(ctx context.Context, rawDataPresetId uuid.UUID) *pgconn.PgError {
	// bsonRawDataPresetId, err := primitive.ObjectIDFromHex(rawDataPresetId)
	// if err == nil {
	// 	_, err := service.RawDataPresetCollection(ctx).DeleteOne(ctx, bson.M{"_id": bsonRawDataPresetId})
	// 	return err
	// } else {
	// 	return err
	// }
	return nil
}

// Internal function
func (service *RawDataPresetService) FindById(ctx context.Context, rawDataPresetId uuid.UUID) (*model.RawDataPreset, *pgconn.PgError) {
	// bsonRawDataPresetId, err := primitive.ObjectIDFromHex(rawDataPresetId)
	// if err != nil {
	// 	return nil, err
	// }
	// var rawDataPreset models.RawDataPreset
	// if err = service.RawDataPresetCollection(ctx).FindOne(ctx, bson.M{"_id": bsonRawDataPresetId}).Decode(&rawDataPreset); err != nil {
	// 	return nil, err
	// }
	// return &rawDataPreset, nil
	return nil, nil
}
