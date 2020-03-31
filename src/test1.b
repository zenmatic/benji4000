def b(x)
    x := x * 2;
    print("x in b:");
    print(x);
    return x;
end

def a(x)
    x := x * x;
    print("x in a:");
    print(x);
    x := b(x);
    print("x in a again:");
    print(x);
    return x;
end

def main()

    # do some math
    x := 2 * (7 - 2);
    x := a(x);

    # show it
    print("x is currently:");
    print(x);

    return x;
end