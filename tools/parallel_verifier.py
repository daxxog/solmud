#!/usr/bin/env python3
"""
Parallel Verification Pipeline
Processes independent files simultaneously for faster completion
"""

import os
import json
import asyncio
import subprocess
from pathlib import Path
from concurrent.futures import ThreadPoolExecutor
from collections import defaultdict

class ParallelVerifier:
    def __init__(self, evidence_dir="bytecode/mapping/evidence/verified"):
        self.evidence_dir = Path(evidence_dir)
        self.max_workers = 4
        self.batch_size = 14
    
    def classify_files(self):
        """Classify files by complexity for parallel processing"""
        files = list(self.evidence_dir.glob("*.md"))
        
        workers = {
            'worker_1': {'complexity': 'network', 'files': []},
            'worker_2': {'complexity': 'objects', 'files': []},
            'worker_3': {'complexity': 'animation', 'files': []},
            'worker_4': {'complexity': 'utilities', 'files': []}
        }
        
        for file_path in files:
            filename = file_path.name
            
            if any(x in filename for x in ['RSSocket', 'OnDemand', 'RSInterface', 'client', 'WorldController']):
                workers['worker_1']['files'].append(filename)
            elif any(x in filename for x in ['Object', 'Model', 'Entity', 'NPC', 'Item']):
                workers['worker_2']['files'].append(filename)
            elif any(x in filename for x in ['Animation', 'Animable', 'SpotAnim']):
                workers['worker_3']['files'].append(filename)
            else:
                workers['worker_4']['files'].append(filename)
        
        return workers
    
    def estimate_processing_times(self, workers):
        """Estimate processing times for each worker"""
        time_per_file = {
            'network': 4.0,     # minutes per complex file
            'objects': 2.5,       # minutes per object file
            'animation': 3.0,     # minutes per animation file
            'utilities': 1.5       # minutes per utility file
        }
        
        for worker_id, worker_data in workers.items():
            complexity = worker_data['complexity']
            file_count = len(worker_data['files'])
            time_per = time_per_file[complexity]
            
            worker_data['estimated_time'] = file_count * time_per
            worker_data['time_per_file'] = time_per
        
        return workers
    
    def simulate_parallel_processing(self, workers):
        """Simulate parallel processing and calculate completion time"""
        # Find the worker with longest processing time (bottleneck)
        max_time = max(worker['estimated_time'] for worker in workers.values())
        total_files = sum(len(worker['files']) for worker in workers.values())
        
        # Sequential time for comparison
        sequential_time = sum(worker['estimated_time'] for worker in workers.values())
        
        # Calculate efficiency gains
        efficiency = (sequential_time / max_time) if max_time > 0 else 0
        
        return {
            'parallel_time': max_time,
            'sequential_time': sequential_time,
            'time_savings': sequential_time - max_time,
            'efficiency_gain': ((efficiency - 1) * 100) if efficiency > 1 else 0,
            'total_files': total_files,
            'workers_used': len(workers)
        }
    
    def generate_processing_plan(self, workers):
        """Generate detailed processing plan for each worker"""
        plan = {}
        
        for worker_id, worker_data in workers.items():
            plan[worker_id] = {
                'complexity': worker_data['complexity'],
                'file_count': len(worker_data['files']),
                'estimated_time': worker_data['estimated_time'],
                'files': worker_data['files'],
                'batches': self.create_batches(worker_data['files'], 3)  # 3 files per batch
            }
        
        return plan
    
    def create_batches(self, files, batch_size):
        """Create batches for a worker"""
        batches = []
        for i in range(0, len(files), batch_size):
            batch = files[i:i + batch_size]
            batches.append({
                'batch_id': len(batches) + 1,
                'files': batch,
                'size': len(batch)
            })
        return batches
    
    def generate_report(self, workers, processing_stats, processing_plan):
        """Generate comprehensive parallel processing report"""
        print("🚀 PARALLEL VERIFICATION PIPELINE REPORT")
        print("=" * 60)
        
        print(f"\n📊 WORKER DISTRIBUTION")
        print("-" * 30)
        total_files = 0
        for worker_id, worker_data in workers.items():
            file_count = len(worker_data['files'])
            total_files += file_count
            time_est = worker_data['estimated_time']
            print(f"{worker_id:10}: {file_count:3} files, {time_est:5.1f}min ({worker_data['complexity']})")
        
        print(f"\n{'Total':10}: {total_files:3} files")
        
        print(f"\n⏱️  PERFORMANCE ANALYSIS")
        print("-" * 35)
        
        parallel = processing_stats['parallel_time']
        sequential = processing_stats['sequential_time']
        savings = processing_stats['time_savings']
        efficiency = processing_stats['efficiency_gain']
        
        print(f"Sequential approach:    {sequential:>6.1f} minutes")
        print(f"Parallel approach:      {parallel:>6.1f} minutes")
        print(f"Time savings:           {savings:>6.1f} minutes ({efficiency:+.1f}%)")
        
        print(f"\n📋 DETAILED PROCESSING PLAN")
        print("-" * 40)
        
        for worker_id, plan_data in processing_plan.items():
            print(f"\n{worker_id.upper()} ({plan_data['complexity'].title()}):")
            for batch in plan_data['batches']:
                files_preview = ', '.join([f.replace('.md', '') for f in batch['files'][:2]])
                if batch['size'] > 2:
                    files_preview += f" (+{batch['size']-2} more)"
                
                print(f"  Batch {batch['batch_id']}: {batch['size']} files - {files_preview}")
        
        return processing_plan
    
    def export_parallel_plan(self, processing_plan, filename="parallel_plan.json"):
        """Export parallel processing plan"""
        plan_data = {
            'created_at': str(Path.cwd()),
            'total_workers': len(processing_plan),
            'total_files': sum(plan['file_count'] for plan in processing_plan.values()),
            'estimated_parallel_time': max(plan['estimated_time'] for plan in processing_plan.values()),
            'workers': processing_plan
        }
        
        with open(filename, 'w') as f:
            json.dump(plan_data, f, indent=2)
        
        print(f"\n💾 Parallel plan exported to: {filename}")
        return plan_data

def main():
    import argparse
    
    parser = argparse.ArgumentParser(description='Parallel Verification Pipeline')
    parser.add_argument('--directory', default='bytecode/mapping/evidence/verified',
                       help='Evidence directory path')
    parser.add_argument('--workers', type=int, default=4, help='Number of parallel workers')
    parser.add_argument('--export', action='store_true', help='Export parallel plan to JSON')
    
    args = parser.parse_args()
    
    verifier = ParallelVerifier(args.directory)
    verifier.max_workers = args.workers
    
    # Classify files by complexity
    workers = verifier.classify_files()
    
    # Estimate processing times
    workers = verifier.estimate_processing_times(workers)
    
    # Simulate parallel processing
    processing_stats = verifier.simulate_parallel_processing(workers)
    
    # Generate processing plan
    processing_plan = verifier.generate_processing_plan(workers)
    
    # Generate report
    verifier.generate_report(workers, processing_stats, processing_plan)
    
    # Export if requested
    if args.export:
        verifier.export_parallel_plan(processing_plan)
    
    return 0

if __name__ == "__main__":
    exit(main())