# try various spacings
a:=1;
b:=2+3;
c := 3 + 4;
# and negative numbers
d := -1;

def main() {
    assert(a, 1);
    assert(b, 5);
    assert(c, 7);
    assert(d, -1);
    assert(b * d, -5);

}
