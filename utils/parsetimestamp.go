package utils

import (
	"strconv"
	"time"
)

type FormattedTimestamp struct {
	Date string
	Time string
}

// ParseTimestamp converts a Unix timestamp to HH:MM and date, separately
func ParseTimestamp(timestamp string) FormattedTimestamp {
	unixTime, err := strconv.ParseInt(timestamp, 10, 64)
	if err != nil {
		return FormattedTimestamp{
			Time: "00:00",
			Date: "01/01/1970",
		}
	}
	timeObj := time.Unix(unixTime/1000, 0)
	location, err := time.LoadLocation("America/Los_Angeles")
	if err != nil {
		return FormattedTimestamp{
			Time: "00:00",
			Date: "01/01/1970",
		}
	}
	localTime := timeObj.In(location)
	formattedTime := localTime.Format("15:04")
	formattedDate := localTime.Format("02/01/2006")
	return FormattedTimestamp{
		Time: formattedTime,
		Date: formattedDate,
	}
}
