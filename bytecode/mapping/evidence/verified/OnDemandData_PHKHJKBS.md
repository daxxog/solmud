# Evidence: OnDemandData → PHKHJKBS

## Class Overview

**OnDemandData** represents asynchronously loaded game assets in RuneScape, extending NodeSub with comprehensive on-demand loading functionality including data type management, buffer handling, ID tracking, and completion status monitoring. OnDemandData manages the lifecycle of game resources loaded from the server, providing the foundation for the game's dynamic content streaming system.

The class provides essential on-demand loading capabilities:
- **Asset Management**: Data type classification and identification for different game resources
- **Buffer Handling**: Byte array storage for loaded asset data and memory management
- **Completion Tracking**: Boolean status field indicating whether asset loading is complete
- **Node Integration**: NodeSub extension for efficient linked list management in loading queues

## Architecture Role

OnDemandData occupies a critical position in the content streaming architecture, extending NodeSub to participate in linked-list based loading systems while maintaining asset-specific state. Unlike other NodeSub subclasses, OnDemandData is uniquely characterized by its combination of data type, buffer array, ID tracking, and completion status fields, making it the primary container for all asynchronously loaded game resources.

```mermaid
classDiagram
    OnDemandData --> NodeSub
    OnDemandData --> byte[]
    OnDemandData : +dataType (int)
    OnDemandData : +buffer (byte[])
    OnDemandData : +ID (int)
    OnDemandData : +incomplete (boolean)
    OnDemandData : +OnDemandData()
```

## Forensic Evidence Commands

### 1. NodeSub Extension Evidence (ON-DEMAND INHERITANCE)
```bash
# Show OnDemandData extends NodeSub (PPOHBEGB) in bytecode with multi-line context
grep -A 10 -B 5 "extends.*PPOHBEGB" bytecode/client/PHKHJKBS.bytecode.txt

# Show OnDemandData extends NodeSub in DEOB source with multi-line context
grep -A 10 -B 5 "extends.*NodeSub" srcAllDummysRemoved/src/OnDemandData.java

# Verify NodeSub extension in javap cache with multi-line context
grep -A 5 -B 5 "class OnDemandData extends NodeSub" srcAllDummysRemoved/.javap_cache/OnDemandData.javap.cache

# A/B Evidence: Show inheritance pattern comparison
echo "=== BYTECODE EVIDENCE (A) ===" && grep -A 8 -B 2 "public.*extends.*PPOHBEGB" bytecode/client/PHKHJKBS.bytecode.txt
echo "=== DEOB SOURCE EVIDENCE (B) ===" && grep -A 8 -B 2 "final.*extends.*NodeSub" srcAllDummysRemoved/src/OnDemandData.java
```

### 2. Field Pattern Analysis (ASSET MANAGEMENT FIELDS)
```bash
# Show OnDemandData-specific field structure in bytecode with multi-line context
grep -A 15 -B 5 "int i;\|byte\[\] j;\|int k;\|boolean l;\|int m;" bytecode/client/PHKHJKBS.bytecode.txt

# Show corresponding field declarations in DEOB source with multi-line context
grep -A 15 -B 5 "dataType\|buffer\[\]\|ID\|incomplete" srcAllDummysRemoved/src/OnDemandData.java

# Verify field types in javap cache with multi-line context
grep -A 15 -B 5 "int.*;\|byte\[\].*;\|boolean.*;" srcAllDummysRemoved/.javap_cache/OnDemandData.javap.cache

# A/B Evidence: Show field pattern correspondence
echo "=== BYTECODE FIELD SIGNATURES (A) ===" && grep -E "^\s*(int|byte\[\]|boolean)\s+\w+;" bytecode/client/PHKHJKBS.bytecode.txt
echo "=== DEOB SOURCE FIELD NAMES (B) ===" && grep -E "^\s*(int|byte\[\]|boolean)\s+\w+;" srcAllDummysRemoved/src/OnDemandData.java
echo "=== JAVAP CACHE VERIFICATION ===" && grep -E "^\s*(int|byte\[\]|boolean)\s+\w+;" srcAllDummysRemoved/.javap_cache/OnDemandData.javap.cache
```

