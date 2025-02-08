package raven

import (
	"benchmark/shared"
	"fmt"
	"log"
	"time"

	ravendb "github.com/ravendb/ravendb-go-client"
)

func simpleRangeQuery(session *ravendb.DocumentSession) time.Duration {
	query := fmt.Sprintf(`
		from index "TimeSeriesIndex"
		where Timestamp > '%s' and Timestamp < '%s'
		and StationId in (%s)
		select Value, StationId, Timestamp
	`, shared.StartTime, shared.EndTime, shared.StationIDs)

	log.Printf("Simple range:")
	return benchmarkQuery(session, query)
}

func aggregationQuery(session *ravendb.DocumentSession) time.Duration {
	query := fmt.Sprintf(`
		from "TSCollections"
		select timeseries (
				from Timeseries
				between '%s' and '%s'
				where Values[0] in (%s)
				group by '1 month'
				select avg()
		)
	`, shared.StartTime, shared.EndTime, shared.StationIDs)

	log.Printf("Aggregation:")
	return benchmarkQuery(session, query)
}

func downsamplingQuery(session *ravendb.DocumentSession) time.Duration {
	query := fmt.Sprintf(`
		from "TSCollections"
		select timeseries (
				from Timeseries
				between '%s' and '%s'
				group by '1 day'
				select avg()
		)
	`, shared.StartTime, shared.EndTime)

	log.Printf("Downsampling:")
	return benchmarkQuery(session, query)
}

func rangeQueryWithFilter(session *ravendb.DocumentSession) time.Duration {
	query := fmt.Sprintf(`
		from index "TimeSeriesIndex"
		where Timestamp > '%s' and Timestamp < '%s'
		and StationId in (%s)
		and Value > %d
		select Value, StationId, Timestamp
	`, shared.StartTime, shared.EndTime, shared.StationIDs, shared.ValueThreshold)

	log.Printf("Range with value filter:")
	return benchmarkQuery(session, query)
}

func latestTimestampQuery(session *ravendb.DocumentSession) time.Duration {
	query := fmt.Sprintf(`
		from index "TimeSeriesIndex"
		where StationId = %s
		order by Timestamp desc
		select Value
		limit 1
	`, shared.LatestTimestampStationId)

	log.Printf("Latest timestamp:")
	return benchmarkQuery(session, query)
}

func runAllTimeseriesQueries(session *ravendb.DocumentSession)  {
	queries := []shared.QueryDefinition{
		{Name: "latest_timestamp", Function: func() time.Duration { return latestTimestampQuery(session) }},
		{Name: "simple_range", Function: func() time.Duration { return simpleRangeQuery(session) }},
		{Name: "range_with_filter", Function: func() time.Duration { return rangeQueryWithFilter(session) }},
		{Name: "aggregation", Function: func() time.Duration { return aggregationQuery(session) }},
		{Name: "downsampling", Function: func() time.Duration { return downsamplingQuery(session) }},
	}

	shared.RunQueries("ravendb", "timeseries", queries)
}