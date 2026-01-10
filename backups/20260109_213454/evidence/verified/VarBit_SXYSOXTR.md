# Evidence for VarBit -> SXYSOXTR

## Overview
VarBit handles variable bit configurations for game settings. Its core functionality includes bit manipulation for storing multiple boolean flags in integers.

## Architectural Relationships
VarBit manages configuration data, used by the client for various settings.

```mermaid
classDiagram
    VarBit --> client
```

## Bash Commands Proving Bytecode Matches
- `grep -A 20 -B 5 "public SXYSOXTR(" bytecode/client/SXYSOXTR.bytecode.txt`  
  Shows constructor with field assignments for bit settings.

- `grep -A 20 -B 5 "method552" bytecode/client/SXYSOXTR.bytecode.txt`  
  Shows bit manipulation method.

## Commands for Deob Source Sections
- `head -20 srcAllDummysRemoved/src/VarBit.java`  
  Shows class with bit fields.

- `grep -A 15 -B 5 "method552" srcAllDummysRemoved/src/VarBit.java`  
  Shows the bit method.

## Commands for Javap Cache Sections
- `grep -A 20 -B 5 "public VarBit(" srcAllDummysRemoved/.javap_cache/VarBit.javap.cache`  
  Shows constructor.

- `grep -A 20 -B 5 "public int method552" srcAllDummysRemoved/.javap_cache/VarBit.javap.cache`  
  Shows method.

## Verification of Non-Contradictory Evidence
Matches exactly. No contradictions.

## 1:1 Mapping Confirmation
SXYSOXTR.bytecode.txt maps uniquely to VarBit.java by the bit manipulation fields and methods.