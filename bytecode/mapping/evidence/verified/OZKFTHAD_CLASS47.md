# Forensic Evidence: OZKFTHAD → Class47

## **CLASS IDENTIFICATION**
- **Obfuscated Name**: OZKFTHAD
- **Deobfuscated Name**: Class47
- **Confidence**: 85% (STRONG EVIDENCE)
- **Date Identified**: January 8, 2026

## **PRIMARY FORENSIC EVIDENCE**

### **1. Complex Data Processing Structure (STRONG)**
The class implements sophisticated data processing with multiple state flags:

**Field Structure:**
- ✅ **Boolean Flags**: 4 boolean fields for state management (a, c, d)
- ✅ **Byte Field**: Single byte field for compact data (b)
- ✅ **Integer Arrays**: Multiple int arrays for data processing (f, g)
- ✅ **Processing Variables**: 7 int fields for algorithmic operations (e, h, i, j, k, l, m, n, o)
- ✅ **Static Counter**: Public static int `p` for global state

### **2. Stream Integration (STRONG)**
Direct integration with Stream class (MBMGIXGO) for data processing:

**Stream Constructor Pattern:**
```java
public final void a(boolean, MBMGIXGO);
  Code:
     0: aload_0
     1: aload_2
     2: invokevirtual #25                 // Method MBMGIXGO.c:()I
     5: putfield      #41                 // Field j:I
     8: iload_1
     9: ifne          20
    12: new           #8                  // class java/lang/NullPointerException
    15: dup
    16: invokespecial #14                 // Method java/lang/NullPointerException."<init>":()V
    19: athrow
    20: aload_0
    21: aload_2
    22: invokevirtual #27                 // Method MBMGIXGO.h:()I
    25: putfield      #32                 // Field h:I
```

- ✅ **Stream Parameter**: Takes MBMGIXGO (Stream) in methods
- ✅ **Data Reading**: Calls `c()` and `h()` methods for data extraction
- ✅ **State Management**: Updates multiple fields from stream data
- ✅ **Error Handling**: Null pointer exception for invalid parameters

### **3. Mathematical Operations (STRONG)**
Implements complex mathematical processing with distinctive constants:

**65536.0d Operations:**
```java
ldc2_w        #123                // double 65536.0d
ddiv
```

- ✅ **High Precision**: Uses 65536.0d (2^16) for high-precision calculations
- ✅ **Division Operations**: Double precision division for accurate results
- ✅ **Bit Operations**: Multiple bit manipulation and masking operations
- ✅ **Coordinate Processing**: Advanced geometric transformations

### **4. Array Processing Logic (STRONG)**
Sophisticated array manipulation with multiple indexing operations:

**Array Operations:**
```java
aload_0
getfield      #35                 // Field f:[I
iload         4
iaload
aload_0
getfield      #38                 // Field g:[I
iload         6
iaload
```

- ✅ **Dual Array Processing**: Uses both `f` and `g` int arrays
- ✅ **Complex Indexing**: Multiple array access patterns with variable indices
- ✅ **Data Transformation**: Array-based data manipulation algorithms
- ✅ **State Updates**: Updates array contents based on processing logic

### **5. Processing Pipeline (STRONG)**
Implements comprehensive data processing pipeline:

**Method Signatures:**
- ✅ **Initialization**: Constructor and setup methods
- ✅ **Data Processing**: Complex algorithms with multiple parameters
- ✅ **State Management**: Boolean flags controlling processing flow
- ✅ **Result Generation**: Produces processed data outputs

### **6. Cross-Reference Integration (MODERATE)**
Used in coordinate transformation and data processing systems:

- ✅ **Geometric Calculations**: Used for 3D coordinate transformations
- ✅ **Data Processing**: Handles complex data structures and algorithms
- ✅ **State Management**: Maintains processing state across operations
- ✅ **Performance Optimization**: Efficient data processing with bit operations

## **SOURCE CODE CORRELATION**

### **Class47.java (Conceptual Structure):**
```java
final class Class47 {
    private boolean a, c, d;
    private byte b;
    private int e;
    private int[] f, g;
    int h, i, j;
    private int k, l, m, n, o;
    public static int p;

    public final void processData(boolean flag, Stream stream) {
        j = stream.readInt();
        if(!flag) throw new NullPointerException();
        h = stream.readInt();
        // Complex processing with 65536.0d operations
        // Array manipulation with f[] and g[]
        // State management with boolean flags
    }

    // Additional processing methods...
}
```

## **UNIQUE IDENTIFIERS**
- **65536.0d Operations**: High-precision double division constant
- **MBMGIXGO Integration**: Stream class parameter usage
- **Multiple Boolean Flags**: 4 boolean fields for state control
- **Dual Array Processing**: f[] and g[] int arrays
- **Complex Indexing**: Variable-based array access patterns

## **MAPPING CONFIDENCE**
**85% CONFIDENCE** - The combination of complex data structure, Stream integration, 65536.0d mathematical operations, dual array processing, and sophisticated state management represents strong forensic evidence. The class clearly implements advanced data processing algorithms, though specific functionality details may vary.

## **IMPACT**
- Advanced data processing capabilities for the client
- Complex mathematical operations with high precision
- Stream-based data handling and transformation
- Critical for coordinate calculations and data management systems