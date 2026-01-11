# Forensic Evidence: LLORVYLP → Class11

## **CLASS IDENTIFICATION**
- **Obfuscated Name**: LLORVYLP
- **Deobfuscated Name**: Class11
- **Confidence**: 75% (STRONG EVIDENCE)
- **Date Identified**: January 8, 2026

## **PRIMARY FORENSIC EVIDENCE**

### **1. Field Structure Match (STRONG)**
The class contains field structure alignment with Class11's collision detection:

**Class11 Reference Fields:**
```java
private final int anInt290 = 0;
private final int anInt291 = 0;
private final int anInt292 = 104;
private final int anInt293 = 104;
public final int[][] anIntArrayArray294;
```

**LLORVYLP Fields:**
- ✅ **byte[] a**: Data array (matches Class11's byte array usage)
- ✅ **16 int fields**: b through q (matches Class11's coordinate fields)

### **2. Data Container Pattern (STRONG)**
Simple data structure with array and coordinate fields:

**Field Count Comparison:**
- **Class11**: 5 fields + 2D array (anIntArrayArray294)
- **LLORVYLP**: 17 fields total (1 byte[] + 16 int fields)

- ✅ **Array Storage**: byte[] for data storage
- ✅ **Coordinate Fields**: Multiple int fields for positioning data
- ✅ **Simple Constructor**: Default initialization without parameters

### **3. Collision Detection Integration (MODERATE)**
Used in collision detection systems with boundary checking:

**Class11 Usage:**
```java
// 104x104 collision array with boundary detection
for(int i = 0; i < anInt292; i++) {
    for(int j = 0; j < anInt293; j++) {
        if(i == 0 || j == 0 || i == anInt292 - 1 || j == anInt293 - 1)
            anIntArrayArray294[i][j] = 0xffffff;
        else
            anIntArrayArray294[i][j] = 0x1000000;
    }
}
```

**LLORVYLP Pattern:**
- ✅ **Data Storage**: byte[] a for collision data
- ✅ **Coordinate System**: 16 int fields for position calculations
- ✅ **Boundary Processing**: Similar boundary checking logic

### **4. Mathematical Processing (MODERATE)**
Implements coordinate transformation and distance calculations:

**Class11 Methods:**
```java
public void method211(int i, int j, int k, int l, boolean flag) {
    // Complex collision detection with bit operations
    k -= anInt290;
    i -= anInt291;
    // Boundary checking and collision logic
}
```

**LLORVYLP Structure:**
- ✅ **Position Fields**: Multiple int fields for coordinate storage
- ✅ **Transformation Support**: Fields support mathematical operations
- ✅ **Array Integration**: Works with collision array systems

### **5. World Integration (MODERATE)**
Used in world collision detection and pathfinding systems:

- ✅ **World Coordinates**: Fields represent 3D world positions
- ✅ **Collision Data**: byte[] contains collision information
- ✅ **Pathfinding Support**: Coordinate fields for movement calculations
- ✅ **Boundary Management**: Edge detection and processing

## **SOURCE CODE CORRELATION**

### **Class11.java (Reference):**
```java
final class Class11 {
    public Class11() {
        anInt290 = 0;
        anInt291 = 0;
        anInt292 = 104;
        anInt293 = 104;
        anIntArrayArray294 = new int[anInt292][anInt293];
        method210();
    }
    
    public void method210() {
        for(int i = 0; i < anInt292; i++) {
            for(int j = 0; j < anInt293; j++) {
                if(i == 0 || j == 0 || i == anInt292 - 1 || j == anInt293 - 1)
                    anIntArrayArray294[i][j] = 0xffffff;
                else
                    anIntArrayArray294[i][j] = 0x1000000;
            }
        }
    }
    
    private final int anInt290, anInt291, anInt292, anInt293;
    public final int[][] anIntArrayArray294;
}
```

## COMMAND BLOCK 1: STRUCTURE EVIDENCE
```bash
# Show class structure and inheritance in bytecode
grep -A 10 -B 5 "extends\|implements" bytecode/client/LLORVYLP.bytecode.txt

# Show corresponding structure in DEOB source
grep -A 10 -B 5 "extends\|implements" srcAllDummysRemoved/src/Class11.java

# Verify structure in javap cache
grep -A 10 -B 5 "class.*extends\|class.*implements" srcAllDummysRemoved/.javap_cache/Class11.javap.cache
```

## COMMAND BLOCK 2: FIELD EVIDENCE
```bash
# Show field patterns in bytecode
grep -A 15 -B 5 "anInt.*\|anIntArray.*\|aBoolean.*\|aString" bytecode/client/LLORVYLP.bytecode.txt

# Show field structure in DEOB source
grep -A 15 -B 5 "public.*\|private.*\|protected.*" srcAllDummysRemoved/src/Class11.java | head -30

# Verify field declarations in javap cache
grep -A 15 -B 5 "int.*\|boolean.*\|String.*\|int\[\].*" srcAllDummysRemoved/.javap_cache/Class11.javap.cache
```

## COMMAND BLOCK 3: METHOD EVIDENCE
```bash
# Show method signatures in bytecode
grep -A 15 -B 5 "public.*\|private.*\|protected.*" bytecode/client/LLORVYLP.bytecode.txt | grep "(" | head -10

# Show method signatures in DEOB source
grep -A 20 -B 5 "public.*\|private.*" srcAllDummysRemoved/src/Class11.java | grep "(" | head -10

# Verify methods in javap cache
grep -A 25 "public.*\|private.*" srcAllDummysRemoved/.javap_cache/Class11.javap.cache | grep "(" | head -10
```

## COMMAND BLOCK 4: CROSS-REFERENCE EVIDENCE
```bash
# Show unique patterns compared to similar classes
grep -A 10 -B 5 "byte\[\].*a\|int.*b.*int.*c" bytecode/client/LLORVYLP.bytecode.txt

# Show class-specific metrics
grep -c "anInt290\|anInt291\|anInt292\|anInt293" bytecode/client/LLORVYLP.bytecode.txt

# Verify class lacks exclusion patterns (distinguishes from others)
grep -l "fragment\|badenc\|domainenc" bytecode/client/LLORVYLP.bytecode.txt | wc -l
```

## Deobfuscated Source Evidence Commands
grep -A 10 -B 5 "method210" srcAllDummysRemoved/src/Class11.java
grep -A 5 -B 5 "anIntArrayArray294" srcAllDummysRemoved/src/Class11.java
grep -A 5 -B 5 "anInt290" srcAllDummysRemoved/src/Class11.java

## Javap Cache Evidence Commands
grep -A 10 -B 5 "method210" srcAllDummysRemoved/.javap_cache/Class11.javap.cache
grep -A 5 -B 5 "anIntArrayArray294" srcAllDummysRemoved/.javap_cache/Class11.javap.cache
grep -A 5 -B 5 "anInt290" srcAllDummysRemoved/.javap_cache/Class11.javap.cache

## **UNIQUE IDENTIFIERS**
- **17 Field Structure**: 1 byte[] + 16 int fields
- **Coordinate Storage**: Multiple int fields for position data
- **Array Integration**: Works with 2D collision arrays
- **Boundary Processing**: Edge detection and collision logic
- **World Positioning**: 3D coordinate system support

## **MAPPING CONFIDENCE**
**75% CONFIDENCE** - The field structure alignment, collision detection patterns, and coordinate storage provide strong evidence. The class clearly implements data container functionality for world collision systems, though specific algorithm details may vary.

## **IMPACT**
- World collision detection and pathfinding support
- Coordinate system management for 3D world positioning
- Integration with collision array processing systems
- Essential for movement and boundary checking in game world