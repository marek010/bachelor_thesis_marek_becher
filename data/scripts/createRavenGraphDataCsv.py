import csv
import json
import hashlib

def get_hash_id(string_id):
    return int(hashlib.md5(string_id.encode()).hexdigest(), 16) % 28147497671000

# For RavenDB the data is in form of documents which represent stations and these have an array of edges to other stations.
def create_raven_graph_data():
    edges_file = 'data/dataset/graph_edges.json'
    nodes_file = 'data/dataset/graph_nodes.json'
    output_file = 'data/raven_graph_data.csv'

    with open(edges_file, 'r') as edges_file:
        edges_data = json.load(edges_file)

    with open(nodes_file, 'r') as nodes_file:
        nodes_data = json.load(nodes_file)

    edges_dict = {}
    for edge in edges_data:
        from_id = get_hash_id(edge['from'])
        to_id = get_hash_id(edge['to'])
        if to_id not in edges_dict:
            edges_dict[to_id] = []
        edges_dict[to_id].append('Stations/' + str(from_id))

    fieldnames = [
      '@id', 'name', 'lat', 'lon', 'region_id',
      'capacity', 'short_name', 'start', 'end', 'edges'
    ]

    with open(output_file, 'w', newline='') as csvfile:
        writer = csv.DictWriter(csvfile, fieldnames=fieldnames)
        writer.writeheader()

        for node in nodes_data:
            node_id = get_hash_id(node['nodeid'])
            row = {
                '@id': 'Stations/' + str(node_id),
                'name': node.get('name', ''),
                'lat': node.get('lat', ''),
                'lon': node.get('lon', ''),
                'region_id': node.get('region_id', ''),
                'capacity': node.get('capacity', ''),
                'short_name': node.get('short_name', ''),
                'start': node.get('start', ''),
                'end': node.get('end', ''),
                'edges': json.dumps(edges_dict.get(node_id, []))
            }
            writer.writerow(row)

    print("Raven Graph Data CSV created successfully!")

if __name__ == '__main__':
    create_raven_graph_data()