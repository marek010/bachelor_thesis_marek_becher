import csv
import json
import hashlib

def get_hash_id(string_id):
    return int(hashlib.md5(string_id.encode()).hexdigest(), 16) % 28147497671000

input_file = 'data/dataset/graph_nodes.json'
output_file = 'data/ts_data.csv'
fieldnames = ['timestamp', 'station_id', 'value']

def create_timeseries_data():
    with open(input_file, 'r') as file:
        data = json.load(file)

    with open(output_file, 'w', newline='') as csvfile:
        writer = csv.DictWriter(csvfile, fieldnames=fieldnames)
        writer.writeheader()

        for station in data:
            station_id = get_hash_id(station['station_id'])
            
            if 'num_bikes_available' in station['ts']:
                for entry in station['ts']['num_bikes_available']:
                    row = {
                        'timestamp': entry['Start'],
                        'station_id': station_id,
                        'value': entry['Value']
                    }
                    writer.writerow(row)

    print("Time Series Data CSV created successfully!")

if __name__ == '__main__':
    create_timeseries_data()