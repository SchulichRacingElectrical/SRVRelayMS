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
	go AwaitThingDataSessions(redisClient, db, conf)
}

func AwaitThingDataSessions(redisClient *redis.Client, db *gorm.DB, conf *config.Configuration) {
	ctx := context.Background()
	subscriber := redisClient.Subscribe(ctx, "THING_CONNECTION")
	defer subscriber.Close()
	connectionChannel := subscriber.Channel()
	sessionService := services.NewSessionService(db, conf)

	for msg := range connectionChannel {
		message := Message{}
		json.Unmarshal([]byte(msg.Payload), &message)
		if message.Active {
			thingId, err := uuid.Parse(message.ThingId)
			if err != nil {
				log.Println("Failed to read thing ID.")
				panic(err)
			}
			generated := true
			session := &model.Session{
				StartTime: time.Now().UnixMilli(),
				EndTime:   nil,
				ThingId:   thingId,
				Name:      uuid.NewString(),
				Generated: &generated,
			}
			perr := sessionService.CreateSession(ctx, session)
			if perr != nil {
				log.Println("Failed to create session.")
				panic(err)
			}
			log.Println("Session created with ID: " + session.Id.String())
			go ThingDataSession(thingId, session, redisClient, db, conf)
		}
	}
}

func ThingDataSession(thingId uuid.UUID, session *model.Session, redisClient *redis.Client, db *gorm.DB, conf *config.Configuration) {
	ctx := context.Background()
	subscriber := redisClient.Subscribe(ctx, "THING_"+thingId.String())
	defer subscriber.Close()
	log.Println("Thing Data Session Started for " + thingId.String())
	thingDataChannel := subscriber.Channel()
	datumService := services.NewDatumService(db, conf)
	sessionService := services.NewSessionService(db, conf)

	// If there is an error anywhere, we want to delete the session
	defer func() {
		if err := recover(); err != nil {
			log.Println("Failed to parse session data, deleting session now...")
			sessionService.DeleteSession(ctx, session.Id)
		}
	}()

	for msg := range thingDataChannel {
		message := Message{}
		json.Unmarshal([]byte(msg.Payload), &message)

		if !message.Active {
			// Get thing data from redis
			thingData, err := redisClient.LRange(ctx, "THING_"+thingId.String(), 0, -1).Result()
			if err != nil {
				panic(err)
			}

			// Delete thing data from redis, next connection will clean up if needed
			redisClient.Del(ctx, "THING_"+thingId.String())

			// Parse thing data from JSON
			var thingDataArray []map[string]float64
			for _, thingDataItem := range thingData {
				var thingDataItemMap map[string]float64
				json.Unmarshal([]byte(thingDataItem), &thingDataItemMap)
				thingDataArray = append(thingDataArray, thingDataItemMap)
			}

			// Exit if no data is found
			if len(thingDataArray) == 0 {
				panic("empty data")
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
			thingDataArray = FillMissingValues(thingDataArray, smallIds, int(timestampInterval))

			// Create map of SmallId to ID and Name
			smallIdToInfoMap := make(map[string]SensorInfo)
			for _, sensor := range sensors {
				smallId := fmt.Sprint(sensor.SmallId)
				smallIdToInfoMap[smallId] = SensorInfo{Id: sensor.Id, Name: sensor.Name}
			}

			// Re-fetch the session in case a user has modified it
			session, perr = sessionService.FindById(ctx, session.Id)

			// Update the session

			endTime := time.Now().UnixMilli()
			session.EndTime = &endTime
			perr = sessionService.UpdateSession(ctx, session)
			if perr != nil {
				panic(err)
			}

			// Save thing data to a .csv file
			ExportToCsv(
				thingDataArray,
				smallIds,
				smallIdToInfoMap,
				int(timestampInterval),
				conf.FilePath+thingId.String(),
				conf.FilePath+thingId.String()+"/"+session.Name+".csv",
			)

			// Process non-linear thing data
			thingDataArray = ReplaceSmallIdsWithIds(thingDataArray, smallIdToInfoMap)
			var datumArray []*model.Datum
			for _, thingDataItem := range thingDataArray {
				for _, smallId := range smallIds {
					strSmallId := strconv.Itoa(smallId)
					sensorId := smallIdToInfoMap[strSmallId].Id
					datumArray = append(datumArray, &model.Datum{
						SessionId: session.Id,
						SensorId:  sensorId,
						Value:     float64(thingDataItem[sensorId.String()]),
						Timestamp: int64(thingDataItem["ts"]),
					})
				}
			}

			// Save thing data in the database
			datumService.CreateMany(ctx, datumArray)

			log.Println("Thing Data Session Ended for " + thingId.String())
			return
		}
	}
}

func FillMissingValues(
	thingDataArray []map[string]float64,
	smallIds []int,
	interval int,
) []map[string]float64 {
	// Get the default values of all the sensors
	currentDataMap := CopyMap(thingDataArray[0])
	for _, smallId := range smallIds {
		if _, ok := currentDataMap[strconv.Itoa(smallId)]; !ok {
			currentDataMap[strconv.Itoa(smallId)] = 0
		}
	}

	// Populate each smallId with its current or previous value
	for _, thingDataItem := range thingDataArray {
		for key := range currentDataMap {
			if _, ok := thingDataItem[key]; !ok {
				thingDataItem[key] = currentDataMap[key]
			} else {
				currentDataMap[key] = thingDataItem[key]
			}
		}
	}

	return thingDataArray
}

func ExportToCsv(
	thingDataArray []map[string]float64,
	smallIds []int,
	smallIdToInfoMap map[string]SensorInfo,
	interval int,
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
	prevTimestamp := 0
	var prevRow []string
	for _, datum := range thingDataArray {
		// Fill gaps after 0
		if prevTimestamp != 0 {
			for int(datum["ts"])-prevTimestamp > interval {
				err = csvWriter.Write(prevRow)
				if err != nil {
					panic(err)
				}
				prevTimestamp += interval
			}
		}

		// Create the csv row and write to the file
		row := CreateCsvRow(datum, smallIds)
		err = csvWriter.Write(row)
		if err != nil {
			panic(err)
		}

		// Set previous values for gap filling
		prevTimestamp = int(datum["ts"])
		prevRow = row
	}
}

func CreateCsvRow(datum map[string]float64, smallIds []int) []string {
	var strArray []string
	timestamp := fmt.Sprintf("%.15f", datum["ts"])
	timestamp = strings.TrimRight(strings.TrimRight(timestamp, "0"), ".")
	strArray = append(strArray, timestamp)
	for _, smallId := range smallIds {
		stringValue := fmt.Sprintf("%.15f", datum[strconv.Itoa(smallId)])
		stringValue = strings.TrimRight(strings.TrimRight(stringValue, "0"), ".")
		strArray = append(strArray, stringValue)
	}
	return strArray
}

func ReplaceSmallIdsWithIds(
	thingDataArray []map[string]float64,
	smallIdToInfoMap map[string]SensorInfo,
) []map[string]float64 {
	for _, thingDataItem := range thingDataArray {
		for smallId, sensorInfo := range smallIdToInfoMap {
			thingDataItem[sensorInfo.Id.String()] = thingDataItem[smallId]
			delete(thingDataItem, smallId)
		}
	}
	return thingDataArray
}

func CopyMap(source map[string]float64) map[string]float64 {
	dest := make(map[string]float64)
	for key, value := range source {
		dest[key] = value
	}
	return dest
}
