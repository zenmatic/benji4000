
def main()
    # declare an array
    a := [ 1, 2, 3, 4, 5 ];
    
    # print it, or a part of it
    print("a=" + a);
    print(a[2]);

    # array element assignment
    a[0] := 13;
    print(a[0]);

    # array element on RHS
    x := a[0];
    print(x);

    # array element dynamic access
    i := 0;
    while (i < 5)
        a[i] := a[i] * 2;
        i := i + 1;
    end
    print(a);
end
