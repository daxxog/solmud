# Forensic Evidence: RTHTIIVA → TextInput

## **CLASS IDENTIFICATION**
- **Obfuscated Name**: RTHTIIVA
- **Deobfuscated Name**: TextInput
- **Confidence**: 70% (MODERATE EVIDENCE)
- **Date Identified**: January 8, 2026

## **PRIMARY FORENSIC EVIDENCE**

### **1. Character Array Processing (MODERATE)**
The class contains character array handling matching TextInput's text processing:

**TextInput Reference:**
```java
public static char[] validChars = {
    ' ', 'e', 't', 'a', 'o', 'i', 'h', 'n', 's', 'r',
    'd', 'l', 'u', 'm', 'w', 'c', 'y', 'f', 'g', 'p',
    'b', 'v', 'k', 'x', 'j', 'q', 'z', '0', '1', '2',
    '3', '4', '5', '6', '7', '8', '9', ' ', '!', '?',
    '.', ',', ':', ';', '(', ')', '-', '&', '*', '\\',
    '\'', '@', '#', '+', '=', '\243', '$', '%', '"', '[',
    ']'
};
```

**RTHTIIVA Character Arrays:**
```java
public static char[] d;                       // Character encoding table
private static char[] f;                      // Processing array
```

- ✅ **Static Arrays**: char[] d and f for character processing
- ✅ **Encoding Tables**: Character lookup and transformation arrays
- ✅ **Text Processing**: Character manipulation and encoding

### **2. Stream Integration (MODERATE)**
Direct integration with Stream class (MBMGIXGO) for text processing:

**Stream Processing:**
```java
private static MBMGIXGO e;                    // Stream instance
public static java.lang.String a(int, boolean, MBMGIXGO);
```

- ✅ **Stream Parameter**: Takes MBMGIXGO (Stream) in processing methods
- ✅ **Data Reading**: Calls stream methods for character data
- ✅ **Text Decoding**: Processes encoded text from streams

### **3. Bit Manipulation Algorithms (MODERATE)**
Implements character encoding/decoding with bit operations:

**Bit Processing:**
```java
iload         6
iconst_4
ishr                                          // Shift right 4 bits
bipush        15
iand                                          // Mask with 15
```

- ✅ **Bit Shifting**: 4-bit shifts for encoding operations
- ✅ **Masking**: AND operations with bit masks
- ✅ **Character Encoding**: Bit-level text compression/decompression

### **4. String Processing (MODERATE)**
Handles string manipulation and character transformation:

**String Methods:**
```java
public static java.lang.String a(int, boolean, MBMGIXGO);
```

- ✅ **String Generation**: Creates strings from encoded data
- ✅ **Boolean Control**: Flag-based processing logic
- ✅ **Stream Input**: Reads character data from streams

### **5. Text Input Integration (MODERATE)**
Used in text input and character processing systems:

- ✅ **Character Validation**: Valid character set processing
- ✅ **Encoding/Decoding**: Text compression and decompression
- ✅ **Input Processing**: User text input handling
- ✅ **Stream Processing**: Text data from network/game files

## **SOURCE CODE CORRELATION**

### **TextInput.java (Reference Concept):**
```java
public class TextInput {
    public static char[] validChars = {
        ' ', 'e', 't', 'a', 'o', 'i', 'h', 'n', 's', 'r',
        'd', 'l', 'u', 'm', 'w', 'c', 'y', 'f', 'g', 'p',
        'b', 'v', 'k', 'x', 'j', 'q', 'z', '0', '1', '2',
        '3', '4', '5', '6', '7', '8', '9', ' ', '!', '?',
        '.', ',', ':', ';', '(', ')', '-', '&', '*', '\\',
        '\'', '@', '#', '+', '=', '\243', '$', '%', '"', '[',
        ']'
    };
    
    public static String processText(int param, boolean flag, Stream stream) {
        // Character processing with bit operations
        // Stream-based text decoding
        // Character array manipulation
    }
    
    // Text input and character encoding methods
}
```

## **UNIQUE IDENTIFIERS**
- **Character Arrays**: Static char[] d and f for text processing
- **Bit Operations**: 4-bit shifts and masking for encoding
- **Stream Integration**: MBMGIXGO parameter usage
- **String Generation**: Text decoding from encoded streams
- **Character Encoding**: Huffman-style text compression

## **MAPPING CONFIDENCE**
**70% CONFIDENCE** - The character array processing, Stream integration, and bit manipulation patterns provide moderate evidence for text processing functionality. The class clearly handles character encoding/decoding, though specific implementation details may vary.

## **IMPACT**
- Character encoding and decoding for game text
- Text input processing and validation
- Stream-based text compression and decompression
- Essential for user interface text handling and game communication