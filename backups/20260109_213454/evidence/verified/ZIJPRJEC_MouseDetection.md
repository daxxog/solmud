# ZIJPRJEC → MouseDetection

## Overview
MouseDetection implements Runnable for background mouse coordinate tracking with synchronization.

## Architectural Relationships
MouseDetection implements Runnable, uses client for mouse access, references RSApplet (KHACHIFW).

```mermaid
classDiagram
    ZIJPRJEC ..|> Runnable
    ZIJPRJEC --> client
    ZIJPRJEC --> KHACHIFW[RSApplet]
    ZIJPRJEC : +run()
```

## Bytecode Matching Commands
To show Runnable implementation:
```
grep -A 30 "public void run" bytecode/client/ZIJPRJEC.bytecode.txt
```

To show synchronization:
```
grep -A 10 -B 5 "monitorenter" bytecode/client/ZIJPRJEC.bytecode.txt
```

## Deobfuscated Source Evidence Commands
For MouseDetection class:
```
grep -A 15 "public class MouseDetection" srcAllDummysRemoved/src/MouseDetection.java
```

For run method:
```
grep -A 10 "public void run" srcAllDummysRemoved/src/MouseDetection.java
```

## Javap Cache Evidence Commands
For class structure:
```
grep -A 15 "public class MouseDetection" srcAllDummysRemoved/.javap_cache/MouseDetection.javap.cache
```

For methods:
```
grep -A 5 "run" srcAllDummysRemoved/.javap_cache/MouseDetection.javap.cache
```

## Verification of Non-Contradictory Evidence
Bytecode matches source/javap in Runnable, coordinate arrays, synchronization. No contradictions. 1:1 mapping confirmed.