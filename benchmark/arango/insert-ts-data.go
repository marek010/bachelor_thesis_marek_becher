package arango

import (
	"benchmark/shared"
	"context"
	"fmt"
	"log"
	"os/exec"
	"time"

	"github.com/arangodb/go-driver"
)

func insertTSData(col driver.Collection) {
		// Clear collection and indexes
    ctx := context.Background()
    
    err := col.Truncate(ctx)
    shared.HandleError(err, "Failed to empty collection")
    
    indexes, _ := col.Indexes(ctx)
    for _, index := range indexes {
        index.Remove(ctx)
    }

		// Insert timeseries data
    log.Println("Starting data insertion...")
    start := time.Now()

    cmd := exec.Command("docker", "exec", "arangodb-container", "arangoimport",
        "--file", "../../../data/ts_data.csv",
        "--type", "csv",
        "--collection", collection,
        "--server.password", password,
        "--server.database", dbname)

    _, err = cmd.CombinedOutput()
    shared.HandleError(err, "Failed to import data")

		elapsedTime := time.Since(start)

    fmt.Printf("Import completed in %s\n", elapsedTime)

		// Save measured duration
		shared.SaveResults([]shared.QueryResult{{
			Database:    "arangodb",
			QueryType:   "timeseries",
			QueryName:   "timeseries_insertion",
			AverageTime: float64(elapsedTime.Milliseconds()),
		}})

		// Create indexes on all fields
		fmt.Print("Creating indexes...")
    indexFields := []string{"timestamp", "station_id", "value"}
    for _, field := range indexFields {
        col.EnsurePersistentIndex(ctx, []string{field}, &driver.EnsurePersistentIndexOptions{})
    }

    fmt.Print(" Success!")
}