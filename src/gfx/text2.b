def main() {
    setVideoMode(0);
    x := 0;
    y := 0;
    while(y < 25) {
        fg := random() * 16;
        bg := random() * 16;
        # todo: make this not a magic number (288 fonts total)    
        ch := 128 + (random() * (288 - 128));
        drawFont(x, y, fg, bg, ch);
        x := x + 1;
        if(x >= 40) {
            x := 0;
            y := y + 1;
        }
    }
    updateVideo();
}
