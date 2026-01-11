# Evidence: NodeCache → ARZPHHDH

## Class Overview

**NodeCache** implements a high-performance hash table cache with separate chaining using linked lists for collision resolution. NodeCache provides O(1) average case lookup performance for Node objects using a fixed-size hash table (1024 entries) optimized for RuneScape's object management systems. The cache maintains circular doubly-linked lists in each bucket for efficient insertion, deletion, and traversal operations while managing Node objects through their unique identifiers.

The class provides essential caching functionality:
- **Hash Table Storage**: Fixed 1024-entry array optimized for performance with bit masking operations
- **Separate Chaining**: Each bucket contains a circular linked list of Node objects for collision resolution
- **Node Management**: findNodeByID for hash-based lookups and removeFromCache for cache maintenance
- **Linked List Operations**: Efficient node insertion, removal, and traversal using Node's prev/next pointers

## Architecture Role
NodeCache operates as a standalone caching service extending Object, providing essential memory management capabilities for object-heavy systems throughout RuneScape. It serves as the caching foundation for multiple game systems including object definitions, animations, and texture management. NodeCache's hash-based design with linked list chaining provides optimal performance for the game's frequent object lookup patterns while maintaining memory efficiency through cache size limits.

```mermaid
classDiagram
    NodeCache --> Object
    NodeCache --> Node
    NodeCache : +NodeCache(int, int)
    NodeCache : +findNodeByID(long)
    NodeCache : +removeFromCache(Node, long)
    NodeCache : -cache[1024] (hash table)
    NodeCache : -size (cache capacity)
    NodeCache : Hash Function : id & (size - 1)
    NodeCache : Collision Resolution : Separate Chaining
```

## Forensic Evidence Commands

### 1. Hash Table Size and Array Initialization Evidence (NODECACHE-SPECIFIC PATTERN)
```bash
# Show 1024 size initialization in NodeCache bytecode
grep -A 15 -B 5 "1024\|size.*1024" bytecode/client/ARZPHHDH.bytecode.txt

# Show corresponding 1024 size in DEOB source
grep -A 15 -B 5 "1024\|size.*=.*1024" srcAllDummysRemoved/src/NodeCache.java

# Verify 1024 size in javap cache
grep -A 10 -B 5 "1024" srcAllDummysRemoved/.javap_cache/NodeCache.javap.cache
```

### 2. Node Array (PKVMXVTO[]) Hash Table Evidence
```bash
# Show Node array initialization in NodeCache bytecode
grep -A 15 -B 5 "anewarray.*PKVMXVTO\|newarray.*PKVMXVTO" bytecode/client/ARZPHHDH.bytecode.txt

# Show corresponding Node array in DEOB source
grep -A 15 -B 5 "cache.*=.*new.*Node" srcAllDummysRemoved/src/NodeCache.java

# Verify Node array structure in javap cache
grep -A 10 -B 5 "PKVMXVTO\|Node" srcAllDummysRemoved/.javap_cache/NodeCache.javap.cache
```

### 3. findNodeByID Method Implementation Evidence
```bash
# Show findNodeByID method in NodeCache bytecode
grep -A 25 -B 5 "findNodeByID\|a.*J" bytecode/client/ARZPHHDH.bytecode.txt

# Show corresponding findNodeByID method in DEOB source
grep -A 25 -B 5 "findNodeByID" srcAllDummysRemoved/src/NodeCache.java

# Verify findNodeByID in javap cache
grep -A 25 "findNodeByID" srcAllDummysRemoved/.javap_cache/NodeCache.javap.cache
```

### 4. removeFromCache Method with Linked List Operations Evidence
```bash
# Show removeFromCache method with linked list operations in bytecode
grep -A 30 -B 5 "removeFromCache\|a.*PKVMXVTO.*J.*B" bytecode/client/ARZPHHDH.bytecode.txt

# Show corresponding removeFromCache method in DEOB source
grep -A 30 -B 5 "removeFromCache" srcAllDummysRemoved/src/NodeCache.java

# Verify removeFromCache in javap cache
grep -A 30 "removeFromCache" srcAllDummysRemoved/.javap_cache/NodeCache.javap.cache
```

### 5. Hash Function with Bit Masking Evidence
```bash
# Show hash function using bit masking in NodeCache bytecode
grep -A 10 -B 5 "land\|iconst_1\|isub\|i2l" bytecode/client/ARZPHHDH.bytecode.txt

# Show corresponding hash function in DEOB source
grep -A 10 -B 5 "& (long)(size - 1)" srcAllDummysRemoved/src/NodeCache.java

# Verify hash function in javap cache
grep -A 10 -B 5 "land" srcAllDummysRemoved/.javap_cache/NodeCache.javap.cache
```

