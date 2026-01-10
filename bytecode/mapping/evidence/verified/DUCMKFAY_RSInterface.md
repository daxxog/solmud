# DUCMKFAY → RSInterface

## Overview
RSInterface handles user interface components, managing widgets, sprites, and interactions for game menus and HUD.

## Architectural Relationships
RSInterface extends Object, uses Sprite (CXGZMTJK), TextDrawingArea (YXVQXWYR), and MRUNodes (GCPOSBWX) for caching.

```mermaid
classDiagram
    DUCMKFAY --> Object
    DUCMKFAY --> CXGZMTJK[Sprite]
    DUCMKFAY --> YXVQXWYR[TextDrawingArea]
    DUCMKFAY --> GCPOSBWX[MRUNodes]
```

## Bytecode Matching Commands
To show class structure:
```
head -n 30 bytecode/client/DUCMKFAY.bytecode.txt
```

To show method signatures:
```
grep -A 10 -B 5 "public static" bytecode/client/DUCMKFAY.bytecode.txt
```

## Deobfuscated Source Evidence Commands
For RSInterface class:
```
grep -A 25 "public class RSInterface" srcAllDummysRemoved/src/RSInterface.java
```

For static methods:
```
grep -A 10 "public static RSInterface" srcAllDummysRemoved/src/RSInterface.java
```

## Javap Cache Evidence Commands
For class structure:
```
grep -A 20 "public class RSInterface" srcAllDummysRemoved/.javap_cache/RSInterface.javap.cache
```

For methods:
```
grep -A 5 "addButton" srcAllDummysRemoved/.javap_cache/RSInterface.javap.cache
```

## Verification of Non-Contradictory Evidence
Bytecode matches source/javap in widget management, static arrays. No contradictions. 1:1 mapping confirmed.