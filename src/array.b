
def main()
    # declare an array
    a := [ 1, 2, 3, 4, 5 ];
    
    # print it, or a part of it
    print("a=" + a);
    print("at index=2, " + a[2]);

    # array element assignment
    a[0] := 13;
    print("After setting index 0, " + a[0]);

    # array element on RHS
    x := a[0];
    print("Afer reading index 0, " + x);

    # array element dynamic access
    i := 0;
    while (i < len(a))
        a[i] := a[i] * 2;
        i := i + 1;
    end
    print("After dynamic access " + a);

    # append to an array by adding a new element at the end
    i := 0;
    while(i < 3)
        a[len(a)] := i;
        i := i + 1;
    end
    print("After adding 3 elements " + a);

    # delete from array
    del a[0];
    print("After removing element 0: " + a + " length=" + len(a));
    del a[3];
    print("After removing element 3: " + a + " length=" + len(a));
end
