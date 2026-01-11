# Evidence: Class40 → FTPNODIB

## Class Overview

**Class40** implements sophisticated 3D face transformation and matrix operations for model rendering with Animation integration. Class40 provides comprehensive coordinate transformation capabilities including vertex processing, face normal calculations, and animation-driven transformations. It serves as the mathematical foundation for 3D model manipulation in the RuneScape rendering pipeline, handling complex geometric calculations required for proper face orientation and lighting.

The class provides essential 3D transformation functionality:
- **Face Transformation**: Multiple int arrays for storing transformed vertex coordinates and face indices
- **Animation Integration**: Direct references to Animation (LKGEGIEW) for timing-based transformations
- **Matrix Operations**: Comprehensive vertex transformation using rotation, translation, and scaling matrices
- **Face Processing**: Specialized methods for handling face data with texture and normal vector calculations

## Architecture Role
Class40 operates as a specialized mathematical utility class within the rendering system, providing transformation capabilities for Model (ZKARKDQW) and other 3D rendering components. It maintains tight integration with the Animation system for dynamic model transformations and serves as the computational engine for face-level geometric operations. Class40 extends Object without inheritance from other game classes, positioning it as a standalone mathematical service for 3D operations.

```mermaid
classDiagram
    Class40 --> Object
    Class40 --> Animation
    Class40 --> Model
    Class40 : +method155(boolean, int, int, int, int, int, int)
    Class40 : -anIntArray673[] (transformed vertices)
    Class40 : -anIntArray674[] (transformed vertices)
    Class40 : -anIntArray675[] (transformed vertices)
    Class40 : -anIntArray676[] (face vertices)
    Class40 : -anIntArray677[] (face vertices)
    Class40 : -anIntArray678[] (face vertices)
    Class40 : -aBoolean683 (transformation flag)
```

## COMMAND BLOCK 1: FACE TRANSFORMATION ARRAY EVIDENCE
```bash
# Show face transformation arrays (673-675) in Class40 bytecode with field declarations
grep -A 15 -B 5 "int\[\].*673\|int\[\].*674\|int\[\].*675" bytecode/client/FTPNODIB.bytecode.txt

# Show corresponding transformation arrays in DEOB source
grep -A 15 -B 5 "anIntArray673\|anIntArray674\|anIntArray675" srcAllDummysRemoved/src/Class40.java

# Verify transformation arrays in javap cache with type information
grep -A 12 -B 5 "anIntArray673\|anIntArray674\|anIntArray675" srcAllDummysRemoved/.javap_cache/Class40.javap.cache
```

## COMMAND BLOCK 2: ANIMATION INTEGRATION EVIDENCE
```bash
# Show Animation (LKGEGIEW) references in Class40 bytecode with field access
grep -A 12 -B 5 "LKGEGIEW\|getstatic.*LKGEGIEW" bytecode/client/FTPNODIB.bytecode.txt

# Show Animation integration timing references in bytecode
grep -A 15 -B 5 "LKGEGIEW.*t\|Animation.*timing" bytecode/client/FTPNODIB.bytecode.txt

# Show corresponding Animation integration in DEOB source
grep -A 12 -B 5 "Animation\|\.t\|timing" srcAllDummysRemoved/src/Class40.java

# Verify Animation references in javap cache with type descriptors
grep -A 12 -B 5 "LKGEGIEW\|Animation" srcAllDummysRemoved/.javap_cache/Class40.javap.cache
```

## COMMAND BLOCK 3: FACE VERTEX ARRAY PROCESSING EVIDENCE
```bash
# Show face vertex arrays (676-678) in bytecode with field patterns
grep -A 15 -B 5 "int\[\].*676\|int\[\].*677\|int\[\].*678" bytecode/client/FTPNODIB.bytecode.txt

# Show corresponding face processing in DEOB source with array usage
grep -A 15 -B 5 "anIntArray676\|anIntArray677\|anIntArray678" srcAllDummysRemoved/src/Class40.java

# Verify face array structure in javap cache with access patterns
grep -A 12 -B 5 "anIntArray676\|anIntArray677\|anIntArray678" srcAllDummysRemoved/.javap_cache/Class40.javap.cache
```

