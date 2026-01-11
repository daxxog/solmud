# Class Mapper Tool

Standalone Go tool for mapping deobfuscated Java class names to obfuscated bytecode names with confidence scoring.

## Purpose

Maps deobfuscated class names (from Mopar 317 deobfuscated client) to obfuscated class names (from original RuneScape 317 bytecode) using intelligent pattern matching.

## Installation

### Build from source

```bash
cd tools
make build
```

### Install to GOPATH

```bash
cd tools
make install
```

## Usage

### Basic usage (run from project root)

```bash
# Build
cd tools && make build
cd ..

# CSV mode (default)
./tools/classmapper/classmapper

# JSON mode
./tools/classmapper/classmapper -mode json

# Custom settings
./tools/classmapper/classmapper -mode json -deob srcAllDummysRemoved/bin -obf bytecode/client -threshold 75
```

## Options

- `-deob`: Path to deobfuscated class files (default: `srcAllDummysRemoved/bin`)
- `-obf`: Path to obfuscated bytecode files (default: `bytecode/client`)
- `-mode`: Output mode: `csv` or `json` (default: `csv`)
- `-threshold`: Minimum confidence threshold (0-100) (default: 65.0)

## Output Formats

### CSV Mode (default)

Outputs to stdout in CSV format with separate sections for high confidence and uncertain matches.

```csv
High Confidence Matches (CSV):
deobfuscated_name,obfuscated_name,confidence_score,superclass_match,interface_count,field_count_deob,field_count_obf,method_count_deob,method_count_obf,notes
Node,PKVMXVTO,100.00,true,0,N/A,N/A,N/A,N/A,Exact match - anchor class

Uncertain Matches (CSV):
deobfuscated_name,obfuscated_name,confidence_score,notes
Class33,BMEXSMOV,68.50,Matched by signature - score 68.50
```

### JSON Mode

Outputs complete JSON structure to stdout with detailed score breakdowns.

```json
{
  "summary": {
    "total_matches": 3,
    "high_confidence_count": 3,
    "medium_confidence_count": 0,
    "low_confidence_count": 0
  },
  "matches": [
    {
      "DeobfuscatedClass": "Node",
      "ObfuscatedClass": "PKVMXVTO",
      "ConfidenceScore": 100,
      "ScoreBreakdown": {
        "InterfaceMatch": 0,
        "SuperclassMatch": 25,
        "FieldCountMatch": 0,
        "FieldSimilarity": 0,
        "MethodCountMatch": 5,
        "MethodSimilarity": 0,
        "ConstructorMatch": 5,
        "AccessMatch": 5,
        "SizePenalty": 0
      },
      "Details": "Exact match - anchor class"
    }
  ]
}
```

## Requirements

- Go 1.21.5 or later
- Java Development Kit (javap must be in PATH)
- Mopar 317 deobfuscated classes
- Original RuneScape 317 bytecode files

## How It Works

1. Parses deobfuscated classes using javap
2. Parses obfuscated classes from existing .bytecode.txt files
3. Uses three-pass resolution:
   - Pass 1: Match known anchor classes
   - Pass 2: Match by inheritance hierarchy
   - Pass 3: Match by method/field signatures
4. Calculates confidence scores based on:
   - Interface implementations
   - Superclass matches
   - Field counts and types
   - Method counts and signatures
   - Access modifiers

## License

Same as parent project.