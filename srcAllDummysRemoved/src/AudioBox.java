/**
 * Interface representing a <bgsound> element from the original JavaScript. Only exposes src and
 * volume setters — stop is handled by setting a silent src.
 */
public interface AudioBox {

    /**
     * Sets the source file path to play (e.g., "jingle1.mid"). Setting to getSilentPath() triggers
     * stop.
     */
    void setSrc(String path);

    /** Sets the volume level (0-127). */
    void setVolume(int volume);

    /**
     * Returns the special "silent" path that triggers stop (e.g., "c:\\silence.mid" or equivalent).
     */
    String getSilentPath();
}
