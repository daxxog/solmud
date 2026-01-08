# Forensic Bytecode Evidence Repository

## Overview

This repository contains detailed forensic evidence documentation for class name mappings between deobfuscated and obfuscated RuneScape 317 bytecode. Each evidence file follows a systematic forensic methodology to establish mapping confidence beyond automated scoring.

## Forensic Methodology Framework

### Evidence Types

#### 1. Anchor Class Evidence
**Definition**: Classes with unique, identifying characteristics that serve as ground truth.
**Forensic Techniques**:
- File loading pattern analysis (e.g., loads specific .dat files)
- Unique implementation signatures
- Critical infrastructure roles

**Strength**: 100% confidence when confirmed
**Examples**: `ItemDef` → `DJRMEMXO` (loads "obj.dat")

#### 2. Inheritance Chain Evidence
**Definition**: Classes identified through confirmed superclass relationships.
**Forensic Techniques**:
- Superclass anchor verification
- Sibling class consistency analysis
- Hierarchy topology validation

**Strength**: 85% confidence when chain is verified
**Examples**: `Player` → `DLZHLHNK` (extends `Entity` → `GQOSZKJC`)

#### 3. Signature Evidence
**Definition**: Classes matched through structural and methodological patterns.
**Forensic Techniques**:
- Field type distribution analysis
- Method signature similarity scoring
- Access modifier pattern matching
- Constructor pattern analysis

**Strength**: 65-84% confidence depending on correlation strength

#### 4. Behavioral Evidence
**Definition**: Classes identified through unique bytecode execution patterns.
**Forensic Techniques**:
- Instruction sequence fingerprinting
- Method call graph analysis
- State manipulation pattern detection
- Cross-reference dependency tracing

**Strength**: Variable confidence based on uniqueness of patterns

### Forensic Validation Techniques

#### Cross-Reference Validation
- Validate mappings through dependency relationships
- Ensure bi-directional consistency in class references
- Identify circular validation loops

#### Pattern Consistency Analysis
- Compare against known obfuscation patterns
- Validate naming conventions compliance
- Check for structural anomalies

#### Multi-Method Corroboration
- Combine multiple forensic techniques
- Weight evidence by reliability
- Require minimum evidence thresholds

## Evidence File Organization

### File Naming Convention
```
evidence/<DEOB_FILENAME>.md
```

### Evidence Categories
- **High Confidence** (≥85%): Anchor and validated inheritance chains
- **Medium Confidence** (65-84%): Strong signature and behavioral matches
- **Low Confidence** (<65%): Tentative matches requiring additional evidence

### Evidence Status Codes
- **VERIFIED**: Multiple independent forensic techniques confirm
- **PROBABLE**: Strong evidence with minor gaps
- **TENTATIVE**: Preliminary evidence pending verification
- **DISPUTED**: Evidence conflicts with other mappings

## Forensic Analysis Tools

### Primary Tool: classmapper
- Automated pattern matching with confidence scoring
- Multi-pass resolution (anchor → inheritance → signature → behavioral)
- Cross-reference validation capabilities

### Supporting Tools
- **Javap Analysis**: Deobfuscated class structure parsing
- **Bytecode Analysis**: Obfuscated class instruction examination
- **Cross-Reference Tracer**: Dependency relationship mapping
- **Pattern Validator**: Consistency checking against known patterns

## Forensic Evidence Standards

### Minimum Evidence Requirements
- **High Confidence**: Minimum 2 independent forensic techniques
- **Medium Confidence**: Minimum 1 strong forensic technique + supporting evidence
- **Low Confidence**: Single forensic technique with preliminary evidence

### Evidence Quality Metrics
- **Reproducibility**: Evidence can be independently verified
- **Specificity**: Evidence uniquely identifies the class
- **Consistency**: Evidence aligns with related mappings
- **Completeness**: Evidence covers critical class characteristics

## Verification Workflow

### 1. Initial Detection
- Run classmapper with comprehensive analysis
- Identify potential matches with confidence scores
- Generate preliminary evidence files

### 2. Evidence Collection
- Extract structural fingerprints
- Analyze behavioral patterns
- Document cross-reference relationships

### 3. Cross-Validation
- Validate against known anchor classes
- Check inheritance chain consistency
- Verify dependency relationships

### 4. Manual Review
- Forensic expert examination of evidence
- Risk assessment and confidence validation
- Final evidence quality determination

### 5. Documentation
- Complete evidence file creation
- Source referencing and tool versioning
- Peer review and final approval

## Current Evidence Coverage

### Statistics
- **Total Classes Analyzed**: `74`
- **High Confidence Mappings**: `36` (automated) + `4` (forensic corrections)
- **Medium Confidence Mappings**: `2` (automated)
- **Low Confidence Mappings**: `0` (automated)
- **Evidence Coverage**: `57%` (after corrections)

### Priority Targets
The following critical mappings are priority evidence collection targets:
1. **ISAACRandomGen** → JOCFVBOI (Security-critical)
2. **MRUNodes** → GCPOSBWX (Memory management)
3. **NodeCache** → ARZPHHDH (Cache infrastructure)
4. **Stream** → [NEEDS IDENTIFICATION] (I/O operations)

## Quality Assurance

### Evidence Review Process
- Initial forensic analyst documentation
- Secondary analyst verification
- Cross-technique consistency check
- Final forensic expert approval

### Continuous Improvement
- Evidence file template refinement
- Forensic technique optimization
- Tool accuracy enhancement
- Pattern library expansion

## Contributing to Evidence Collection

### Evidence Submission Guidelines
1. Follow the evidence file template exactly
2. Document all forensic techniques used
3. Include source references and tool versions
4. Provide confidence scoring justification
5. Cross-validate against related mappings

### Evidence Review Criteria
- Forensic method validity
- Evidence completeness
- Cross-reference consistency
- Confidence scoring accuracy
- Documentation quality

---
**Last Updated**: `2026-01-08`
**Evidence Version**: `1.0`
**Maintained By**: `opencode`