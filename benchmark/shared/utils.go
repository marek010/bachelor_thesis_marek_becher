package shared

import "log"

const (
	NumRepetitions = 10
	StartTime      = "2024-05-16T17:10:01"
	EndTime        = "2024-05-31T16:45:06"
	ValueThreshold = 40
	StationIDs = "10799515428237,9575977656619,11288438576825,20053112967470,11086574027625"
	LatestTimestampStationId = "10799515428237"
)

func HandleError(err error, message string) {
	if err != nil {
			log.Fatalf("%s: %v", message, err)
	}
}