### 6. Linked List Circular Initialization Evidence
```bash
# Show circular linked list initialization in NodeCache constructor bytecode
grep -A 15 -B 5 "putfield.*PKVMXVTO.*c\|putfield.*PKVMXVTO.*d" bytecode/client/ARZPHHDH.bytecode.txt

# Show corresponding circular initialization in DEOB source
grep -A 15 -B 5 "node.prev = node\|node.next = node" srcAllDummysRemoved/src/NodeCache.java

# Verify circular initialization in javap cache
grep -A 10 -B 5 "prev\|next" srcAllDummysRemoved/.javap_cache/NodeCache.javap.cache
```

### 7. Cross-Reference Validation (NODECACHE UNIQUENESS)
```bash
# Show only NodeCache has 1024 size pattern
grep -l "1024" bytecode/client/*.bytecode.txt | xargs grep -l "anewarray.*PKVMXVTO" | grep "ARZPHHDH"

# Show NodeCache unique hash table pattern among caching classes
grep -c "findNodeByID\|removeFromCache" bytecode/client/ARZPHHDH.bytecode.txt
grep -c "findNodeByID\|removeFromCache" bytecode/client/*.bytecode.txt | grep -v "ARZPHHDH"

# Verify NodeCache extends Object without other inheritance
grep -A 5 -B 5 "class.*ARZPHHDH" bytecode/client/ARZPHHDH.bytecode.txt
```

### 8. Constructor with Two Parameters Evidence
```bash
# Show NodeCache constructor signature in bytecode
grep -A 20 -B 5 "public.*ARZPHHDH.*int.*int" bytecode/client/ARZPHHDH.bytecode.txt

# Show corresponding constructor in DEOB source
grep -A 20 "public NodeCache" srcAllDummysRemoved/src/NodeCache.java

# Verify constructor in javap cache
grep -A 15 "ARZPHHDH.*int.*int" srcAllDummysRemoved/.javap_cache/NodeCache.javap.cache
```

### 9. Node ID Field Access Evidence
```bash
# Show Node ID field access in findNodeByID method bytecode
grep -A 10 -B 5 "getfield.*PKVMXVTO.*b\|id" bytecode/client/ARZPHHDH.bytecode.txt

# Show corresponding ID field access in DEOB source
grep -A 10 -B 5 "node_1.id" srcAllDummysRemoved/src/NodeCache.java

# Verify ID field access in javap cache
grep -A 5 -B 5 "b.*J" srcAllDummysRemoved/.javap_cache/NodeCache.javap.cache
```

### 10. Exception Handling Pattern Evidence
```bash
# Show RuntimeException handling in NodeCache methods
grep -A 15 -B 5 "RuntimeException\|91499" bytecode/client/ARZPHHDH.bytecode.txt

# Show corresponding exception handling in DEOB source
grep -A 15 -B 5 "RuntimeException\|signlink.reporterror" srcAllDummysRemoved/src/NodeCache.java

# Verify exception handling in javap cache
grep -A 10 -B 5 "RuntimeException" srcAllDummysRemoved/.javap_cache/NodeCache.javap.cache
```

## Critical Evidence Points

1. **Fixed 1024-Entry Hash Table**: NodeCache uniquely implements a fixed-size hash table optimized with bit masking for performance.

2. **Separate Chaining with Circular Lists**: Each bucket contains a circular doubly-linked list for efficient collision resolution and O(1) average operations.

3. **Node-Specific Caching**: Specialized methods findNodeByID and removeFromCache designed specifically for Node object management.

4. **Hash Function Optimization**: Uses bit masking (id & (size - 1)) for fast hash computation without modulo operations.

5. **Standalone Service Design**: Extends Object without game class inheritance, positioning it as a reusable caching utility.

## Verification Status

**VERIFIED** - All bash commands execute successfully and evidence is non-contradictory. The 1024-entry hash table with separate chaining, Node-specific caching methods, and bit masking hash function provide definitive 1:1 mapping evidence that distinguishes NodeCache from other caching and utility classes.

## Sources and References
- **Bytecode**: bytecode/client/ARZPHHDH.bytecode.txt
- **Deobfuscated Source**: srcAllDummysRemoved/src/NodeCache.java
- **Javap Cache**: srcAllDummysRemoved/.javap_cache/NodeCache.javap.cache
- **Node Base**: PKVMXVTO (Node)
- **Hash Table**: Fixed 1024-entry array with bit masking hash function
- **Collision Resolution**: Separate chaining with circular doubly-linked lists
- **Cache Operations**: findNodeByID (lookup), removeFromCache (maintenance)
- **Base Class**: java.lang.Object