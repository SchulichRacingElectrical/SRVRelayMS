package redisHandler

import (
	"context"
	"database-ms/config"
	"log"
	"sort"
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

			// Fill missing timestamps
			timeLinearThingData := fillMissingTimestamps(thingDataArray)

			// Convert to 2D array
			timeLinearThingData2DArray := mapArrayTo2DArray(timeLinearThingData)
			log.Println(timeLinearThingData2DArray)

			// // Get sensor list
			// sensorService := services.NewSensorService(dbSession, conf)
			// sensors, err := sensorService.FindByThingId(ctx, thingId)
			// if err != nil {
			// 	panic(err)
			// }

			// // Create map of SmallId to ID and Name
			// smallIdToInfoMap := make(map[string]struct {
			// 	ID   primitive.ObjectID
			// 	Name string
			// })
			// for _, sensor := range sensors {
			// 	smallId := strconv.Itoa(*sensor.SmallId)
			// 	smallIdToInfoMap[smallId] = struct {
			// 		ID   primitive.ObjectID
			// 		Name string
			// 	}{ID: sensor.ID, Name: sensor.Name}
			// }

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
	currentDataMap := copyMap(thingDataArray[0])

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

func fillMissingTimestamps(thingDataArray []map[string]int) []map[string]int {
	lastTimestamp := thingDataArray[len(thingDataArray)-1]["ts"]
	output := make([]map[string]int, lastTimestamp+1)

	// Copy first map with 0 values
	currentDataMap := copyMapWithDefaultValues(thingDataArray[0])
	currentTimestamp := 0

	for _, thingDataItem := range thingDataArray {
		// if thingDataItem has a higher timestamp than currentTimestamp,
		// add currentDataMap to output and increment currentTimestamp
		// until thingDataItem has a timestamp equal to currentTimestamp,
		// then add thingDataItem to output
		if thingDataItem["ts"] > currentTimestamp {
			for i := currentTimestamp; i < thingDataItem["ts"]; i++ {
				currentDataMap["ts"] = i
				output[i] = copyMap(currentDataMap)
			}
			currentTimestamp = thingDataItem["ts"]
		}
		output[thingDataItem["ts"]] = copyMap(thingDataItem)

		// Update currentDataMap
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

	return output
}

func copyMap(source map[string]int) map[string]int {
	dest := make(map[string]int)
	for key, value := range source {
		dest[key] = value
	}
	return dest
}

func copyMapWithDefaultValues(source map[string]int) map[string]int {
	dest := make(map[string]int)
	for key := range source {
		dest[key] = 0
	}
	return dest
}

func mapArrayTo2DArray(mapArray []map[string]int) [][]int {
	// Create sorted array of smallIds to keep them ordered
	var smallIds []int
	for k := range mapArray[0] {
		if k == "ts" {
			continue
		}
		smallId, _ := strconv.Atoi(k)
		smallIds = append(smallIds, smallId)
	}
	sort.Ints(smallIds)

	// Create 2D array
	output := make([][]int, len(mapArray))
	for i := range mapArray {
		output[i] = make([]int, len(mapArray[i]))
		output[i][0] = mapArray[i]["ts"]
		for j, smallId := range smallIds {
			output[i][j+1] = mapArray[i][strconv.Itoa(smallId)]
		}
	}
	return output
}
