import csv
import logging
from datetime import datetime, timedelta
from ravendb import DocumentStore

logging.basicConfig(level=logging.INFO)
logger = logging.getLogger(__name__)

url = "http://127.0.0.1:8080"
dbname = "benchmark_db"
time_series_name = "Timeseries"

store = DocumentStore(urls=[url], database=dbname)
store.initialize()

with store.open_session() as session:
    doc_id = "TSCollection/1854354353"
    empty_doc = {"@id": doc_id}
    session.store(empty_doc, doc_id)
    session.save_changes()

csv_file_path = "../data/ts_data.csv"
entries = []
timestamp_counts = {}

# Read CSV data
with open(csv_file_path, mode='r') as file:
    csv_reader = csv.reader(file)
    next(csv_reader)
    for row in csv_reader:
        timestamp_str, station_id, value = row
        timestamp = datetime.fromisoformat(timestamp_str)
        value = float(value)

        if timestamp in timestamp_counts:
            timestamp_counts[timestamp] += 1
            timestamp += timedelta(milliseconds=timestamp_counts[timestamp])
        else:
            timestamp_counts[timestamp] = 0
        
        entries.append((timestamp, [station_id, value]))

# Bulk insert
try:
    with store.bulk_insert(dbname) as bulk_insert:
        with bulk_insert.time_series_for(doc_id, time_series_name) as ts_bulk_insert:
            for entry in entries:
                ts_bulk_insert.append(*entry)
except Exception as e:
    print(f"An error occurred: {e}")