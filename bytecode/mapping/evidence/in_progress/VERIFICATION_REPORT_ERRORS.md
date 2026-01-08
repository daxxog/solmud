# FORENSIC MAPPING VERIFICATION REPORT

## **🔍 VERIFICATION SUMMARY**

**Date**: January 8, 2026
**Purpose**: Double-check forensic claims before updating Go tooling
**Result**: **CRITICAL ERRORS FOUND** - Multiple incorrect mappings identified

## **❌ CRITICAL ERRORS IDENTIFIED**

### **ERROR 1: NodeCache Mapping (WRONG)**
- **My Claim**: `NodeCache → LHGXPZPG` (error code "91809")
- **Verification Result**: **FALSE**
- **Correct Mapping**: `NodeCache → ARZPHHDH` (error code "91499")

**Evidence**:
```bash
$ grep -l "91499" *.bytecode.txt
ARZPHHDH.bytecode.txt  # TRUE NodeCache

$ grep -l "91809" *.bytecode.txt  
LHGXPZPG.bytecode.txt  # DIFFERENT class
```

**Impact**: I incorrectly changed a verified correct mapping from previous session
**Error Code Mismatch**:
- LHGXPZPG: "91809"
- ARZPHHDH: "91499" (matches NodeCache.java line 49)
- NodeCache.java: `signlink.reporterror("91499, " + node + ", " + l + ", " + (byte)7 + ", " + runtimeexception.toString());`

### **ERROR 2: Class32 Structure Mismatch (WRONG)**
- **My Claim**: VBAXKVMG has 19 arrays matching Class32 structure
- **Verification Result**: **FALSE** - Complete structure mismatch

**Class32.java Actual Structure**:
```java
// Constructor initializes these arrays:
anIntArray583 = new int[256];      // array a
anIntArray585 = new int[257];      // array b
aBooleanArray589 = new boolean[256];
aBooleanArray590 = new boolean[16];
aByteArray591 = new byte[256];
aByteArray592 = new byte[4096];
anIntArray593 = new int[16];
aByteArray594 = new byte[18002];
aByteArray595 = new byte[18002];
aByteArrayArray596 = new byte[6][258];
anIntArrayArray597 = new int[6][258];
anIntArrayArray598 = new int[6][258];
anIntArrayArray599 = new int[6][258];
anIntArray600 = new int[6];
```

**Expected Unique Identifiers**:
- Array sizes: 256, 257, 258, 6, 16, 4096, 18002
- Multidimensional arrays: `[6][258]` pattern (4 occurrences)
- Field names: anIntArray583, anIntArray585, etc.

**VBAXKVMG Analysis Status**: Structure does NOT match Class32.java
**Impact**: VBAXKVMG is likely a different compression/utility class

### **ERROR 3: StreamLoader Magic Constants (OVERSTATED)**
- **My Claim**: Magic constants 44820, -29508, 891 provide irrefutable evidence
- **Verification Result**: **PARTIALLY CORRECT** - Hash algorithm confirmed, constants purpose unclear

**Verified Evidence**:
```java
// StreamLoader.java line 46:
i = (i * 61 + s.charAt(j)) - 32;  // CORRECT - base-61 hash algorithm
```

**Bytecode Verification**:
```
21: bipush 61          // hash *= 61
31: bipush 32          // - 32 offset
```

**Constants Status**:
- 44820: Present in bytecode but purpose unclear
- -29508: Present but not in source code
- 891: Present in bytecode but purpose unclear

**Impact**: Hash algorithm is correct evidence, but magic constants are overstated

## **✅ VERIFIED CORRECT MAPPINGS**

### **Mapping 1: Skills → YUXCUCXD (CORRECT)**
**Status**: **VERIFIED CORRECT**

**Evidence Match**:
```java
// Skills.java line 10-13:
public static final String[] skillNames = {
    "attack", "defence", "strength", "hitpoints", "ranged", "prayer", "magic", 
    "cooking", "woodcutting", "fletching", "fishing", "firemaking", "crafting", 
    "smithing", "mining", "herblore", "agility", "thieving", "slayer", 
    "farming", "runecraft", "-unused-", "-unused-", "-unused-", "-unused-"
};
```

