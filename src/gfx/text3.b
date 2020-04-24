def main() {
    i := 128;
    x:=0;
    y:=0;    
    while(i < 512) {
        drawFont(x, y, 1, 0, i);
        i:=i+1;
        x:=x+1;
        if(x >= 40) {
            x:=0;
            y:=y+1;
        }
    }
}
