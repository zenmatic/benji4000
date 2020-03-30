def b(x)
    let x = x * 2;
    print("x in b:");
    print(x);
    return x;
end

def a(x)
    let x = x * x;
    print("x in a:");
    print(x);
    let x = b(x);
    print("x in a again:");
    print(x);
    return x;
end

def main()

    # do some math
    let x = 2 * (7 - 2);
    let x = a(x);

    # show it
    print("x is currently:");
    print(x);

    return x;
end