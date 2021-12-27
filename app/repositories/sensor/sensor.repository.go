package sensor

import (
	"context"

	model "database-ms/app/models"
	"database-ms/config"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// SensorRepository, used to perform DB operations
// Interface contains basic operations on sensor document
// So that, DB operation can be performed easily
type SensorRepository interface {

	// Create, will perform db operation to save sensor
	// Returns new id and error if occurs
	Create(context.Context, *model.Sensor) error

	// FindByThingId, find all sensor belonging to a thing
	// It will retrun error if occurs
	FindByThingId(context.Context, string) ([]*model.Sensor, error)

	// FindOneById, find the sensor by the provided id
	// return matched sensor and error if any
	FindOneById(context.Context, string) (*model.Sensor, error)

	// FindOneByThingIdAndSid
	FindOneByThingIdAndSid(context.Context, string, int) (*model.Sensor, error)

	// FindByThingIdAndLastUpdate
	FindByThingIdAndLastUpdate(context.Context, string, int64) ([]*model.Sensor, error)

	// FindOne
	FindOne(context.Context, interface{}) (*model.Sensor, error)

	// Update
	Update(context.Context, string, interface{}) error

	// Delete
	Delete(context.Context, string) error

	IsSensorAlreadyExisits(context.Context, bson.ObjectId, int) bool
}

type SensorRepositoryImp struct {
	db     *mgo.Session
	config *config.Configuration
}

func New(db *mgo.Session, c *config.Configuration) SensorRepository {

	return &SensorRepositoryImp{db: db, config: c}

}

func (service *SensorRepositoryImp) Create(ctx context.Context, sensor *model.Sensor) error {

	return service.collection().Insert(sensor)

}

func (service *SensorRepositoryImp) FindOne(ctx context.Context, query interface{}) (*model.Sensor, error) {

	var sensor model.Sensor
	err := service.collection().Find(query).One(&sensor)
	return &sensor, err

}

func (service *SensorRepositoryImp) FindByThingId(ctx context.Context, thingId string) ([]*model.Sensor, error) {

	var sensors []*model.Sensor
	query := bson.M{"thingId": bson.ObjectIdHex(thingId)}
	err := service.collection().Find(query).All(&sensors)

	return sensors, err

}

func (service *SensorRepositoryImp) FindOneById(ctx context.Context, id string) (*model.Sensor, error) {

	query := bson.M{"_id": bson.ObjectIdHex(id)}
	user, err := service.FindOne(ctx, query)

	return user, err

}

func (service *SensorRepositoryImp) FindOneByThingIdAndSid(ctx context.Context, thingId string, sid int) (*model.Sensor, error) {

	query := bson.M{"thingId": bson.ObjectIdHex(thingId), "sid": sid}
	user, err := service.FindOne(ctx, query)

	return user, err

}

func (service *SensorRepositoryImp) FindByThingIdAndLastUpdate(ctx context.Context, thingId string, lastUpdate int64) ([]*model.Sensor, error) {

	var sensors []*model.Sensor
	query := bson.M{
		"thingId": bson.ObjectIdHex(thingId),
		"lastUpdate": bson.M{
			"$gte": lastUpdate,
		},
	}
	err := service.collection().Find(query).All(&sensors)

	return sensors, err

}

func (service *SensorRepositoryImp) Update(ctx context.Context, id string, update interface{}) error {
	return nil
}

func (service *SensorRepositoryImp) Delete(ctx context.Context, id string) error {
	// TODO Delete from Thing document
	return nil
}

func (service *SensorRepositoryImp) IsSensorAlreadyExisits(ctx context.Context, thingId bson.ObjectId, sid int) bool {

	query := bson.M{"thingId": thingId, "sid": sid}
	_, err := service.FindOne(ctx, query)
	if err != nil {
		return false
	}
	return true

}

func (service *SensorRepositoryImp) collection() *mgo.Collection {

	return service.db.DB(service.config.MongoDbName).C("Sensor")

}
