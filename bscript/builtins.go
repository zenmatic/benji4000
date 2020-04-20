package bscript

import (
	"bufio"
	"fmt"
	"math"
	"math/rand"
	"os"
	"strings"
)

func print(ctx *Context, arg ...interface{}) (interface{}, error) {
	fmt.Println(EvalString(arg[0]))
	return nil, nil
}

func input(ctx *Context, arg ...interface{}) (interface{}, error) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print(arg[0])
	text, err := reader.ReadString('\n')
	return strings.TrimSpace(text), err
}

func length(ctx *Context, arg ...interface{}) (interface{}, error) {
	a, ok := arg[0].(*[]interface{})
	if !ok {
		s, ok := arg[0].(string)
		if !ok {
			return nil, fmt.Errorf("argument to len() should be an array or a string")
		}
		return float64(len(s)), nil
	}
	return float64(len(*a)), nil
}

func substr(ctx *Context, arg ...interface{}) (interface{}, error) {
	s, ok := arg[0].(string)
	if !ok {
		return nil, fmt.Errorf("argument 1 to substr() should be a string")
	}
	index, ok := arg[1].(float64)
	if !ok {
		return nil, fmt.Errorf("argument 2 to substr() should be a number")
	}
	length := len(s)
	if len(arg) > 2 {
		f, ok := arg[2].(float64)
		if !ok {
			return nil, fmt.Errorf("argument 3 to substr() should be a number")
		}
		length = int(f)
	}
	start := int(math.Min(math.Max(index, 0), float64(len(s))))
	end := int(math.Min(math.Max(float64(start+length), 0), float64(len(s))))
	return string(s[start:end]), nil
}

func replace(ctx *Context, arg ...interface{}) (interface{}, error) {
	s, ok := arg[0].(string)
	if !ok {
		return nil, fmt.Errorf("argument 1 to replace() should be a string")
	}
	oldstring, ok := arg[1].(string)
	if !ok {
		return nil, fmt.Errorf("argument 2 to replace() should be a string")
	}
	newstring, ok := arg[2].(string)
	if !ok {
		return nil, fmt.Errorf("argument 3 to replace() should be a string")
	}

	return strings.ReplaceAll(s, oldstring, newstring), nil
}

func keys(ctx *Context, arg ...interface{}) (interface{}, error) {
	m, ok := arg[0].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("argument to key() should be a map")
	}
	keys := make([]interface{}, len(m))
	index := 0
	for k := range m {
		keys[index] = k
		index++
	}
	return &keys, nil
}

func setVideoMode(ctx *Context, arg ...interface{}) (interface{}, error) {
	mode, ok := arg[0].(float64)
	if !ok {
		return nil, fmt.Errorf("First parameter should be the number of the video mode")
	}
	ctx.VideoMode = int(mode)
	return nil, nil
}

func setPixel(ctx *Context, arg ...interface{}) (interface{}, error) {
	x, ok := arg[0].(float64)
	if !ok {
		return nil, fmt.Errorf("First parameter should be a number")
	}
	y, ok := arg[1].(float64)
	if !ok {
		return nil, fmt.Errorf("Second parameter should be a number")
	}
	color, ok := arg[2].(float64)
	if !ok {
		return nil, fmt.Errorf("Third parameter should be a number")
	}
	return nil, ctx.SetPixel(int(x), int(y), uint8(0), uint8(color), uint8(0))
}

func drawLine(ctx *Context, arg ...interface{}) (interface{}, error) {
	x, ok := arg[0].(float64)
	if !ok {
		return nil, fmt.Errorf("First parameter should be a number")
	}
	y, ok := arg[1].(float64)
	if !ok {
		return nil, fmt.Errorf("Second parameter should be a number")
	}
	x2, ok := arg[2].(float64)
	if !ok {
		return nil, fmt.Errorf("Third parameter should be a number")
	}
	y2, ok := arg[3].(float64)
	if !ok {
		return nil, fmt.Errorf("Fourth parameter should be a number")
	}
	color, ok := arg[4].(float64)
	if !ok {
		return nil, fmt.Errorf("Fifth parameter should be a number")
	}
	return nil, ctx.DrawLine(int(x), int(y), int(x2), int(y2), uint8(color))
}

