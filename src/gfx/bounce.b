# constants
const ne = 1;
const se = 2;
const nw = 3;
const sw = 4;

# global vars
x := 80;
y := 100;
dir := sw;
speed := 0.025;

def step() {
    if(dir = se) {
        x := x + speed;
        y := y + speed;
    }
    if(dir = ne) {
        x := x + speed;
        y := y - speed;
    }
    if(dir = sw) {
        x := x - speed;
        y := y + speed;
    }
    if(dir = nw) {
        x := x - speed;
        y := y - speed;
    }
}

def boundsCheck() {
    if(x >= 150) {
        if(dir = se) {
            dir := sw;
        }
        if(dir = ne) {
            dir := nw;
        }
    }
    if(x <= 10) {
        if(dir = nw) {
            dir := ne;
        }
        if(dir = sw) {
            dir := se;
        }
    }
    if(y >= 190) {
        if(dir = se) {
            dir := ne;
        }
        if(dir = sw) {
            dir := nw;
        }
    }
    if(y <= 10) {
        if(dir = ne) {
            dir := se;
        }
        if(dir = nw) {
            dir := sw;
        }
    }
}

def drawBorder() {
    drawLine(10, 10, 150, 10, 3);
    drawLine(10, 10, 10, 190, 3);
    drawLine(10, 190, 150, 190, 3);
    drawLine(150, 10, 150, 190, 3);
    drawLine(10, 10, 150, 190, 8);
    drawLine(10, 190, 150, 10, 8);
}

def main() {
    setVideoMode(2);
    while(dir != 0) {
        clearVideo();        
        drawBorder();
        fillCircle(x, y, 10, 4);
        updateVideo();
        step();
        boundsCheck();
    }
}
