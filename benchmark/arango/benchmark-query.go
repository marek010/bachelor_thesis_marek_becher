package arango

import (
	"benchmark/shared"
	"context"
	"log"
	"time"

	"github.com/arangodb/go-driver"
)

// Exectues the query repeatedly and returns the average execution time
func benchmarkQuery(db driver.Database, query string) time.Duration {
    var total time.Duration
    
    for i := 0; i < shared.NumRepetitions; i++ {
        start := time.Now()
        
        cursor, err := db.Query(context.Background(), query, nil)
        shared.HandleError(err, "Error executing query")
        cursor.Close()
        
        elapsed := time.Since(start)
        total += elapsed
        log.Printf("Run %d duration: %s", i+1, elapsed)
    }

    avg := total / time.Duration(shared.NumRepetitions)
    log.Printf("Average Execution Time: %s\n\n", avg)
    return avg
}