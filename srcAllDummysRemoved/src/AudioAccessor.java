/**
 * Interface matching the public methods exposed by the original loader applet for JavaScript
 * (LiveConnect) to poll and play audio.
 */
public interface AudioAccessor {

    /**
     * Returns the path to the latest MIDI file (e.g., "jingle2.mid") or "none". Clears the
     * underlying field after reading.
     */
    String getmidi();

    /** Returns current MIDI volume (0-127). */
    int getmidivol();

    /** Returns MIDI fade flag (0 or 1). */
    int getmidifade();

    /**
     * Returns the path to the latest WAV sound effect or "none". Clears the underlying field after
     * reading.
     */
    String getwave();

    /** Returns current WAV volume (0-127). */
    int getwavevol();
}
