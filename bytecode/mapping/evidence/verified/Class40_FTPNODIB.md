# FTPNODIB → Class40

## Overview
Class40 handles 3D face transformations and matrices for model rendering, with animation integration.

## Architectural Relationships
Class40 extends Object, references Animation (LKGEGIEW) for timing.

```mermaid
classDiagram
    FTPNODIB --> Object
    FTPNODIB --> LKGEGIEW[Animation]
```

## Bytecode Matching Commands
To show matrix operations:
```
tail -n +28 bytecode/client/FTPNODIB.bytecode.txt | head -n 30
```

To show animation refs:
```
grep -A 5 -B 5 "LKGEGIEW" bytecode/client/FTPNODIB.bytecode.txt
```

## Deobfuscated Source Evidence Commands
For Class40 class:
```
grep -A 15 "public final class Class40" srcAllDummysRemoved/src/Class40.java
```

For transformation methods:
```
grep -A 10 "public static void method155" srcAllDummysRemoved/src/Class40.java
```

## Javap Cache Evidence Commands
For class structure:
```
grep -A 10 "public final class Class40" srcAllDummysRemoved/.javap_cache/Class40.javap.cache
```

For methods:
```
grep -A 5 "method155" srcAllDummysRemoved/.javap_cache/Class40.javap.cache
```

## Verification of Non-Contradictory Evidence
Bytecode matches source/javap in matrix calculations, animation refs. No contradictions. 1:1 mapping confirmed.