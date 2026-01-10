# client → client

## Overview
The main client class for RuneScape 317. This is the core game client that handles rendering, networking, game logic, and user interface. It extends RSApplet and implements the main game loop.

## Bytecode Evidence Commands
```bash
head -10 ../../client/client.bytecode.txt
grep -A 5 -B 5 "main" ../../client/client.bytecode.txt
```

## Deobfuscated Source Evidence Commands  
```bash
head -15 ../../../srcAllDummysRemoved/src/client.java
grep -A 10 -B 5 "main" ../../../srcAllDummysRemoved/src/client.java
```

## Javap Cache Evidence Commands
```bash
head -15 ../../../srcAllDummysRemoved/.javap_cache/client.javap.cache
```

## Verification
This mapping is trivial as both classes have identical names "client". This represents the main game client class.

