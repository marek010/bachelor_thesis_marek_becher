package postgres

import (
	"database/sql"
	"log"
)


func insertGraphData(db *sql.DB) {
	log.Print("Inserting graph data...\n")

	// Load graph extension "AGE" and creaet a graph
	db.Exec("LOAD 'age';")
	db.Exec("SET search_path = ag_catalog, '$user', public;")
	db.Exec("SELECT drop_graph('graph', cascade);")
	db.Exec("SELECT create_graph('graph');")
	db.Exec("SELECT create_vlabel('graph','Station');")

	// Load edges and nodes from CSV files into the graph
	db.Exec("SELECT load_labels_from_file('graph','Station','/data/graph_nodes.csv');")
	db.Exec("SELECT create_elabel('graph','has_bikeride');")
	db.Exec("SELECT load_edges_from_file('graph', 'has_bikeride','/data/graph_edges.csv');")

	log.Print("Success!")
}	