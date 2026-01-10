# Stream_MBMGIXGO

## Overview
Stream handles data buffering and reading for RuneScape, providing methods for reading various data types from byte arrays.

## Architectural Relationships
Core utility for data streams. No diagram.

## Bytecode Matching Commands
To show constructor:
```
grep -A 10 "public MBMGIXGO" bytecode/client/MBMGIXGO.bytecode.txt
```

## Deobfuscated Source Sections
```
grep -A 10 "public Stream" srcAllDummysRemoved/src/Stream.java
```

## Javap Cache Sections
```
grep -A 10 "public Stream" srcAllDummysRemoved/.javap_cache/Stream.javap.cache
```

## Verification of Non-Contradictory Evidence
Aligns on buffer access. No contradictions. 1:1 confirmed. Note: This mapping is disputed due to potential overlap with Stream.md (same class, conflicting evidence). Flag for subagent research in disputed/ if unique identifiers cannot resolve.