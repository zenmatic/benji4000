# benji4000
The "other" lost personal computer of the 80s. 
It was similar to the better known C64 but programmed in a higher-level (and much more imaginary) language.

# To build the code
`go build`

# To run a bscript file

For example:
`./benji4000 -source=src/adventure.b`

# bscript
The programming language of benji. Execution starts by calling the function named "main".

## Features:
- single line comments: `# this is a comment`
- variable declarations: `a := 1;` Global variables are declared outside of any function. Variable values can be a number, a string, an array or a map.
- constants: `const PI=3.14159;`
- strings: `a := "hello";`
- control flow: `if(a = 1) { doSomething(); } else { doSomethingElse(); }`
- loop: `while(a < 10) { a := a + 1; }`
- arrays: `a := [1, 2, 3];`
- maps: `a := { "a": 1, "b": 2 };` Map keys are always strings, values can be anything (including other maps.)
- function definitions: `def hello(x) { print(x); }`
- function calls: `f(g(123));`
- builtin functions:
   - length: the length of a string, array or map
   - keys: returns a map's keys as an array (always strings)
   - substr: substring, for example: `substr("Hello World", 1, 2)` prints "el"
   - print: print strings + variables
   - input: ask for user input
   - debug: print closures and stack trace
   - assert: assertion testing
- first class functions: `def f(x) { return 2; } x := f;`
   - functions as parameters
   - anonymous functions: `def f(x) { return (n) => { return x + n; }; }`
   
## Coming soon:
- boolean operators (and, or, not)

## bscript syntax highlighting
The vscode directory contains a plugin for syntax highlighting for .b files.

  
