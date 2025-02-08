package shared

import (
	"encoding/json"
	"fmt"
	"os"
)

type QueryResult struct {
    Database    string  `json:"database"`
    QueryType   string  `json:"queryType"`
    QueryName   string  `json:"queryName"`
    AverageTime float64 `json:"averageTime"`
}

type BenchmarkResults struct {
    Results []QueryResult `json:"results"`
}

func SaveResults(newResults []QueryResult) {
		// Read results json
    data, err := os.ReadFile("results/benchmark_results.json")
    existing := BenchmarkResults{}
    if err == nil {
        json.Unmarshal(data, &existing)
    }

    seen := make(map[string]bool)
    final := BenchmarkResults{}

    // Add new results or overwrite if already exists
    for _, r := range existing.Results {
        key := fmt.Sprintf("%s-%s", r.Database, r.QueryName)
        seen[key] = true
        final.Results = append(final.Results, r)
    }

    for _, r := range newResults {
        key := fmt.Sprintf("%s-%s", r.Database, r.QueryName)
        if seen[key] {
            for i, e := range final.Results {
                if e.Database == r.Database && e.QueryName == r.QueryName {
                    final.Results[i] = r
                    break
                }
            }
        } else {
            final.Results = append(final.Results, r)
        }
    }

    data, err = json.MarshalIndent(final, "", "  ")
    HandleError(err, "Error marshaling results")

    err = os.WriteFile("results/benchmark_results.json", data, 0644)
    HandleError(err, "Error writing results file")
}