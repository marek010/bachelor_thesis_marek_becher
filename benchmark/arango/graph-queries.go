package arango

import (
	"benchmark/shared"
	"fmt"
	"log"
	"time"

	"github.com/arangodb/go-driver"
)

const (
    startStationId = "5230109759983"
    endStationId   = "32116902472201"
)

func getAllNeighbours(db driver.Database) time.Duration {
    query := fmt.Sprintf(`
        FOR a IN stations
        FILTER a._key == '%s'
        FOR b IN 1..1 OUTBOUND a GRAPH '%s'
        RETURN b.name`, startStationId, graphName)

		log.Printf("Get neighbours:")
    return benchmarkQuery(db, query)
}

func search3HopPath(db driver.Database) time.Duration {
    query := fmt.Sprintf(`
        FOR n IN stations
        FILTER n._key == '%s'
        FOR v IN 1..3 OUTBOUND n GRAPH '%s'
        FILTER v._key == '%s'
        RETURN n.name`, startStationId, graphName, endStationId)

		log.Printf("Search 3 hop path:")
    return benchmarkQuery(db, query)
}

func highestDegree(db driver.Database) time.Duration {
    query := fmt.Sprintf(`
        FOR station IN stations
        LET degree = LENGTH(
            FOR v IN 1..1 ANY station GRAPH '%s'
            RETURN v
        )
        SORT degree DESC
        LIMIT 1
        RETURN {
            station: station.name,
            degree: degree
        }`, graphName)

		log.Printf("Highest degree:")
    return benchmarkQuery(db, query)
}

func runAllGraphQueries(db driver.Database) {
    queries := []shared.QueryDefinition{
        {Name: "all_neighbours", Function: func() time.Duration { return getAllNeighbours(db) }},
        {Name: "three_hop_path", Function: func() time.Duration { return search3HopPath(db) }},
        {Name: "highest_degree", Function: func() time.Duration { return highestDegree(db) }},
    }
    shared.RunQueries("arangodb", "graph", queries)
}