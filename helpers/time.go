package helpers

import (
	"strconv"
	"strings"
	"time"
)

func Convert2DateTime(dateTimeStr string) time.Time {
	parsedDatetime := strings.Split(dateTimeStr, "T")
	parsedDate := strings.Split(parsedDatetime[0], "-")
	parsedTime := strings.Split(parsedDatetime[1], ":")

	year, _ := strconv.Atoi(parsedDate[0])
	month, _ := strconv.Atoi(parsedDate[1])
	day, _ := strconv.Atoi(parsedDate[2])

	hour, _ := strconv.Atoi(parsedTime[0])
	min, _ := strconv.Atoi(parsedTime[1])
	sec, _ := strconv.Atoi(parsedTime[2])

	return time.Date(year, time.Month(month), day, hour, min, sec, 0, time.UTC)
}
