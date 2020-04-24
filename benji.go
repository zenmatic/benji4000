package main

import (
	"flag"

	"github.com/uzudil/benji4000/bscript"
	"github.com/uzudil/benji4000/gfx"
)

func repl(video *gfx.Gfx) {
	bscript.Repl(video)
}

func main() {
	var source string
	flag.StringVar(&source, "source", "", "the bscript file to run")
	showAst := flag.Bool("ast", false, "print AST and not execute?")
	flag.Parse()

	video := gfx.NewGfx()

	if source != "" {
		go func() {
			bscript.Run(source, showAst, nil, video)
		}()
	} else {
		go repl(video)
	}

	video.Render.MainLoop()
}
