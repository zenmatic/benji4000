def fib(x)
    if (x = 0)
        return 0;
    end
    if (x = 1)
        return 1;
    end
    return fib(x - 1) + fib(x - 2);
end

def fibseq(x)
    while(x >= 0)
        print(fib(x));
        let x = x - 1;
    end
end

def main()
    fibseq(10);
end
