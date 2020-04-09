package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/alecthomas/repr"
	"github.com/uzudil/benji4000/bscript"
)

func main() {
	var source string
	flag.StringVar(&source, "source", "src/test1.b", "the bscript file to run")
	showAst := flag.Bool("ast", false, "print AST and not execute?")
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
	} else {
		// run it
		err = ast.Evaluate()
		if err != nil {
			fmt.Println("Error running program:", err)
		}
	}
}
