package subscriber

import (
	"context"
	"database-ms/app/model"
	"database-ms/app/services"
	"database-ms/config"
	"fmt"
	"log"
	"math"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"

	"encoding/csv"
	"encoding/json"

	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Message struct {
	Active  bool   `json:"active"`
	ThingId string `json:"THING"`
}

type SensorInfo struct {
	Id   uuid.UUID
	Name string
}

func Initialize(conf *config.Configuration, db *gorm.DB) {
	redisClient := redis.NewClient(&redis.Options{
		Addr:     conf.RedisUrl + ":" + conf.RedisPort,
		Password: conf.RedisPassword,
	})

	go awaitThingDataSessions(redisClient, db, conf)
}

func awaitThingDataSessions(redisClient *redis.Client, db *gorm.DB, conf *config.Configuration) {
	ctx := context.Background()
	subscriber := redisClient.Subscribe(ctx, "THING_CONNECTION")
	defer subscriber.Close()
	connectionChannel := subscriber.Channel()
	sessionService := services.NewSessionService(db, conf)

	for msg := range connectionChannel {
		message := Message{}
		json.Unmarshal([]byte(msg.Payload), &message)
		if message.Active {
			thingObjId, err := uuid.Parse(message.ThingId)
			if err != nil {
				panic(err)
			}
			session := &model.Session{
				StartTime: time.Now().UnixMilli(),
				EndTime:   nil,
				ThingId:   thingObjId,
				Name:      uuid.NewString(),
			}
			perr := sessionService.CreateSession(ctx, session)
			if perr != nil {
				panic(err)
			}
			log.Println("Session created with ID: " + session.Id.String())
			thingId, err := uuid.Parse(message.ThingId)
			if err != nil {
				// Failed - Do something
			}
			go thingDataSession(thingId, session, redisClient, db, conf)
		}
	}
}

func thingDataSession(thingId uuid.UUID, session *model.Session, redisClient *redis.Client, db *gorm.DB, conf *config.Configuration) {
	ctx := context.Background()
	subscriber := redisClient.Subscribe(ctx, "THING_"+thingId.String())
	defer subscriber.Close()
	log.Println("Thing Data Session Started for " + thingId.String())
	thingDataChannel := subscriber.Channel()

	//datumService := services.NewDatumService(db, conf)
	sessionService := services.NewSessionService(db, conf)

	for msg := range thingDataChannel {
		message := Message{}
		json.Unmarshal([]byte(msg.Payload), &message)

		if !message.Active {
			// Get thing data from redis
			thingData, err := redisClient.LRange(ctx, "THING_"+thingId.String(), 0, -1).Result()
			if err != nil {
				panic(err)
			}

			// Delete thing data from redis
			err = redisClient.Del(ctx, "THING_"+thingId.String()).Err()
			if err != nil {
				panic(err)
			}

			// Parse thing data from JSON
			var thingDataArray []map[string]float64
			for _, thingDataItem := range thingData {
				var thingDataItemMap map[string]float64
				json.Unmarshal([]byte(thingDataItem), &thingDataItemMap)
				thingDataArray = append(thingDataArray, thingDataItemMap)
			}

			// Sort the thing data by timestamp
			sort.Slice(thingDataArray, func(i, j int) bool {
				return thingDataArray[i]["ts"] < thingDataArray[j]["ts"]
			})

			// Get sensor list
			sensorService := services.NewSensorService(db, conf)
			sensors, perr := sensorService.FindByThingId(ctx, thingId)
			if perr != nil {
				panic(err)
			}
			sort.Slice(sensors, func(i, j int) bool {
				return sensors[i].Name < sensors[j].Name
			})

			// Get small ids in the respective order of the sensors
			var smallIds []int
			highestFrequency := 0.0
			for _, sensor := range sensors {
				smallIds = append(smallIds, sensor.SmallId)
				if sensor.Frequency > int32(highestFrequency) {
					highestFrequency = float64(sensor.Frequency)
				}
			}

			// Get the timestamp interval
			timestampInterval := math.Round(1000 / highestFrequency)

			// Process thing data to fill missing sensor values and missing timestamps
			thingDataArray = fillMissingValues(thingDataArray, smallIds, int(timestampInterval))

			// Create map of SmallId to ID and Name
			smallIdToInfoMap := make(map[string]SensorInfo)
			for _, sensor := range sensors {
				smallId := fmt.Sprint(sensor.SmallId)
				smallIdToInfoMap[smallId] = SensorInfo{Id: sensor.Id, Name: sensor.Name}
			}

			// Re-fetch the session in case a user has modified it
			session, perr = sessionService.FindById(ctx, session.Id)
			if perr != nil {
				panic(err)
			}

			// Update the session
			endTime := session.StartTime + int64(thingDataArray[len(thingDataArray)-1]["ts"])
			session.EndTime = &endTime
			perr = sessionService.UpdateSession(ctx, session)
			if perr != nil {
				panic(err)
			}

			// Save thing data to csv
			exportToCsv(
				thingDataArray,
				smallIds,
				smallIdToInfoMap,
				conf.FilePath+thingId.String(),
				conf.FilePath+thingId.String()+"/"+session.Name+".csv",
			)

			// // Process non-linear thing data
			// thingDataArray = replaceSmallIdsWithIds(thingDataArray, smallIdToInfoMap)
			// datumArray := make([]*model.Datum, len(thingDataArray)*len(smallIds))
			// for i, thingDataItem := range thingDataArray {
			// 	for j, smallId := range smallIds {
			// 		strSmallId := strconv.Itoa(smallId)
			// 		sensorId := smallIdToInfoMap[strSmallId].Id
			// 		datumArray[i*len(smallIds)+j] = &model.Datum{
			// 			SessionId: session.Id,
			// 			SensorId:  sensorId,
			// 			Value:     float64(thingDataItem[sensorId.String()]),
			// 			Timestamp: int64(thingDataItem["ts"]),
			// 		}
			// 	}
			// }

			// // Save thing data in the database
			// datumService.CreateMany(ctx, datumArray)

			log.Println("Thing Data Session Ended for " + thingId.String())
			return
		}
	}
}

func fillMissingValues(
	thingDataArray []map[string]float64,
	smallIds []int,
	interval int,
) []map[string]float64 {
	// Get the default values of all the sensors
	currentDataMap := copyMap(thingDataArray[0])
	for _, smallId := range smallIds {
		if _, ok := currentDataMap[strconv.Itoa(smallId)]; !ok {
			currentDataMap[strconv.Itoa(smallId)] = 0
		}
	}

	// Fill from timestamp 0 to the first timestamp
	firstTimeStamp := int(currentDataMap["ts"])
	if firstTimeStamp != 0 {
		currentTimeStamp := 0
		for currentTimeStamp != firstTimeStamp {
			newMap := copyMap(currentDataMap)
			newMap["ts"] = float64(currentTimeStamp)
			thingDataArray = append([]map[string]float64{newMap}, thingDataArray...)
			currentTimeStamp += interval
		}
	}

	// Fill timestamp gaps
	println(len(thingDataArray))
	var filledDataArray []map[string]float64
	prevTimestamp := float64(0)
	prevMap := thingDataArray[0]
	for _, datum := range thingDataArray {
		currentTimestamp := datum["ts"]
		if prevTimestamp != 0 && int(currentTimestamp-prevTimestamp) != interval {
			for currentTimestamp-prevTimestamp != 0 {
				fillItem := copyMap(datum)
				fillItem["ts"] = currentTimestamp
				filledDataArray = append(filledDataArray, fillItem)
			}
		} else {
			filledDataArray = append(filledDataArray, datum)
		}
		prevMap = datum
		prevTimestamp = currentTimestamp
	}
	_ = prevMap
	println(len(filledDataArray))

	// Populate each smallId with its current or previous value
	for _, thingDataItem := range filledDataArray {
		for key := range currentDataMap {
			if _, ok := thingDataItem[key]; !ok {
				thingDataItem[key] = currentDataMap[key]
			} else {
				currentDataMap[key] = thingDataItem[key]
			}
		}
	}

	return filledDataArray
}

func copyMap(source map[string]float64) map[string]float64 {
	dest := make(map[string]float64)
	for key, value := range source {
		dest[key] = value
	}
	return dest
}

func exportToCsv(
	thingDataArray []map[string]float64,
	smallIds []int,
	smallIdToInfoMap map[string]SensorInfo,
	filePath string,
	fileName string,
) {
	// Attempt to create the directory
	err := os.MkdirAll(filePath, 0777)
	if err != nil {
		panic(err)
	}

	// Attempt to create the file
	csvFile, err := os.Create(fileName)
	if err != nil {
		panic(err)
	}
	defer csvFile.Close()

	// Create the CSV writer
	csvWriter := csv.NewWriter(csvFile)
	defer csvWriter.Flush()

	// Write header
	var header []string
	header = append(header, "Timestamp")
	for _, smallId := range smallIds {
		strSmallId := strconv.Itoa(smallId)
		sensorName := smallIdToInfoMap[strSmallId].Name
		header = append(header, sensorName)
	}
	err = csvWriter.Write(header)
	if err != nil {
		panic(err)
	}

	// Write data
	for _, datum := range thingDataArray {
		var strArray []string
		timestamp := fmt.Sprintf("%.15f", datum["ts"])
		timestamp = strings.TrimRight(strings.TrimRight(timestamp, "0"), ".")
		strArray = append(strArray, timestamp)
		for _, smallId := range smallIds {
			stringValue := fmt.Sprintf("%.15f", datum[strconv.Itoa(smallId)])
			stringValue = strings.TrimRight(strings.TrimRight(stringValue, "0"), ".")
			strArray = append(strArray, stringValue)
		}
		err = csvWriter.Write(strArray)
		if err != nil {
			panic(err)
		}
	}
}

func replaceSmallIdsWithIds(thingDataArray []map[string]float64, smallIdToInfoMap map[string]SensorInfo) []map[string]float64 {
	for _, thingDataItem := range thingDataArray {
		for smallId, sensorInfo := range smallIdToInfoMap {
			thingDataItem[sensorInfo.Id.String()] = thingDataItem[smallId]
			delete(thingDataItem, smallId)
		}
	}

	return thingDataArray
}
