# SpotAnim_MUDLUUBC

## Overview
SpotAnim manages spot animations, loading from cache and handling animation frames for effects in RuneScape.

## Architectural Relationships
SpotAnim uses Stream for data. Part of animation system. No diagram.

## Bytecode Matching Commands
To show unpackConfig:
```
grep -A 15 "public static void a" bytecode/client/MUDLUUBC.bytecode.txt
```

## Deobfuscated Source Sections
```
grep -A 15 "public static void unpackConfig" srcAllDummysRemoved/src/SpotAnim.java
```

## Javap Cache Sections
```
grep -A 15 "public static void unpackConfig" srcAllDummysRemoved/.javap_cache/SpotAnim.javap.cache
```

## Verification of Non-Contradictory Evidence
Matches on cache loading. No contradictions. 1:1 confirmed.