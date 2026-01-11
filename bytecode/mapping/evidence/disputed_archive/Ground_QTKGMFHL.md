# Dispute: Class30_Sub1 → QTKGMFHL

## Issue
The mapping shows contradictory evidence between the obfuscated bytecode and deobfuscated source.

### Bytecode Evidence (QTKGMFHL)
- Constructor: Parametrized constructor(int, int, int) with sipush -589, iconst_5 anewarray OPNPFUJE, newarray int, and field assignments
- Fields: Many fields including arrays, objects, booleans

### Deobfuscated Source Evidence (Class30_Sub1)
- Constructor: Default constructor that sets anInt1294 to -1
- Fields: Exactly 12 public int fields, no other types

### Contradiction
The constructor is completely different. The field count and types do not match. QTKGMFHL has arrays and objects, Class30_Sub1 has only ints.

## Recommendation
This mapping should be disputed and re-investigated. The class QTKGMFHL does not match the structure of Class30_Sub1.

## Evidence Commands
To verify the contradiction:
```
grep -A 30 "QTKGMFHL(int" bytecode/client/QTKGMFHL.bytecode.txt
grep -A 10 "Class30_Sub1();" srcAllDummysRemoved/.javap_cache/Class30_Sub1.javap.cache
```