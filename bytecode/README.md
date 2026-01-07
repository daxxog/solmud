# RS317 Original Game (OG) Client Bytecode

This directory contains the decompiled bytecode from the original RuneScape 317 client JAR file (`rs317og.jar`).

## Origin

The bytecode was extracted from `rs317og.jar` using the Java `javap` tool with the following command:
```bash
javap -c -p -classpath rs317og.jar <classname>
```

Where:
- `-c`: Print disassembled code (bytecode instructions)
- `-p`: Show private members
- `-classpath rs317og.jar`: Load classes from the JAR file

## Directory Structure

```
bytecode/
├── client/           # Main client classes (default package)
│   ├── client.bytecode.txt
│   ├── KHACHIFW.bytecode.txt
│   └── ... (71 more files)
└── sign/             # Signlink package classes
    └── signlink.bytecode.txt
```

## Purpose

These files provide insight into the original RuneScape 317 client's implementation, including:
- Class structure and inheritance
- Field definitions and access modifiers
- Method signatures and bytecode instructions
- Internal class names (obfuscated)

## Note on Obfuscation

The original client was obfuscated, so class names, method names, and field names are typically single characters or short meaningless strings. This makes the code difficult to understand without additional reverse engineering work.

## Total Classes

- 73 classes in the main client package
- 1 class in the sign package
- **Total: 74 classes**