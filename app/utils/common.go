package utils

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"gopkg.in/mgo.v2/bson"
)

func UnitMilli(t time.Time) int64 {
	return t.Round(time.Millisecond).UnixNano() / (int64(time.Millisecond) / int64(time.Nanosecond))
}

func CurrentTimeInMilli() int64 {
	return UnitMilli(time.Now())
}

func ToMap(s interface{}) (map[string]interface{}, error) {
	var stringInterfaceMap map[string]interface{}
	itr, _ := bson.Marshal(s)
	err := bson.Unmarshal(itr, &stringInterfaceMap)
	return stringInterfaceMap, err

}

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
