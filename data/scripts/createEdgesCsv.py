import csv
import json
import hashlib

def get_hash_id(station_id: str) -> int:
    return int(hashlib.md5(station_id.encode()).hexdigest(), 16) % 28147497671000

def create_edges_csv():
    with open('data/dataset/graph_edges.json') as infile:
        edges = json.load(infile)

    with open('data/graph_edges.csv', 'w', newline='') as outfile:
        writer = csv.DictWriter(outfile, fieldnames=[
            'start_id', 
            'start_vertex_type', 
            'end_id', 
            'end_vertex_type'
        ])
        writer.writeheader()

        for edge in edges:
            row = {
                'start_id': get_hash_id(edge['from']),
                'start_vertex_type': 'Station',
                'end_id': get_hash_id(edge['to']),
                'end_vertex_type': 'Station'
            }
            writer.writerow(row)

    print("Edges CSV created successfully!")

if __name__ == '__main__':
    create_edges_csv()