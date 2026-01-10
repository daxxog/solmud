# SWTXAYDT → Animable_Sub4

## Overview
Animable_Sub4 handles projectile or particle animations with physics-based movement, gravity, and SpotAnim rendering.

## Architectural Relationships
Animable_Sub4 extends Animable (XHHRODPC), uses SpotAnim (MUDLUUBC), Model (ZKARKDQW), Animation (LKGEGIEW).

```mermaid
classDiagram
    SWTXAYDT --> XHHRODPC[Animable]
    SWTXAYDT --> MUDLUUBC[SpotAnim]
    SWTXAYDT --> ZKARKDQW[Model]
    SWTXAYDT --> LKGEGIEW[Animation]
    SWTXAYDT : +getRotatedModel()
    SWTXAYDT : +method456(int)
```

## Bytecode Matching Commands
To show class and inheritance:
```
grep -A 10 "final class SWTXAYDT" bytecode/client/SWTXAYDT.bytecode.txt
```

To show physics methods:
```
grep -A 10 -B 5 "Math\." bytecode/client/SWTXAYDT.bytecode.txt
```

## Deobfuscated Source Evidence Commands
For Animable_Sub4 class:
```
grep -A 20 "public final class Animable_Sub4" srcAllDummysRemoved/src/Animable_Sub4.java
```

For method456:
```
grep -A 10 "public void method456" srcAllDummysRemoved/src/Animable_Sub4.java
```

## Javap Cache Evidence Commands
For class structure:
```
grep -A 20 "public final class Animable_Sub4" srcAllDummysRemoved/.javap_cache/Animable_Sub4.javap.cache
```

For methods:
```
grep -A 5 "getRotatedModel" srcAllDummysRemoved/.javap_cache/Animable_Sub4.javap.cache
```

## Verification of Non-Contradictory Evidence
Bytecode matches source/javap in projectile physics, animation refs. No contradictions. 1:1 mapping confirmed.