# RuneScape Build 317 Audio Playback System

**A Historical Analysis of How MIDI Music and Sound Effects Were Handled in 2005–2006**

## Overview

In RuneScape Classic and early RS2 (build ~317, circa 2005–2006), **the client did not play audio itself**.  
The system was a hybrid architecture involving three layers:

1. **Main gamepack** (obfuscated JAR): Received audio data from server/cache and saved it to temporary files in the cache directory.
2. **Loader applet** (separate small JAR): Bootstrapped the gamepack and exposed public getter methods to read file paths.
3. **JavaScript on the launch webpage**: Polled the loader applet via LiveConnect and triggered browser-native playback using `<bgsound>` tags (IE-specific).

This design offloaded audio playback to the browser, bypassing applet sandbox limitations.

## 1. Main Gamepack (gamepack_317.jar)

**Role**: Dumps MIDI music and WAV sound effects to temporary files.

**Key Class**: `sign/signlink.class`

**Relevant Methods** (decompiled):

```java
public static final synchronized void midisave(byte[] abyte0, int i) {
    if (i > 0x1e8480)
        return;
    if (savereq != null)
        return;
    midipos = (midipos + 1) % 5;
    savelen = i;
    savebuf = abyte0;
    midiplay = true;
    savereq = "jingle" + midipos + ".mid";
}

public static final synchronized boolean wavesave(byte[] abyte0, int i) {
    if (i > 0x1e8480)
        return false;
    if (savereq != null)
        return false;
    wavepos = (wavepos + 1) % 5;
    savelen = i;
    savebuf = abyte0;
    waveplay = true;
    savereq = "sound" + wavepos + ".wav";
    return true;
}
```

**Flow**:
- Server sends music packet (opcode 74) → calls `midisave(data, length)`.
- Saves to cache dir (e.g., `./cache/jingle2.mid`).
- Sets `signlink.midi = full path`.

**No playback code** — only file writing and field setting.

## 2. Loader Applet (loader317.jar)

**Role**: Downloads and launches the main gamepack; exposes public getters for JavaScript.

**Key Methods** (decompiled):

```java
public final String getmidi() {
    if (signlink.midi == null) {
        return "none";
    } else {
        String var1 = signlink.midi;
        signlink.midi = null;  // Clears after read
        return var1;
    }
}

public final int getmidivol() {
    return signlink.midivol;
}

public final int getmidifade() {
    return signlink.midifade;
}

public final String getwave() {
    if (signlink.wave == null) {
        return "none";
    } else {
        String var1 = signlink.wave;
        signlink.wave = null;  // Clears after read
        return var1;
    }
}

public final int getwavevol() {
    return signlink.wavevol;
}
```

These getters allow JavaScript to poll and consume paths without race conditions.

## 3. Webpage JavaScript (from preserved page)

**Role**: Polls the loader applet via LiveConnect and plays audio natively.

**HTML Tags** (used for playback):

```html
<bgsound id="midibox" loop="0" volume="0" autostart="yes">
<bgsound id="wavebox" loop="0" volume="0" autostart="yes">
```

**Full Audio Polling & Playback Logic** (decompiled from preserved page):

```javascript
var fade = 200;
var fademidi = null;
var fadevol = 0;

function givetrue() {
    return true;
}

function musicloop() {
    var midi = document.app.getmidi();
    if (midi != "none") {
        var midivol = document.app.getmidivol();
        var midifade = document.app.getmidifade();
        if (midi == "voladjust" && fademidi != null) {
            fadevol = midivol;
            return;
        }
        fademidi = null;
        if (midi == "stop") {
            midibox.src = "c:\\silence.mid";
            fade = 200;
            return;
        }
        if (midi == "voladjust") {
            midibox.volume = midivol;
            if (midifade == 1) fade = -(midivol / 25);
            else fade = 200;
            return;
        }
        if (midifade == 1) {
            fademidi = midi;
            fadevol = midivol;
        } else {
            midibox.src = midi;
            midibox.volume = midivol;
            fade = 200;
        }
    }
    if (fademidi != null) {
        fade = fade + 1;
        midibox.volume = -(fade * 25);
        if (fade >= 144) {
            midibox.src = fademidi;
            midibox.volume = fadevol;
            fademidi = null;
            fade = -(fadevol / 25);
        }
    }
    var wave = document.app.getwave();
    if (wave != "none") {
        var wavevol = document.app.getwavevol();
        wavebox.src = wave;
        wavebox.volume = wavevol;
    }
}

function startloop() {
    onerror = givetrue;
    setInterval("musicloop()", 50);
}
```

**How `musicloop()` is Called**:
- `<body onload="startloop();">` triggers `startloop()` on page load.
- `startloop()` sets up error handling and starts the polling loop:

```javascript
function startloop() {
    onerror = givetrue;
    setInterval("musicloop()", 50);  // Polls every 50ms
}
```

## How Everything Works Together

1. **Server → Gamepack**:
   - Sends music/SFX data via packets.
   - Gamepack saves to temp files and sets `signlink.midi` / `signlink.wave`.

2. **Loader Applet**:
   - Bootstraps gamepack.
   - Exposes public getters (`getmidi()`, `getwave()`, etc.).

3. **Webpage JavaScript**:
   - `onload="startloop();"` runs on page load.
   - `startloop()` starts a 50ms interval calling `musicloop()`.
   - `musicloop()` polls `document.app.getmidi()` / `getwave()` via LiveConnect.
   - Updates `<bgsound src="...">` with file paths.
   - Handles volume/fade using `getmidivol()` / `getmidifade()`.

4. **Browser (IE)**:
   - Plays local `.mid`/`.wav` files via `<bgsound>`.

This hybrid system explains why audio "just worked" in the official client (in IE) but never in decompiled/standalone versions — the playback was **entirely in JavaScript on the webpage**.
