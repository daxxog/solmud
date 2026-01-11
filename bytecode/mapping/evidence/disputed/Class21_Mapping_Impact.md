# Class21 Mapping Ripple Effects Analysis

## Current Problematic Mappings
1. **Class21 → OZKFTHAD** (90% confidence) - INCORRECT
2. **Class11 → LLORVYLP** (100% confidence) - INCORRECT  

## Correct Mappings Based on Forensic Evidence
1. **Class21 → LLORVYLP** (should be 100% confidence) - CORRECT
2. **Class11 → ?** (requires investigation)
3. **OZKFTHAD → ?** (requires investigation)

## Impact Analysis

### Classes Requiring Re-mapping:
1. **Class11** - Currently incorrectly mapped to LLORVYLP
   - Should map to a complex collision detection OG class
   - OZKFTHAD is a candidate but needs verification

2. **OZKFTHAD** - Currently incorrectly mapped to Class21  
   - Should map to a complex DEOB class with multiple methods
   - Class11 is a candidate but needs verification

### Verification Commands for Ripple Effects:

```bash
# Check Class11 complexity (should be complex like collision detection)
grep -c "public.*(" srcAllDummysRemoved/src/Class11.java  # Should be >10
grep -c "int\[\]" srcAllDummysRemoved/src/Class11.java     # Should be 1 (104x104 array)

# Check OZKFTHAD complexity (should be complex with methods)  
grep -c "public.*(" bytecode/client/OZKFTHAD.bytecode.txt  # Should be >1
grep -c "Field " bytecode/client/OZKFTHAD.bytecode.txt     # Should be >5

# Check Class21 simplicity (should be simple data container)
grep -c "public.*(" srcAllDummysRemoved/src/Class21.java   # Should be 1
grep -c "public.*\[\]" srcAllDummysRemoved/src/Class21.java # Should be 1
grep -c "public int" srcAllDummysRemoved/src/Class21.java   # Should be 16

# Check LLORVYLP simplicity (should be simple data container)
grep -c "public.*(" bytecode/client/LLORVYLP.bytecode.txt  # Should be 1  
grep -c "Field " bytecode/client/LLORVYLP.bytecode.txt     # Should be 17
```

## Immediate Action Required
1. **Priority 1**: Fix Class21 → LLORVYLP mapping (irrefutable evidence)
2. **Priority 2**: Investigate Class11 → OZKFTHAD potential mapping
3. **Priority 3**: Verify no other mappings are affected

## Verification Status
✅ Class21 ↔ LLORVYLP: **FORENSIC PROVEN** (1:1 exact match)
❓ Class11 ↔ OZKFTHAD: **NEEDS INVESTIGATION** (both complex classes)
❓ OZKFTHAD correct mapping: **UNKNOWN** (complex class pattern)
❓ Class11 correct mapping: **UNKNOWN** (collision detection pattern)
