// util.go
package util

import (
	"strconv"
	"time"
)

// return 1441007112776 in millisecond
func GetTimestampInMilli() int64 {
	return int64(time.Now().UnixNano() / (1000 * 1000)) // ms
}

// return 1441007112776 in millisecond
func GetTimestampInMilliString() string {
	return strconv.FormatInt(GetTimestampInMilli(), 10)
}
