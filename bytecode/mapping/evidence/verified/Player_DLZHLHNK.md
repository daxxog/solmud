# DLZHLHNK → Player

## Overview
Player manages player character data, extending Entity for base functionality, with fields for equipment, skills, and appearance.

## Architectural Relationships
Player extends Entity (GQOSZKJC), uses EntityDef (CKDEJADD), Model (ZKARKDQW), MRUNodes (GCPOSBWX), and Animation (LKGEGIEW).

```mermaid
classDiagram
    DLZHLHNK --> GQOSZKJC[Entity]
    DLZHLHNK --> CKDEJADD[EntityDef]
    DLZHLHNK --> ZKARKDQW[Model]
    DLZHLHNK --> GCPOSBWX[MRUNodes]
    DLZHLHNK --> LKGEGIEW[Animation]
```

## Bytecode Matching Commands
To show inheritance and fields:
```
head -n 25 bytecode/client/DLZHLHNK.bytecode.txt
```

To show method signatures:
```
grep -A 10 -B 5 "public.*ZKARKDQW" bytecode/client/DLZHLHNK.bytecode.txt
```

## Deobfuscated Source Evidence Commands
For Player class:
```
grep -A 20 "public final class Player" srcAllDummysRemoved/src/Player.java
```

For getRotatedModel:
```
grep -A 10 "public Model getRotatedModel" srcAllDummysRemoved/src/Player.java
```

## Javap Cache Evidence Commands
For class and inheritance:
```
grep -A 15 "public final class Player" srcAllDummysRemoved/.javap_cache/Player.javap.cache
```

For methods:
```
grep -A 5 "getRotatedModel" srcAllDummysRemoved/.javap_cache/Player.javap.cache
```

## Verification of Non-Contradictory Evidence
Bytecode matches source/javap in Entity extension, Player fields. No contradictions. 1:1 mapping confirmed.