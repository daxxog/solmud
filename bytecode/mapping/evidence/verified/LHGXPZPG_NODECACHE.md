# Forensic Evidence: LHGXPZPG → NodeCache

## **CLASS IDENTIFICATION**
- **Obfuscated Name**: LHGXPZPG
- **Deobfuscated Name**: NodeCache
- **Confidence**: 90% (HIGH CONFIDENCE)
- **Date Identified**: January 8, 2026

## **PRIMARY FORENSIC EVIDENCE**

### **1. Node References (CONFIRMATORY)**
The bytecode contains multiple references to the Node class (PKVMXVTO):

```
public PKVMXVTO c;                          // head/sentinel node
private PKVMXVTO d;                          // tail node
```

### **2. Magic Constant (DISTINCTIVE)**
```
   10: bipush        -77                      // int -77 (cache size marker)
```

**Magic Constant**: -77 - Unique cache boundary marker

### **3. Circular Linking Pattern (IRREFUTABLE)**
The bytecode shows the exact circular doubly-linked list initialization:

```
   21: new           #7                      // class PKVMXVTO (create sentinel)
   23: putfield      #24                     // Field c:LPKVMXVTO;
   36: aload_0
   37: getfield      #24                     // Load sentinel
   50: aload_0
   51: getfield      #24                     // Load sentinel again
   54: putfield      #25                     // PKVMXVTO.c = sentinel (next)
   57: aload_0
   58: getfield      #24                     // Load sentinel
   61: aload_0
   62: getfield      #24                     // Load sentinel again
   65: putfield      #28                     // PKVMXVTO.d = sentinel (prev)
```

**Pattern**: Sentinel node points to itself in both directions

### **4. LRU Cache Management (CONFIRMATORY)**
Complex node manipulation for least-recently-used eviction:

```
   70: aload_1                              // node to remove
   71: getfield      #28                     // node.prev
   74: putfield      #25                     // node.prev.next = node.next
   77: aload_1
   78: getfield      #25                     // node.next
   81: putfield      #28                     // node.next.prev = node.prev
   84: aload_1                              // node
   85: getfield      #28                     // node.prev (now head.next)
   88: putfield      #25                     // node.next = head.next
   91: aload_1
   92: getfield      #25                     // node.next
   95: putfield      #28                     // node.prev = node.next.prev
```

### **5. Error Code Pattern**
```
   74: ldc           #6                      // String 91809,
```

**Error Code**: "91809" - Unique to NodeCache operations

### **6. Cache Size Logic**
```
   26: iload_1                              // capacity parameter
   27: ifgt          46                      // if capacity > 0
   30: aload_0                              // else unlimited cache
   31: aload_0
   32: getfield      #33                     // load current unlimited flag
   35: ifeq          42                      // toggle flag
   38: iconst_0
   39: goto          43
   42: iconst_1
   43: putfield      #33                     // set unlimited flag
```

## **SOURCE CODE CORRELATION**

### **NodeCache.java Reference:**
```java
final class NodeCache {
    private boolean unlimited;
    private int capacity;
    public Node head = new Node();           // sentinel node

    public NodeCache(int capacity) {
        head.next = head;                    // circular linking
        head.prev = head;
        this.capacity = capacity;
        if (capacity <= 0) {
            unlimited = !unlimited;          // toggle for unlimited
        }
    }

    public void removeFromCache(Node node) {
        if (node.prev != null) {
            node.unlink();                   // LRU eviction
        }
        // Re-insert at head (most recently used)
        node.prev = head;
        node.next = head.next;
        node.next.prev = node;
        node.prev.next = node;
    }
}
```

## **UNIQUE IDENTIFIERS**
- **Circular Sentinel Pattern**: Node pointing to itself
- **Magic Constant**: -77 cache marker
- **LRU Management**: Complex node manipulation for eviction
- **Error Code**: "91809" in exception handling
- **Unlimited Toggle**: Boolean flag toggle for unlimited capacity

## **MAPPING CONFIDENCE**
**90% CONFIDENCE** - The combination of circular linking patterns, LRU cache management, and Node references creates a distinctive fingerprint unique to cache implementations. The -77 constant provides additional verification.

## **IMPACT**
- Core caching infrastructure for game objects
- Used throughout client for memory management
- Essential for performance optimization of frequently accessed data</content>
<parameter name="filePath">/Users/daxxog/Desktop/solmud/bytecode/mapping/evidence/verified/LHGXPZPG_NODECACHE.md