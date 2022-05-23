package redisHandler

import (
	"context"
	"database-ms/app/models"
	"database-ms/app/services"
	"database-ms/config"
	"log"
	"strconv"

	"encoding/json"

	"github.com/go-redis/redis/v8"
	"gopkg.in/mgo.v2"
)

type Message struct {
	Active  bool   `json:"active"`
	ThingId string `json:"THING"`
}

func Initialize(conf *config.Configuration, dbSession *mgo.Session) {
	redisClient := redis.NewClient(&redis.Options{
		Addr:     conf.RedisUrl + ":" + conf.RedisPort,
		Password: conf.RedisPassword,
	})

	go awaitThingDataSessions(redisClient, dbSession, conf)
}

func awaitThingDataSessions(redisClient *redis.Client, dbSession *mgo.Session, conf *config.Configuration) {
	ctx := context.Background()
	subscriber := redisClient.Subscribe(ctx, "THING_CONNECTION")

	for {
		msg, err := subscriber.ReceiveMessage(ctx)
		if err != nil {
			panic(err)
		}
		message := Message{}
		json.Unmarshal([]byte(msg.Payload), &message)
		if message.Active {
			go thingDataSession(message.ThingId, redisClient, dbSession, conf)
		}
	}
}

func thingDataSession(thingId string, redisClient *redis.Client, dbSession *mgo.Session, conf *config.Configuration) {
	ctx := context.Background()
	subscriber := redisClient.Subscribe(ctx, "THING_"+thingId)
	log.Println("Thing Data Session Started for " + thingId)

	for {
		msg, err := subscriber.ReceiveMessage(ctx)
		if err != nil {
			panic(err)
		}
		message := Message{}
		json.Unmarshal([]byte(msg.Payload), &message)

		if !message.Active {
			// Get thing data from redis
			thingData, err := redisClient.LRange(ctx, "THING_"+thingId, 0, -1).Result()
			if err != nil {
				panic(err)
			}

			// Delete thing data from redis
			err = redisClient.Del(ctx, "THING_"+thingId).Err()
			if err != nil {
				panic(err)
			}

			// Parse thing data from JSON
			var thingDataArray []map[string]int
			for _, thingDataItem := range thingData {
				var thingDataItemMap map[string]int
				json.Unmarshal([]byte(thingDataItem), &thingDataItemMap)
				thingDataArray = append(thingDataArray, thingDataItemMap)
			}
			log.Println(thingDataArray)

			// Process thing data to fill missing sensor values
			thingDataArray = fillMissingValues(thingDataArray)

			// Get sensor list
			sensorService := services.NewSensorService(dbSession, conf)
			sensors, err := sensorService.FindByThingId(ctx, thingId)
			if err != nil {
				panic(err)
			}

			// Replace SmallId with ID
			thingDataArray = replaceSmallIdWithId(thingDataArray, sensors)

			// Save thing data to csv

			// Save thing data to mongo
		}
	}
}

func fillMissingValues(thingDataArray []map[string]int) []map[string]int {
	// use first map from thingDataArray to initialize currentDataMap,
	// then iterate through thingDataArray and fill missing values in each map,
	// updating currentDataMap as we go

	// creates a copy of the first map
	currentDataMap := make(map[string]int)
	for key, value := range thingDataArray[0] {
		currentDataMap[key] = value
	}

	// iterate through the rest of the maps
	for _, thingDataItem := range thingDataArray {
		for key := range currentDataMap {
			// if key doesn't exist on the current map, add it from currentDataMap
			// if the key does exist, add the value from the current map to currentDataMap
			if _, ok := thingDataItem[key]; !ok {
				thingDataItem[key] = currentDataMap[key]
			} else {
				currentDataMap[key] = thingDataItem[key]
			}
		}
	}

	return thingDataArray
}

func replaceSmallIdWithId(thingDataArray []map[string]int, sensors []*models.Sensor) []map[string]int {
	// Create map of SmallId to ID
	smallIdToIdMap := make(map[string]string)
	for _, sensor := range sensors {
		smallId := strconv.Itoa(*sensor.SmallId)
		smallIdToIdMap[smallId] = sensor.ID.Hex()
	}

	// Replace SmallId with ID
	for _, thingDataItem := range thingDataArray {
		for smallId, id := range smallIdToIdMap {
			thingDataItem[id] = thingDataItem[smallId]
			delete(thingDataItem, smallId)
		}
	}

	return thingDataArray
}
