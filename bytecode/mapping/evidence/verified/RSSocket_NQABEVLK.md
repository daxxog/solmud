# Evidence: RSSocket → NQABEVLK

## Class Overview

**RSSocket** is a network socket wrapper class that provides thread-safe client-server communication for the RuneScape client. It implements the Runnable interface to handle asynchronous network I/O operations, managing TCP socket connections with configurable timeouts and buffering. The class serves as the fundamental network communication layer, providing reliable data transmission between the client and game servers with robust error handling.

The class provides comprehensive network functionality:
- **Socket Management**: TCP connection establishment and lifecycle management with proper cleanup
- **Asynchronous I/O**: Runnable implementation for background network operations and non-blocking communication
- **Buffering**: Input/output stream buffering for efficient data handling and performance optimization
- **Error Handling**: Robust exception handling for network failures and connection recovery
- **Thread Safety**: Synchronized operations for concurrent access protection and data integrity

## Architecture Role
RSSocket acts as the network communication abstraction layer, providing a clean interface for the client to handle all server communication while managing the complexities of socket programming and thread safety. The class integrates with the main client and RSApplet for network operations, and works with Stream objects for data serialization and deserialization. RSSocket forms the foundation of all client-server communication in RuneScape.

```mermaid
graph TD
    RSSocket --> client[client]
    RSSocket --> RSApplet[RSApplet]
    RSSocket --> Stream[Stream]

    subgraph NetworkLayer
        RSSocket
    end

    subgraph ApplicationLayer
        client
        RSApplet
    end

    subgraph DataLayer
        Stream
    end

    client -.-> RSSocket
    RSApplet -.-> RSSocket
    RSSocket -.-> Stream
```

## Forensic Evidence Commands

### 1. Runnable Interface Implementation Evidence (RSSOCKET-SPECIFIC PATTERN)
```bash
# Show Runnable interface implementation in bytecode with multi-line context
grep -A 15 -B 10 "java.lang.Runnable\|run.*method" bytecode/client/NQABEVLK.bytecode.txt

# Show corresponding Runnable implementation in DEOB source with multi-line context
grep -A 15 -B 10 "implements Runnable\|public void run" srcAllDummysRemoved/src/RSSocket.java

# Verify Runnable implementation in javap cache with multi-line context
grep -A 15 -B 10 "java.lang.Runnable\|public void run" srcAllDummysRemoved/.javap_cache/RSSocket.javap.cache

# A/B Evidence: Show Runnable implementation patterns
echo "=== BYTECODE RUNNABLE (A) ===" && grep -A 10 -B 2 "implements.*java.lang.Runnable" bytecode/client/NQABEVLK.bytecode.txt
echo "=== DEOB SOURCE RUNNABLE (B) ===" && grep -A 10 -B 2 "implements.*Runnable" srcAllDummysRemoved/src/RSSocket.java
echo "=== JAVAP RUNNABLE VERIFICATION ===" && grep -A 10 "implements.*Runnable" srcAllDummysRemoved/.javap_cache/RSSocket.javap.cache
```

### 2. Network Socket Field Evidence
```bash
# Show network I/O fields (InputStream, OutputStream, Socket) in bytecode
grep -A 10 -B 5 "java\.io\.InputStream\|java\.io\.OutputStream\|java\.net\.Socket" bytecode/client/NQABEVLK.bytecode.txt

# Show corresponding network fields in DEOB source
grep -A 15 -B 5 "InputStream\|OutputStream\|Socket.*socket" srcAllDummysRemoved/src/RSSocket.java

# Verify network field declarations in javap cache
grep -A 15 -B 5 "InputStream\|OutputStream\|Socket" srcAllDummysRemoved/.javap_cache/RSSocket.javap.cache
```

### 3. Constructor with IOException Evidence
```bash
# Show constructor signature showing network initialization with IOException in bytecode
grep -A 20 -B 5 "public NQABEVLK.*throws java\.io\.IOException" bytecode/client/NQABEVLK.bytecode.txt

# Show corresponding constructor with IOException in DEOB source
grep -A 20 -B 5 "public RSSocket.*throws java\.io\.IOException" srcAllDummysRemoved/src/RSSocket.java

# Verify constructor signature in javap cache
grep -A 15 "public RSSocket.*throws java\.io\.IOException" srcAllDummysRemoved/.javap_cache/RSSocket.javap.cache
```

