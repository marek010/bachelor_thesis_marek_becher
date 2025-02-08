import csv
import json
import hashlib

def get_hash_id(station_id):
    return int(hashlib.md5(station_id.encode()).hexdigest(), 16) % 28147497671000

def create_nodes_csv():
    with open('data/dataset/graph_nodes.json', 'r') as file:
        data = json.load(file)

    fieldnames = [
        'id', 'name', 'lat', 'lon', 'region_id', 
        'capacity', 'station_id', 'short_name', 'start', 'end'
    ]

    with open('data/graph_nodes.csv', 'w', newline='') as csvfile:
        writer = csv.DictWriter(csvfile, fieldnames=fieldnames)
        writer.writeheader()

        for station in data:
            row = {
                'id': get_hash_id(station['nodeid']),
                'name': station.get('name', ''),
                'lat': station.get('lat', ''),
                'lon': station.get('lon', ''),
                'region_id': station.get('region_id', ''),
                'capacity': station.get('capacity', ''),
                'station_id': station.get('station_id', ''),
                'short_name': station.get('short_name', ''),
                'start': station.get('start', ''),
                'end': station.get('end', '')
            }
            writer.writerow(row)

    print("Nodes CSV created successfully!")

if __name__ == '__main__':
    create_nodes_csv()