**Bytecode Verification**:
```
16: ldc #3    // String "attack"
20: ldc #20   // String "defence"
... (all 21 skill names present)
```

**Confidence**: **100%** - Literal strings are irrefutable

### **Mapping 2: TextClass → ZTQFNQRH (CORRECT)**
**Status**: **VERIFIED CORRECT**

**Evidence Match**:
```java
// TextClass.java line 15:
l *= 37L;  // Base-37 multiplication
```

**Bytecode Verification**:
```
22: ldc2_w #44  // long 37l
25: lmul        // hash *= 37
```

**Additional Verification**:
- Character ranges: A-Z (65-90), a-z (97-122), 0-9 (48-57) ✅
- Invalid name boundary: `0x5b5b57f8a98a5dd1L` = 6582952005840035281L ✅
- Error code: "81570" in source line 50 ✅

**Confidence**: **100%** - Cryptographic algorithm is irrefutable

### **Mapping 3: StreamLoader → XTGLDHGX (MOSTLY CORRECT)**
**Status**: **VERIFIED CORRECT** (with caveats)

**Evidence Match**:
```java
// StreamLoader.java line 46:
i = (i * 61 + s.charAt(j)) - 32;  // Base-61 hash algorithm
```

**Bytecode Verification**:
```
239: bipush 61   // hash *= 61
245: bipush 32   // - 32 offset
```

**Array Structure**:
```java
// StreamLoader.java lines 25-28:
anIntArray728 = new int[dataSize];  // hashes
anIntArray729 = new int[dataSize];  // sizes
anIntArray730 = new int[dataSize];  // compressed sizes
anIntArray731 = new int[dataSize];  // offsets
```

**Confidence**: **90%** - Hash algorithm is correct evidence, magic constants overstated

## **📊 VERIFICATION SUMMARY**

| Mapping | Status | Confidence | Issues |
|----------|---------|-------------|---------|
| Skills → YUXCUCXD | ✅ CORRECT | 100% | None |
| TextClass → ZTQFNQRH | ✅ CORRECT | 100% | None |
| StreamLoader → XTGLDHGX | ✅ CORRECT | 90% | Magic constants overstated |
| NodeCache → LHGXPZPG | ❌ WRONG | N/A | Wrong class, corrected mapping changed |
| Class32 → VBAXKVMG | ❌ WRONG | N/A | Structure mismatch, wrong class |

## **🔧 CORRECTIVE ACTIONS REQUIRED**

### **IMMEDIATE (Critical)**
1. **REVERT NodeCache mapping**: Change back to `NodeCache → ARZPHHDH`
2. **REMOVE Class32 mapping**: Delete `Class32 → VBAXKVMG` from tool
3. **UPDATE evidence documentation**: Correct wrong claims in evidence files

### **REQUIRED (Documentation)**
1. **CORRECT LHGXPZPG evidence**: Identify what LHGXPZPG actually is
2. **FIND true Class32**: Search for arrays with sizes 256, 257, 258, 6, 16, 4096, 18002
3. **STREAMLOADER refinement**: Clarify magic constants purpose

### **RECOMMENDED (Process)**
1. **VALIDATION CHECKLIST**: Create systematic verification for future mappings
2. **ERROR CODE DATABASE**: Build lookup table for error code patterns
3. **STRUCTURE FINGERPRINTING**: Better array size/structure matching

## **🎯 UPDATED TOOLING PLAN**

**Valid Mappings for Go Tool**:
```go
"Skills":         "YUXCUCXD", // VERIFIED: literal skill names
"TextClass":      "ZTQFNQRH", // VERIFIED: base-37 hashing algorithm  
"StreamLoader":   "XTGLDHGX", // VERIFIED: base-61 hash algorithm
// REMOVE:
// "NodeCache":   "LHGXPZPG", // WRONG - revert to ARZPHHDH
// "Class32":     "VBAXKVMG", // WRONG - structure mismatch
```

**Status**: Requires forensic correction before tool update</content>
<parameter name="filePath">/Users/daxxog/Desktop/solmud/bytecode/mapping/evidence/VERIFICATION_REPORT_ERRORS.md