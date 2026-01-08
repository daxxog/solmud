# Class Mapping

This directory contains mappings between deobfuscated class names (from Mopar 317 deobfuscated client) and obfuscated class names (from original RuneScape 317 bytecode).

## Files

- `class_mapping.csv` - CSV format with high confidence (≥85.0) and uncertain (<85.0) matches (Phase 1)
- `class_mapping.json` - Complete JSON with detailed score breakdowns (Phase 1)
- `class_mapping_phase1.csv` - Phase 1 results (18 matches)
- `class_mapping_phase1.json` - Phase 1 detailed results
- `class_mapping_phase2_final.json` - Phase 2 final results (26 mappings)
- `class_mapping_phase2.json` - Phase 2 results (25 matches, enhanced matching)

## How to Regenerate

To regenerate these files from project root:
```bash
# Generate both CSV and JSON
./tools/classmapper/classmapper -mode csv > bytecode/mapping/class_mapping.csv
./tools/classmapper/classmapper -mode json > bytecode/mapping/class_mapping.json
```

## CSV Format

Columns:
- deobfuscated_name - Deobfuscated class name
- obfuscated_name - Obfuscated class name
- confidence_score - Confidence score (0-100)
- superclass_match - Whether superclass matches
- interface_count - Number of matching interfaces
- field_count_deob - Field count in deobfuscated class
- field_count_obf - Field count in obfuscated class
- method_count_deob - Method count in deobfuscated class
- method_count_obf - Method count in obfuscated class
- notes - Match reasoning

## JSON Format

Complete structure with:
- `summary` - Total matches, high/medium/low confidence counts
- `matches` - Array of match objects with:
  - `DeobfuscatedClass` - Deobfuscated class name
  - `ObfuscatedClass` - Obfuscated class name
  - `ConfidenceScore` - Confidence score (0-100)
  - `ScoreBreakdown` - Detailed scoring breakdown:
    - `InterfaceMatch` - Interface matching points
    - `SuperclassMatch` - Superclass matching points
    - `FieldCountMatch` - Field count matching points
    - `FieldSimilarity` - Field type similarity
    - `MethodCountMatch` - Method count matching points
    - `MethodSimilarity` - Method signature similarity
    - `ConstructorMatch` - Constructor matching points
    - `AccessMatch` - Access modifier matching points
    - `SizePenalty` - Size difference penalty
  - `Details` - Human-readable match reasoning

## Confidence Scoring

Matches are scored on a 0-100 scale using:

- **Interface matching** (20 points) - Exact interface matches
- **Superclass matching** (25 points) - Superclass matches
- **Field count/type similarity** (15 points) - Field matching
- **Method count/signature similarity** (25 points) - Method matching
- **Constructor matching** (5 points) - Constructor count
- **Access modifier matching** (5 points) - Access modifier matches
- **Size penalty** (-10 points) - Applied for significant size differences

**Thresholds**:
- **High confidence**: ≥85.0
- **Medium confidence**: 65.0-84.9
- **Low confidence**: <65.0

## Known Anchor Mappings

The tool uses these pre-configured anchor classes for exact matches (25 total):

**Core Infrastructure:**
- `Node` → `PKVMXVTO`
- `NodeSub` → `PPOHBEGB`
- `RSApplet` → `KHACHIFW`
- `client` → `client`
- `sign/signlink` → `sign/signlink`

**Game Entities & Rendering:**
- `Animable` → `XHHRODPC`
- `Entity` → `GQOSZKJC`
- `Model` → `ZKARKDQW`
- `Stream` → `AFCKELYG`

**Data Loading (File Pattern Matches):**
- `IDK` → `TAVAECED` (loads "idk.dat")
- `VarBit` → `SXYSOXTR` (loads "varbit.dat")
- `ItemDef` → `DJRMEMXO` (loads "obj.dat")
- `ObjectDef` → `YZDBYLRM` (loads "loc.dat")
- `EntityDef` → `CKDEJADD` (loads "npc.dat")
- `Animation` → `LKGEGIEW` (loads "seq.dat")
- `Flo` → `MNHKFPQO` (loads "flo.dat")
- `Varp` → `VGXVBFVC` (loads "varp.dat")
- `SpotAnim` → `MUDLUUBC` (loads "spotanim.dat")

**Inheritance Chain Matches:**
- `Ground` → `DYMVKFXP`
- `Class30_Sub1` → `QTKGMFHL`
- `OnDemandData` → `PHKHJKBS`
- `Player` → `DLZHLHNK`
- `NPC` → `CWNCPMLX`
- `Animable_Sub4` → `SWTXAYDT`
- `Animable_Sub5` → `WBWOBAFW`
- `Animable_Sub3` → `OJEALINP`
- `Item` → `HNKCWGJM`

These are matched with 100.0 confidence scores (anchors) and 85.0 confidence scores (inheritance chains).

## Current Status (Data Cleanup Complete)

**Overall Coverage**: 86.3% (63/73 obfuscated classes mapped)

**Coverage Breakdown**:
- 100% Confidence: 53 classes (anchor classes & exact matches) - 72.6%
- 85% Confidence: 10 classes (inherited class hierarchy) - 13.7%
- Overall Quality: 100% (all mappings meet minimum 85% standard)

**Phase Progress**:
- **Phase 1**: 25 mappings (34% coverage) - Anchor and inheritance chain matching
- **Phase 2**: 26 mappings (35% coverage) - Enhanced scoring, cross-references, validation
- **Phase 3**: 12 additional mappings (17% coverage) - Advanced pattern matching through Phase 13
- **Data Cleanup**: Removed low-quality mappings, established accurate statistics
- **Current Total**: 63 high-confidence mappings (86.3% coverage)

## Remaining Unmapped Classes (10/73)

**High-Priority Targets** (3 classes):
- Class18, Class33, DummyClass - Small utility classes

**Complex Analysis Required** (3 classes):
- Censor, ObjectManager, Class36/Class40 - Large, complex systems

**Special Cases** (3 classes):
- GUI, NodeList, Class21 - May be version-specific or conditional

**Complex Obfuscated Classes** (9 remaining):
- BISVHPUN, CRRWDRTI, DUCMKFAY, FTPNODIB, OIBEELAZ, RJXWGZGD, VBAXKVMG, XPBACSMK, ZARDZRHZ

## Next Steps

**Phase 4**: Target manual mapping of small utility classes
- Map Class18 → OIBEELAZ (90%+ expected)
- Map DummyClass → RJXWGZGD (95%+ expected)
- Map Class33 → XPBACSMK/ZARDZRHZ (85-90% expected)
- Target: 66/73 (90.4% coverage)

**Phase 5**: Deep analysis of complex systems
- Analyze Censor → CRRWDRTI (text filtering systems)
- Analyze ObjectManager → VBAXKVMG (game object management)
- Analyze Class36/Class40 → FTPNODIB/DUCMKFAY (animation/graphics)
- Target: 70/73 (95.9% coverage)

**Phase 6**: Special case investigation
- Investigate GUI, NodeList, Class21 for mapping feasibility
- Final documentation and validation
- Target: 70/73 (95.9% coverage with clear explanation of 3 unmapped classes)

## Notes

Generated by classmapper tool
Source: Mopar 317 deobfuscated (srcAllDummysRemoved/bin/)
Target: RuneScape 317 obfuscated (bytecode/client/)
Phase 2 Date: $(date)