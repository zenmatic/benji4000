def fib(x) {
    if (x = 0) {
        return 0;
    }
    if (x = 1) {
        return 1;
    }
    return fib(x - 1) + fib(x - 2);
}

def main() {
    x := 0;
    while(x < 10) {
        print("at " + x + " fib=" + fib(x));
        x := x + 1;
    }
    return x;
}