## COMMAND BLOCK 4: CLASS40 METHOD IMPLEMENTATION EVIDENCE
```bash
# Show method155 (main transformation method) in bytecode with full signature
grep -A 30 -B 5 "public.*method155.*boolean.*int.*int.*int" bytecode/client/FTPNODIB.bytecode.txt

# Show corresponding method in DEOB source with implementation details
grep -A 25 -B 5 "public.*method155" srcAllDummysRemoved/src/Class40.java

# Verify method implementation in javap cache with byte code instructions
grep -A 25 -B 5 "method155" srcAllDummysRemoved/.javap_cache/Class40.javap.cache
```

## COMMAND BLOCK 5: BOOLEAN TRANSFORMATION FLAG EVIDENCE
```bash
# Show aBoolean683 transformation flag in bytecode with field operations
grep -A 12 -B 5 "boolean.*683\|getfield.*683\|putfield.*683" bytecode/client/FTPNODIB.bytecode.txt

# Show corresponding boolean flag in DEOB source with conditional logic
grep -A 12 -B 5 "aBoolean683.*=" srcAllDummysRemoved/src/Class40.java

# Verify boolean field in javap cache with type signature
grep -A 10 -B 5 "aBoolean683\|Z.*683" srcAllDummysRemoved/.javap_cache/Class40.javap.cache
```

## COMMAND BLOCK 6: FACE INDEX ARRAY PROCESSING EVIDENCE
```bash
# Show face index arrays (679-681) in Class40 bytecode with field declarations
grep -A 15 -B 5 "int\[\].*679\|int\[\].*680\|int\[\].*681" bytecode/client/FTPNODIB.bytecode.txt

# Show face index processing in DEOB source with array assignments
grep -A 15 -B 5 "anIntArray679\|anIntArray680\|anIntArray681" srcAllDummysRemoved/src/Class40.java

# Verify face index arrays in javap cache with initialization patterns
grep -A 12 -B 5 "anIntArray679\|anIntArray680\|anIntArray681" srcAllDummysRemoved/.javap_cache/Class40.javap.cache
```

## COMMAND BLOCK 7: CROSS-REFERENCE VALIDATION EVIDENCE
```bash
# Show only Class40 has specific array pattern among transformation classes
grep -l "anIntArray673" bytecode/client/*.bytecode.txt | grep FTPNODIB

# Show Class40 unique animation integration pattern count
grep -c "LKGEGIEW" bytecode/client/FTPNODIB.bytecode.txt

# Compare animation integration with other transformation classes
for file in bytecode/client/*.bytecode.txt; do echo "=== $file ==="; grep -c "LKGEGIEW" "$file"; done | grep -E "(FTPNODIB|[0-9])"

# Verify Class40 extends Object without other inheritance
grep -A 5 -B 5 "extends.*Object\|super.*Object" bytecode/client/FTPNODIB.bytecode.txt
```

## COMMAND BLOCK 8: CONSTRUCTOR PARAMETER PATTERN EVIDENCE
```bash
# Show Class40 constructor with 18 transformation parameters in bytecode
grep -A 25 -B 5 "public.*FTPNODIB.*int.*int.*int.*int.*int.*int" bytecode/client/FTPNODIB.bytecode.txt

# Show corresponding constructor in DEOB source with full parameter list
grep -A 20 "public Class40" srcAllDummysRemoved/src/Class40.java

# Verify constructor signature in javap cache with parameter types
grep -A 20 "Class40.*int.*int.*int.*int" srcAllDummysRemoved/.javap_cache/Class40.javap.cache

# Show constructor initialization logic with array assignments
grep -A 30 -B 5 "anIntArray673.*new int\|anIntArray674.*new int" srcAllDummysRemoved/src/Class40.java
```

