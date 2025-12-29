/**
 * Polling thread that emulates the original JavaScript musicloop() function.
 * Runs every 50ms to check for new MIDI/WAV files and handle playback.
 */
public class MusicLoop implements Runnable {
    private static final int POLL_INTERVAL_MS = 50;          // Original setInterval(..., 50)
    private static final int FADE_INITIAL_VALUE = 200;       // Starting fade value
    private static final int FADE_STEPS = 144;               // Fade completion threshold
    private static final int FADE_VOLUME_MULTIPLIER = 25;    // Volume change per step

    private final AudioAccessor accessor;
    private final AudioBox midibox;
    private final AudioBox wavebox;
    private volatile boolean running = true;

    public MusicLoop(AudioAccessor accessor, MidiAudioBox midibox, WaveAudioBox wavebox) {
        this.accessor = accessor;
        this.midibox = midibox;
        this.wavebox = wavebox;
    }

    @Override
    public void run() {
        while (running) {
            try {
                Thread.sleep(POLL_INTERVAL_MS);
                musicloop();
            } catch (InterruptedException e) {
                Thread.currentThread().interrupt();
                break;
            } catch (Exception e) {
                e.printStackTrace();
            }
        }
    }

    private void musicloop() {
        int fade = FADE_INITIAL_VALUE;
        String fademidi = null;
        int fadevol = 0;

        // MIDI
        String midi = accessor.getmidi();
        if (!"none".equals(midi)) {
            int midivol = accessor.getmidivol();
            int midifade = accessor.getmidifade();

            if ("voladjust".equals(midi) && fademidi != null) {
                fadevol = midivol;
            } else {
                fademidi = null;

                if ("stop".equals(midi)) {
                    midibox.setSrc(midibox.getSilentPath());
                    fade = FADE_INITIAL_VALUE;
                } else if ("voladjust".equals(midi)) {
                    midibox.setVolume(midivol);
                    if (midifade == 1) fade = -(midivol / FADE_VOLUME_MULTIPLIER);
                    else fade = FADE_INITIAL_VALUE;
                } else if (midifade == 1) {
                    fademidi = midi;
                    fadevol = midivol;
                } else {
                    midibox.setSrc(midi);
                    midibox.setVolume(midivol);
                    fade = FADE_INITIAL_VALUE;
                }
            }
        }

        // Fade logic
        if (fademidi != null) {
            fade += 1;
            midibox.setVolume(-(fade * FADE_VOLUME_MULTIPLIER));
            if (fade >= FADE_STEPS) {
                midibox.setSrc(fademidi);
                midibox.setVolume(fadevol);
                fademidi = null;
                fade = -(fadevol / FADE_VOLUME_MULTIPLIER);
            }
        }

        // WAV
        String wave = accessor.getwave();
        if (!"none".equals(wave)) {
            int wavevol = accessor.getwavevol();
            wavebox.setSrc(wave);
            wavebox.setVolume(wavevol);
        }
    }

    public void stop() {
        running = false;
        midibox.setSrc(midibox.getSilentPath());
        wavebox.setSrc(wavebox.getSilentPath());
    }
}
