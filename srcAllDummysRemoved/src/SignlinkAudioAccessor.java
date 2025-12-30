import sign.signlink;

public class SignlinkAudioAccessor implements AudioAccessor {

    @Override
    public final String getmidi() {
        String path = signlink.midi;
        if (path == null || path.isEmpty()) {
            return "none";
        }
        signlink.midi = null; // Clear after read
        return path;
    }

    @Override
    public final int getmidivol() {
        return signlink.midivol;
    }

    @Override
    public final int getmidifade() {
        return signlink.midifade;
    }

    @Override
    public final String getwave() {
        String path = signlink.wave;
        if (path == null || path.isEmpty()) {
            return "none";
        }
        signlink.wave = null; // Clear after read
        return path;
    }

    @Override
    public final int getwavevol() {
        return signlink.wavevol;
    }
}