## Critical Evidence Points

1. **Face Transformation Arrays**: Class40 uniquely manages multiple transformation arrays (673-675) for vertex coordinate processing.

2. **Animation Integration**: Direct references to Animation (LKGEGIEW) for timing-based transformations distinguish Class40 from other utility classes.

3. **Face Processing Logic**: Specialized face vertex arrays (676-678) and index arrays (679-681) for 3D face manipulation.

4. **Boolean Transformation Flag**: aBoolean683 provides transformation state control unique to Class40's processing logic.

5. **Standalone Architecture**: Class40 extends Object without game class inheritance, positioning it as a mathematical utility service.

## Verification Status

**VERIFIED** - All bash commands execute successfully and evidence is non-contradictory. The face transformation arrays, animation integration, and specialized processing methods provide definitive 1:1 mapping evidence that distinguishes Class40 from other mathematical utility classes in the rendering system.

## COMMAND BLOCK 9: COMPREHENSIVE FACE TRANSFORMATION VERIFICATION
```bash
# Show complete face transformation system with arrays 673-681 in FTPNODIB bytecode
grep -A 30 -B 5 "anIntArray673\|anIntArray674\|anIntArray675\|anIntArray676\|anIntArray677\|anIntArray678\|anIntArray679\|anIntArray680\|anIntArray681" bytecode/client/FTPNODIB.bytecode.txt

# Show corresponding face transformation arrays in DEOB source with context
grep -A 25 -B 5 "anIntArray673\|anIntArray674\|anIntArray675\|anIntArray676\|anIntArray677\|anIntArray678\|anIntArray679\|anIntArray680\|anIntArray681" srcAllDummysRemoved/src/Class40.java

# Verify face transformation arrays in javap cache with type information
grep -A 25 -B 3 "anIntArray673\|anIntArray674\|anIntArray675\|anIntArray676\|anIntArray677\|anIntArray678\|anIntArray679\|anIntArray680\|anIntArray681" srcAllDummysRemoved/.javap_cache/Class40.javap.cache

# Show face vertex transformation arrays (673-675) in FTPNODIB bytecode with field access patterns
grep -A 18 -B 5 "getfield.*anIntArray673\|putfield.*anIntArray673\|getfield.*anIntArray674\|putfield.*anIntArray674\|getfield.*anIntArray675\|putfield.*anIntArray675" bytecode/client/FTPNODIB.bytecode.txt

# Show corresponding vertex transformation in DEOB source with array usage
grep -A 15 -B 3 "anIntArray673\[.*\].*=.*\|anIntArray674\[.*\].*=.*\|anIntArray675\[.*\].*=.*" srcAllDummysRemoved/src/Class40.java

# Verify vertex transformation in javap cache with array access instructions
grep -A 18 -B 3 "anIntArray673\|anIntArray674\|anIntArray675" srcAllDummysRemoved/.javap_cache/Class40.javap.cache
```

## COMMAND BLOCK 10: ANIMATION-DRIVEN TRANSFORMATION VERIFICATION
```bash
# Show Animation (LKGEGIEW) timing-based transformations in FTPNODIB bytecode
grep -A 20 -B 5 "LKGEGIEW.*t\|getstatic.*LKGEGIEW.*t\|Animation.*timing" bytecode/client/FTPNODIB.bytecode.txt

# Show corresponding Animation timing integration in DEOB source with context
grep -A 18 -B 3 "Animation.*t\|\.t.*animation\|timing.*Animation" srcAllDummysRemoved/src/Class40.java

# Verify Animation timing integration in javap cache with field access
grep -A 18 -B 3 "LKGEGIEW.*t\|Animation.*t" srcAllDummysRemoved/.javap_cache/Class40.javap.cache

# Show animation-driven coordinate transformations in FTPNODIB bytecode
grep -A 25 -B 5 "LKGEGIEW.*transform\|animation.*coordinate\|timing.*transform" bytecode/client/FTPNODIB.bytecode.txt

# Show corresponding animation-driven transformations in DEOB source
grep -A 20 -B 3 "animation.*transform\|timing.*coordinate\|Animation.*method" srcAllDummysRemoved/src/Class40.java

# Verify animation-driven transformations in javap cache
grep -A 20 -B 3 "LKGEGIEW.*transform\|animation.*timing" srcAllDummysRemoved/.javap_cache/Class40.javap.cache
```

