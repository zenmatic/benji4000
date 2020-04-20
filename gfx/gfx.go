package gfx

import (
	"fmt"
	"log"
	"runtime"
	"strings"
	"sync"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
)

const (
	Width  = 320
	Height = 200
	scale  = 2

	vertexShaderSource = `
		#version 410
		layout (location = 0) in vec3 aPos;
		layout (location = 1) in vec3 aColor;
		layout (location = 2) in vec2 aTexCoord;

		out vec3 ourColor;
		out vec2 TexCoord;

		void main() {
			gl_Position = vec4(aPos, 1.0);
			ourColor = aColor;
			TexCoord = aTexCoord;
		}
	` + "\x00"

	fragmentShaderSource = `
		#version 410
		out vec4 FragColor;
  
		in vec3 ourColor;
		in vec2 TexCoord;

		uniform sampler2D ourTexture;

		void main()
		{
			FragColor = texture(ourTexture, TexCoord);
		}
	` + "\x00"
)

var (
	screen = []float32{
		// xyz		color		texture coords
		-1, 1, 0, 1, 1, 1, 0, 0,
		-1, -1, 0, 1, 1, 1, 0, 1,
		1, -1, 0, 1, 1, 1, 1, 1,
		1, 1, 0, 1, 1, 1, 1, 0,
		-1, 1, 0, 1, 1, 1, 0, 0,
		1, -1, 0, 1, 1, 1, 1, 1,
	}
)

// Gfx is the "video card" api
type Gfx struct {
	// the video memory
	PixelMemory [Width * Height * 3]byte
	// the color palette (16 rgb values)
	Colors [16 * 3]uint8
	// The background color (used by some modes)
	BackgroundColor uint8

	Lock    sync.Mutex
	Window  *glfw.Window
	Program uint32
	Vao     uint32
}

// NewGfx lets you create a new Gfx video card
func NewGfx() *Gfx {
	gfx := &Gfx{
		PixelMemory: [320 * 200 * 3]byte{},
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

		Lock:    sync.Mutex{},
		Window:  initGlfw(),
		Program: initOpenGL(),
		Vao:     makeVao(),
	}

	runtime.LockOSThread()

	return gfx
}

// initGlfw initializes glfw and returns a Window to use.
func initGlfw() *glfw.Window {
	if err := glfw.Init(); err != nil {
		panic(err)
	}
	glfw.WindowHint(glfw.Resizable, glfw.True)
	glfw.WindowHint(glfw.ContextVersionMajor, 4)
	glfw.WindowHint(glfw.ContextVersionMinor, 1)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)

	window, err := glfw.CreateWindow(Width*scale, Height*scale, "Benji4000", nil, nil)
	if err != nil {
		panic(err)
	}
	window.MakeContextCurrent()

	return window
}

// initOpenGL initializes OpenGL and returns an intiialized program.
func initOpenGL() uint32 {
	if err := gl.Init(); err != nil {
		panic(err)
	}
	version := gl.GoStr(gl.GetString(gl.VERSION))
	log.Println("OpenGL version", version)

	vertexShader, err := compileShader(vertexShaderSource, gl.VERTEX_SHADER)
	if err != nil {
		panic(err)
	}

	fragmentShader, err := compileShader(fragmentShaderSource, gl.FRAGMENT_SHADER)
	if err != nil {
		panic(err)
	}

	prog := gl.CreateProgram()
	gl.AttachShader(prog, vertexShader)
	gl.AttachShader(prog, fragmentShader)
	gl.LinkProgram(prog)
	return prog
}

func compileShader(source string, shaderType uint32) (uint32, error) {
	shader := gl.CreateShader(shaderType)

	csources, free := gl.Strs(source)
	gl.ShaderSource(shader, 1, csources, nil)
	free()
	gl.CompileShader(shader)

	var status int32
	gl.GetShaderiv(shader, gl.COMPILE_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetShaderiv(shader, gl.INFO_LOG_LENGTH, &logLength)

		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetShaderInfoLog(shader, logLength, nil, gl.Str(log))

		return 0, fmt.Errorf("failed to compile %v: %v", source, log)
	}

	return shader, nil
}

// makeVao initializes and returns a vertex array from the points provided.
func makeVao() uint32 {
	var vbo uint32
	gl.GenBuffers(1, &vbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.BufferData(gl.ARRAY_BUFFER, 4*len(screen), gl.Ptr(screen), gl.STATIC_DRAW)

	var vao uint32
	gl.GenVertexArrays(1, &vao)
	gl.BindVertexArray(vao)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	var offset int

	// position attribute
	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 8*4, gl.PtrOffset(offset))
	gl.EnableVertexAttribArray(0)
	offset += 3 * 4

	// color attribute
	gl.VertexAttribPointer(1, 3, gl.FLOAT, false, 8*4, gl.PtrOffset(offset))
	gl.EnableVertexAttribArray(1)
	offset += 3 * 4

	// texture coord attribute
	gl.VertexAttribPointer(2, 2, gl.FLOAT, false, 8*4, gl.PtrOffset(offset))
	gl.EnableVertexAttribArray(2)
	offset += 2 * 4

	return vao
}

func (gfx *Gfx) MainLoop() {
	defer glfw.Terminate()

	// texture setup
	var texture uint32
	gl.GenTextures(1, &texture)
	gl.BindTexture(gl.TEXTURE_2D, texture)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.REPEAT)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.REPEAT)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
	// gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGB, Width, Height, 0, gl.RGB, gl.UNSIGNED_BYTE, gl.Ptr(&pixels[0]))
	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGB, Width, Height, 0, gl.RGB, gl.UNSIGNED_BYTE, nil)
	gl.GenerateMipmap(gl.TEXTURE_2D)

	// bind to shader
	gl.UseProgram(gfx.Program)
	gl.Uniform1i(gl.GetUniformLocation(gfx.Program, gl.Str("ourTexture\x00")), 0)

	var lastTime, delta float64
	var nbFrames int
	for !gfx.Window.ShouldClose() {
		currentTime := glfw.GetTime()
		delta = currentTime - lastTime
		nbFrames++
		if delta >= 1.0 { // If last cout was more than 1 sec ago
			gfx.Window.SetTitle(fmt.Sprintf("FPS: %.2f", float64(nbFrames)/delta))
			nbFrames = 0
			lastTime = currentTime
		}

		// blip the video ram to texture (random for now)
		// for i := 0; i < len(pixels); i++ {
		// 	pixels[i] = byte(rand.Intn(255))
		// }
		gfx.Lock.Lock()
		// need to do this so go.Ptr() works. This could be a bug in go: https://github.com/golang/go/issues/14210
		pixels := gfx.PixelMemory
		gl.TexSubImage2D(gl.TEXTURE_2D, 0, 0, 0, Width, Height, gl.RGB, gl.UNSIGNED_BYTE, gl.Ptr(&pixels[0]))
		gfx.Lock.Unlock()

		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
		gl.ActiveTexture(gl.TEXTURE0)
		gl.BindTexture(gl.TEXTURE_2D, texture)
		gl.UseProgram(gfx.Program)

		gl.BindVertexArray(gfx.Vao)
		gl.DrawArrays(gl.TRIANGLES, 0, int32(len(screen)/8))

		glfw.PollEvents()
		gfx.Window.SwapBuffers()
	}
}

/*
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
*/
