package postgres

import (
	"benchmark/shared"
	"database/sql"
	"log"
	"time"
)

// Exectues the query repeatedly and returns the average execution time
func benchmarkQuery(db *sql.DB, query string, args ...interface{}) time.Duration {
	var total time.Duration
	
	for i := 0; i < shared.NumRepetitions; i++ {
			// The graph extension age needs to be loaded before every query
			db.Exec("LOAD 'age';")
			db.Exec("SET search_path = ag_catalog, '$user', public;")
			start := time.Now()
			
			rows, err := db.Query(query, args...)
			shared.HandleError(err, "Error executing query")
			rows.Close()
			
			elapsed := time.Since(start)
			total += elapsed
			log.Printf("Run %d duration: %s", i+1, elapsed)
	}

	avg := total / time.Duration(shared.NumRepetitions)
	log.Printf("Average Execution Time: %s\n\n", avg)
	return avg
}