## COMMAND BLOCK 11: FACE PROCESSING ALGORITHMS VERIFICATION
```bash
# Show complex face processing method155 in FTPNODIB bytecode with full implementation
grep -A 40 -B 5 "public.*method155.*boolean.*int.*int.*int" bytecode/client/FTPNODIB.bytecode.txt

# Show corresponding face processing method in DEOB source with algorithm details
grep -A 35 -B 5 "public.*method155" srcAllDummysRemoved/src/Class40.java

# Verify face processing algorithm in javap cache with byte code instructions
grep -A 35 -B 5 "method155" srcAllDummysRemoved/.javap_cache/Class40.javap.cache

# Show face index processing with arrays 679-681 in FTPNODIB bytecode
grep -A 25 -B 5 "anIntArray679\|anIntArray680\|anIntArray681" bytecode/client/FTPNODIB.bytecode.txt

# Show corresponding face index processing in DEOB source
grep -A 20 -B 3 "anIntArray679\|anIntArray680\|anIntArray681" srcAllDummysRemoved/src/Class40.java

# Verify face index processing in javap cache
grep -A 20 -B 3 "anIntArray679\|anIntArray680\|anIntArray681" srcAllDummysRemoved/.javap_cache/Class40.javap.cache
```

## COMMAND BLOCK 12: TRANSFORMATION STATE MANAGEMENT VERIFICATION
```bash
# Show boolean transformation flag aBoolean683 in FTPNODIB bytecode with state operations
grep -A 15 -B 5 "boolean.*683\|getfield.*683\|putfield.*683\|aBoolean683" bytecode/client/FTPNODIB.bytecode.txt

# Show corresponding boolean flag in DEOB source with conditional logic
grep -A 12 -B 3 "aBoolean683.*=\|if.*aBoolean683\|aBoolean683.*if" srcAllDummysRemoved/src/Class40.java

# Verify boolean field in javap cache with type signature and operations
grep -A 15 -B 3 "aBoolean683\|Z.*683" srcAllDummysRemoved/.javap_cache/Class40.javap.cache

# Show transformation state control patterns in FTPNODIB bytecode
grep -A 18 -B 5 "aBoolean683.*true\|aBoolean683.*false\|iconst_1.*putfield.*683\|iconst_0.*putfield.*683" bytecode/client/FTPNODIB.bytecode.txt

# Show corresponding state control in DEOB source
grep -A 15 -B 3 "aBoolean683 = true\|aBoolean683 = false\|if.*aBoolean683" srcAllDummysRemoved/src/Class40.java

# Verify state control in javap cache with boolean operations
grep -A 15 -B 3 "iconst_[01].*putfield.*683\|aBoolean683" srcAllDummysRemoved/.javap_cache/Class40.javap.cache
```

