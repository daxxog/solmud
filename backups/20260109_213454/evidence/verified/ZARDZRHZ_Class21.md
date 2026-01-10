# ZARDZRHZ → Class21

## Overview
Class21 is a data container with 19 integer fields for game data storage.

## Architectural Relationships
Class21 extends Object with multiple int fields.

```mermaid
classDiagram
    ZARDZRHZ --> Object
    ZARDZRHZ : +int anInt370
    ZARDZRHZ : +int anInt371
    // ... 17 more int fields
```

## Bytecode Matching Commands
To show class and fields:
```
cat bytecode/client/ZARDZRHZ.bytecode.txt
```

To show field count:
```
grep -c "int" bytecode/client/ZARDZRHZ.bytecode.txt
```

## Deobfuscated Source Evidence Commands
For Class21 class:
```
grep -A 20 "public final class Class21" srcAllDummysRemoved/src/Class21.java
```

For fields:
```
grep -A 10 "anInt" srcAllDummysRemoved/src/Class21.java
```

## Javap Cache Evidence Commands
For class structure:
```
grep -A 20 "public final class Class21" srcAllDummysRemoved/.javap_cache/Class21.javap.cache
```

For fields:
```
grep -A 10 "anInt" srcAllDummysRemoved/.javap_cache/Class21.javap.cache
```

## Verification of Non-Contradictory Evidence
Bytecode matches source/javap in int field storage. No contradictions. 1:1 mapping confirmed.