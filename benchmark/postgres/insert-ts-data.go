package postgres

import (
	"benchmark/shared"
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/lib/pq"
)

func insertTSData(db *sql.DB) {
	// Create table and turn it into hypertable with Timescale extension
	log.Print("Creating table...")

	db.Exec("DROP TABLE IF EXISTS ts_table")
	db.Exec("CREATE TABLE IF NOT EXISTS ts_table (timestamp TIMESTAMPTZ NOT NULL, station_id BIGINT NOT NULL, value INT);")
	db.Exec("SELECT create_hypertable('ts_table', by_range('timestamp', INTERVAL '1 day'));")

	log.Print(" Success! \n")

	// Insert timeseries data using COPY
	fmt.Printf("Inserting data...")

	csvFilePath := "/data/ts_data.csv"
	
  copyCommand := fmt.Sprintf(`
    COPY ts_table(timestamp, station_id, value)
    FROM '%s'
    WITH (FORMAT csv, HEADER true);
  `, csvFilePath)

  startTime := time.Now()
  _, err := db.Exec(copyCommand)
  shared.HandleError(err, "Error inserting data")
  elapsedTime := time.Since(startTime)

  fmt.Printf("Elapsed Time: %v\n", elapsedTime)

	// Save measured duration
	shared.SaveResults([]shared.QueryResult{{
		Database:    "postgresql",
		QueryType:   "timeseries",
		QueryName:   "timeseries_insertion",
		AverageTime: float64(elapsedTime.Milliseconds()),
	}})

	// Determine table size
	log.Print("Determining table size...")
	var tableSize string
	err = db.QueryRow("SELECT pg_size_pretty(hypertable_size('ts_table'))").Scan(&tableSize)
	shared.HandleError(err, "Error querying hypertable size")

	fmt.Printf("Hypertable size: %s\n", tableSize)

	// Create indexes
	log.Print("Creating Indexes...")
	db.Exec("CREATE INDEX ON ts_table (station_id, timestamp DESC);")
	db.Exec("CREATE INDEX ON ts_table (value, timestamp DESC);")
	log.Print(" Success! \n")
}
