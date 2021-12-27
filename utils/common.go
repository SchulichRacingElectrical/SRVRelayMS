package utils

import (
	"time"

	"gopkg.in/mgo.v2/bson"
)

// UnixMilli use to get ms of given time
// @params - time
// return - ms of the given time
func UnitMilli(t time.Time) int64 {
	return t.Round(time.Millisecond).UnixNano() / (int64(time.Millisecond) / int64(time.Nanosecond))
}

// CurrentTimeInMill uset to get current time in ms
// Used to obtain current timestamp
// return = curren timestamp in ms
func CurrentTimeInMilli() int64 {
	return UnitMilli(time.Now())
}

type CustomBson struct{}

type BsonWrapper struct {
	Set interface{} `json:"$set,omitempty" bson:"$set,omitempty"`
}

func ToMap(s interface{}) (map[string]interface{}, error) {
	var stringInterfaceMap map[string]interface{}
	itr, _ := bson.Marshal(s)
	err := bson.Unmarshal(itr, &stringInterfaceMap)
	return stringInterfaceMap, err

}

func (customBson *CustomBson) Set(data interface{}) (map[string]interface{}, error) {
	s := BsonWrapper{Set: data}
	return ToMap(s)
}
