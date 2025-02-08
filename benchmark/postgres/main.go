package postgres

import (
	"benchmark/shared"
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

// Database and Benchmark Configuration
const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "root"
	dbname   = "benchmark_db"
)

func PostgresMain(insertData *bool, runQueries *bool, dataType *string) {
	// Connect to the database
	db, err := connectToDb()
	shared.HandleError(err, "Error opening database")
	defer db.Close()

	// Run the insertion or query function based on the command line arguments
	if *insertData {
			switch *dataType {
			case "graph":
					insertGraphData(db)
			case "timeseries":
					insertTSData(db)
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

func connectToDb() (*sql.DB, error) {
	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
	host, port, user, password, dbname)
	return sql.Open("postgres", connStr)
}
