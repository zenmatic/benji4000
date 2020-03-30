package main

import (
	"fmt"
	"os"

	"github.com/uzudil/benji4000/bscript"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Specify a source file to run")
		return
	}
	r, err := os.Open(os.Args[1])
	if err != nil {
		fmt.Println("Unable to open source file:", err)
		return
	}
	defer r.Close()

	ast := &bscript.Program{}
	bscript.Parser.Parse(r, ast)

	// print the ast
	// fmt.Println(">>>>>>>>>>>>>>>>>>>>>>>>>>")
	// repr.Println(ast)
	// fmt.Println(">>>>>>>>>>>>>>>>>>>>>>>>>>")

	// run it
	fmt.Println("Running program:")
	err = ast.Evaluate()
	if err != nil {
		fmt.Println("Error running program:", err)
	}
}
