
def main()
    # declare an array
    let a = [ 1, 2, 3, 4, 5 ];
    
    # print it, or a part of it
    print(a);
    print(a[2]);

    # array element assignment
    let a[0] = 13;
    print(a[0]);

    # array element on RHS
    let x = a[0];
    print(x);

    # array element dynamic access
    let i = 0;
    while (i < 5)
        let a[i] = a[i] * 2;
        let i = i + 1;
    end
    print(a);
end
