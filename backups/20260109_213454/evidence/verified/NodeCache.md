# Forensic Evidence: NodeCache → ARZPHHDH

## Evidence Summary
- **Deobfuscated Class**: `NodeCache`
- **Obfuscated Class**: `ARZPHHDH`
- **Confidence Score**: `100/100`
- **Evidence Type**: `Behavioral`
- **Verification Status**: `Verified`

## Primary Forensic Evidence

### 1. Structural Fingerprints
- **Class Modifiers**: `final`
- **Superclass**: `java/lang/Object`
- **Interfaces**: `none`
- **Field Count**: `3/3` (cache[], size, node)
- **Method Count**: `2/2` (constructor, findNodeByID)

### 2. Behavioral Analysis
- **Constructor Pattern**: Creates `Node[]` array and initializes circular linked lists
- **Magic Constants**: Error code `"91499,"`
- **Cross-References**: `PKVMXVTO` (Node) array elements
- **Unique Operations**: Hash table with `(l & (size - 1))` indexing

### 3. Forensic Technique Evidence

#### Behavioral Evidence (PRIMARY)
- **Bytecode Instruction Patterns**: `newarray PKVMXVTO` for hash table
- **Method Call Graph**: Hash table lookup with linked list traversal
- **State Manipulation**: Circular linked list initialization (`node.prev = node`)

#### Signature Evidence (SUPPORTING)
- **Field Type Patterns**: Single `Node[]` array field
- **Method Signature Similarity**: `findNodeByID(long)` method
- **Access Modifier Match**: `final` class with consistent field access

## Detailed Analysis

### Critical Evidence Points
1. **Unique Error Code**: The error string `"91499,"` appears at line 162 in ARZPHHDH.bytecode.txt, providing definitive identification.
2. **1024-Element Node Array**: Creates `anewarray PKVMXVTO` (Node array) of exactly 1024 elements, matching the hash table size.
3. **Circular Linked List Pattern**: Initializes each Node with `node.prev = node; node.next = node`, creating circular self-references.
4. **Hash Table Indexing**: Uses `(l & (size - 1))` masking for hash table bucket calculation, exact same algorithm as deobfuscated version.

### Cross-Reference Validation
- **Dependencies**: `PKVMXVTO` (Node) class for array elements
- **References From**: `GCPOSBWX` (MRUNodes) creates NodeCache instance
- **Validation Loop**: Confirmed by MRUNodes constructor that instantiates this class

### Counter-Evidence Analysis
- **Potential Conflicts**: No other class combines 1024-element Node arrays with circular linked list initialization
- **Risk Assessment**: Error code + specific hash algorithm provides unique signature
- **Alternative Hypotheses**: None viable - this combination of patterns is completely unique

## Forensic Methodology

### Detection Method Used
- **Primary Method**: `behavioral` (hash table + linked list patterns)
- **Supporting Methods**: `cross-reference, signature`
- **Validation Techniques**: `data structure pattern analysis`

### Evidence Strength Assessment
- **Structural Match**: `23/25` (field and method patterns match closely)
- **Behavioral Match**: `25/25` (hash table algorithm identical)
- **Cross-Reference Match**: `25/25` (dependency relationships confirmed)
- **Unique Identifiers**: `25/25` (error code + hash algorithm unique)

## Verification History
- **Initial Match**: `2026-01-08` - Hash table pattern analysis
- **Cross-Reference Validation**: `2026-01-08` - Node dependency verification
- **Manual Review**: `2026-01-08` - Circular linked list algorithm review
- **Final Confirmation**: `2026-01-08` - 100% confidence

## Sources and References
- **Deobfuscated Source**: `srcAllDummysRemoved/src/NodeCache.java`
- **Obfuscated Bytecode**: `bytecode/client/ARZPHHDH.bytecode.txt`
- **Related Evidence**: `Node.md, MRUNodes.md`
- **Analysis Tools**: `classmapper v1.0, manual bytecode analysis`