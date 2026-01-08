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

### Folder Structure
```
evidence/
├── verified/           # 100% confirmed mappings with irrefutable evidence
├── disputed/          # Incorrect claims moved here with dispute documentation
├── in_progress/       # Active research and verification work
├── conflict_resolution/ # Previous conflict resolutions and analyses
└── evidence_template.md # Standard template for evidence documentation
```

### File Naming Convention
```
verified/<DEOB_FILENAME>.md              # Confirmed mappings
disputed/<OBF_CLASS>_DISPUTED_NOT_<DEOB>.md # Disputed claims
in_progress/<WORK_FILE>.md              # Active research
conflict_resolution/<TOPIC>.md          # Conflict analyses
```

### Evidence Categories
- **High Confidence** (≥85%): Anchor and validated inheritance chains
- **Medium Confidence** (65-84%): Strong signature and behavioral matches
- **Low Confidence** (<65%): Tentative matches requiring additional evidence

### Evidence Status Codes
- **VERIFIED**: Multiple independent forensic techniques confirm (in `verified/`)
- **DISPUTED**: Incorrect claims verified through cross-validation (in `disputed/`)
- **IN_PROGRESS**: Active research and verification (in `in_progress/`)
- **RESOLVED**: Conflicts successfully resolved (in `conflict_resolution/`)

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
- **NEW**: Cross-check error codes against source files
- **NEW**: Verify array structures match exactly

### 4. Manual Review
- Forensic expert examination of evidence
- Risk assessment and confidence validation
- Final evidence quality determination
- **NEW**: Dispute identification and correction

### 5. Documentation
- Complete evidence file creation
- Source referencing and tool versioning
- Peer review and final approval
- **NEW**: Dispute documentation for incorrect claims

## Dispute Resolution Process

### When Disputes Are Identified
1. **Error Code Verification**: Cross-reference error codes with source files
2. **Structure Analysis**: Verify array sizes and structures match exactly
3. **Source Bytecode Comparison**: Compare bytecode patterns with source code
4. **Conflict Documentation**: Move disputed claims to `disputed/` folder with correction notes
5. **Tool Correction**: Update Go tooling with verified correct mappings
6. **Evidence Reorganization**: Ensure correct evidence stays in `verified/`

### Dispute Status Codes
- **INCORRECT**: Claim proven wrong through verification
- **CONFLICTING**: Multiple valid interpretations
- **UNCERTAIN**: Insufficient evidence for determination

### Moving Evidence Between Folders
```bash
# Move incorrect claim:
mv verified/WRONG_MAPPING.md disputed/WRONG_DISPUTED.md

# Add dispute notice:
# Add ⚠️ header explaining why claim is wrong
# Include correct mapping reference
# Document verification evidence that disproves claim
```

## Current Evidence Coverage

### Statistics (as of January 8, 2026 - CORRECTED)
- **Total Classes**: `74`
- **Verified Mappings**: `44` (100% confidence)
- **High Confidence Mappings**: `1` (90% confidence)
- **Medium Confidence Mappings**: `1` (70% confidence)
- **Total Coverage**: `60.8%` (45/74 classes)
- **High-Confidence Rate**: `97.8%` (44/45 mapped classes)

### Recent Corrections (January 8, 2026)
- **Removed**: `LHGXPZPG → NodeCache` (INCORRECT - moved to `disputed/`)
- **Removed**: `VBAXKVMG → Class32` (INCORRECT - moved to `disputed/`)
- **Added**: `Skills → YUXCUCXD` (VERIFIED - literal skill names)
- **Added**: `TextClass → ZTQFNQRH` (VERIFIED - base-37 hashing)
- **Added**: `StreamLoader → XTGLDHGX` (VERIFIED - base-61 hashing)

### Priority Research Targets
The following critical mappings are priority research targets:
1. **Class32** → [NEEDS IDENTIFICATION] (Bzip2 decompression - search arrays: 256, 257, 258, 6, 16, 4096, 18002)
2. **LHGXPZPG** → [NEEDS IDENTIFICATION] (Error code "91809" - NOT NodeCache)
3. **Object4** → FEHPTPDG (Animable references, coordinate systems)
4. **Object1/2/3** → OIBEELAZ (Boolean flag systems, object management)

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
**Evidence Version**: `1.1` (CORRECTED)
**Maintained By**: `opencode`
**Recent Corrections**: `2` incorrect mappings identified and moved to `disputed/`