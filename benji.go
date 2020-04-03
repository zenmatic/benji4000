package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/alecthomas/repr"
	"github.com/uzudil/benji4000/bscript"
)

func main2() {
	a := map[int]int{10: 1, 11: 2, 12: 3}
	b := a
	fmt.Printf("b addr=%p value=%v\n", &b, b)
	fmt.Printf("a addr=%p value=%v\n", &a, a)

	delete(a, 11)
	fmt.Printf("b addr=%p value=%v\n", &b, b)
	fmt.Printf("a addr=%p value=%v\n", &a, a)

}

func main() {
	var source string
	flag.StringVar(&source, "source", "src/test1.b", "the bscript file to run")
	showAst := flag.Bool("ast", false, "print AST?")
	evalAst := flag.Bool("eval", true, "eval code?")
	flag.Parse()

	r, err := os.Open(source)
	if err != nil {
		fmt.Println("Unable to open source file:", err)
		return
	}
	defer r.Close()

	ast := &bscript.Program{}
	bscript.Parser.Parse(r, ast)
	if *showAst {
		// print the ast
		fmt.Println(">>>>>>>>>>>>>>>>>>>>>>>>>>")
		repr.Println(ast)
		fmt.Println(">>>>>>>>>>>>>>>>>>>>>>>>>>")
	}

	// run it
	if *evalAst {
		err = ast.Evaluate()
		if err != nil {
			fmt.Println("Error running program:", err)
		}
	}
}