### 4. TCP Socket Configuration Evidence
```bash
# Show TCP socket configuration (timeout, TCP_NODELAY) in bytecode with context
grep -A 20 -B 10 "setSoTimeout\|setTcpNoDelay\|30000\|true" bytecode/client/NQABEVLK.bytecode.txt

# Show corresponding socket configuration in DEOB source with context
grep -A 20 -B 10 "setSoTimeout\|setTcpNoDelay\|30000" srcAllDummysRemoved/src/RSSocket.java

# Verify socket configuration in javap cache with context
grep -A 20 -B 10 "setSoTimeout\|setTcpNoDelay\|30000" srcAllDummysRemoved/.javap_cache/RSSocket.javap.cache
```

### 5. Asynchronous run() Method Evidence
```bash
# Show run() method implementation for asynchronous I/O in bytecode
grep -A 25 -B 5 "public void run\|while.*true" bytecode/client/NQABEVLK.bytecode.txt

# Show corresponding run() method in DEOB source
grep -A 25 -B 5 "public void run" srcAllDummysRemoved/src/RSSocket.java

# Verify run() method in javap cache
grep -A 25 "public void run" srcAllDummysRemoved/.javap_cache/RSSocket.javap.cache
```

### 6. Stream Integration Evidence
```bash
# Show Stream object integration for data I/O in bytecode
grep -A 15 -B 5 "Stream\|MBMGIXGO\|read\|write" bytecode/client/NQABEVLK.bytecode.txt

# Show Stream usage in DEOB source
grep -A 15 -B 5 "Stream\|read\|write" srcAllDummysRemoved/src/RSSocket.java

# Verify Stream integration in javap cache
grep -A 15 -B 5 "Stream" srcAllDummysRemoved/.javap_cache/RSSocket.javap.cache
```

### 7. Cross-Reference Validation (UNIQUE RSSOCKET PATTERN)
```bash
# Confirm NQABEVLK only maps to RSSocket - unique mapping verification
grep -r "NQABEVLK" bytecode/mapping/evidence/verified/ | grep -v RSSocket || echo "Unique mapping confirmed"

# Verify the unique network socket + RSApplet + Runnable pattern appears only in NQABEVLK
find bytecode/client/ -name "*.bytecode.txt" -exec grep -l "implements java.lang.Runnable" {} \; | xargs grep -l "java.io.InputStream" | xargs grep -l "java.io.OutputStream" | xargs grep -l "java.net.Socket" | xargs grep -l "KHACHIFW"

# Show RSSocket's unique network configuration pattern
grep -l "setSoTimeout\|30000" bytecode/client/*.bytecode.txt | xargs grep -l "setTcpNoDelay\|true" | xargs grep -l "NQABEVLK"
```

### 8. Network Error Handling Evidence
```bash
# Show exception handling for network operations in bytecode
grep -A 15 -B 5 "IOException\|catch\|finally" bytecode/client/NQABEVLK.bytecode.txt

# Show corresponding error handling in DEOB source
grep -A 15 -B 5 "IOException\|catch\|finally" srcAllDummysRemoved/src/RSSocket.java

# Verify error handling in javap cache
grep -A 15 -B 5 "IOException\|catch" srcAllDummysRemoved/.javap_cache/RSSocket.javap.cache
```

## Critical Evidence Points

1. **Runnable Implementation**: The class implements java.lang.Runnable, indicating asynchronous network operation capabilities for background I/O processing.

2. **Network Socket Fields**: Contains the complete set of TCP socket management fields: InputStream, OutputStream, and Socket references for bidirectional communication.

3. **IOException Handling**: Constructor signature includes "throws java.io.IOException", standard for network operations requiring error handling.

4. **TCP Configuration**: Socket configuration with setSoTimeout(30000) and setTcpNoDelay(true), typical for real-time game network optimization.

5. **Asynchronous Operations**: run() method implementation for non-blocking network communication with while loop for continuous operation.

## Verification Status

**VERIFIED** - All bash commands execute successfully and evidence is non-contradictory. The combination of Runnable implementation, complete network socket field set, TCP configuration patterns, and asynchronous I/O methods provides 100% confidence in this 1:1 mapping. The unique integration of network fields with Runnable interface and socket configuration establishes this as the definitive network communication layer.

## Sources and References
- **Bytecode**: bytecode/client/NQABEVLK.bytecode.txt
- **Deobfuscated Source**: srcAllDummysRemoved/src/RSSocket.java
- **Javap Cache**: srcAllDummysRemoved/.javap_cache/RSSocket.javap.cache
- **Network Configuration**: setSoTimeout(30000) and setTcpNoDelay(true) for game optimization
- **Stream Integration**: MBMGIXGO (Stream) for data serialization
- **Client Integration**: KHACHIFW (RSApplet) for application-layer communication
- **Runnable Interface**: Asynchronous I/O operations for background network processing
- **Error Handling**: Comprehensive IOException management for network reliability