package arango

import (
	"benchmark/shared"
	"fmt"
	"log"
	"time"

	"github.com/arangodb/go-driver"
)

func latestTimestampQuery(db driver.Database) time.Duration {
    query := fmt.Sprintf(`
        FOR doc IN ts_table
        FILTER doc.station_id == %s
        SORT doc.timestamp DESC
        LIMIT 1
        RETURN doc`, shared.LatestTimestampStationId)

    log.Printf("Latest timestamp:")
    return benchmarkQuery(db, query)
}

func simpleRangeQuery(db driver.Database) time.Duration {
    query := fmt.Sprintf(`
        FOR doc IN ts_table
        FILTER doc.timestamp >= DATE_ISO8601('%s') 
        AND doc.timestamp <= DATE_ISO8601('%s')
        AND doc.station_id IN [%s]
        RETURN doc`, shared.StartTime, shared.EndTime, shared.StationIDs)

    log.Printf("Simple range:")
    return benchmarkQuery(db, query)
}

func aggregationQuery(db driver.Database) time.Duration {
    query := fmt.Sprintf(`
        FOR doc IN ts_table
        FILTER doc.timestamp >= DATE_ISO8601('%s') 
        AND doc.timestamp <= DATE_ISO8601('%s')
        AND doc.station_id IN [%s]
        COLLECT station_id = doc.station_id 
        AGGREGATE avgValue = AVG(doc.value)
        RETURN { 
            station_id: station_id, 
            average_value: avgValue 
        }`, shared.StartTime, shared.EndTime, shared.StationIDs)

    log.Printf("Aggregation:")
    return benchmarkQuery(db, query)
}

func downsamplingQuery(db driver.Database) time.Duration {
    query := fmt.Sprintf(`
        FOR doc IN ts_table
        FILTER IN_RANGE(doc.timestamp, '%s', '%s', true, true)
        COLLECT day = DATE_TRUNC(doc.timestamp, "day") 
        AGGREGATE average_value = AVERAGE(doc.value)
        RETURN { 
            day_interval: day, 
            average_value: average_value 
        }`, shared.StartTime, shared.EndTime)

    log.Printf("Downsampling:")
    return benchmarkQuery(db, query)
}

func rangeQueryWithFilter(db driver.Database) time.Duration {
    query := fmt.Sprintf(`
        FOR doc IN ts_table
        FILTER doc.timestamp >= DATE_ISO8601('%s') 
        AND doc.timestamp <= DATE_ISO8601('%s')
        AND doc.station_id IN [%s]
        AND doc.value > %d
        RETURN doc`, shared.StartTime, shared.EndTime, shared.StationIDs, shared.ValueThreshold)

    log.Printf("Range with value filter:")
    return benchmarkQuery(db, query)
}

func runAllTimeseriesQueries(db driver.Database) {
    queries := []shared.QueryDefinition{
        {Name: "latest_timestamp", Function: func() time.Duration { return latestTimestampQuery(db) }},
        {Name: "simple_range", Function: func() time.Duration { return simpleRangeQuery(db) }},
        {Name: "range_with_filter", Function: func() time.Duration { return rangeQueryWithFilter(db) }},
        {Name: "aggregation", Function: func() time.Duration { return aggregationQuery(db) }},
        {Name: "downsampling", Function: func() time.Duration { return downsamplingQuery(db) }},
    }

    shared.RunQueries("arangodb", "timeseries", queries)
}