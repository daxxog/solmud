import javax.sound.sampled.*;
import java.io.File;

/**
 * Implementation of AudioBox for WAV sound effects playback.
 * Plays the sound once (no looping).
 * No special silent path needed — setting a new src stops previous.
 */
public class WaveAudioBox implements AudioBox {
    private Clip clip;
    private String currentSrc;
    private int currentVolume = 127;

    @Override
    public void setSrc(String path) {
        stopInternal(); // Stop previous WAV first
        this.currentSrc = path;
        try {
            File file = new File(path);
            if (!file.exists()) return;

            AudioInputStream audioInputStream = AudioSystem.getAudioInputStream(file);
            clip = AudioSystem.getClip();
            clip.open(audioInputStream);

            setVolume(currentVolume);
            clip.start();
        } catch (Exception e) {
            e.printStackTrace();
        }
    }

    @Override
    public void setVolume(int volume) {
        currentVolume = Math.max(0, Math.min(127, volume));
        if (clip != null && clip.isOpen()) {
            FloatControl gainControl = (FloatControl) clip.getControl(FloatControl.Type.MASTER_GAIN);
            if (gainControl != null) {
                float range = gainControl.getMaximum() - gainControl.getMinimum();
                float gain = gainControl.getMinimum() + (currentVolume / 127f) * range;
                gainControl.setValue(gain);
            }
        }
    }

    @Override
    public String getSilentPath() {
        return "";  // unused in WaveAudioBox
    }

    private void stopInternal() {
        if (clip != null) {
            clip.stop();
            clip.close();
            clip = null;
        }
    }
}
