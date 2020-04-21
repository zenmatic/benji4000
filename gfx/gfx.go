package gfx

import (
	"math"
)

const (
	Width  = 320
	Height = 200

	// GfxTextMode is a 40x25 char text mode, 16 color background, 16 color foreground
	GfxTextMode = 0

	// GfxHiresMode is a 320x200 pixels, 1 background color, each 8x8 block 16 colors graphics mode
	GfxHiresMode = 1

	// GfxMultiColorMode is a 160x200 pixels (double wide) 16 colors graphics mode
	GfxMultiColorMode = 2
)

// Gfx is the "video card" api
type Gfx struct {
	// the video mode
	VideoMode int
	// video memory
	VideoMemory [Width * Height]byte
	// text memory
	TextMemory [Width / 8 * Height / 8]byte
	// the actual renderer
	Render *Render
	// Color definitions
	Colors [16 * 3]uint8
	// the global background color
	BackgroundColor int
}

// NewGfx lets you create a new Gfx video card
func NewGfx() *Gfx {
	return &Gfx{
		VideoMode:   GfxTextMode,
		VideoMemory: [Width * Height]byte{},
		TextMemory:  [Width / 8 * Height / 8]byte{},
		Render:      NewRender(),
		Colors: [16 * 3]uint8{
			// C64 colors :-)
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
	}
}

func (gfx *Gfx) DrawCircle(x, y, r int, fg uint8) error {
	return gfx.circle(x, y, r, fg, false)
}

func (gfx *Gfx) FillCircle(x, y, r int, fg uint8) error {
	return gfx.circle(x, y, r, fg, true)
}

// 90 degrees in radians
const rad90 = 0.5 * math.Pi

// there is probably a more efficient way to do this
func (gfx *Gfx) circle(x, y, r int, fg uint8, filled bool) error {
	var px, py float64
	circleSteps := int(math.Min(float64(r*2), 100))
	for a := 0; a <= circleSteps; a++ {
		rad := (float64(a) / float64(circleSteps)) * rad90
		dx := float64(r) * math.Cos(rad)
		dy := float64(r) * math.Sin(rad)
		if filled {
			gfx.FillRect(int(float64(x)+dx), int(float64(y)+dy), int(float64(x)+px), int(float64(y)-dy), fg)
			gfx.FillRect(int(float64(x)-px), int(float64(y)+dy), int(float64(x)-dx), int(float64(y)-dy), fg)
		} else if !(px == 0 && py == 0) {
			gfx.DrawLine(int(float64(x)+dx), int(float64(y)+dy), int(float64(x)+px), int(float64(y)+py), fg)
			gfx.DrawLine(int(float64(x)-dx), int(float64(y)+dy), int(float64(x)-px), int(float64(y)+py), fg)
			gfx.DrawLine(int(float64(x)-dx), int(float64(y)-dy), int(float64(x)-px), int(float64(y)-py), fg)
			gfx.DrawLine(int(float64(x)+dx), int(float64(y)-dy), int(float64(x)+px), int(float64(y)-py), fg)
		}
		px = dx
		py = dy
	}
	return nil
}

func (gfx *Gfx) DrawRect(x, y, x2, y2 int, fg uint8) error {
	gfx.DrawLine(x, y, x2, y, fg)
	gfx.DrawLine(x, y2, x2, y2, fg)
	gfx.DrawLine(x, y, x, y2, fg)
	gfx.DrawLine(x2, y, x2, y2, fg)
	return nil
}

func (gfx *Gfx) FillRect(x, y, x2, y2 int, fg uint8) error {
	dx := 1
	if x > x2 {
		dx = -1
	}
	dy := 1
	if y > y2 {
		dy = -1
	}
	for xx := x; xx != x2; xx += dx {
		for yy := y; yy != y2; yy += dy {
			gfx.SetPixel(xx, yy, 0, fg, 0)
		}
	}
	return nil
}

func (gfx *Gfx) DrawLine(x, y, x2, y2 int, fg uint8) error {
	sx := float64(x)
	sy := float64(y)
	ex := float64(x2)
	ey := float64(y2)

	if math.Abs(float64(x)-float64(x2)) > math.Abs(float64(y)-float64(y2)) {
		// walk along x
		if x > x2 {
			sx, ex = ex, sx
			sy, ey = ey, sy
		}
		dy := (ey - sy) / (ex - sx)
		yy := sy
		for xx := sx; xx <= ex; xx++ {
			gfx.SetPixel(int(xx), int(yy), 0, fg, 0)
			yy += dy
		}
	} else {
		// walk along y
		if y > y2 {
			sy, ey = ey, sy
			sx, ex = ex, sx
		}
		dx := (ex - sx) / (ey - sy)
		xx := sx
		for yy := sy; yy <= ey; yy++ {
			gfx.SetPixel(int(xx), int(yy), 0, fg, 0)
			xx += dx
		}
	}

	return nil
}

func (gfx *Gfx) SetPixel(x, y int, ch, fg, bg uint8) error {
	switch {
	case gfx.VideoMode == GfxTextMode:
	case gfx.VideoMode == GfxHiresMode:
		if x >= 0 && y >= 0 && x < Width && y < Height {
			// set the pixel asked for
			gfx.VideoMemory[y*Width+x] = byte(fg)

			// set other pixels (if >0) in this 8x8 area
			bx := (x / 8) * 8
			by := (y / 8) * 8
			for xx := 0; xx < 8; xx++ {
				for yy := 0; yy < 8; yy++ {
					addr := (by+yy)*Width + (bx + xx)
					if gfx.VideoMemory[addr] > 0 {
						gfx.VideoMemory[addr] = byte(fg)
					}
				}
			}
		}
	case gfx.VideoMode == GfxMultiColorMode:
		if x >= 0 && y >= 0 && x < Width/2 && y < Height {
			gfx.VideoMemory[y*Width+x*2] = byte(fg)
			gfx.VideoMemory[y*Width+x*2+1] = byte(fg)
		}
	}
	return nil
}

func (gfx *Gfx) ClearVideo() error {
	for i := range gfx.VideoMemory {
		gfx.VideoMemory[i] = 0
	}
	return nil
}

func (gfx *Gfx) UpdateVideo() error {
	gfx.Render.Lock.Lock()
	for index, colorIndex := range gfx.VideoMemory {
		r, g, b := gfx.Colors[colorIndex*3], gfx.Colors[colorIndex*3+1], gfx.Colors[colorIndex*3+2]
		gfx.Render.PixelMemory[index*3] = r
		gfx.Render.PixelMemory[index*3+1] = g
		gfx.Render.PixelMemory[index*3+2] = b
	}
	gfx.Render.Lock.Unlock()
	// runtime.Gosched()
	return nil
}
