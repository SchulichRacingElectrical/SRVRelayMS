package services

import (
	"context"
	"database-ms/app/model"
	"database-ms/config"

	"github.com/google/uuid"
	"github.com/jackc/pgconn"
	"gorm.io/gorm"
)

type ChartPresetServiceInterface interface {
	// Public
	FindByThingId(context.Context, uuid.UUID) ([]*model.ChartPreset, *pgconn.PgError)
	Create(context.Context, *model.ChartPreset) *pgconn.PgError
	Update(context.Context, *model.ChartPreset) *pgconn.PgError
	Delete(context.Context, uuid.UUID) *pgconn.PgError

	// Private
	FindById(context.Context, uuid.UUID) (*model.ChartPreset, *pgconn.PgError)
}

type ChartPresetService struct {
	db     *gorm.DB
	config *config.Configuration
}

func NewChartPresetService(db *gorm.DB, c *config.Configuration) ChartPresetServiceInterface {
	return &ChartPresetService{db: db, config: c}
}

// PUBLIC FUNCTIONS

func (service *ChartPresetService) FindByThingId(ctx context.Context, thingId uuid.UUID) ([]*model.ChartPreset, *pgconn.PgError) {
	// bsonThingId, err := primitive.ObjectIDFromHex(thingId)
	// if err != nil {
	// 	return nil, err
	// }
	// var chartPresets []*models.ChartPreset
	// cursor, err := service.ChartPresetCollection(ctx).Find(ctx, bson.D{{"thingId", bsonThingId}})
	// if err = cursor.All(ctx, &chartPresets); err != nil {
	// 	return nil, err
	// }
	// if chartPresets == nil {
	// 	chartPresets = []*models.ChartPreset{}
	// } else {
	// 	// Create a list for all the charts to query
	// 	chartPresetIds := []primitive.ObjectID{}
	// 	for _, preset := range chartPresets {
	// 		chartPresetIds = append(chartPresetIds, preset.ID)
	// 		preset.Charts = []models.Chart{}
	// 	}

	// 	// Fetch all the charts
	// 	cursor, err := service.ChartCollection(ctx).Find(ctx, bson.M{"chartPresetId": bson.M{"$in": chartPresetIds}})
	// 	if err == nil {
	// 		var charts []models.Chart
	// 		if err = cursor.All(ctx, &charts); err != nil {
	// 			return nil, err
	// 		}
	// 		chartMap := make(map[primitive.ObjectID][]models.Chart)
	// 		for _, chart := range charts {
	// 			chartMap[chart.ChartPresetID] = append(chartMap[chart.ChartPresetID], chart)
	// 		}
	// 		for _, preset := range chartPresets {
	// 			preset.Charts = chartMap[preset.ID]
	// 		}
	// 	} else {
	// 		return nil, err
	// 	}
	// }
	// return chartPresets, nil
	return nil, nil
}

func (service *ChartPresetService) Create(ctx context.Context, chartPreset *model.ChartPreset) *pgconn.PgError {
	// client, err := databases.GetDBClient(service.config.AtlasUri, ctx)
	// if err != nil {
	// 	return err
	// }

	// callback := func (sessCtx mongo.SessionContext) (interface{}, error) {
	// 	db := client.Database(service.config.MongoDbName)

	// 	// Create the preset
	// 	result, err := db.Collection("ChartPreset").InsertOne(ctx, chartPreset)
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// 	chartPreset.ID = (result.InsertedID).(primitive.ObjectID)

	// 	// Remove duplicate sensor Ids from each chart
	// 	for _, chart := range chartPreset.Charts {
	// 		sensorIdMap := make(map[primitive.ObjectID]int)
	// 		for _, sensorId := range chart.SensorIds {
	// 			sensorIdMap[sensorId] = 0
	// 		}
	// 		chart.SensorIds = []primitive.ObjectID{}
	// 		for id, _ := range sensorIdMap {
	// 			chart.SensorIds = append(chart.SensorIds, id)
	// 		}
	// 	}

	// 	// Create chart objects
	// 	chartsInterface := make([]interface{}, len(chartPreset.Charts))
	// 	for i, chart := range chartPreset.Charts {
	// 		chart.ChartPresetID = (result.InsertedID).(primitive.ObjectID)
	// 		chartsInterface[i] = chart
	// 	}
	// 	res, err := db.Collection("Chart").InsertMany(ctx, chartsInterface)
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// 	chartIds := []primitive.ObjectID{}
	// 	for _, id := range res.InsertedIDs {
	// 		chartIds = append(chartIds, (id).(primitive.ObjectID))
	// 	}

	// 	return nil, nil
	// }

	// _, err = databases.WithTransaction(client, ctx, callback)
	// return err
	return nil
}

