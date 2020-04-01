# global vars demo

value := 10;


def foo()
    print("in foo, value=" + value);
    value := 30;
end

def foo2(x)
    print("in foo2, x=" + x);
    x := x * 100;
    print("in foo2, updated x=" + x);
    return x;
end

def foo3(value)
    print("in foo3, shadowed value=" + value);
    value := 1000;
    print("in foo3, updated shadowed value=" + value);
end

def main()
    print("value=" + value);

    value := value * 2;
    print("value=" + value);

    foo();
    print("after foo, value=" + value);

    new_value := foo2(value);
    print("after foo2, value=" + value);
    print("after foo2, new_value=" + new_value);

    foo3(5);
    print("after foo3, value=" + value);
end