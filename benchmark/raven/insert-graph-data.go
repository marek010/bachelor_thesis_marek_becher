package raven

import (
	"benchmark/shared"
	"encoding/csv"
	"encoding/json"
	"log"
	"os"

	ravendb "github.com/ravendb/ravendb-go-client"
)

func NewStationsDegreeIndex() *ravendb.IndexCreationTask {
    index := ravendb.NewIndexCreationTask("StationsDegreeIndex")
    index.Map = `
        map("Stations", function (station) {
            var results = [];
            var degree = station.edges ? station.edges.length : 0;
            
            results.push({
                StationId: id(station),
                Degree: degree
            });

            if (station.edges) {
                station.edges.forEach(edgeId => {
                    if (edgeId && edgeId !== id(station)) {
                        results.push({
                            StationId: edgeId,
                            Degree: 1
                        });
                    }
                });
            }
            return results;
        });`

    index.Reduce = `
        reduce(groupBy(x => x.StationId)
            .aggregate(g => ({
                StationId: g.key,
                Degree: g.values.reduce((sum, val) => sum + val.Degree, 0)
            }))
        )`

    return index
}

func insertGraphData(session *ravendb.DocumentSession, store *ravendb.DocumentStore) {
    // Open csv
    csvPath := "../data/raven_graph_data.csv"
    file, err := os.Open(csvPath)
    shared.HandleError(err, "Error opening CSV file")
    defer file.Close()

    reader := csv.NewReader(file)
    headers, err := reader.Read()
    shared.HandleError(err, "Error reading CSV header")

    records, err := reader.ReadAll()
    shared.HandleError(err, "Error opening CSV records")

    for _, record := range records {
        doc := &map[string]interface{}{}
        *doc = make(map[string]interface{})
        
        var documentId string
        
        for i, value := range record {
            if headers[i] == "@id" {
                documentId = value
                continue
            }
            if headers[i] == "edges" {
                var edges []string
                err := json.Unmarshal([]byte(value), &edges)
                shared.HandleError(err, "Error parsing edges")
                (*doc)[headers[i]] = edges
            } else {
                (*doc)[headers[i]] = value
            }
        }

				// Store the documents
        err = session.StoreWithID(doc, documentId)
        shared.HandleError(err, "Error storing entry")

        metadata, err := session.Advanced().GetMetadataFor(doc)
        shared.HandleError(err, "Error getting metadata")
        metadata.Put("@collection", "Stations")
    }

    session.SaveChanges()
    log.Printf("Data import successful!")

		// Create index
    log.Print("Creating Indexes...")
    index := NewStationsDegreeIndex()
    err = index.Execute(store, nil, dbname)
    shared.HandleError(err, "Error creating index")
    log.Printf("Success!")
}