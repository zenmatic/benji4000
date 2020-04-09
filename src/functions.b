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

def multiply(n) {
    def inner(m) {
        return n * m;
    }
    return inner;
}

def map(fx, array) {
    i := 0;
    while(i < len(array)) {
        array[i] := fx(array[i]);
        i := i + 1;
    }
}

def double2() {
    return (x) => { 
        return x * 2; 
    };
}


def main() {
    # call a function
    x := double(2);
    print("x=" + x);

    # call an embedded function
    y := triple(2);
    print("y=" + y);

    # function references
    five_x := multiply(5);
    print("5 * 6=" + five_x(6));

    # double parens
    print("4 * 6=" + multiply(4)(6));

    # map
    a := [1, 2, 3, 4, 5];
    print("a=" + a);
    map(double, a);
    print("doubled a=" + a);

    # anonymous functions (function literals)
    map((n) => {
        return n + 1;
    }, a);
    print("added 1 to a=" + a);

    # anon function returned
    anon := double2();
    print("2 * 5=" + anon(5));
}