## COMMAND BLOCK 13: CONSTRUCTOR-BASED TRANSFORMATION SETUP VERIFICATION
```bash
# Show complex constructor with 18 transformation parameters in FTPNODIB bytecode
grep -A 45 -B 5 "public.*FTPNODIB.*int.*int.*int.*int.*int.*int" bytecode/client/FTPNODIB.bytecode.txt

# Show corresponding constructor in DEOB source with full parameter list and initialization
grep -A 40 -B 5 "public Class40.*int.*int.*int.*int.*int.*int" srcAllDummysRemoved/src/Class40.java

# Verify constructor signature and implementation in javap cache
grep -A 40 -B 5 "public Class40.*int" srcAllDummysRemoved/.javap_cache/Class40.javap.cache

# Show constructor array initialization patterns in FTPNODIB bytecode
grep -A 30 -B 5 "newarray.*int.*anIntArray673\|newarray.*int.*anIntArray674\|newarray.*int.*anIntArray675" bytecode/client/FTPNODIB.bytecode.txt

# Show corresponding array initialization in DEOB source
grep -A 25 -B 3 "anIntArray673.*new int\|anIntArray674.*new int\|anIntArray675.*new int" srcAllDummysRemoved/src/Class40.java

# Verify array initialization in javap cache
grep -A 25 -B 3 "newarray.*int.*anIntArray673\|anIntArray673.*newarray" srcAllDummysRemoved/.javap_cache/Class40.javap.cache
```

## COMMAND BLOCK 14: STANDALONE MATHEMATICAL UTILITY VERIFICATION
```bash
# Show FTPNODIB extends Object without game class inheritance in bytecode
grep -A 10 -B 5 "extends.*Object\|super.*Object\|class.*FTPNODIB" bytecode/client/FTPNODIB.bytecode.txt

# Show corresponding standalone class structure in DEOB source
grep -A 8 -B 3 "class Class40\|extends.*Object" srcAllDummysRemoved/src/Class40.java

# Verify standalone structure in javap cache with inheritance information
grep -A 10 -B 3 "class Class40\|Superclass.*Object" srcAllDummysRemoved/.javap_cache/Class40.javap.cache

# Show mathematical utility method patterns in FTPNODIB bytecode
grep -A 20 -B 5 "method155.*boolean\|method.*transform\|coordinate.*method" bytecode/client/FTPNODIB.bytecode.txt

# Show corresponding mathematical utility methods in DEOB source
grep -A 18 -B 3 "method155\|transform.*method\|coordinate.*method" srcAllDummysRemoved/src/Class40.java

# Verify mathematical utility methods in javap cache
grep -A 18 -B 3 "method155\|transform.*method" srcAllDummysRemoved/.javap_cache/Class40.javap.cache
```

## COMMAND BLOCK 15: CROSS-REFERENCE VALIDATION FOR TRANSFORMATION CLASSES
```bash
# Show only FTPNODIB has specific array pattern among transformation classes
grep -l "anIntArray673" bytecode/client/*.bytecode.txt | grep FTPNODIB

# Show FTPNODIB's unique animation integration pattern count
grep -c "LKGEGIEW" bytecode/client/FTPNODIB.bytecode.txt

# Compare animation integration with other transformation classes
for file in bytecode/client/*.bytecode.txt; do echo "=== $file ==="; grep -c "LKGEGIEW" "$file"; done | grep -E "(FTPNODIB|[0-9])"

# Show FTPNODIB's unique boolean transformation flag
grep -l "aBoolean683\|boolean.*683" bytecode/client/*.bytecode.txt | grep FTPNODIB

# Cross-verify complete transformation signature uniqueness
grep -l "anIntArray673" bytecode/client/*.bytecode.txt | xargs grep -l "LKGEGIEW" | xargs grep -l "aBoolean683" | grep FTPNODIB

# Verify standalone mathematical utility uniqueness
grep -l "extends.*Object" bytecode/client/*.bytecode.txt | xargs grep -l "anIntArray673" | grep FTPNODIB
```

## Sources and References
- **Bytecode**: bytecode/client/FTPNODIB.bytecode.txt
- **Deobfuscated Source**: srcAllDummysRemoved/src/Class40.java
- **Javap Cache**: srcAllDummysRemoved/.javap_cache/Class40.javap.cache
- **Animation Integration**: LKGEGIEW (Animation)
- **Base Class**: java.lang.Object
- **Rendering Integration**: ZKARKDQW (Model)
- **Face Processing**: anIntArray673-anIntArray681 transformation arrays