# first-class functions demo

def double(n) {
    return n * 2;
}

def triple(n) {

    def mul(m) {
        return n * m;
    }

    return mul(3);
}


def main() {
    # call a function
    x := double(2);
    print("x=" + x);

    # call an embedded function
    y := triple(2);
    print("y=" + y);
}