func drawCircle(ctx *Context, arg ...interface{}) (interface{}, error) {
	x, ok := arg[0].(float64)
	if !ok {
		return nil, fmt.Errorf("First parameter should be a number")
	}
	y, ok := arg[1].(float64)
	if !ok {
		return nil, fmt.Errorf("Second parameter should be a number")
	}
	r, ok := arg[2].(float64)
	if !ok {
		return nil, fmt.Errorf("Third parameter should be a number")
	}
	color, ok := arg[3].(float64)
	if !ok {
		return nil, fmt.Errorf("Fourth parameter should be a number")
	}
	return nil, ctx.DrawCircle(int(x), int(y), int(r), uint8(color))
}

func fillCircle(ctx *Context, arg ...interface{}) (interface{}, error) {
	x, ok := arg[0].(float64)
	if !ok {
		return nil, fmt.Errorf("First parameter should be a number")
	}
	y, ok := arg[1].(float64)
	if !ok {
		return nil, fmt.Errorf("Second parameter should be a number")
	}
	r, ok := arg[2].(float64)
	if !ok {
		return nil, fmt.Errorf("Third parameter should be a number")
	}
	color, ok := arg[3].(float64)
	if !ok {
		return nil, fmt.Errorf("Fourth parameter should be a number")
	}
	return nil, ctx.FillCircle(int(x), int(y), int(r), uint8(color))
}

func clearVideo(ctx *Context, arg ...interface{}) (interface{}, error) {
	return nil, ctx.ClearVideo()
}

func updateVideo(ctx *Context, arg ...interface{}) (interface{}, error) {
	if ctx.Video == nil {
		panic("Video card not initialized")
	}
	return nil, ctx.UpdateVideo()
}

func random(ctx *Context, arg ...interface{}) (interface{}, error) {
	return rand.Float64(), nil
}

func debug(ctx *Context, arg ...interface{}) (interface{}, error) {
	message, ok := arg[0].(string)
	if !ok {
		return nil, fmt.Errorf("argument to debug() should be a string")
	}
	ctx.debug(message)
	return nil, nil
}

func assert(ctx *Context, arg ...interface{}) (interface{}, error) {
	a := arg[0]
	b := arg[1]
	msg := "Incorrect value"
	if len(arg) > 2 {
		msg = arg[2].(string)
	}

	var res bool

	// for arrays, compare the values
	arr, ok := a.(*[]interface{})
	if ok {
		// array
		brr, ok := b.(*[]interface{})
		if !ok {
			res = true
		} else {
			if len(*arr) == len(*brr) {
				res = false
				for i := range *arr {
					if (*arr)[i] != (*brr)[i] {
						res = true
						break
					}
				}
			} else {
				res = true
			}
		}
	} else {
		// map
		amap, ok := a.(map[string]interface{})
		if ok {
			bmap, ok := b.(map[string]interface{})
			if !ok {
				res = true
			} else {
				if len(amap) == len(bmap) {
					res = false
					for k := range amap {
						if amap[k] != bmap[k] {
							res = true
							break
						}
					}
				} else {
					res = true
				}
			}
		} else {
			// default is to compare equality
			res = a != b
		}
	}

	if res {
		debug(ctx, fmt.Sprintf("Assertion failure: %s: %v != %v", msg, a, b))
		return nil, fmt.Errorf("%s Assertion failure: %s: %v != %v", ctx.Pos, msg, a, b)
	}
	return nil, nil
}

func Builtins() map[string]Builtin {
	return map[string]Builtin{
		"print":        print,
		"input":        input,
		"len":          length,
		"keys":         keys,
		"substr":       substr,
		"replace":      replace,
		"debug":        debug,
		"assert":       assert,
		"setVideoMode": setVideoMode,
		"setPixel":     setPixel,
		"random":       random,
		"updateVideo":  updateVideo,
		"clearVideo":   clearVideo,
		"drawLine":     drawLine,
		"drawCircle":   drawCircle,
		"fillCircle":   fillCircle,
	}
}
