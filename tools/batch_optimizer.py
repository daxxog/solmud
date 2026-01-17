#!/usr/bin/env python3
"""
Dynamic Batch Optimizer
Analyzes class complexity and creates optimal batch sizes for efficient processing
"""

import os
import re
import json
from pathlib import Path
from collections import defaultdict

class BatchOptimizer:
    def __init__(self, evidence_dir="bytecode/mapping/evidence/verified"):
        self.evidence_dir = Path(evidence_dir)
        
        # Complexity classification based on class naming patterns
        self.complexity_patterns = {
            'complex': {
                'patterns': [r'RSSocket', r'OnDemand', r'RSInterface', r'client', r'WorldController'],
                'batch_size': 2,
                'estimated_time': 45  # minutes per batch
            },
            'medium': {
                'patterns': [r'Object\d', r'Animation', r'Model', r'Entity', r'Class\d\d', r'Animable'],
                'batch_size': 4,
                'estimated_time': 30  # minutes per batch
            },
            'simple': {
                'patterns': [r'SizeConstants', r'TextInput', r'Node', r'Stream', r'Decompressor'],
                'batch_size': 6,
                'estimated_time': 20  # minutes per batch
            }
        }
    
    def classify_file(self, filename):
        """Classify a file by complexity based on naming patterns"""
        for complexity, config in self.complexity_patterns.items():
            for pattern in config['patterns']:
                if re.search(pattern, filename):
                    return complexity
        return 'medium'  # Default to medium
    
    def analyze_files(self):
        """Analyze all files and classify by complexity"""
        if not self.evidence_dir.exists():
            print(f"Error: Directory {self.evidence_dir} not found")
            return None
        
        files = list(self.evidence_dir.glob("*.md"))
        classification = {'complex': [], 'medium': [], 'simple': []}
        
        for file_path in files:
            complexity = self.classify_file(file_path.name)
            classification[complexity].append(file_path.name)
        
        return classification
    
    def create_batches(self, classification):
        """Create optimal batches based on complexity classification"""
        batches = []
        
        for complexity, files in classification.items():
            if not files:
                continue
                
            batch_size = self.complexity_patterns[complexity]['batch_size']
            
            for i in range(0, len(files), batch_size):
                batch_files = files[i:i + batch_size]
                batch = {
                    'complexity': complexity,
                    'files': batch_files,
                    'size': len(batch_files),
                    'estimated_time': self.complexity_patterns[complexity]['estimated_time'],
                    'batch_id': len(batches) + 1
                }
                batches.append(batch)
        
        return batches
    
    def calculate_optimizations(self, current_batches, optimized_batches):
        """Calculate time savings from optimization"""
        current_time = sum(batch['estimated_time'] for batch in current_batches)
        optimized_time = sum(batch['estimated_time'] for batch in optimized_batches)
        
        return {
            'current_time': current_time,
            'optimized_time': optimized_time,
            'time_savings': current_time - optimized_time,
            'improvement_percentage': ((current_time - optimized_time) / current_time * 100) if current_time > 0 else 0
        }
    
    def generate_report(self, classification, batches, optimizations):
        """Generate comprehensive optimization report"""
        print("🔍 DYNAMIC BATCH OPTIMIZATION REPORT")
        print("=" * 60)
        
        print("\n📊 CLASSIFICATION SUMMARY")
        print("-" * 30)
        total_files = sum(len(files) for files in classification.values())
        for complexity, files in classification.items():
            percentage = (len(files) / total_files) * 100 if total_files > 0 else 0
            print(f"{complexity.capitalize():8}: {len(files):3} files ({percentage:5.1f}%)")
        
        print(f"\n{'Total':8}: {total_files} files (100.0%)")
        
        print("\n📋 OPTIMIZED BATCH DISTRIBUTION")
        print("-" * 40)
        
        complexity_counts = defaultdict(int)
        for batch in batches:
            complexity_counts[batch['complexity']] += 1
            
        for complexity, count in complexity_counts.items():
            batch_size = self.complexity_patterns[complexity]['batch_size']
            print(f"{complexity.capitalize():8}: {count:2} batches × {batch_size} files/batch")
        
        print(f"\nTotal: {len(batches)} batches")
        
        print("\n⏱️  TIME OPTIMIZATION ANALYSIS")
        print("-" * 35)
        
        if optimizations:
            current = optimizations['current_time']
            optimized = optimizations['optimized_time']
            savings = optimizations['time_savings']
            improvement = optimizations['improvement_percentage']
            
            print(f"Current approach:     {current:>3.0f} minutes")
            print(f"Optimized approach:   {optimized:>3.0f} minutes")
            print(f"Time savings:          {savings:>3.0f} minutes ({improvement:.1f}%)")
            
        print("\n📝 DETAILED BATCH BREAKDOWN")
        print("-" * 40)
        
        for batch in batches:
            files_str = ", ".join([f.replace('.md', '') for f in batch['files'][:2]])
            if len(batch['files']) > 2:
                files_str += f" (+{len(batch['files'])-2} more)"
            
            print(f"Batch {batch['batch_id']:2d} ({batch['complexity']:7}): {batch['size']} files, {batch['estimated_time']:2d}min")
            print(f"         {files_str}")
        
        return batches
    
    def export_batch_plan(self, batches, output_file="batch_plan.json"):
        """Export batch plan for processing automation"""
        plan = {
            'created_at': str(Path.cwd()),
            'total_batches': len(batches),
            'total_files': sum(batch['size'] for batch in batches),
            'estimated_total_time': sum(batch['estimated_time'] for batch in batches),
            'batches': batches
        }
        
        with open(output_file, 'w') as f:
            json.dump(plan, f, indent=2)
        
        print(f"\n💾 Batch plan exported to: {output_file}")
        return plan

def main():
    import argparse
    
    parser = argparse.ArgumentParser(description='Dynamic Batch Optimizer')
    parser.add_argument('--directory', default='bytecode/mapping/evidence/verified', 
                       help='Evidence directory path')
    parser.add_argument('--min-batch-size', type=int, default=2, help='Minimum batch size')
    parser.add_argument('--max-batch-size', type=int, default=6, help='Maximum batch size')
    parser.add_argument('--export', action='store_true', help='Export batch plan to JSON')
    
    args = parser.parse_args()
    
    optimizer = BatchOptimizer(args.directory)
    
    # Analyze and classify files
    classification = optimizer.analyze_files()
    if not classification:
        return 1
    
    # Create optimized batches
    batches = optimizer.create_batches(classification)
    
    # Calculate optimizations (vs current uniform batching)
    current_batches = [{'estimated_time': 35}] * ((sum(len(files) for files in classification.values()) // 4) + 1)
    optimizations = optimizer.calculate_optimizations(current_batches, batches)
    
    # Generate report
    optimizer.generate_report(classification, batches, optimizations)
    
    # Export if requested
    if args.export:
        optimizer.export_batch_plan(batches)
    
    return 0

if __name__ == "__main__":
    exit(main())