package main

import (
	"flag"

	"github.com/uzudil/benji4000/bscript"
)

func main() {
	var source string
	flag.StringVar(&source, "source", "src/test1.b", "the bscript file to run")
	showAst := flag.Bool("ast", false, "print AST and not execute?")
	repl := flag.Bool("repl", false, "Enter interactive mode?")
	flag.Parse()

	if *repl {
		bscript.Repl()
	} else {
		bscript.Run(source, showAst, nil)
	}
}
