def b(x) {
    debug("in b");
    x := x * 2;
    print("x in b:");
    print(x);
    return x;
}

def a(x) {
    debug("in a");
    x := x * x;
    print("x in a:");
    print(x);
    x := b(x);
    print("x in a again:");
    print(x);
    return x;
}

def main() {

    # do some math
    x := 2 * (7 - 2);
    print("x=" + x);
    debug("before");
    a(x);
    debug("after");

    # show it
    print("x is currently:");
    print(x);

    return x;
}
