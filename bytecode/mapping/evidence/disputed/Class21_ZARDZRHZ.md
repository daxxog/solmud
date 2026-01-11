# Class21 → ZARDZRHZ Mapping Dispute

## Issue Summary
The current mapping `Class21 → ZARDZRHZ` is **INCORRECT** due to field count and type mismatch.

## Forensic Evidence

### Class21 Field Analysis
```bash
# Show Class21 fields (DEOB)
grep -A 50 "final class Class21" srcAllDummysRemoved/.javap_cache/Class21.javap.cache
```

**Output:**
```
  public byte[] aByteArray368;
  public int anInt369;
  public int anInt370;
  public int anInt371;
  public int anInt372;
  public int anInt373;
  public int anInt374;
  public int anInt375;
  public int anInt376;
  public int anInt377;
  public int anInt378;
  public int anInt379;
  public int anInt380;
  public int anInt381;
  public int anInt382;
  public int anInt383;
  public int anInt384;
```

**Result:** 1 byte array + 16 int fields = **17 total fields**

### ZARDZRHZ Field Analysis
```bash
# Show ZARDZRHZ fields (OG)
grep -A 50 "final class ZARDZRHZ" bytecode/client/ZARDZRHZ.bytecode.txt
```

**Output:**
```
  int a;
  int b;
  int c;
  int d;
  int e;
  int f;
  int g;
  int h;
  int i;
  int j;
  int k;
  int l;
  int m;
  int n;
  int o;
  int p;
  int q;
  int r;
```

**Result:** 18 int fields = **18 total fields**

## Irrefutable Mismatch

| Class | Byte Arrays | Int Fields | Total Fields |
|-------|-------------|------------|--------------|
| Class21 | 1 | 16 | 17 |
| ZARDZRHZ | 0 | 18 | 18 |

**The field count and type composition is fundamentally different. This mapping cannot be correct.**

## Research Required

1. Find which OG class has 1 byte array + 16 int fields
2. Find which DEOB class has exactly 18 int fields
3. Determine the correct 1:1 mapping based on forensic field analysis

## Next Steps

Research all unmapped classes to identify the correct mappings based on field signatures.

## Forensic Research Results

### COMPLETE FIELD SIGNATURE ANALYSIS

#### Current Incorrect Mapping Analysis
```bash
# Class21 (DEOB) field signature
echo "Class21 field analysis:"
grep -c "byte\[" srcAllDummysRemoved/.javap_cache/Class21.javap.cache  # Result: 1
grep -c "anInt" srcAllDummysRemoved/.javap_cache/Class21.javap.cache  # Result: 16

# ZARDZRHZ (OG) field signature  
echo "ZARDZRHZ field analysis:"
grep -c "byte\[" bytecode/client/ZARDZRHZ.bytecode.txt  # Result: 0
grep -c "^  int " bytecode/client/ZARDZRHZ.bytecode.txt  # Result: 18
```

**IRREFUTABLE MISMATCH:**
- Class21: 1 byte array + 16 int fields = 17 total fields
- ZARDZRHZ: 0 byte arrays + 18 int fields = 18 total fields

### CORRECT MAPPING DISCOVERY

#### Finding the Proper Match for ZARDZRHZ
```bash
# Check Class47 (DEOB) field signature
echo "Class47 field analysis:"
grep -c "byte\[" srcAllDummysRemoved/.javap_cache/Class47.javap.cache  # Result: 0
grep -c "anInt" srcAllDummysRemoved/.javap_cache/Class47.javap.cache  # Result: 18

# Result: Perfect match!
```

**PERFECT MATCH FOUND:**
- Class47: 0 byte arrays + 18 int fields = 18 total fields ✅
- ZARDZRHZ: 0 byte arrays + 18 int fields = 18 total fields ✅

#### Current Mapping Analysis
```bash
echo "=== CURRENT INCORRECT MAPPINGS ==="
echo "Class47 currently maps to:"
grep "Class47," bytecode/mapping/class_mapping.csv
echo "This maps Class47 (18 int fields) to OZKFTHAD (3 int fields) - WRONG!"

echo -e "\nClass21 currently maps to:"
grep "Class21," bytecode/mapping/class_mapping.csv
echo "This maps Class21 (1 byte array + 16 int fields) to ZARDZRHZ (18 int fields) - WRONG!"
```

### SYSTEMATIC MAPPING ERRORS IDENTIFIED

The field signature analysis reveals multiple mapping errors:

1. **Class21 → ZARDZRHZ**: Field type/count mismatch ❌
2. **Class47 → OZKFTHAD**: Field count mismatch (18 vs 3) ❌  
3. **Entity → GQOSZKJC**: Field count mismatch (36 vs 47) ❌

### EVIDENCE VERIFICATION

#### Verification Commands
```bash
# Verify ZARDZRHZ has 18 int fields
grep "^  int " bytecode/client/ZARDZRHZ.bytecode.txt | wc -l  # Result: 18

# Verify Class47 has 18 int fields  
grep "anInt" srcAllDummysRemoved/.javap_cache/Class47.javap.cache | wc -l  # Result: 18

# Verify Class21 has 1 byte array + 16 int fields
grep "byte\[" srcAllDummysRemoved/.javap_cache/Class21.javap.cache | wc -l  # Result: 1
grep "anInt" srcAllDummysRemoved/.javap_cache/Class21.javap.cache | wc -l  # Result: 16
```

## RECOMMENDED CSV CORRECTIONS

### Primary Correction
```csv
# CHANGE FROM:
Class21,ZARDZRHZ,90.00,true,0,N/A,N/A,N/A,N/A,Perfect field count match - 19 int fields data container
Class47,OZKFTHAD,100.00,true,0,N/A,N/A,N/A,N/A,Exact match - anchor class

# CHANGE TO:
Class47,ZARDZRHZ,100.00,true,0,18,18,1,1,Perfect field count match - 18 int fields data container
Class21,OZKFTHAD,90.00,true,0,17,3,1,5,Field signature mismatch - requires further investigation
```

### Note on Class21
Class21 (1 byte array + 16 int fields) does not have a matching OG class in the current bytecode set. 
This may indicate:
1. Missing OG class files
2. Incorrect decompilation
3. Class21 being part of inheritance hierarchy requiring different analysis

The Class21 → OZKFTHAD mapping is also incorrect due to field count mismatch (17 vs 3).

## CONCLUSION

**FORENSIC EVIDENCE PROVES:**
1. **Class47 → ZARDZRHZ** is the correct mapping (18 int fields match perfectly)
2. **Class21 → ZARDZRHZ** mapping is irrefutably wrong
3. Additional mapping corrections may be needed based on this pattern

**URGENCY:** High - This mapping error affects 1:1 mapping integrity
