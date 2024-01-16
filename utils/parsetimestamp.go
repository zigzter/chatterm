package utils

import (
	"strconv"
	"time"
)

func ParseTimestamp(timestamp string) string {
	unixTime, err := strconv.ParseInt(timestamp, 10, 64)
	if err != nil {
		return "00:00"
	}
	timeObj := time.Unix(unixTime/1000, 0)
	// TODO: Custom timezone
	location, err := time.LoadLocation("America/Los_Angeles")
	if err != nil {
		return "00:00"
	}
	localTime := timeObj.In(location)
	formattedTime := localTime.Format("15:04")
	return formattedTime
}
