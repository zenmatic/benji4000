# multi-dimensional array testing
a := [
    [1, 2, 3, 4, 5],
    ["a", "b", "c", "d", "e"]
];

def foo() {
    b := a[1];
    del b[2];
}

def main() {

    print(a);

    b := a[0];
    print(b);

    del b[1]; # maybe del should be a function
    print(b);
    print(a);

    # foo();
    # print(a);

    # need support for: print(a[0][1]);
}