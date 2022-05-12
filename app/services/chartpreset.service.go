package services

import (
	"context"
	"database-ms/app/models"
	"database-ms/config"
	"database-ms/databases"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"gopkg.in/mgo.v2"
)

type ChartPresetServiceInterface interface {
	Create(context.Context, *models.ChartPreset) error
	FindByThingId(context.Context, string) ([]*models.ChartPreset, error)
	Update(context.Context, *models.ChartPreset) error
	Delete(context.Context, string) error
	FindById(context.Context, string) (*models.ChartPreset, error)
	IsPresetUnique(context.Context, *models.ChartPreset) bool
	IsPresetValid(context.Context, *models.ChartPreset) bool
}

type ChartPresetService struct {
	db			*mgo.Session
	config	*config.Configuration
}

func NewChartPresetService(db *mgo.Session, c *config.Configuration) ChartPresetServiceInterface {
	return &ChartPresetService{db: db, config: c}
}

func (service *ChartPresetService) Create(ctx context.Context, chartPreset *models.ChartPreset) error {
	client, err := databases.GetDBClient(service.config.AtlasUri, ctx)
	if err != nil {
		return err
	}

	callback := func (sessCtx mongo.SessionContext) (interface{}, error) {
		db := client.Database(service.config.MongoDbName)

		// Create the preset
		result, err := db.Collection("ChartPreset").InsertOne(ctx, chartPreset)
		if err != nil {
			return nil, err
		}

		// Remove duplicate sensor Ids from each chart
		for _, chart := range chartPreset.Charts {
			sensorIdMap := make(map[primitive.ObjectID]int)
			for _, sensorId := range chart.SensorIds {
				sensorIdMap[sensorId] = 0
			}
			chart.SensorIds = []primitive.ObjectID{}
			for id, _ := range sensorIdMap {
				chart.SensorIds = append(chart.SensorIds, id)
			}
		}

		// Create chart objects
		chartsInterface := make([]interface{}, len(chartPreset.Charts))
		for i, chart := range chartPreset.Charts {
			chart.ChartPresetID = (result.InsertedID).(primitive.ObjectID)
			chartsInterface[i] = chart
		}
		res, err := db.Collection("Chart").InsertMany(ctx, chartsInterface)
		if err != nil {
			return nil, err
		}
		chartIds := []primitive.ObjectID{}
		for _, id := range res.InsertedIDs {
			chartIds = append(chartIds, (id).(primitive.ObjectID))
		}
		chartPreset.ChartIds = chartIds

		return nil, nil
	}

	_, err = databases.WithTransaction(client, ctx, callback)
	return err
}

func (service *ChartPresetService) FindByThingId(ctx context.Context, thingId string) ([]*models.ChartPreset, error) {
	bsonThingId, err := primitive.ObjectIDFromHex(thingId)
	if err != nil {
		return nil, err
	}
	var chartPresets []*models.ChartPreset
	cursor, err := service.ChartPresetCollection(ctx).Find(ctx, bson.D{{"thingId", bsonThingId}})
	if err = cursor.All(ctx, &chartPresets); err != nil {
		return nil, err
	}
	if chartPresets == nil {
		chartPresets = []*models.ChartPreset{}
	} else {
		// Get all the charts
		chartIdMap := make(map[primitive.ObjectID]int)
		for _, preset := range chartPresets {
			chartIdMap[preset.ID] = 0
		}
		chartIds := []primitive.ObjectID{}
		for id, _ := range chartIdMap {
			chartIds = append(chartIds, id)
		}

		// Query all the required charts then attach the each preset
		cursor, err := service.ChartCollection(ctx).Find(ctx, bson.M{"_id": bson.M{"$in": chartIds}})
		if err == nil {
			var charts []*models.Chart
			if err = cursor.All(ctx, &charts); err != nil {
				return nil, err
			}
			chartMap := make(map[primitive.ObjectID]*models.Chart)
			for _, chart := range charts {
				chartMap[chart.ID] = chart
			}
			for _, chartPreset := range chartPresets {
				for _, id := range chartPreset.ChartIds {
					chartPreset.Charts = append(chartPreset.Charts, chartMap[id])
				}
			}
		} else {
			return nil, err
		}

		return nil, nil
	}
	return chartPresets, nil
}

