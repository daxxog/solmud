#!/usr/bin/env python3
"""
Automated Quality Gate Validator
Pre-validates evidence files before forensic analysis to reduce rework
"""

import os
import re
import sys
from pathlib import Path

class QualityGateValidator:
    def __init__(self, evidence_dir="bytecode/mapping/evidence/verified"):
        self.evidence_dir = Path(evidence_dir)
        self.quality_gates = {
            'bash_commands': {'min_score': 3, 'weight': 20},
            'multi_line_context': {'min_score': 3, 'weight': 20},
            'deob_diagrams': {'min_score': 2, 'weight': 15},
            'relative_paths': {'min_score': 4, 'weight': 10},
            'overview_section': {'min_score': 3, 'weight': 15},
            'architecture_docs': {'min_score': 2, 'weight': 10},
            'command_verification': {'min_score': 2, 'weight': 10}
        }
    
    def analyze_file(self, file_path):
        """Analyze a single evidence file for quality gate compliance"""
        with open(file_path, 'r', encoding='utf-8') as f:
            content = f.read()
        
        scores = {}
        
        # Bash commands presence
        bash_patterns = ['grep -A', 'grep -B', '`', '```bash']
        scores['bash_commands'] = sum(1 for pattern in bash_patterns if pattern in content)
        
        # Multi-line context requirements
        multi_line_patterns = [r'grep -A \d+', r'grep -B \d+', r'-A \d+ -B \d+']
        scores['multi_line_context'] = sum(1 for pattern in multi_line_patterns if re.search(pattern, content))
        
        # DEOB diagrams (no OG names)
        has_mermaid = 'classDiagram' in content or 'flowchart' in content
        has_no_og = 'OG_' not in content or content.count('OG_') <= 2  # Allow minimal references
        scores['deob_diagrams'] = 2 if has_mermaid and has_no_og else (1 if has_mermaid else 0)
        
        # Relative paths only
        has_absolute = '/Users/daxxog/Desktop' in content
        scores['relative_paths'] = 5 if not has_absolute else 0
        
        # Overview section
        overview_keywords = ['## Overview', '## Class Overview', '### Purpose', '### Core Functionality']
        scores['overview_section'] = sum(1 for keyword in overview_keywords if keyword in content)
        
        # Architecture documentation
        arch_keywords = ['## Architecture', '## Relationships', '### Dependencies', 'classDiagram']
        scores['architecture_docs'] = sum(1 for keyword in arch_keywords if keyword in content)
        
        # Command verification
        verify_keywords = ['## Verification', '### Commands', '### Evidence', '### Testing']
        scores['command_verification'] = sum(1 for keyword in verify_keywords if keyword in content)
        
        return scores
    
    def calculate_quality_score(self, scores):
        """Calculate overall quality score from individual scores"""
        total_score = 0
        max_possible = 0
        
        for gate, data in scores.items():
            gate_config = self.quality_gates.get(gate, {})
            min_required = gate_config.get('min_score', 2)
            weight = gate_config.get('weight', 10)
            
            # Score is percentage of minimum requirement met
            gate_score = min(100, (data / min_required) * 100) if min_required > 0 else 0
            weighted_score = (gate_score / 100) * weight
            
            total_score += weighted_score
            max_possible += weight
        
        return (total_score / max_possible) * 100 if max_possible > 0 else 0
    
    def validate_file(self, file_path):
        """Validate a single file and return results"""
        scores = self.analyze_file(file_path)
        quality_score = self.calculate_quality_score(scores)
        
        # Check if passes quality gates
        passed_gates = []
        failed_gates = []
        
        for gate, data in scores.items():
            min_required = self.quality_gates.get(gate, {}).get('min_score', 2)
            if data >= min_required:
                passed_gates.append(gate)
            else:
                failed_gates.append(f"{gate} ({data}/{min_required})")
        
        return {
            'file': file_path.name,
            'quality_score': quality_score,
            'passed_gates': passed_gates,
            'failed_gates': failed_gates,
            'individual_scores': scores
        }
    
    def validate_batch(self, file_limit=10, min_score=60):
        """Validate a batch of files and provide summary"""
        if not self.evidence_dir.exists():
            print(f"Error: Directory {self.evidence_dir} not found")
            return False
        
        files = list(self.evidence_dir.glob("*.md"))[:file_limit]
        if not files:
            print(f"No .md files found in {self.evidence_dir}")
            return False
        
        print(f"Validating {len(files)} files with minimum score: {min_score}")
        print("=" * 60)
        
        passed_files = []
        failed_files = []
        
        for file_path in files:
            result = self.validate_file(file_path)
            
            status = "✅ PASS" if result['quality_score'] >= min_score else "❌ FAIL"
            print(f"{status} {result['file']}: {result['quality_score']:.1f}/100")
            
            if result['quality_score'] >= min_score:
                passed_files.append(result)
            else:
                failed_files.append(result)
                if result['failed_gates']:
                    print(f"   Failed gates: {', '.join(result['failed_gates'])}")
        
        # Summary
        print("=" * 60)
        print(f"VALIDATION SUMMARY:")
        print(f"Passed: {len(passed_files)}/{len(files)} ({len(passed_files)/len(files)*100:.1f}%)")
        print(f"Failed: {len(failed_files)}/{len(files)} ({len(failed_files)/len(files)*100:.1f}%)")
        
        if failed_files:
            print("\nCommon failure patterns:")
            from collections import Counter
            failure_counts = Counter()
            for result in failed_files:
                for gate in result['failed_gates']:
                    gate_name = gate.split(' ')[0]
                    failure_counts[gate_name] += 1
            
            for gate, count in failure_counts.most_common(3):
                print(f"  {gate}: {count} files")
        
        return len(passed_files) == len(files)

def main():
    import argparse
    
    parser = argparse.ArgumentParser(description='Automated Quality Gate Validator')
    parser.add_argument('--min-score', type=int, default=60, help='Minimum quality score to pass')
    parser.add_argument('--file-limit', type=int, default=10, help='Number of files to validate')
    parser.add_argument('--directory', default='bytecode/mapping/evidence/verified', help='Evidence directory')
    
    args = parser.parse_args()
    
    validator = QualityGateValidator(args.directory)
    success = validator.validate_batch(args.file_limit, args.min_score)
    
    sys.exit(0 if success else 1)

if __name__ == "__main__":
    main()