# benji4000
The "other" lost personal computer of the 80s. 
It was similar to the better known C64 but programmed in a higher-level (and much more imaginary) language.

# To build the code
`go build`

# To run a bscript file

For example:
`./benji4000 -source=src/adventure.b`

# bscript
The programming language of benji. 

## Features:
- single line comments: `# this is a comment`
- variable declarations: `a := 1;`
- strings: `a := "hello";`
- control flow: `if(a = 1) { doSomething(); } else { doSomethingElse(); }`
- loop: `while(a < 10) { a := a + 1; }`
- arrays: `a := [1, 2, 3];`
- maps: `a := { "a": 1, "b": 2 };`
- builtin functions:
   - length: the length of a string, array or map
   - keys: returns a map's keys as an array (always strings)
   - substr: substring, for example: `substr("Hello World", 1, 2)` prints "el"
   - print: print strings + variables
   - input: ask for user input
   
  ## Coming soon:
  - first-class functions
  - closures
  - boolean operators (and, or, not)
  
