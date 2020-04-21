r := 10;
dir := 0.3;

def main() {
    setVideoMode(1);
    while(dir != 0) {
        clearVideo();
        drawCircle(160, 100, r + 20, 5);
        fillCircle(160, 100, r + 10, 7);
        r := r + dir;
        if(r >= 150) {
            dir := -1 * dir;
        }
        if(r <= 5) {
            dir := -1 * dir;
        }
        updateVideo();
    }
}
