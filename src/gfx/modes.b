def fillScreen(w, h) {
    index := 0;
    c := 0;
    while(index < w * h) {
        y := index / w;
        x := index % w;
        setPixel(x, y, c);
        c := random() * 16;
        index := index + 1;
    }
}

def main() {
    setVideoMode(1); # high res mode
    fillScreen(320, 200);

    input(">>>");

    setVideoMode(2); # multi-color mode
    fillScreen(160, 200);
}
