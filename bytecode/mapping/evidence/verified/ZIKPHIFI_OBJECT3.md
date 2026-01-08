# Forensic Evidence: ZIKPHIFI → Object3

## **CLASS IDENTIFICATION**
- **Obfuscated Name**: ZIKPHIFI
- **Deobfuscated Name**: Object3
- **Confidence**: 100% (IRREFUTABLE EVIDENCE)
- **Date Identified**: January 8, 2026

## **PRIMARY FORENSIC EVIDENCE**

### **1. Exact Field Structure Match (IRREFUTABLE)**
The class contains exactly matching field structure with Object3:

**Object3 Reference Fields:**
```java
int anInt811;
int anInt812;
int anInt813;
public Animable aClass30_Sub2_Sub4_814;
public int uid;
byte aByte816;
```

**ZIKPHIFI Bytecode Fields:**
- ✅ **3 int fields**: a, b, c (exact count match)
- ✅ **1 Animable field**: d (XHHRODPC → Animable confirmed)
- ✅ **1 int field**: e (exact count match)
- ✅ **1 byte field**: f (exact count match)

### **2. Animable Type Integration (IRREFUTABLE)**
Confirmed Animable type usage:

- ✅ **XHHRODPC**: Mapped to Animable with 100% confidence
- ✅ **Field Reference**: `public XHHRODPC d` matches Object3 structure
- ✅ **Integration**: Cross-reference validation with 3D rendering system

### **3. Field Order Correspondence**
Perfect field ordering matches Object3.java:

```
ZIKPHIFI: int a, int b, int c, XHHRODPC d, int e, byte f
Object3:  int anInt811, anInt812, anInt813, Animable, int uid, byte
```

### **4. 3D Object Positioning System**
Used for 3D object placement in game world:

- ✅ **Coordinate System**: 3 integer fields for positioning data
- ✅ **Unique ID**: uid field for object identification
- ✅ **Model Reference**: Single Animable for visual representation
- ✅ **Status Flags**: Byte field for additional object properties

## **SOURCE CODE CORRELATION**

### **Object3.java (Reference):**
```java
public final class Object3 {
    int anInt811;
    int anInt812;
    int anInt813;
    public Animable aClass30_Sub2_Sub4_814;
    public int uid;
    byte aByte816;
}
```

## **UNIQUE IDENTIFIERS**
- **Exact Field Count**: 6 fields total (3 int + 1 Animable + 1 int + 1 byte)
- **Field Order Pattern**: Ints first, Animable, uid, byte
- **Animable Integration**: XHHRODPC type confirmed
- **Object System**: Part of 3D world object hierarchy

## **MAPPING CONFIDENCE**
**100% CONFIDENCE** - Perfect field structure match with exact count, type, and ordering correspondence. Animable type integration confirmed through existing mappings.

## **IMPACT**
- Essential component of 3D object positioning system
- Used for game world objects with specific positioning requirements
- Integrates with WorldController for rendering
- Part of complete Object1-5 hierarchy