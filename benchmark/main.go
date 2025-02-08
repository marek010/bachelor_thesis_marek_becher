package main

import (
	arango "benchmark/arango"
	postgres "benchmark/postgres"
	raven "benchmark/raven"
	"flag"
)

func main() {
	// Parse command line arguments to determine which database, datatype and method to use
	insertData := flag.Bool("insert", false, "Run data insertion")
	runQueries := flag.Bool("query", false, "Run queries")
	dataTypeFlag := flag.String("datatype", "", "Graph or Timeseries data") // "graph" or "timeseries"
	databaseFlag := flag.String("db", "", "Database to use") // "arango", "raven" or "postgres"
	flag.Parse()

	if *databaseFlag == "arango" {
		arango.ArangoMain(insertData, runQueries, dataTypeFlag)
	} else if *databaseFlag == "raven" {
		raven.RavenMain(insertData, runQueries, dataTypeFlag)
	} else if *databaseFlag == "postgres" {
		postgres.PostgresMain(insertData, runQueries, dataTypeFlag)
	}
}