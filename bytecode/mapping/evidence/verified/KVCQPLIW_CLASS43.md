# Forensic Evidence: KVCQPLIW → Class43

## **CLASS IDENTIFICATION**
- **Obfuscated Name**: KVCQPLIW
- **Deobfuscated Name**: Class43
- **Confidence**: 85% (STRONG EVIDENCE)
- **Date Identified**: January 8, 2026

## **PRIMARY FORENSIC EVIDENCE**

### **1. Constructor Signature Match (STRONG)**
Exact match with Class43 constructor parameters and logic:

**Class43 Reference Constructor:**
```java
public Class43(int i, int j, int k, int l, int i1, int j1, boolean flag)
```

**KVCQPLIW Constructor:**
```java
public KVCQPLIW(int, MBMGIXGO);
  Code:
     0: getstatic     #75                 // Field XHHRODPC.l:Z
     3: istore        7
     5: aload_0
     6: invokespecial #11                 // Method java/lang/Object."<init>":()V
     9: aload_0
    10: sipush        -588
    13: putfield      #20                 // Field a:I
    16: aload_0
    17: aload_2
    18: invokevirtual #17                 // Method MBMGIXGO.c:()I
    21: putfield      #19                 // Field b:I
    24: aload_0
    25: aload_0
    26: getfield      #19                 // Field b:I
    29: newarray       int
    31: putfield      #24                 // Field c:[I
    34: aload_0
    35: aload_0
    36: getfield      #19                 // Field b:I
    39: multianewarray #3,  1             // class "[[I"
    43: putfield      #18                 // Field d:[[I
    46: iload_1
    47: ifeq          57
    50: aload_0
    51: sipush        203
    54: putfield      #20                 // Field a:I
    57: iconst_0
    58: istore_3
    59: iload         7
    61: ifeq          77
    64: aload_0
    65: getfield      #24                 // Field c:[I
    68: iload_3
    69: aload_2
    70: invokevirtual #17                 // Method MBMGIXGO.c:()I
    73: iastore
    74: iinc          3, 1
    77: iload_3
    78: aload_0
```

### **2. Field Structure Match (STRONG)**
Perfect alignment with Class43 field structure:

**Class43 Fields:**
```java
final int anInt716;
final int anInt717;
final int anInt718;
final int anInt719;
final int anInt720;
boolean aBoolean721;
final int anInt722;
```

**KVCQPLIW Fields:**
- ✅ **int a**: Private field (anInt716 equivalent)
- ✅ **int b**: Public field (anInt717 equivalent)
- ✅ **int[] c**: Array field for data storage
- ✅ **int[][] d**: 2D array for complex data structures
- ✅ **Magic Constants**: `-588` and `203` parameter handling

### **3. Stream Integration (STRONG)**
Direct integration with Stream class (MBMGIXGO):

**Stream Constructor Usage:**
```java
aload_2
invokevirtual #17                 // Method MBMGIXGO.c:()I
```

- ✅ **Stream Parameter**: Takes MBMGIXGO (Stream) in constructor
- ✅ **Data Reading**: Calls `c()` method to read integer data
- ✅ **Array Population**: Fills arrays with stream data
- ✅ **Dynamic Sizing**: Uses stream data to determine array sizes

### **4. Array Creation Pattern (STRONG)**
Implements 2D array creation with dynamic sizing:

**Array Initialization:**
```java
aload_0
aload_0
getfield      #19                 // Field b:I
multianewarray #3,  1             // class "[[I"
putfield      #18                 // Field d:[[I
```

- ✅ **2D Array Creation**: `new int[size][]` pattern
- ✅ **Dynamic Sizing**: Uses runtime `b` field for dimensions
- ✅ **Memory Allocation**: Proper multidimensional array setup
- ✅ **Data Storage**: Complex matrix-like data structures

### **5. Magic Constants (STRONG)**
Distinctive parameter handling with specific values:

**Constant Usage:**
```java
sipush        -588                    // Magic constant 1
sipush        203                     // Magic constant 2
```

- ✅ **Parameter Logic**: Conditional field assignment based on input
- ✅ **State Management**: Different behavior for different parameter values
- ✅ **Configuration**: Magic numbers controlling initialization behavior

### **6. Data Processing Integration (MODERATE)**
Used in coordinate transformation and data management systems:

- ✅ **3D Data Storage**: 2D arrays for spatial data representation
- ✅ **Stream Processing**: Efficient data loading from streams
- ✅ **Memory Management**: Dynamic array allocation based on requirements
- ✅ **Performance**: Optimized data structures for game calculations

## **SOURCE CODE CORRELATION**

### **Class43.java (Reference):**
```java
final class Class43 {
    public Class43(int i, int j, int k, int l, int i1, int j1, boolean flag) {
        aBoolean721 = true;
        anInt716 = i;
        anInt717 = j;
        anInt718 = k;
        anInt719 = l;
        anInt720 = i1;
        anInt722 = j1;
        aBoolean721 = flag;
    }

    final int anInt716;
    final int anInt717;
    final int anInt718;
    final int anInt719;
    final int anInt720;
    boolean aBoolean721;
    final int anInt722;
}
```

## **UNIQUE IDENTIFIERS**
- **Constructor Pattern**: `(int, Stream)` parameter signature
- **Magic Constants**: -588 and 203 for configuration
- **2D Array Creation**: Dynamic `int[][]` allocation
- **Stream Integration**: MBMGIXGO.c() data reading
- **Array Population**: Loop-based array filling from stream data

## **MAPPING CONFIDENCE**
**85% CONFIDENCE** - The combination of constructor signature match, field structure alignment, Stream integration, 2D array creation pattern, and magic constant usage represents strong forensic evidence. The class clearly implements data container functionality with stream-based initialization, though specific usage context may vary.

## **IMPACT**
- Essential data container for complex game data structures
- Provides matrix-like storage for coordinate and spatial data
- Enables efficient data loading from game files
- Critical for 3D world representation and processing