import json
import matplotlib.pyplot as plt
import pandas as pd

databases = ['postgresql', 'arangodb', 'ravendb']
colors = ['#4C78A8', '#00B050', '#D6277F']
results_path = 'benchmark/results'

def visualize_benchmarks(filename):
    df = pd.DataFrame(json.load(open(filename))['results'])
    
    for query_type in ['timeseries', 'graph']:
        type_df = df[df['queryType'] == query_type]
        
        for query_name in type_df['queryName'].unique():
            fig, ax = plt.subplots(figsize=(7, 5))
            
            # Order the bars according to the databases list
            data = (type_df[type_df['queryName'] == query_name]
                   .set_index('database')
                   .reindex(databases)
                   .reset_index())
            
            bars = ax.bar(data['database'], 
                         data['averageTime'], 
                         color=colors)
            
            ax.set_title(query_name.replace('_', ' ').title())
            ax.set_xlabel('Database')
            
            if query_name == "timeseries_insertion":
                ax.set_ylabel('Insertion Time (ms)')
            else:
                ax.set_ylabel('Average Query Time (ms)')

            ax.grid(axis='y', linestyle='solid', alpha=0.5)
            
            for bar, time in zip(bars, data['averageTime']):
                ax.text(bar.get_x() + bar.get_width()/2, 
                       time,
                       f'{time:.1f}', 
                       ha='center', 
                       va='bottom')
            
            plt.tight_layout()
            plt.savefig(f'{results_path}/results_images/{query_type}_{query_name}_comparison.png')
            plt.close()

if __name__ == '__main__':
    visualize_benchmarks(f'{results_path}/benchmark_results.json')