// shared/query_runner.go
package shared

import (
	"fmt"
	"time"
)

type QueryDefinition struct {
    Name     string
    Function func() time.Duration
}

// Runs all queries and saves the results in a json file
func RunQueries(dbName string, queryType string, queries []QueryDefinition) []QueryResult {
	var results []QueryResult
	var durations = make(map[string]time.Duration)

	// Execute queries and collect results
	for _, q := range queries {
			duration := q.Function()
			durations[q.Name] = duration
			
			result := QueryResult{
					Database:    dbName,
					QueryType:   queryType,
					QueryName:   q.Name,
					AverageTime: float64(duration.Nanoseconds()) / 1_000_000,
			}
			results = append(results, result)
	}

	// Print results
	fmt.Print("\n Average Execution Times For Queries:\n")
	fmt.Println("-----------------------------------------")
	for _, q := range queries {
			fmt.Printf("%s: %s\n", q.Name, durations[q.Name])
	}

	SaveResults(results)
	return results
}