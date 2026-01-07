import java.io.File;

import javax.sound.midi.*;

/**
 * Implementation of AudioBox for MIDI music playback. Plays the track once (matches original
 * loop="0"). Stop is triggered by setting src to the silent path.
 */
public class MidiAudioBox implements AudioBox {
    private Sequencer sequencer;
    private String currentSrc;
    private int currentVolume = 127;

    @Override
    public void setSrc(String path) {
        if (getSilentPath().equals(path)) {
            stopInternal();
            return;
        }

        stopInternal(); // Stop any current playback first

        this.currentSrc = path;
        try {
            File file = new File(path);
            if (!file.exists()) return;

            Sequence sequence = MidiSystem.getSequence(file);
            sequencer = MidiSystem.getSequencer();
            sequencer.open();
            sequencer.setSequence(sequence);
            sequencer.setLoopCount(0); // Play once only — matches loop="0"

            setVolume(currentVolume);
            sequencer.start();
        } catch (Exception e) {
            e.printStackTrace();
        }
    }

    @Override
    public void setVolume(int volume) {
        currentVolume = Math.max(0, Math.min(127, volume));
        if (sequencer instanceof Synthesizer) {
            Synthesizer synth = (Synthesizer) sequencer;
            for (MidiChannel channel : synth.getChannels()) {
                if (channel != null) {
                    channel.controlChange(7, currentVolume);
                }
            }
        }
    }

    @Override
    public String getSilentPath() {
        return "c:\\silence.mid"; // Path matches original JavaScript implementation
    }

    private void stopInternal() {
        if (sequencer != null) {
            try {
                sequencer.stop();
                sequencer.close();
            } catch (Exception ignored) {
            }
            sequencer = null;
        }
    }
}
