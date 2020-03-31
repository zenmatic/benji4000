# map demo
def main()
    let map = { "a": 1, "b": 2, "c": 3 };
    let keys = keys(map);
    let i = 0;
    while(i < len(k))
        let k = keys[i];

        print("key=");
        print(k);
        print("value=");
        print(map[k]);

        map[k] = map[k] * 2;

        let i = i + 1;
    end

    print(map);
end
