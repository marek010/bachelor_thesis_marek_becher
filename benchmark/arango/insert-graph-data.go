package arango

import (
	"benchmark/shared"
	"context"
	"log"
	"os/exec"

	"github.com/arangodb/go-driver"
)

const (
    vertexCollection = "stations"
    edgeCollection  = "edge_collection"
)

func insertGraphData(db driver.Database) {
		// Clear existing vertex and edge collections
    ctx := context.Background()
    
    vertices, err := db.Collection(ctx, vertexCollection)
    shared.HandleError(err, "Error getting vertex collection")
    
    edges, err := db.Collection(ctx, edgeCollection)
    shared.HandleError(err, "Error getting edge collection")

    vertices.Truncate(ctx)
    edges.Truncate(ctx)

		// Import vertices and edges
    importNodes := exec.Command("docker", "exec", "arangodb-container", "arangoimport",
        "--file", "/data/graph_nodes.csv",
        "--type", "csv",
        "--collection", vertexCollection,
        "--server.password", password,
        "--server.database", dbname,
        "--translate", "id=_key")

    importEdges := exec.Command("docker", "exec", "arangodb-container", "arangoimport",
        "--file", "/data/arango_graph_edges.csv",
        "--type", "csv",
        "--collection", edgeCollection,
        "--server.password", password,
        "--server.database", dbname)

    _, err = importNodes.CombinedOutput()
    shared.HandleError(err, "Error importing vertices")

    _, err = importEdges.CombinedOutput()
    shared.HandleError(err, "Error importing edges")

		// Create a graph with the imported data
    exists, err := db.GraphExists(ctx, graphName)
    shared.HandleError(err, "Error checking graph existence")
    
    if exists {
        log.Print("Graph already exists")
        return
    }

    opts := driver.CreateGraphOptions{
        EdgeDefinitions: []driver.EdgeDefinition{{
            Collection: edgeCollection,
            From:      []string{vertexCollection},
            To:        []string{vertexCollection},
        }},
    }

    _, err = db.CreateGraphV2(ctx, graphName, &opts)
    shared.HandleError(err, "Error creating graph")

    log.Print("Graph created successfully!")
}