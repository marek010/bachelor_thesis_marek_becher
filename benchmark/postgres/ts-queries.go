package postgres

import (
	"benchmark/shared"
	"database/sql"
	"fmt"
	"log"
	"time"
)

func simpleRangeQuery(db *sql.DB) time.Duration {
	query := fmt.Sprintf(`
			SELECT * FROM ts_table
			WHERE timestamp BETWEEN '%s' AND '%s'
			AND station_id IN (%s)
	`, shared.StartTime, shared.EndTime, shared.StationIDs)

	log.Printf("Simple range:")
	return benchmarkQuery(db, query)
}

func aggregationQuery(db *sql.DB) time.Duration {
	query := fmt.Sprintf(`
			SELECT avg(value) AS average_value
			FROM ts_table
			WHERE timestamp BETWEEN '%s' AND '%s'
			AND station_id IN (%s)
	`, shared.StartTime, shared.EndTime, shared.StationIDs)

	log.Printf("Aggregation:")
	return benchmarkQuery(db, query)
}

func downsamplingQuery(db *sql.DB) time.Duration {
	query := fmt.Sprintf(`
			SELECT time_bucket('1 day', timestamp) AS day_interval,
			avg(value) AS average_value
			FROM ts_table
			WHERE timestamp BETWEEN '%s' AND '%s'
			GROUP BY day_interval
			ORDER BY day_interval
	`, shared.StartTime, shared.EndTime)

	log.Printf("Downsampling:")
	return benchmarkQuery(db, query)
}

func rangeQueryWithFilter(db *sql.DB) time.Duration {
	query := fmt.Sprintf(`
			SELECT * FROM ts_table
			WHERE timestamp BETWEEN '%s' AND '%s'
			AND station_id IN (%s)
			AND value > %d
	`, shared.StartTime, shared.EndTime, shared.StationIDs, shared.ValueThreshold)

	log.Printf("Range with value filter:")
	return benchmarkQuery(db, query)
}

func latestTimestampQuery(db *sql.DB) time.Duration {
	query := fmt.Sprintf(`
			SELECT *
			FROM ts_table
			WHERE station_id = %s
			ORDER BY timestamp DESC
			LIMIT 1
	`, shared.LatestTimestampStationId)

	log.Printf("Latest timestamp:")
	return benchmarkQuery(db, query)
}

func runAllTimeseriesQueries(db *sql.DB)  {
	queries := []shared.QueryDefinition{
		{Name: "latest_timestamp", Function: func() time.Duration { return latestTimestampQuery(db) }},
		{Name: "simple_range", Function: func() time.Duration { return simpleRangeQuery(db) }},
		{Name: "range_with_filter", Function: func() time.Duration { return rangeQueryWithFilter(db) }},
		{Name: "aggregation", Function: func() time.Duration { return aggregationQuery(db) }},
		{Name: "downsampling", Function: func() time.Duration { return downsamplingQuery(db) }},
	}

	shared.RunQueries("postgresql", "timeseries", queries)
}
