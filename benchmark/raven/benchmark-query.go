package raven

import (
	"benchmark/shared"
	"log"
	"time"

	ravendb "github.com/ravendb/ravendb-go-client"
)

// Exectues the query repeatedly and returns the average execution time
func benchmarkQuery(session *ravendb.DocumentSession, baseQuery string) time.Duration {
		var total time.Duration
    
    for i := 0; i < shared.NumRepetitions; i++ {
			query := session.Advanced().RawQuery(baseQuery).NoCaching()
			start := time.Now()

			var results []map[string]interface{}
			err := query.GetResults(&results)
			shared.HandleError(err, "Error getting results")

			elapsed := time.Since(start)
        total += elapsed
        log.Printf("Run %d duration: %s", i+1, elapsed)
    }

    avg := total / time.Duration(shared.NumRepetitions)
    log.Printf("Average Execution Time: %s\n\n", avg)
    return avg
}