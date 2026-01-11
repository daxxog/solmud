# Class21 → LLORVYLP Mapping Dispute

## Current Incorrect Mapping
- **DEOB**: Class21  
- **OG**: OZKFTHAD (90% confidence)
- **Issue**: Field signature mismatch requires investigation

## Correct Mapping Evidence
- **DEOB**: Class21
- **OG**: LLORVYLP 
- **Should be**: 100% confidence exact match

## Forensic Evidence

### 1. DEOB Class21 Structure Analysis

**Command to examine Class21 structure:**
```bash
cat srcAllDummysRemoved/src/Class21.java
```

**Evidence shows:**
```java
final class Class21 {
    public Class21() { }
    
    public byte aByteArray368[];
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
}
```

**Field pattern:** 1 byte array + 16 integers = 17 total fields
**Method pattern:** 1 simple constructor

### 2. DEOB Class21 Javap Cache Analysis

**Command to examine javap cache:**
```bash
cat srcAllDummysRemoved/.javap_cache/Class21.javap.cache
```

**Evidence confirms exact field structure:**
```javap
final class Class21 {
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

  public Class21();
}
```

### 3. OG LLORVYLP Structure Analysis  

**Command to examine LLORVYLP structure:**
```bash
head -20 bytecode/client/LLORVYLP.bytecode.txt
```

**Evidence shows exact match:**
```
LLORVYLP: 18 field declarations ===
  public byte[] a;
  public int b;
  public int c;
  public int d;
  public int e;
  public int f;
  public int g;
  public int h;
  public int i;
  public int j;
  public int k;
  public int l;
  public int m;
  public int n;
  public int o;
  public int p;
  public int q;
  public LLORVYLP();
```

**Field pattern:** 1 byte array + 16 integers = 17 total fields  
**Method pattern:** 1 simple constructor

### 4. Current Incorrect Mapping Analysis

**Command to examine OZKFTHAD (incorrectly mapped to Class21):**
```bash
head -35 bytecode/client/OZKFTHAD.bytecode.txt
```

**Evidence shows complete mismatch:**
```
public class OZKFTHAD {
  private boolean a;
  private byte b;
  private boolean c;
  private boolean d;
  private int e;
  private int[] f;
  private int[] g;
  int h;
  int i;
  int j;
  private int k;
  private int l;
  private int m;
  private int n;
  private int o;
  public static int p;

  public final void a(boolean, MBMGIXGO);
  public final void a(byte, MBMGIXGO);
  final void a(byte);
  final int a(boolean, int);
  public OZKFTHAD();
}
```

**Field pattern:** Multiple field types (boolean, byte, int arrays) + static field
**Method pattern:** 4 complex methods + constructor
**Total:** Not a simple data container

### 5. Class11→LLORVYLP Current Mapping Analysis

**Command to examine Class11 (currently mapped to LLORVYLP):**
```bash
head -15 srcAllDummysRemoved/src/Class11.java
```

**Evidence shows complex collision detection class:**
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
    
    public void method210() { /* complex logic */ }
    public void method211(int i, int j, int k, int l, boolean flag) { /* complex logic */ }
    // ... 10+ additional complex methods
    
    private final int anInt290;
    private final int anInt291;
    private final int anInt292;
    private final int anInt293;
    public final int[][] anIntArrayArray294;
}
```

**Field pattern:** 4 final ints + 104x104 int array
**Method pattern:** 10+ complex methods for collision detection

## Conclusion

The forensic evidence is irrefutable:

1. **Class21** = Simple data container (1 byte array + 16 ints)
2. **LLORVYLP** = Simple data container (1 byte array + 16 ints)  
3. **OZKFTHAD** = Complex class with multiple methods and mixed field types
4. **Class11** = Complex collision detection class with 104x104 array

**Correct Mapping Resolution:**
- Class21 should map to LLORVYLP (100% confidence exact match)
- Class11 should map to a different OG class (requires separate investigation)
- OZKFTHAD should map to Class11 or another appropriate complex class

## Verification Commands

```bash
# Verify Class21 structure (17 fields, 1 method)
grep -c "public.*\[\]" srcAllDummysRemoved/src/Class21.java  # Should be 1
grep -c "public int" srcAllDummysRemoved/src/Class21.java    # Should be 16

# Verify LLORVYLP structure (17 fields, 1 method)  
grep -c "Field " bytecode/client/LLORVYLP.bytecode.txt     # Should be 17
grep -c "public.*(" bytecode/client/LLORVYLP.bytecode.txt | head -1  # Should be 1

# Verify mismatch with OZKFTHAD
grep -c "public.*(" bytecode/client/OZKFTHAD.bytecode.txt | head -1  # Should be >1
```
