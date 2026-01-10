# BISVHPUN → NodeList

## Overview
NodeList manages a doubly-linked list of NodeSub objects using a sentinel node pattern. It provides methods for inserting, removing, and iterating through nodes efficiently.

## Architectural Relationships
NodeList extends Object and uses NodeSub (PPOHBEGB) for list elements. It maintains a sentinel head node and current node for traversal.

```mermaid
classDiagram
    BISVHPUN --> PPOHBEGB[NodeSub]
    BISVHPUN --> Object
```

## Bytecode Matching Commands
To show class structure and sentinel initialization:
```
grep -A 20 -B 5 "public final class BISVHPUN" bytecode/client/BISVHPUN.bytecode.txt
```

To show insertion method:
```
grep -A 15 -B 5 "public final void a" bytecode/client/BISVHPUN.bytecode.txt
```

## Deobfuscated Source Evidence Commands
For NodeList constructor:
```
grep -A 10 "public NodeList" srcAllDummysRemoved/src/NodeList.java
```

For insertHead method:
```
grep -A 10 "public void insertHead" srcAllDummysRemoved/src/NodeList.java
```

## Javap Cache Evidence Commands
For class structure:
```
grep -A 10 "public final class NodeList" srcAllDummysRemoved/.javap_cache/NodeList.javap.cache
```

For methods:
```
grep -A 5 "insertHead" srcAllDummysRemoved/.javap_cache/NodeList.javap.cache
```

## Verification of Non-Contradictory Evidence
Bytecode matches source/javap in sentinel pattern, NodeSub usage. No contradictions. 1:1 mapping confirmed.