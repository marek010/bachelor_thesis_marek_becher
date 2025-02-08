package raven

import (
	"benchmark/shared"
	"fmt"
	"log"
	"os/exec"
	"time"

	ravendb "github.com/ravendb/ravendb-go-client"
)

func NewTimeSeriesIndex() *ravendb.IndexCreationTask {
	index := ravendb.NewIndexCreationTask("TimeSeriesIndex")
	index.Map = `
		timeSeries.map('TSCollections', 'Timeseries', function (segment) {
				return segment.Entries.map(entry => {
						return {
								Value: entry.Values[1]
								StationId: entry.Values[0]
								Timestamp: entry.Timestamp
						}
				});
		});
	`

	return index
}

func insertTSData(store *ravendb.DocumentStore) {
	fmt.Printf("Inserting data...")

	startTime := time.Now()

	//Use python script to insert data as go driver does not support timeseries bulk insertion
	cmd := exec.Command("python3", "raven/insert-ts-data-script.py")
	_, err := cmd.CombinedOutput()
	shared.HandleError(err, "Error running python insertion script")
	cmd.Run()

	elapsedTime := time.Since(startTime)

	// Save measured duration
	shared.SaveResults([]shared.QueryResult{{
		Database:    "ravendb",
		QueryType:   "timeseries",
		QueryName:   "timeseries_insertion",
		AverageTime: float64(elapsedTime.Milliseconds()),
	}})

	log.Print(" Success! \n")
	fmt.Printf("Elapsed Time: %v\n", elapsedTime)

	// Add index for timeseries data
	log.Print("Creating Index...")
	index := NewTimeSeriesIndex()
	err = index.Execute(store, nil, dbname)
	shared.HandleError(err, "Error creating index")
	log.Printf("Success!")
}
