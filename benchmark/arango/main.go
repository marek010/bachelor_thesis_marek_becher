package arango

import (
	"benchmark/shared"
	"context"
	"fmt"

	"github.com/arangodb/go-driver"
	"github.com/arangodb/go-driver/http"
)

const (
    host       = "localhost"
    port       = 8529
    user       = "root"
    password   = "root"
    dbname     = "benchmark_db"
    collection = "ts_table"
    graphName  = "benchmark_graph"
)

func ArangoMain(insertData *bool, runQueries *bool, dataType *string) {
    db, col := connectToDb()

		// Run the insertion or query function based on the command line arguments
		if *insertData {
			switch *dataType {
			case "graph":
					insertGraphData(db)
			case "timeseries":
					insertTSData(col)
			}
		} else if *runQueries {
				switch *dataType {
				case "graph":
						runAllGraphQueries(db)
				case "timeseries":
						runAllTimeseriesQueries(db)
				}
		}
}

func connectToDb() (driver.Database, driver.Collection) {
    conn, err := http.NewConnection(http.ConnectionConfig{
        Endpoints: []string{fmt.Sprintf("http://%s:%d", host, port)},
    })
    shared.HandleError(err, "Error creating connection")

    client, err := driver.NewClient(driver.ClientConfig{
        Connection:     conn,
        Authentication: driver.BasicAuthentication(user, password),
    })
    shared.HandleError(err, "Error creating client")

    ctx := context.Background()
    db, err := client.Database(ctx, dbname)
    shared.HandleError(err, "Error connecting to database")

    col, err := db.Collection(ctx, collection)
    shared.HandleError(err, "Error accessing collection")

    return db, col
}