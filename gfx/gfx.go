package gfx

import (
	"fmt"

	"github.com/veandco/go-sdl2/sdl"
)

const SCALE = 2.0
const WIDTH = 320
const HEIGHT = 200

// GfxTextMode is a 40x25 char text mode, 16 color background, 16 color foreground
const GfxTextMode = 0

// GfxHiresMode is a 320x200 pixels, 1 background color, each 8x8 block 16 colors graphics mode
const GfxHiresMode = 1

// GfxMultiColorMode is a 160x200 pixels (double wide) 16 colors graphics mode
const GfxMultiColorMode = 2

// Gfx is the "video card" api
type Gfx struct {
	// the graphics mode
	Mode uint8
	// the video memory
	PixelMemory [HEIGHT * WIDTH]uint8
	// character memory
	TextMemory [HEIGHT / 8 * WIDTH / 8]uint8
	// the color palette (16 rgb values)
	Colors [16 * 3]uint8
	// The background color (used by some modes)
	BackgroundColor uint8

	Window   *sdl.Window
	Renderer *sdl.Renderer
}

// NewGfx lets you create a new Gfx video card
func NewGfx() *Gfx {
	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		panic(err)
	}

	window, renderer, err := sdl.CreateWindowAndRenderer(int32(WIDTH*SCALE), int32(HEIGHT*SCALE), sdl.RENDERER_ACCELERATED|sdl.RENDERER_PRESENTVSYNC)
	if err != nil {
		panic(err)
	}

	renderer.SetScale(SCALE, SCALE)
	renderer.SetLogicalSize(WIDTH, HEIGHT)

	gfx := &Gfx{
		Mode:        GfxTextMode,
		PixelMemory: [HEIGHT * WIDTH]uint8{},
		TextMemory:  [HEIGHT / 8 * WIDTH / 8]uint8{},
		Colors: [16 * 3]uint8{
			0x00, 0x00, 0x00,
			0xff, 0xff, 0xff,
			0x88, 0x20, 0x00,
			0x68, 0xd0, 0xa8,
			0xa8, 0x38, 0xa0,
			0x50, 0xb8, 0x18,
			0x18, 0x10, 0x90,
			0xf0, 0xe8, 0x58,
			0xa0, 0x48, 0x00,
			0x47, 0x2b, 0x1b,
			0xc8, 0x78, 0x70,
			0x48, 0x48, 0x48,
			0x80, 0x80, 0x80,
			0x98, 0xff, 0x98,
			0x50, 0x90, 0xd0,
			0xb8, 0xb8, 0xb8,
		},
		BackgroundColor: 0,
		Window:          window,
		Renderer:        renderer,
	}

	return gfx
}

func (gfx *Gfx) MainLoop() {
	defer sdl.Quit()
	defer gfx.Window.Destroy()
	defer gfx.Renderer.Destroy()

	running := true
	var current_time, fps_last_time uint32 = 0, sdl.GetTicks()
	var c int
	var fps float32
	for running {
		current_time = sdl.GetTicks()
		c++

		// FPS counter math:
		if current_time-fps_last_time > 500 {
			fps = float32(c*1000) / float32(current_time-fps_last_time)
			gfx.Window.SetTitle(fmt.Sprintf("fps: %.2f", fps))
			fps_last_time = current_time
			c = 0
		}

		gfx.Renderer.SetDrawColor(0, 0, 0, 0)
		gfx.Renderer.Clear()

		switch {
		case gfx.Mode == GfxTextMode:
			gfx.renderTextMode()
		case gfx.Mode == GfxHiresMode:
			gfx.renderHiresMode()
		case gfx.Mode == GfxMultiColorMode:
			gfx.renderMultiColorMode()
		}
		gfx.Renderer.Present()

		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch event.(type) {
			case *sdl.QuitEvent:
				println("Quit")
				running = false
				break
			}
		}
	}
}

func (gfx *Gfx) renderTextMode() {
}

func (gfx *Gfx) renderHiresMode() {
	for index := 0; index < HEIGHT*(WIDTH/8); index++ {
		// convert index to pixel location
		y := index / (WIDTH / 8)
		x := index % (WIDTH / 8)
		b := gfx.PixelMemory[index]

		// get corresponding color info from text memory
		fg := gfx.TextMemory[(y/8)*WIDTH/8+x]
		gfx.Renderer.SetDrawColor(gfx.Colors[fg*3], gfx.Colors[fg*3+1], gfx.Colors[fg*3+2], 255)

		// draw pixel for each set bit at this location
		for bit := 0; bit < 8; bit++ {
			if b&(1<<bit) > 0 {
				gfx.Renderer.DrawPoint(int32(x*8+bit), int32(y))
			}
		}
	}
}

func (gfx *Gfx) renderMultiColorMode() {
	for index := 0; index < HEIGHT*(WIDTH/4); index++ {
		// convert index to pixel location
		y := index / (WIDTH / 4)
		x := index % (WIDTH / 4)
		b := gfx.PixelMemory[index]

		fg := b & 0x0f
		gfx.Renderer.SetDrawColor(gfx.Colors[fg*3], gfx.Colors[fg*3+1], gfx.Colors[fg*3+2], 255)
		gfx.Renderer.DrawPoint(int32(x*4), int32(y))
		gfx.Renderer.DrawPoint(int32(x*4+1), int32(y))

		fg = ((b >> 4) & 0x0f)
		gfx.Renderer.SetDrawColor(gfx.Colors[fg*3], gfx.Colors[fg*3+1], gfx.Colors[fg*3+2], 255)
		gfx.Renderer.DrawPoint(int32(x*4+2), int32(y))
		gfx.Renderer.DrawPoint(int32(x*4+3), int32(y))
	}
}

func (gfx *Gfx) SetPixel(x, y int, ch uint8, fg uint8, bg uint8) {

	switch {
	case gfx.Mode == GfxTextMode:
		panic("Implement me!")
	case gfx.Mode == GfxHiresMode:
		if x < 0 || x >= WIDTH || y < 0 || y >= HEIGHT {
			return
		}

		// pixel on/off bit
		index := x/8 + (y * WIDTH / 8)
		gfx.PixelMemory[index] |= (1 << (x % 8))

		// color info is stored in the text memory
		gfx.TextMemory[(y/8)*WIDTH/8+(x/8)] = fg
	case gfx.Mode == GfxMultiColorMode:
		if x < 0 || x >= WIDTH/2 || y < 0 || y >= HEIGHT {
			return
		}

		// pixel color: 2 double-wide pixels are stored in 1 byte
		index := x/2 + (y * WIDTH / 4)
		if x%2 == 0 {
			gfx.PixelMemory[index] |= fg
		} else {
			gfx.PixelMemory[index] |= (fg << 4)
		}
	}
}
