# Benchmark Instructions

Welcome to the Timeseries + Graph benchmarks for PostgreSQL, ArangoDB, and RavenDB. Follow the steps below to set up your environment, prepare the data, execute benchmarks, and visualize the results. The benchmark includes 5 Timeseries Queries and 3 Graph Queries. Note that all commands should be executed from the base directory unless stated otherwise.

## Prerequisites

Ensure you have the following installed:

- **Docker**
- **Go** (version 1.23.3)
- **Python 3**
  
## Data Preparation

### 1. Download the Dataset

Download the dataset from [Zenodo](https://zenodo.org/records/13846868) and save the files `graph_edges.json` and `graph_nodes.json` inside the `data/dataset` folder (this folder might need to be created first).


### 2. Run Preparation Scripts

Generate the necessary CSV files for the databases:

```bash
# Create CSV for timeseries benchmark
python3 data/scripts/createTimeSeriesDataCsv.py

# Create CSV files for nodes and edges (PostgreSQL and ArangoDB)
python3 data/scripts/createArangoEdgesCsv.py
python3 data/scripts/createEdgesCsv.py
python3 data/scripts/createNodesCsv.py

# Create CSV file for graph data (RavenDB)
python3 data/scripts/createRavenGraphDataCsv.py
```

---

## Database Setup

### PostgreSQL (AGE + Timescale)

#### 1. Build the Docker Image (it includes Postgres16 and the Extensions AGE and Timescale)

```bash
docker build -f other/postgres.Dockerfile -t postgres-timescale-age .
```

#### 2. Run the Docker Container

```bash
docker run --name postgres-timescale-age-container \
    -e POSTGRES_PASSWORD=root \
    -d \
    -p 5432:5432 \
    -v ./data:/data \
    postgres-timescale-age
```

#### 3. Enable Extensions

```bash
docker exec -it postgres-timescale-age-container psql -U postgres
```

Inside the PostgreSQL shell:

```sql
ALTER SYSTEM SET shared_preload_libraries = 'timescaledb';
```

Restart the container:

```bash
docker restart postgres-timescale-age-container
```

Reconnect to the database:

```bash
docker exec -it postgres-timescale-age-container psql -U postgres
```

Inside the PostgreSQL shell:

```sql
CREATE DATABASE benchmark_db;
```

Exit and reconnect to the new database:

```bash
\c benchmark_db
```

Enable the required extensions:

```sql
CREATE EXTENSION IF NOT EXISTS timescaledb;
CREATE EXTENSION IF NOT EXISTS age;
```

---

### ArangoDB

#### 1. Download and Run the Docker Container

```bash
docker pull arangodb:3.12.3
docker run -d -p 8529:8529 -e ARANGO_ROOT_PASSWORD=root --name arangodb-container -v ./data:/data arangodb:3.12.3
```

#### 2. Create Database and Collections

```bash
docker exec -it arangodb-container arangosh --server.endpoint tcp://localhost:8529 --server.password root
```

Inside the shell:

```javascript
db._createDatabase("benchmark_db");
db._useDatabase("benchmark_db");
db._create("ts_table");
db._create("stations");
db._createEdgeCollection("edge_collection");
```

---

### RavenDB

1. Download and install [RavenDB](https://ravendb.net/download) (or download the [docker image](https://hub.docker.com/r/ravendb/ravendb/)) **Version 5.4**
2. Start the database or Docker container and follow the setup wizard selecting all default options but **enabling experimental features**.
3. In the Database [Web Interface](http://127.0.0.1:8080) create a new database named `benchmark_db` (leave all other settings as default).

---

## Running Benchmarks

Navigate to the `benchmark` directory:

```bash
cd benchmark
```

Replace `{database}` with either `postgres`, `raven`, or `arango` to run the queries in the corresponding database.

### Insert Timeseries Data

```bash
go run . -insert -datatype timeseries -db {database}
```

### Run Timeseries Queries

```bash
go run . -query -datatype timeseries -db {database}
```

### Insert Graph Data

```bash
go run . -insert -datatype graph -db {database}
```

### Run Graph Queries

```bash
go run . -query -datatype graph -db {database}
```

---

## Visualizing Benchmark Results

After running an insertion or query, the results are stored in a JSON file in **`benchmark/results`**. Create a folder named **`results_images`** inside the folder **`results`** After running all insertions and queries for all databases you can visualize the results in bar-charts by running the following script:

```bash
python3 benchmark/results/visualize_benchmark_results.py
```

This script will generate bar charts as PNG files comparing the benchmark results. They will be stored in the folder **`benchmark/results/results_images`**

---

## Additional Remarks

- The database setup and benchmark scripts were tested on macOS. Minor adjustments may be required for other operating systems.
- The storage size of the time- eries data needs to be determined manually in the arangoDB and RavenDB web interfaces, for PostgreSQL the result is printed after the insertion

