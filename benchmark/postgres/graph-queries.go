package postgres

import (
	"benchmark/shared"
	"database/sql"
	"fmt"
	"log"
	"time"
)

const (
	startStationId = "5230109759983"
	endStationId   = "32116902472201"
)

func getAllNeighbours(db *sql.DB) time.Duration {
	query := fmt.Sprintf(`
			SELECT *
			FROM cypher('graph', $$
					MATCH (a:Station {id: '%s'})-->(b:Station)
					RETURN b.id
			$$) AS (neighbor_id TEXT);`,
			startStationId)

	log.Printf("Get neighbours:")
	return benchmarkQuery(db, query)
}

func search3HopPath(db *sql.DB) time.Duration {
	query := fmt.Sprintf(`
			SELECT *
			FROM cypher('graph', $$
					MATCH (n:Station {id: '%s'})
					WHERE exists((n)-[*1..3]->(:Station {id: '%s'}))
					RETURN n.name
			$$) as (name agtype);`,
			startStationId, endStationId)

	log.Printf("Search 3 hop path:")
	return benchmarkQuery(db, query)
}

func highestDegree(db *sql.DB) time.Duration {
    query := `
        SELECT result.station_name, result.degree
        FROM cypher('graph', $$
            MATCH (station:Station)-[:has_bikeride*1..1]-(:Station)
            RETURN station.name AS station_name, COUNT(*) AS degree
            ORDER BY COUNT(*) DESC
            LIMIT 1
        $$) AS result(station_name agtype, degree agtype);`

    log.Printf("Highest degree:")
    return benchmarkQuery(db, query)
}

func runAllGraphQueries(db *sql.DB) {
    queries := []shared.QueryDefinition{
        {Name: "all_neighbours", Function: func() time.Duration { return getAllNeighbours(db) }},
        {Name: "three_hop_path", Function: func() time.Duration { return search3HopPath(db) }},
        {Name: "highest_degree", Function: func() time.Duration { return highestDegree(db) }},
    }

    shared.RunQueries("postgresql", "graph", queries)
}