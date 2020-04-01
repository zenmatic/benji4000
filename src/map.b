# map demo

def foo(m)
    m["added in foo"] := "heck yeah";
end

def make_map()
    return { "xxx": 123, "yyy": "zaza" };
end

def main()
    # create a map
    map := { "a": 1, "b": 2, "c": 3 };
    print("map is " + map);
    print("value for 'a' is: " + map["a"]);
    print("value for 'b' is: " + map["b"]);
    print("value for 'c' is: " + map["c"]);
    print("value for non-existing key: " + map["d"]);

    # iterate the keys
    keys := keys(map);
    i := 0;
    while(i < len(keys))
        print("i=" + i + " key=" + keys[i]);
        i := i + 1;
    end

    # update same value
    map["b"] := map["b"] * 2;
    print("map after updating existing key, using the key: " + map);

    # assignment to map: update an existing key
    map["a"] := 22;
    print("map after updating existing key: " + map);

    # add a new key
    map["zorro"] := "hello world";
    print("map after adding a key: " + map);

    # delete a key + value from the map
    del map["b"];
    print("map after deleting a key: " + map);

    # pass by reference
    foo(map);
    print("map after pass by reference: " + map);

    # return by reference
    new_map := make_map();
    print("new map created in a function: " + new_map);
end
