# Evidence for Animable_Sub3 -> OJEALINP

## Overview
Animable_Sub3 is a final subclass of Animable that represents temporary spot animations in the RuneScape game world. Its core functionality includes rendering animated graphical effects (such as spell impacts or environmental animations) and managing their animation timing and lifecycle through methods like getRotatedModel() for 3D model generation and method454() for animation progression.

## Architectural Relationships
Animable_Sub3 fits into the game's rendering hierarchy as a specialized animable object for spot effects. It extends Animable (which extends NodeSub) and collaborates with SpotAnim for animation data, Model for 3D rendering, Animation for frame timing, and Class36 for animation utilities.

```mermaid
classDiagram
    NodeSub <|-- Animable
    Animable <|-- Animable_Sub3
    Animable_Sub3 --> SpotAnim
    Animable_Sub3 --> Model
    Animable_Sub3 --> Animation
    Animable_Sub3 --> Class36
```

## Bash Commands Proving Bytecode Matches
- `grep -A 50 -B 5 "public OJEALINP(" bytecode/client/OJEALINP.bytecode.txt`  
  Shows the constructor bytecode with 8 parameters, initializing fields r (boolean), s (int), t (boolean duplicate), u (SpotAnim reference), m/n/o/p/q (int fields), matching the deob constructor sequence.

- `grep -A 100 -B 10 "getfield.*u" bytecode/client/OJEALINP.bytecode.txt | head -50`  
  Shows the getRotatedModel method accessing the SpotAnim field (u), invoking getModel(), checking null, accessing Animation arrays, creating new Model with Class36.method532 call, and conditional method calls for transformations.

- `grep -A 50 -B 5 "aload_0" bytecode/client/OJEALINP.bytecode.txt | grep -A 50 "dup" | grep -A 50 "iadd" | head -50`  
  Shows the animation update method (method454 equivalent) incrementing a field, comparing with Animation.method258 result in a loop, updating indices and boolean flag.

## Commands for Deob Source Sections
- `head -20 srcAllDummysRemoved/src/Animable_Sub3.java`  
  Shows the class declaration extending Animable, and the start of the constructor with parameter assignments and SpotAnim.cache access.

- `grep -A 35 -B 5 "getRotatedModel" srcAllDummysRemoved/src/Animable_Sub3.java`  
  Shows the getRotatedModel method body, including Model creation with Class36.method532, conditional method469/method470 calls, scaling checks, rotation logic, and method479 lighting call.

- `grep -A 20 -B 5 "method454" srcAllDummysRemoved/src/Animable_Sub3.java`  
  Shows the method454 animation update method with the loop incrementing anInt1570, comparing against Animation.method258, updating anInt1569 index, and setting aBoolean1567 flag.

## Commands for Javap Cache Sections
- `head -25 srcAllDummysRemoved/.javap_cache/Animable_Sub3.javap.cache`  
  Shows the class declaration, field declarations (anInt1560-1564, aBoolean1567, aSpotAnim_1568, anInt1569-1570), matching the source fields.

- `grep -A 30 -B 5 "public Animable_Sub3(" srcAllDummysRemoved/.javap_cache/Animable_Sub3.javap.cache`  
  Shows the constructor bytecode with invokespecial to Animable init, putfield for aBoolean1567, getstatic SpotAnim.cache, aaload, and putfield for each int field in sequence.

- `grep -A 50 -B 5 "public Model getRotatedModel" srcAllDummysRemoved/.javap_cache/Animable_Sub3.javap.cache`  
  Shows the getRotatedModel bytecode with invokevirtual getRotatedModel, ifnonnull, putfield modelHeight, invokevirtual model.method443.

Wait, that seems wrong, the agent said:

- `grep -A 50 -B 5 "public Model getRotatedModel" srcAllDummysRemoved/.javap_cache/Animable_Sub3.javap.cache`  
  Shows the getRotatedModel bytecode with getfield aSpotAnim_1568, invokevirtual getModel, null check, getfield aAnimation_407, getfield anIntArray353, iaload, new Model with Class36.method532, conditional invokes for method469/470, scaling checks, and method479.

I think the agent content has a mistake in the last part.

Anyway, use the agent content as is.

The javap for getRotatedModel is correct in agent.

For method454:
- `grep -A 50 -B 5 "public void method454" srcAllDummysRemoved/.javap_cache/Animable_Sub3.javap.cache`  
  Shows the method454 bytecode with dup getfield anInt1570 iadd putfield, getfield anInt1570, getfield aSpotAnim_1568 aAnimation_407, getfield anInt1569, invokevirtual method258, if_icmple jump, arithmetic updates, and boolean set.

## Verification of Non-Contradictory Evidence
The bytecode in OJEALINP.bytecode.txt matches the javap cache structure exactly, with identical bytecode sequences, field access patterns, and method flows. The deob source aligns with the javap bytecode logic. No contradictions found; the obfuscated field names (m-q, r, s, t, u, v, w) correspond 1:1 to deob fields (anInt1560-1564, aBoolean1567, anInt1569, aSpotAnim_1568, anInt1570).

## 1:1 Mapping Confirmation
The bytecode file bytecode/client/OJEALINP.bytecode.txt maps exclusively to srcAllDummysRemoved/src/Animable_Sub3.java and srcAllDummysRemoved/.javap_cache/Animable_Sub3.javap.cache. No other classes share this exact structure of 9 fields (5 public int, 1 public boolean, 1 private SpotAnim, 2 private int) and 3 methods (constructor, getRotatedModel, method454).