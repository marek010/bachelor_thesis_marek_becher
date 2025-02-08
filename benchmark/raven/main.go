package raven

import (
	"benchmark/shared"

	ravendb "github.com/ravendb/ravendb-go-client"
)

const (
	url      = "http://127.0.0.1:8080"
	dbname   = "benchmark_db"
)

func RavenMain(insertData *bool, runQueries *bool, dataType *string) {
	store := connectToDb()
	defer store.Close()

	session, err := store.OpenSession("")
	shared.HandleError(err, "Error opening session")
	defer session.Close()
	session.Advanced().SetMaxNumberOfRequestsPerSession(60)

	// Run the insertion or query function based on the command line arguments
	if *insertData {
		switch *dataType {
		case "graph":
				insertGraphData(session, store)
		case "timeseries":
				insertTSData(store)
		}
	} else if *runQueries {
			switch *dataType {
			case "graph":
					runAllGraphQueries(session)
			case "timeseries":
					runAllTimeseriesQueries(session)
			}
	}
}

func connectToDb() *ravendb.DocumentStore {
	store := ravendb.NewDocumentStore([]string{url}, dbname)
	store.Initialize()
	return store
}