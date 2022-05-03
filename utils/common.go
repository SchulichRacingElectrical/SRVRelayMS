package utils

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
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

// Removes duplicates from a int slice
// source: https://www.golangprograms.com/remove-duplicate-values-from-slice.html
// @params - integer slice
// return - integers slice without duplicates
func Unique(intSlice []int) []int {
	keys := make(map[int]bool)
	list := []int{}
	for _, entry := range intSlice {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}

func IdInSlice(a primitive.ObjectID, list []primitive.ObjectID) bool {
	for _, b := range list {
		if a == b {
			return true
		}
	}
	return false
}
