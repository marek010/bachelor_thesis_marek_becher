import csv
import json
import hashlib

def get_hash_id(station_id):
    return int(hashlib.md5(station_id.encode()).hexdigest(), 16) % 28147497671000

def create_edges_csv():
    with open('data/dataset/graph_edges.json') as infile:
        edges = json.load(infile)

    with open('data/arango_graph_edges.csv', 'w', newline='') as outfile:
        writer = csv.DictWriter(outfile, fieldnames=['_from', '_to'])
        writer.writeheader()

        for edge in edges:
            row = {
                '_from': f"stations/{get_hash_id(edge['from'])}",
                '_to': f"stations/{get_hash_id(edge['to'])}"
            }
            writer.writerow(row)

    print("Arango Edges CSV created successfully!")

if __name__ == '__main__':
    create_edges_csv()