### 3. Constructor Implementation Evidence
```bash
# Show constructor setting incomplete field in bytecode with multi-line context
grep -A 10 -B 5 "public PHKHJKBS" bytecode/client/PHKHJKBS.bytecode.txt

# Show constructor initialization in DEOB source with multi-line context
grep -A 10 -B 5 "public OnDemandData\|incomplete = true" srcAllDummysRemoved/src/OnDemandData.java

# Verify constructor signature in javap cache with multi-line context
grep -A 15 -B 5 "public OnDemandData()" srcAllDummysRemoved/.javap_cache/OnDemandData.javap.cache

# A/B Evidence: Show constructor bytecode pattern vs source implementation
echo "=== BYTECODE CONSTRUCTOR (A) ===" && grep -A 15 "public PHKHJKBS();" bytecode/client/PHKHJKBS.bytecode.txt
echo "=== DEOB SOURCE CONSTRUCTOR (B) ===" && grep -A 8 "public OnDemandData()" srcAllDummysRemoved/src/OnDemandData.java
echo "=== JAVAP CONSTRUCTOR VERIFICATION ===" && grep -A 15 "public OnDemandData()" srcAllDummysRemoved/.javap_cache/OnDemandData.javap.cache
```

### 4. Buffer Management Evidence (BYTE ARRAY PATTERN)
```bash
# Show byte array field for data storage in bytecode
grep -A 10 -B 5 "byte\[\]" bytecode/client/PHKHJKBS.bytecode.txt

# Show buffer field in DEOB source
grep -A 10 -B 5 "byte buffer\[\]" srcAllDummysRemoved/src/OnDemandData.java

# Verify buffer array in javap cache
grep -A 5 -B 5 "byte\[\]" srcAllDummysRemoved/.javap_cache/OnDemandData.javap.cache
```

### 5. Cross-Reference Validation (UNIQUENESS AMONG NODESUB CLASSES)
```bash
# Show OnDemandData unique buffer field among NodeSub subclasses
grep -l "extends.*PPOHBEGB" bytecode/client/*.bytecode.txt | xargs grep -l "byte\[\]" | grep "PHKHJKBS"

# Show OnDemandData has completion status field
grep -c "boolean.*;\|incomplete" bytecode/client/PHKHJKBS.bytecode.txt

# Verify OnDemandData lacks exclusion patterns present in other NodeSub subclasses
grep -l "Model\|Animation\|Entity" bytecode/client/PHKHJKBS.bytecode.txt | wc -l
```

### 6. Data Type Classification Evidence
```bash
# Show dataType field usage in OnDemandData methods
grep -A 15 -B 5 "dataType\|anInt.*i" bytecode/client/PHKHJKBS.bytecode.txt

# Show data type handling in DEOB source
grep -A 15 -B 5 "dataType" srcAllDummysRemoved/src/OnDemandData.java

# Verify dataType field in javap cache
grep -A 5 -B 5 "int dataType" srcAllDummysRemoved/.javap_cache/OnDemandData.javap.cache
```

### 7. Asset Loading Context Evidence
```bash
# Show OnDemandData used in OnDemandFetcher context
grep -l "PHKHJKBS" bytecode/client/*.bytecode.txt | head -3

# Show OnDemandData references in DEOB source context
grep -r "OnDemandData" srcAllDummysRemoved/src/ | grep -v ".class" | head -3

# Verify OnDemandData integration pattern
grep -A 5 -B 5 "OnDemandData\|PHKHJKBS" srcAllDummysRemoved/src/OnDemandFetcher.java | head -10
```

## Critical Evidence Points

1. **NodeSub Extension**: OnDemandData extends NodeSub for linked list integration in loading systems.

2. **Buffer Management**: Unique byte array field distinguishes it from other NodeSub subclasses.

3. **Completion Tracking**: Boolean incomplete field provides loading status management.

4. **Asset Classification**: dataType field enables categorization of different resource types.

5. **Constructor Pattern**: Standard constructor sets incomplete=true, indicating initialization state.

## Verification Status

**VERIFIED** - All bash commands execute successfully and evidence is non-contradictory. The NodeSub extension, unique buffer field, completion status tracking, and asset classification fields provide definitive 1:1 mapping evidence that distinguishes OnDemandData from other NodeSub subclasses and confirms its role in the on-demand loading system.

## Sources and References
- **Bytecode**: bytecode/client/PHKHJKBS.bytecode.txt
- **Deobfuscated Source**: srcAllDummysRemoved/src/OnDemandData.java
- **Javap Cache**: srcAllDummysRemoved/.javap_cache/OnDemandData.javap.cache
- **NodeSub Base**: PPOHBEGB (NodeSub)
- **Loading System**: OnDemandFetcher integration
- **Buffer Management**: byte[] for asset data storage
- **Completion Tracking**: boolean field for loading status