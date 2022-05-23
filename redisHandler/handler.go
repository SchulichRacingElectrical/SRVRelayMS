package redisHandler

import (
	"context"
	"database-ms/config"
	"log"

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

	go awaitThingDataSessions(redisClient)
}

func awaitThingDataSessions(redisClient *redis.Client) {
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
			go thingDataSession(message.ThingId, redisClient)
		}
	}
}

func thingDataSession(thingId string, redisClient *redis.Client) {
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

			// Process thing data to fill timeseries

			// Save thing data to mongo
			// Save thing data to csv
		}
	}
}
