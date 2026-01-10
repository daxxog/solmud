# CWNCPMLX → NPC

## Overview
NPC represents non-player characters in RuneScape, extending Entity for shared properties like position and animation, with additional fields for NPC-specific data.

## Architectural Relationships
NPC extends Entity (GQOSZKJC), uses EntityDef (CKDEJADD), Model (ZKARKDQW), and Animation (LKGEGIEW) for rendering and behavior.

```mermaid
classDiagram
    CWNCPMLX --> GQOSZKJC[Entity]
    CWNCPMLX --> CKDEJADD[EntityDef]
    CWNCPMLX --> ZKARKDQW[Model]
    CWNCPMLX --> LKGEGIEW[Animation]
```

## Bytecode Matching Commands
To show inheritance and fields:
```
head -n 20 bytecode/client/CWNCPMLX.bytecode.txt
```

To show method signatures:
```
grep -A 10 -B 5 "public.*ZKARKDQW" bytecode/client/CWNCPMLX.bytecode.txt
```

## Deobfuscated Source Evidence Commands
For NPC class:
```
grep -A 10 -B 5 "public final class NPC" srcAllDummysRemoved/src/NPC.java
```

For getRotatedModel:
```
grep -A 10 -B 5 "public Model getRotatedModel" srcAllDummysRemoved/src/NPC.java
```

## Javap Cache Evidence Commands
For class and inheritance:
```
grep -A 10 -B 5 "public final class NPC" srcAllDummysRemoved/.javap_cache/NPC.javap.cache
```

For methods:
```
grep -A 10 -B 5 "public Model getRotatedModel" srcAllDummysRemoved/.javap_cache/NPC.javap.cache
```

## Independent Verification Commands
To demonstrate unique 1:1 correspondence beyond inheritance:
```
# Count classes extending Entity (GQOSZKJC) with a method returning Model (ZKARKDQW)
grep -l "extends.*GQOSZKJC" bytecode/client/*.bytecode.txt | xargs grep -l "ZKARKDQW" | wc -l
```

## Verification of Non-Contradictory Evidence
Bytecode matches source/javap in Entity extension, NPC fields. No contradictions. 1:1 mapping confirmed.