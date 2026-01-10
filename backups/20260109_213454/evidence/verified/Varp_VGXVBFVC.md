# Varp_VGXVBFVC

## Overview

The Varp class manages variable player parameters in RuneScape, such as settings or quest states. It unpacks configuration data from "varp.dat" using a StreamLoader, stores an array of Varp instances, and reads values from streams to populate fields like anInt709 and aBoolean713.

## Architectural Relationships

Varp is a data container class that interacts with Stream for data parsing and StreamLoader for resource loading. It is likely used by the client to handle dynamic player variables. No mermaid diagram needed as it has no complex class relationships.

## Bytecode Matching Commands

To show the unpackConfig method signature and initialization:

```
grep -A 15 "public static void a.*XTGLDHGX" bytecode/client/VGXVBFVC.bytecode.txt
```

To show the readValues method and loop structure:

```
grep -A 25 "public void a.*MBMGIXGO.*boolean.*int" bytecode/client/VGXVBFVC.bytecode.txt
```

## Deobfuscated Source Sections

For unpackConfig method:

```
grep -A 25 "public static void unpackConfig" srcAllDummysRemoved/src/Varp.java
```

For readValues method:

```
grep -A 35 "private void readValues" srcAllDummysRemoved/src/Varp.java
```

## Javap Cache Sections

For unpackConfig:

```
grep -A 25 "public static void unpackConfig" srcAllDummysRemoved/.javap_cache/Varp.javap.cache
```

For readValues:

```
grep -A 35 "private void readValues" srcAllDummysRemoved/.javap_cache/Varp.javap.cache
```

## Verification of Non-Contradictory Evidence

Bytecode, javap cache, and source match exactly in method signatures, field access (e.g., e for anInt702, f for anIntArray703), and control flow. No contradictions. 1:1 mapping confirmed via class_mapping.csv.