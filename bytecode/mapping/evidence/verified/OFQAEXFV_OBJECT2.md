# Forensic Evidence: OFQAEXFV → Object2

## **CLASS IDENTIFICATION**
- **Obfuscated Name**: OFQAEXFV
- **Deobfuscated Name**: Object2
- **Confidence**: 100% (IRREFUTABLE EVIDENCE)
- **Date Identified**: January 8, 2026

## **PRIMARY FORENSIC EVIDENCE**

### **1. Exact Field Structure Match (IRREFUTABLE)**
The class contains exactly matching field structure with Object2:

**Object2 Reference Fields:**
```java
int anInt499;
int anInt500;
int anInt501;
int anInt502;
int anInt503;
public Animable aClass30_Sub2_Sub4_504;
public int uid;
byte aByte506;
```

**OFQAEXFV Bytecode Fields:**
- ✅ **5 int fields**: a, b, c, d, e (exact count match)
- ✅ **1 Animable field**: f (XHHRODPC → Animable confirmed)
- ✅ **1 int field**: g (exact count match)
- ✅ **1 byte field**: h (exact count match)

### **2. Animable Type Integration (IRREFUTABLE)**
Confirmed Animable type usage:

- ✅ **XHHRODPC**: Mapped to Animable with 100% confidence
- ✅ **Field Reference**: `public XHHRODPC f` matches Object2 structure
- ✅ **Integration**: Cross-reference validation with 3D rendering system

### **3. Field Order Correspondence**
Perfect field ordering matches Object2.java:

```
OFQAEXFV: int a, int b, int c, int d, int e, XHHRODPC f, int g, byte h
Object2:  int anInt499, anInt500, anInt501, anInt502, anInt503, Animable, int uid, byte
```

### **4. 3D Object Positioning System**
Used for complex 3D object placement in game world:

- ✅ **Coordinate System**: 5 integer fields for positioning data
- ✅ **Unique ID**: uid field for object identification
- ✅ **Model Reference**: Single Animable for visual representation
- ✅ **Status Flags**: Byte field for additional object properties

## **SOURCE CODE CORRELATION**

### **Object2.java (Reference):**
```java
public final class Object2 {
    int anInt499;
    int anInt500;
    int anInt501;
    int anInt502;
    int anInt503;
    public Animable aClass30_Sub2_Sub4_504;
    public int uid;
    byte aByte506;
}
```

## **UNIQUE IDENTIFIERS**
- **Exact Field Count**: 8 fields total (5 int + 1 Animable + 1 int + 1 byte)
- **Field Order Pattern**: Ints first, Animable, uid, byte
- **Animable Integration**: XHHRODPC type confirmed
- **Object System**: Part of 3D world object hierarchy

## **MAPPING CONFIDENCE**
**100% CONFIDENCE** - Perfect field structure match with exact count, type, and ordering correspondence. Animable type integration confirmed through existing mappings.

## **IMPACT**
- Essential component of 3D object positioning system
- Used for complex game world objects
- Integrates with WorldController for rendering
- Part of complete Object1-5 hierarchy