func (service *ChartPresetService) Update(ctx context.Context, updatedChartPreset *models.ChartPreset) error {
	client, err := databases.GetDBClient(service.config.AtlasUri, ctx)
	if err != nil {
		return err
	}

	callback := func (sessCtx mongo.SessionContext) (interface{}, error) {
		db := client.Database(service.config.MongoDbName)

		// Delete the existing charts
		if _, err := db.Collection("Chart").DeleteMany(ctx, bson.M{"chartPresetId": updatedChartPreset.ID}); err != nil {
			return nil, err
		}

		// Remove duplicate sensor Ids from each chart - TODO: Create generic function for removing duplicates
		for _, chart := range updatedChartPreset.Charts {
			sensorIdMap := make(map[primitive.ObjectID]int)
			for _, sensorId := range chart.SensorIds {
				sensorIdMap[sensorId] = 0
			}
			chart.SensorIds = []primitive.ObjectID{}
			for id, _ := range sensorIdMap {
				chart.SensorIds = append(chart.SensorIds, id)
			}
		}

		// Create chart objects
		chartsInterface := make([]interface{}, len(updatedChartPreset.Charts))
		for i, chart := range updatedChartPreset.Charts {
			chart.ChartPresetID = updatedChartPreset.ID
			chartsInterface[i] = chart
		}
		res, err := db.Collection("Chart").InsertMany(ctx, chartsInterface)
		if err != nil {
			return nil, err
		}
		chartIds := []primitive.ObjectID{}
		for _, id := range res.InsertedIDs {
			chartIds = append(chartIds, (id).(primitive.ObjectID))
		}
		updatedChartPreset.ChartIds = chartIds

		// Update the chart preset
		if _, err := db.Collection("ChartPreset").ReplaceOne(ctx, bson.M{"_id": updatedChartPreset.ID}, updatedChartPreset); err != nil {
			return nil, err
		}

		// Fetch the new charts
		if cursor, err := db.Collection("Chart").Find(ctx, bson.M{"chartPresetId": updatedChartPreset.ID}); err != nil {
			return nil, err
		} else {
			charts := []*models.Chart{}
			if err = cursor.All(ctx, &charts); err != nil {
				return nil, err
			} else {
				updatedChartPreset.Charts = charts
			}
		}

		return nil, nil
	}

	_, err = databases.WithTransaction(client, ctx, callback)
	return err
}

func (service *ChartPresetService) Delete(ctx context.Context, chartPresetId string) error {
	bsonChartPresetId, err := primitive.ObjectIDFromHex(chartPresetId)
	if err != nil {
		return err
	}

	client, err := databases.GetDBClient(service.config.AtlasUri, ctx)
	if err != nil {
		return err
	}

	callback := func (sessCtx mongo.SessionContext) (interface{}, error) {
		db := client.Database(service.config.MongoDbName)

		// Delete the associated charts
		if _, err := db.Collection("Chart").DeleteMany(ctx, bson.M{"chartPresetId": bsonChartPresetId}); err != nil {
			return nil, err
		}

		// Delete the preset
		if _, err := db.Collection("ChartPreset").DeleteOne(ctx, bson.M{"chartPresetId": bsonChartPresetId}); err != nil {
			return nil, err
		}

		return nil, nil
	}

	_, err = databases.WithTransaction(client, ctx, callback)
	return err
}

func (service *ChartPresetService) FindById(ctx context.Context, chartPresetId string) (*models.ChartPreset, error) { // Should this return a copy or pointer?
	bsonChartPresetId, err := primitive.ObjectIDFromHex(chartPresetId)
	if err != nil {
		return nil, err
	}
	var chartPreset models.ChartPreset
	if err = service.ChartPresetCollection(ctx).FindOne(ctx, bson.M{"_id": bsonChartPresetId}).Decode(&chartPreset); err != nil {
		return nil, err
	}
	return &chartPreset, nil
}

func (service *ChartPresetService) IsPresetValid(ctx context.Context, chartPreset *models.ChartPreset) bool {
	sensorMap := make(map[primitive.ObjectID]int)
	for _, chart := range chartPreset.Charts {
		for _, sensorId := range chart.SensorIds {
			sensorMap[sensorId] = 0
		}
	}
	sensorIds := []primitive.ObjectID{}
	for sensorId, _ := range sensorMap {
		sensorIds = append(sensorIds, sensorId)
	}
	_, err := service.SensorCollection(ctx).Find(ctx, bson.M{"_id": bson.M{"$in": sensorIds}})
	return err == nil
}

func (service *ChartPresetService) IsPresetUnique(ctx context.Context, chartPreset *models.ChartPreset) bool {
	var existing models.ChartPreset
	if err := service.ChartPresetCollection(ctx).FindOne(ctx, bson.M{"name": chartPreset.Name}).Decode(&existing); err != nil {
		return true
	} else {
		return existing.ID != chartPreset.ID
	}
}

func (service *ChartPresetService) ChartPresetCollection(ctx context.Context) *mongo.Collection {
	dbClient, err := databases.GetDBClient(service.config.AtlasUri, ctx)
	if err != nil {
		panic(err)
	}
	return dbClient.Database(service.config.MongoDbName).Collection("ChartPreset")
}

func (service *ChartPresetService) ChartCollection(ctx context.Context) * mongo.Collection {
	dbClient, err := databases.GetDBClient(service.config.AtlasUri, ctx)
	if err != nil {
		panic(err)
	}
	return dbClient.Database(service.config.MongoDbName).Collection("Chart")
}

func (service *ChartPresetService) SensorCollection(ctx context.Context) *mongo.Collection {
	dbClient, err := databases.GetDBClient(service.config.AtlasUri, ctx)
	if err != nil {
		panic(err)
	}
	return dbClient.Database(service.config.MongoDbName).Collection("Sensor")
}
