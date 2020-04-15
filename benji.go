package main

import (
	"github.com/uzudil/benji4000/bscript"
	"github.com/uzudil/benji4000/gfx"
)

func repl(video *gfx.Gfx) {
	bscript.Repl(video)
}

func main() {
	// var source string
	// flag.StringVar(&source, "source", "src/test1.b", "the bscript file to run")
	// showAst := flag.Bool("ast", false, "print AST and not execute?")
	// repl := flag.Bool("repl", false, "Enter interactive mode?")
	// flag.Parse()

	video := gfx.NewGfx()
	// video.Mode = gfx.GfxHiresMode
	// for i := 0; i < 320; i++ {
	// 	for t := 0; t < 100; t++ {
	// 		video.SetPixel(i+t, i, 0, uint8(rand.Intn(16)), 0)
	// 	}
	// }

	// video.Mode = gfx.GfxMultiColorMode
	// for i := 0; i < 160; i++ {
	// 	for t := 0; t < 100; t++ {
	// 		video.SetPixel(i+t, i, 0, uint8(rand.Intn(16)), 0)
	// 	}
	// }

	go repl(video)

	video.MainLoop()

	// if *repl {
	// 	bscript.Repl()
	// } else {
	// 	bscript.Run(source, showAst, nil)
	// }
}