func (service *ChartPresetService) Update(ctx context.Context, updatedChartPreset *model.ChartPreset) *pgconn.PgError {
	// client, err := databases.GetDBClient(service.config.AtlasUri, ctx)
	// if err != nil {
	// 	return err
	// }

	// // TODO: Don't delete charts if we don't need to?
	// callback := func(sessCtx mongo.SessionContext) (interface{}, error) {
	// 	db := client.Database(service.config.MongoDbName)

	// 	// Delete the existing charts
	// 	if _, err := db.Collection("Chart").DeleteMany(ctx, bson.M{"chartPresetId": updatedChartPreset.ID}); err != nil {
	// 		return nil, err
	// 	}

	// 	// Remove duplicate sensor Ids from each chart - TODO: Create generic function for removing duplicates
	// 	for _, chart := range updatedChartPreset.Charts {
	// 		sensorIdMap := make(map[primitive.ObjectID]int)
	// 		for _, sensorId := range chart.SensorIds {
	// 			sensorIdMap[sensorId] = 0
	// 		}
	// 		chart.SensorIds = []primitive.ObjectID{}
	// 		for id, _ := range sensorIdMap {
	// 			chart.SensorIds = append(chart.SensorIds, id)
	// 		}
	// 	}

	// 	// Create chart objects
	// 	chartsInterface := make([]interface{}, len(updatedChartPreset.Charts))
	// 	for i, chart := range updatedChartPreset.Charts {
	// 		chart.ChartPresetID = updatedChartPreset.ID
	// 		chartsInterface[i] = chart
	// 	}
	// 	if len(chartsInterface) == 0 {
	// 		return nil, bson.ErrDecodeToNil
	// 	}
	// 	res, err := db.Collection("Chart").InsertMany(ctx, chartsInterface)
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// 	chartIds := []primitive.ObjectID{}
	// 	for _, id := range res.InsertedIDs {
	// 		chartIds = append(chartIds, (id).(primitive.ObjectID))
	// 	}

	// 	// Update the chart preset
	// 	if _, err := db.Collection("ChartPreset").ReplaceOne(ctx, bson.M{"_id": updatedChartPreset.ID}, updatedChartPreset); err != nil {
	// 		return nil, err
	// 	}

	// 	// Fetch the new charts
	// 	if cursor, err := db.Collection("Chart").Find(ctx, bson.M{"chartPresetId": updatedChartPreset.ID}); err != nil {
	// 		return nil, err
	// 	} else {
	// 		charts := []models.Chart{}
	// 		if err = cursor.All(ctx, &charts); err != nil {
	// 			return nil, err
	// 		} else {
	// 			updatedChartPreset.Charts = charts
	// 		}
	// 	}

	// 	return nil, nil
	// }

	// _, err = databases.WithTransaction(client, ctx, callback)
	// return err
	return nil
}

func (service *ChartPresetService) Delete(ctx context.Context, chartPresetId uuid.UUID) *pgconn.PgError {
	// bsonChartPresetId, err := primitive.ObjectIDFromHex(chartPresetId)
	// if err != nil {
	// 	return err
	// }

	// client, err := databases.GetDBClient(service.config.AtlasUri, ctx)
	// if err != nil {
	// 	return err
	// }

	// callback := func(sessCtx mongo.SessionContext) (interface{}, error) {
	// 	db := client.Database(service.config.MongoDbName)

	// 	// Delete the associated charts
	// 	if _, err := db.Collection("Chart").DeleteMany(ctx, bson.M{"chartPresetId": bsonChartPresetId}); err != nil {
	// 		return nil, err
	// 	}

	// 	// Delete the preset
	// 	if _, err := db.Collection("ChartPreset").DeleteOne(ctx, bson.M{"chartPresetId": bsonChartPresetId}); err != nil {
	// 		return nil, err
	// 	}

	// 	return nil, nil
	// }

	// _, err = databases.WithTransaction(client, ctx, callback)
	// return err
	return nil
}

// PRIVATE FUNCTIONS

func (service *ChartPresetService) FindById(ctx context.Context, chartPresetId uuid.UUID) (*model.ChartPreset, *pgconn.PgError) { // Should this return a copy or pointer?
	// bsonChartPresetId, err := primitive.ObjectIDFromHex(chartPresetId)
	// if err != nil {
	// 	return nil, err
	// }
	// var chartPreset models.ChartPreset
	// if err = service.ChartPresetCollection(ctx).FindOne(ctx, bson.M{"_id": bsonChartPresetId}).Decode(&chartPreset); err != nil {
	// 	return nil, err
	// }
	// return &chartPreset, nil
	return nil, nil
}
