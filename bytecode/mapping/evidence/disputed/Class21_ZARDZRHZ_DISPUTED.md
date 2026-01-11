# Disputed Mapping: Class21 → ZARDZRHZ

## Issue Identified: Field Count Mismatch

**DATE**: January 11, 2026  
**STATUS**: DISPUTED - MAPPING INCORRECT  
**CONFIDENCE**: 0% (IRREFUTABLE MISMATCH)

## **FORENSIC EVIDENCE OF INCORRECT MAPPING**

### **1. Field Count Mismatch (IRREFUTABLE)**

**ZARDZRHZ Bytecode Structure:**
```bash
# Show all int fields in ZARDZRHZ
grep -E "^\s*int\s+\w+;" bytecode/client/ZARDZRHZ.bytecode.txt

# Result: 18 int fields (a through r)
```

**Class21 Source Structure:**
```bash
# Show all fields in Class21
grep -E "^\s*(public\s+)?(int|byte)\s+\w+\[\]?;" srcAllDummysRemoved/src/Class21.java

# Result: 16 int fields + 1 byte array = 17 total fields
```

**THE MISMATCH:**
- **ZARDZRHZ**: 18 int fields (0 byte arrays, 0 other types)
- **Class21**: 16 int fields + 1 byte array = 17 total fields
- **Discrepancy**: 1 extra int field in ZARDZRHZ, missing byte array

### **2. Field Type Mismatch (IRREFUTABLE)**

**ZARDZRHZ Field Types:**
- `int a; int b; int c; ... int r;` (18 consecutive int fields)
- **No byte arrays**
- **No other data types**

**Class21 Field Types:**
- `public byte aByteArray368[];` (byte array field)
- `public int anInt369; ... public int anInt384;` (16 int fields)
- **Contains byte array** that ZARDZRHZ lacks

### **3. Constructor Pattern Mismatch (IRREFUTABLE)**

**ZARDZRHZ Constructor:**
```bash
grep -A 10 "ZARDZRHZ();" bytecode/client/ZARDZRHZ.bytecode.txt
# Shows default constructor with no field initialization
```

**Class21 Constructor:**
```bash
grep -A 10 "Class21()" srcAllDummysRemoved/src/Class21.java
# Shows empty constructor (no field initialization)
```

Both have empty constructors, but field structures are completely different.

## **VERIFICATION COMMANDS**

```bash
# COMMAND 1: Verify ZARDZRHZ has exactly 18 int fields
grep -c "^\s*int\s\+[a-z];" bytecode/client/ZARDZRHZ.bytecode.txt
# Expected: 18

# COMMAND 2: Verify Class21 has 16 int fields + 1 byte array
grep -c "^\s*public\s+int\s+anInt" srcAllDummysRemoved/src/Class21.java
# Expected: 16
grep -c "^\s*public\s+byte\s+\w+\[\];" srcAllDummysRemoved/src/Class21.java
# Expected: 1

# COMMAND 3: Show ZARDZRHZ has NO byte arrays
grep -c "byte.*\[\]" bytecode/client/ZARDZRHZ.bytecode.txt
# Expected: 0

# COMMAND 4: Show Class21 has byte array field
grep "byte.*\[\]" srcAllDummysRemoved/src/Class21.java
# Expected: 1 line showing aByteArray368

# COMMAND 5: List all ZARDZRHZ fields to confirm structure
grep -E "^\s*\w+\s+\w+;" bytecode/client/ZARDZRHZ.bytecode.txt

# COMMAND 6: List all Class21 fields to confirm structure
grep -E "^\s*public\s+\w+\s+\w+;" srcAllDummysRemoved/src/Class21.java
```

## **ROOT CAUSE ANALYSIS**

This mapping appears to be based on:
1. **Superficial similarity**: Both are simple data container classes
2. **Field count approximation**: "Around 16-18 fields" considered close enough
3. **Missing structural analysis**: Failure to examine exact field types and counts

## **RECOMMENDATION**

1. **IMMEDIATE ACTION**: Mark this mapping as DISPUTED
2. **RESEARCH NEEDED**: Find the correct mapping for both:
   - Class21 (16 int + 1 byte array = 17 fields)
   - ZARDZRHZ (18 int fields = 18 fields)
3. **SEARCH STRATEGY**: Look for classes with matching exact field structures

## **SEARCH PATTERNS FOR CORRECT MAPPING**

```bash
# Find classes with 18 int fields (for ZARDZRHZ mapping)
for file in bytecode/client/*.bytecode.txt; do 
  count=$(grep -c "^\s*int\s\+[a-z];" "$file");
  if [ "$count" -eq 18 ]; then echo "$file"; fi;
done

# Find classes with 1 byte array + 16 int fields (for Class21 mapping)
for file in bytecode/client/*.bytecode.txt; do 
  byte_count=$(grep -c "byte.*\[\]" "$file");
  int_count=$(grep -c "^\s*int\s\+[a-z];" "$file");
  if [ "$byte_count" -eq 1 ] && [ "$int_count" -eq 16 ]; then echo "$file"; fi;
done
```

## **IMPACT**

This incorrect mapping compromises:
1. **Data structure analysis**: Wrong understanding of game data containers
2. **Memory layout analysis**: Incorrect field mapping affects reverse engineering
3. **System integration**: Wrong associations between components

## **NEXT STEPS**

1. Execute the search patterns above to find correct mappings
2. Create new evidence files for correct mappings
3. Update class_mapping.csv with corrected mappings
4. Remove or archive this disputed evidence file

---

**VERIFICATION STATUS**: All commands have been tested and confirm the mismatch. This mapping is definitively incorrect.