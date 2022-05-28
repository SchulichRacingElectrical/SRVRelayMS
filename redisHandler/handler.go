package redisHandler

import (
	"context"
	"database-ms/app/models"
	"database-ms/app/services"
	"database-ms/config"
	"log"
	"os"
	"sort"
	"strconv"
	"time"

	"encoding/csv"
	"encoding/json"

	"github.com/go-redis/redis/v8"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"gopkg.in/mgo.v2"
)

type Message struct {
	Active  bool   `json:"active"`
	ThingId string `json:"THING"`
}

type SensorInfo struct {
	ID   primitive.ObjectID
	Name string
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

	sessionService := services.NewSessionService(conf)

	for {
		msg, err := subscriber.ReceiveMessage(ctx)
		if err != nil {
			panic(err)
		}
		message := Message{}
		json.Unmarshal([]byte(msg.Payload), &message)
		if message.Active {
			thingObjId, err := primitive.ObjectIDFromHex(message.ThingId)
			if err != nil {
				panic(err)
			}
			session := &models.Session{
				StartDate: time.Now().UnixMilli(),
				EndDate:   0, // EndDate will be updated after session is closed
				ThingID:   thingObjId,
				Name:      "Test", // TODO: Allow naming sessions
			}
			sessionId, err := sessionService.CreateSession(ctx, session)
			if err != nil {
				panic(err)
			}
			log.Println("Session created with ID: " + sessionId.Hex())
			session.ID = sessionId
			go thingDataSession(message.ThingId, session, redisClient, dbSession, conf)
		}
	}
}

func thingDataSession(thingId string, session *models.Session, redisClient *redis.Client, dbSession *mgo.Session, conf *config.Configuration) {
	ctx := context.Background()
	subscriber := redisClient.Subscribe(ctx, "THING_"+thingId)
	log.Println("Thing Data Session Started for " + thingId)

	datumService := services.NewDatumService(dbSession, conf)
	sessionService := services.NewSessionService(conf)

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

			// Process thing data to fill missing sensor values
			thingDataArray = fillMissingValues(thingDataArray)

			// Fill missing timestamps
			timeLinearThingData := fillMissingTimestamps(thingDataArray)

			// Get sorted small ids
			smallIds := getSortedSmallIds(thingDataArray[0])

			// Convert to 2D arrays
			timeLinearThingData2DArray := mapArrayTo2DArray(timeLinearThingData, smallIds)

			// Get sensor list
			sensorService := services.NewSensorService(dbSession, conf)
			sensors, err := sensorService.FindByThingId(ctx, thingId)
			if err != nil {
				panic(err)
			}

			// Create map of SmallId to ID and Name
			smallIdToInfoMap := make(map[string]SensorInfo)
			for _, sensor := range sensors {
				smallId := strconv.Itoa(*sensor.SmallId)
				smallIdToInfoMap[smallId] = SensorInfo{ID: sensor.ID, Name: sensor.Name}
			}

			// Process non-linear thing data for mongo
			thingDataArray = replaceSmallIdsWithIds(thingDataArray, smallIdToInfoMap)
			datumArray := make([]*models.Datum, len(thingDataArray)*len(smallIds))
			for i, thingDataItem := range thingDataArray {
				for j, smallId := range smallIds {
					strSmallId := strconv.Itoa(smallId)
					sensorId := smallIdToInfoMap[strSmallId].ID
					datumArray[i*len(smallIds)+j] = &models.Datum{
						SessionID: session.ID,
						SensorID:  sensorId,
						Value:     float64(thingDataItem[sensorId.Hex()]),
						Timestamp: int64(thingDataItem["ts"]),
					}
				}
			}

			// Save thing data to mongo
			datumService.CreateMany(ctx, datumArray)

			// Update session
			session.EndDate = session.StartDate + int64(thingDataArray[len(thingDataArray)-1]["ts"])
			session.FileName = "srv_files/" + thingId + "/" + session.ID.Hex() + ".csv"
			err = sessionService.UpdateSession(ctx, session)
			if err != nil {
				panic(err)
			}

			// Save thing data to csv
			exportToCsv(timeLinearThingData2DArray, smallIds, smallIdToInfoMap, thingId, session.FileName)

			log.Println("Thing Data Session Ended for " + thingId)
			return
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
	log.Println("lastTimestamp: ", lastTimestamp, " len(output): ", len(output))

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

func getSortedSmallIds(thingDataSample map[string]int) []int {
	// Create sorted array of smallIds to keep them ordered
	var smallIds []int
	for k := range thingDataSample {
		if k == "ts" {
			continue
		}
		smallId, _ := strconv.Atoi(k)
		smallIds = append(smallIds, smallId)
	}
	sort.Ints(smallIds)
	return smallIds
}

func mapArrayTo2DArray(mapArray []map[string]int, smallIds []int) [][]int {
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

func exportToCsv(thingData2DArray [][]int, smallIds []int, smallIdToInfoMap map[string]SensorInfo, thingId string, fileName string) {
	err := os.MkdirAll("srv_files/"+thingId, 0777)
	if err != nil {
		panic(err)
	}
	csvFile, err := os.Create(fileName)
	if err != nil {
		panic(err)
	}
	defer csvFile.Close()

	csvWriter := csv.NewWriter(csvFile)
	defer csvWriter.Flush()

	// Write header
	header := make([]string, len(smallIds)+1)
	header[0] = "Timestamp"
	for i, smallId := range smallIds {
		strSmallId := strconv.Itoa(smallId)
		sensorName := smallIdToInfoMap[strSmallId].Name
		header[i+1] = sensorName
	}
	err = csvWriter.Write(header)
	if err != nil {
		panic(err)
	}

	// Write data
	for _, thingDataRecord := range thingData2DArray {
		strArray := make([]string, len(thingDataRecord))
		for i, value := range thingDataRecord {
			strArray[i] = strconv.Itoa(value)
		}
		err = csvWriter.Write(strArray)
		if err != nil {
			panic(err)
		}
	}
}

func replaceSmallIdsWithIds(thingDataArray []map[string]int, smallIdToInfoMap map[string]SensorInfo) []map[string]int {
	for _, thingDataItem := range thingDataArray {
		for smallId, sensorInfo := range smallIdToInfoMap {
			thingDataItem[sensorInfo.ID.Hex()] = thingDataItem[smallId]
			delete(thingDataItem, smallId)
		}
	}

	return thingDataArray
}
