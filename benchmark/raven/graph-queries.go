package raven

import (
	"benchmark/shared"
	"fmt"
	"log"
	"time"

	ravendb "github.com/ravendb/ravendb-go-client"
)

const (
	startStationId = "Stations/5230109759983"
	endStationId   = "Stations/32116902472201"
)

func getAllNeighbours(session *ravendb.DocumentSession) time.Duration {
	query := fmt.Sprintf(`
			match (Stations where id() = "%s") - [edges] -> (Stations as neighbour)
			select neighbour.name
	`, startStationId)

	log.Printf("Get neighbours:")
	return benchmarkQuery(session, query)
}

func search3HopPath(session *ravendb.DocumentSession) time.Duration {
	query := fmt.Sprintf(`
		match (Stations as a where id() = "%s")-
			recursive(1,3,all) {
			[edges]->(Stations as b)
		}
    where b.id = "%s"
    select b
	`, startStationId, endStationId)

	log.Printf("Search 3 hop path:")
	return benchmarkQuery(session, query)
}

func highestDegree(session *ravendb.DocumentSession) time.Duration {
	query := `
		from index 'StationsDegreeIndex'
		order by Degree as double desc
		select StationId, Degree
		limit 1
	`

	log.Printf("Highest degree:")
	return benchmarkQuery(session, query)
}

func runAllGraphQueries(session *ravendb.DocumentSession){
	queries := []shared.QueryDefinition{
		{Name: "all_neighbours", Function: func() time.Duration { return getAllNeighbours(session) }},
		{Name: "three_hop_path", Function: func() time.Duration { return search3HopPath(session) }},
		{Name: "highest_degree", Function: func() time.Duration { return highestDegree(session) }},
	}

	shared.RunQueries("ravendb", "graph", queries)
}