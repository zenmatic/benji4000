def a() {
    def b() {
        debug("in b");
    }

    def c() {
        debug("In c");
        b();
    }

    debug("In a");
    c();
}

def main() {
    a();
    // b(); this should not work
}
