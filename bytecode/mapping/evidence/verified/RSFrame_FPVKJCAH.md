# Evidence: RSFrame → FPVKJCAH

## Class Overview

**RSFrame** extends Java AWT Frame to provide the main game window for RuneScape, with specialized graphics translation and event handling. RSFrame acts as the container for the game applet, managing window properties, graphics context transformation, and delegating rendering operations to the embedded RSApplet. The class provides the bridge between Java's windowing system and the game's rendering pipeline.

The class provides comprehensive window management:
- **Graphics Translation**: Automatic translation of graphics coordinates to account for window borders and title bar
- **Applet Integration**: Seamless delegation of paint and update operations to the contained RSApplet
- **Window Configuration**: Pre-configured window properties including title, resizability, and visibility
- **Event Handling**: Specialized update and paint methods that maintain game rendering consistency

## Architecture Role
RSFrame occupies the top-level position in the client rendering hierarchy, serving as the main application window that contains the RSApplet game instance. Unlike pure utility classes, RSFrame extends Java's AWT Frame and provides the essential window management layer. The class acts as a transparent wrapper that transforms system-level graphics events into game-appropriate rendering operations while maintaining the game's coordinate system.

```mermaid
classDiagram
    RSFrame --> Frame
    RSFrame --> RSApplet
    RSFrame --> Graphics
    RSFrame : +RSFrame(RSApplet, int, int)
    RSFrame : +getGraphics() : Graphics
    RSFrame : +update(Graphics)
    RSFrame : +paint(Graphics)
    RSFrame : -rsApplet (KHACHIFW)
    RSFrame : -a (boolean)
```

## Forensic Evidence Commands

### 1. RSFrame Structure and Frame Extension
```bash
# Show RSFrame extends java.awt.Frame in bytecode
grep -A 15 -B 5 "public final class FPVKJCAH extends java.awt.Frame" bytecode/client/FPVKJCAH.bytecode.txt

# Show RSFrame class definition in DEOB source
grep -A 10 -B 5 "final class RSFrame extends Frame" srcAllDummysRemoved/src/RSFrame.java

# Verify Frame extension in javap cache
grep -A 5 -B 5 "class RSFrame extends Frame" srcAllDummysRemoved/.javap_cache/RSFrame.javap.cache
```

### 2. RSApplet Integration Pattern
```bash
# Show RSApplet (KHACHIFW) field usage in constructor bytecode
grep -A 20 -B 5 "KHACHIFW\|putfield.*b" bytecode/client/FPVKJCAH.bytecode.txt

# Show RSApplet field initialization in DEOB source
grep -A 10 -B 5 "rsApplet.*=" srcAllDummysRemoved/src/RSFrame.java

# Verify RSApplet integration in javap cache
grep -A 5 -B 5 "KHACHIFW" srcAllDummysRemoved/.javap_cache/RSFrame.javap.cache
```

### 3. Window Configuration and Initialization
```bash
# Show window setup (title, resizability) in constructor bytecode
grep -A 15 -B 5 "setTitle\|setResizable\|setVisible" bytecode/client/FPVKJCAH.bytecode.txt

# Show corresponding window configuration in DEOB source
grep -A 10 -B 5 "setTitle\|setResizable\|setVisible" srcAllDummysRemoved/src/RSFrame.java

# Verify window setup methods in javap cache
grep -A 10 -B 5 "setTitle\|setResizable" srcAllDummysRemoved/.javap_cache/RSFrame.javap.cache
```

### 4. Graphics Translation System
```bash
# Show getGraphics method with coordinate translation in bytecode
grep -A 15 -B 5 "public java.awt.Graphics getGraphics" bytecode/client/FPVKJCAH.bytecode.txt

# Show corresponding getGraphics with translation in DEOB source
grep -A 15 -B 5 "public Graphics getGraphics" srcAllDummysRemoved/src/RSFrame.java

# Verify graphics translation in javap cache
grep -A 15 -B 5 "getGraphics\|translate" srcAllDummysRemoved/.javap_cache/RSFrame.javap.cache
```

### 5. Cross-Reference Validation (RSFRAME vs APPLET DISTINCTION)
```bash
# Show RSFrame extends Frame (distinguishes from RSApplet)
grep -c "java.awt.Frame" bytecode/client/FPVKJCAH.bytecode.txt

# Show RSFrame has graphics delegation methods (unique to window wrapper)
grep -A 10 -B 2 "update.*Graphics\|paint.*Graphics" bytecode/client/FPVKJCAH.bytecode.txt

# Verify RSFrame's window-specific field patterns
grep -A 5 -B 2 "boolean.*a" srcAllDummysRemoved/.javap_cache/RSFrame.javap.cache
```

### 6. Rendering Delegation Pattern
```bash
# Show update method delegation to RSApplet in bytecode
grep -A 10 -B 5 "public final void update" bytecode/client/FPVKJCAH.bytecode.txt

# Show update delegation in DEOB source
grep -A 10 -B 5 "public void update" srcAllDummysRemoved/src/RSFrame.java

# Verify paint method delegation in javap cache
grep -A 10 -B 5 "paint.*Graphics" srcAllDummysRemoved/.javap_cache/RSFrame.javap.cache
```

### 7. Window Sizing and Border Handling
```bash
# Show window sizing with border compensation in bytecode
grep -A 10 -B 5 "setSize\|i.*8\|j.*28" bytecode/client/FPVKJCAH.bytecode.txt

# Show border compensation in DEOB source
grep -A 10 -B 5 "setSize.*+ 8\|setSize.*+ 28" srcAllDummysRemoved/src/RSFrame.java

# Verify sizing constants in javap cache
grep -A 10 -B 2 "8.*28\|28.*8" srcAllDummysRemoved/.javap_cache/RSFrame.javap.cache
```

## Critical Evidence Points

1. **AWT Frame Extension**: RSFrame extends java.awt.Frame, providing native window system integration.

2. **Graphics Translation**: Specialized getGraphics method with translate(4, 24) for border compensation.

3. **Applet Delegation**: update and paint methods delegate to contained RSApplet for game rendering.

4. **Window Configuration**: Pre-configured with "Jagex" title, non-resizable, and proper sizing.

## Verification Status

**VERIFIED** - All bash commands execute successfully and evidence is non-contradictory. The AWT Frame extension, graphics translation system, applet delegation pattern, and window configuration provide definitive 1:1 mapping evidence that distinguishes RSFrame from RSApplet and other rendering classes.

## Sources and References
- **Bytecode**: bytecode/client/FPVKJCAH.bytecode.txt
- **Deobfuscated Source**: srcAllDummysRemoved/src/RSFrame.java
- **Javap Cache**: srcAllDummysRemoved/.javap_cache/RSFrame.javap.cache
- **Applet Integration**: KHACHIFW (RSApplet)
- **Base Window Class**: java.awt.Frame
- **Graphics System**: java.awt.Graphics
- **Window Configuration**: setTitle, setResizable, setSize methods