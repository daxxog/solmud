# DYMVKFXP → Ground

## Overview
Ground represents ground tiles in the 3D world, storing terrain properties and object placements.

## Architectural Relationships
Ground extends Object, references WorldController (NYFUGYQS) for scene integration.

```mermaid
classDiagram
    DYMVKFXP --> Object
    DYMVKFXP --> NYFUGYQS[WorldController]
```

## Bytecode Matching Commands
To show class structure:
```
head -n 25 bytecode/client/DYMVKFXP.bytecode.txt
```

To show field assignments:
```
grep -A 10 -B 5 "putfield" bytecode/client/DYMVKFXP.bytecode.txt
```

## Deobfuscated Source Evidence Commands
For Ground class:
```
grep -A 10 -B 5 "public final class Ground" srcAllDummysRemoved/src/Ground.java
```

For constructor:
```
grep -A 15 -B 5 "public Ground" srcAllDummysRemoved/src/Ground.java
```

## Javap Cache Evidence Commands
For class structure:
```
grep -A 10 -B 5 "public final class Ground" srcAllDummysRemoved/.javap_cache/Ground.javap.cache
```

For fields:
```
grep -A 5 -B 5 "Object" srcAllDummysRemoved/.javap_cache/Ground.javap.cache
```

## Verification of Non-Contradictory Evidence
Bytecode matches source/javap in terrain data storage. No contradictions. 1:1 mapping confirmed.