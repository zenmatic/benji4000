package bscript

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"strings"
)

func print(arg ...interface{}) (interface{}, error) {
	fmt.Println(EvalString(arg[0]))
	return nil, nil
}

func input(arg ...interface{}) (interface{}, error) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print(arg[0])
	text, err := reader.ReadString('\n')
	return strings.TrimSpace(text), err
}

func length(arg ...interface{}) (interface{}, error) {
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

func substr(arg ...interface{}) (interface{}, error) {
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

func keys(arg ...interface{}) (interface{}, error) {
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

func Builtins() map[string]Function {
	return map[string]Function{
		"print":  print,
		"input":  input,
		"len":    length,
		"keys":   keys,
		"substr": substr,
	}
}
