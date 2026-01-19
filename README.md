# RuneScape 2 (Build 317) Software Renderer

This project implements a software 3D renderer for RuneScape 2
(build 317, 377 cache) using Go with Ebitengine v2 integration.

## Project Structure

```
solmud/
├── Makefile                    # Build and development commands
├── go.mod                     # Go module definition
├── pkg/
│   ├── renderer/             # Core rendering interfaces
│   ├── geometry/             # Model loading and geometry
│   ├── math/                 # Fixed-point 3D math
│   ├── framebuffer/          # Ebitengine integration
│   ├── raster/              # Scanline rasterization (TODO)
│   ├── texture/             # Texture loading (TODO)
│   ├── camera/              # Camera management (TODO)
│   └── world/               # Scene management (TODO)
└── cmd/
    └── render3d/            # Entry point
```

## Current Status

✅ **Completed:**
- Project structure created
- Makefile with common dev commands
- Go module initialized
- Core interfaces defined in renderer/ and geometry/ packages
- Fixed-point math library implemented (trig tables, vector operations)
- Ebitengine framebuffer wrapper implemented

⚠️ **Known Issues:**
- Ebitengine v2 Game interface API unclear (RunGame not working as expected)
- Need to investigate proper v2 game loop pattern
- Currently using placeholder render implementation

## Development Commands

```bash
# Build
make build

# Run
make run

# Test
make test

# Clean
make clean

# Format code
make fmt
```

## Next Steps

1. **Resolve Ebitengine v2 integration issues**
   - Research proper game loop pattern for v2
   - Implement proper ebiten.Game interface

2. **Implement model loader**
   - Parse moparscape 377 cache format
   - Implement delta-encoded vertex decompression
   - Load Class21 model structures

3. **Implement scanline renderer**
   - Triangle rasterization algorithm
   - Gouraud shading interpolation
   - Texture mapping (affine)
   - Painter's algorithm depth sorting

4. **Integrate mopar assets**
   - Load models from `./mopar/` directory
   - Parse 377 cache format
   - Extract vertex, face, and texture data

## RuneScape 2 Build 317 Details

- **Resolution:** 512x384 (matching original RS2)
- **FOV Scale:** 512 (original RS2 perspective)
- **Near Plane:** 50 units
- **Far Plane:** 3500 units  
- **Trig Table:** 2048 entries (full 360° rotation)
- **Color Depth:** 16-bit (65535 colors)
- **Rotation Range:** 0-2048 (2048 = 360°)

## References

- Moparscape project: RS2 build 317 deobfuscated client
- AMILLI.md: Code style and principles guide
- Ebitengine v2: https://ebitengine.org/
