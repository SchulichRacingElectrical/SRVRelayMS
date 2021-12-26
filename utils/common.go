package utils

import "time"

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
