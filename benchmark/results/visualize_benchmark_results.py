import json
import matplotlib.pyplot as plt
import pandas as pd
import numpy as np

databases = ['postgresql', 'arangodb', 'ravendb']
colors = ['#4C78A8', '#00B050', '#D6277F']
results_folder_path = 'benchmark/results'
timeseries_queries = ['simple_range', 'range_with_filter', 'latest_timestamp', 'aggregation', 'downsampling']
graph_queries = ['all_neighbours', 'three_hop_path', 'highest_degree']

def setup_chart(width, height, log_scale=False):
    fig, ax = plt.subplots(figsize=(width, height))
    if log_scale:
        ax.set_yscale('log')
        ax.yaxis.set_major_formatter(plt.FuncFormatter(lambda y, _: f'{int(y)}'))
    return fig, ax

def format_chart(ax, title, ylabel):
    ax.set_title(title, fontsize=14, fontweight='bold')
    ax.set_ylabel(ylabel, fontsize=14)
    ax.grid(axis='y', linestyle='solid', alpha=0.5)
    ax.tick_params(axis='x', labelsize=14)

# Value numbers on top of bars
def add_value_labels(ax, bars):
    for bar in bars:
        ax.text(bar.get_x() + bar.get_width()/2, bar.get_height(), f'{bar.get_height():.1f}', ha='center', va='bottom', fontsize=12)

def create_x_axis_labels(queries, base):
    return [f'{base + (i+1)/12:.1f} {q.replace("_", " ").title()}' for i, q in enumerate(queries)]

def create_result_visualizations(filename):
    df = pd.DataFrame(json.load(open(filename))['results'])
    width = 0.25
    
    # Time-Series Query Chart
    ts_subset = df[(df['queryType'] == 'timeseries') & (df['queryName'] != 'insert')]
    fig, ax = setup_chart(14, 8, log_scale=True)
    x = np.arange(len(timeseries_queries))
    
    for i, db in enumerate(databases):
        db_data = ts_subset[ts_subset['database'] == db]
        heights = [db_data.loc[db_data['queryName'] == q, 'averageTime'].iloc[0] for q in timeseries_queries]
        bars = ax.bar(x + i*width, heights, width, label=db, color=colors[i])
        add_value_labels(ax, bars)
    
    ax.set_xticks(x + width)
    ax.set_xticklabels(create_x_axis_labels(timeseries_queries, 1.0), ha='center')
    ax.legend(loc='upper left', fontsize=14)
    format_chart(ax, 'Time-Series Queries', 'Average Query Latency in ms (Log Scale)')
    plt.tight_layout()
    plt.savefig(f'{results_folder_path}/results_images/timeseries_queries_comparison.png', dpi=300)
    plt.close()
    
    # Graph Query Chart
    graph_subset = df[(df['queryType'] == 'graph') & (df['queryName'] != 'insert')]
    fig, ax = setup_chart(14, 8, log_scale=True)
    x = np.arange(len(graph_queries))
    
    for i, db in enumerate(databases):
        db_data = graph_subset[graph_subset['database'] == db]
        heights = [db_data.loc[db_data['queryName'] == q, 'averageTime'].iloc[0] for q in graph_queries]
        bars = ax.bar(x + i*width, heights, width, label=db, color=colors[i])
        add_value_labels(ax, bars)
    
    ax.set_xticks(x + width)
    ax.set_xticklabels(create_x_axis_labels(graph_queries, 2.0), ha='center')
    ax.legend(loc='upper left', fontsize=14)
    format_chart(ax, 'Graph Queries', 'Average Query Latency in ms (Log Scale)')
    plt.tight_layout()
    plt.savefig(f'{results_folder_path}/results_images/graph_queries_comparison.png', dpi=300)
    plt.close()
    
    # Data Insertion Chart
    insert_df = df[df['queryName'] == 'timeseries_insertion']
    fig, ax = setup_chart(8, 5.5)
    heights = [insert_df[insert_df['database'] == db]['averageTime'].iloc[0] for db in databases]
    bars = ax.bar(databases, heights, color=colors)

    add_value_labels(ax, bars)
    format_chart(ax, 'Data Insertion Performance', 'Insertion Time (ms)')
    ax.set_xlabel('Database', fontsize=12)
    plt.tight_layout()

    plt.savefig(f'{results_folder_path}/results_images/insertion_comparison.png', dpi=300)
    plt.close()

if __name__ == '__main__':
    create_result_visualizations(f'{results_folder_path}/benchmark_results.json')