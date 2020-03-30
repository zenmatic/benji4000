package bscript

import (
	"math"

	"github.com/alecthomas/participle/lexer"
	"github.com/alecthomas/repr"

	"fmt"
)

type Evaluatable interface {
	Evaluate(ctx *Context) (interface{}, error)
}

type Function func(args ...interface{}) (interface{}, error)

// Context for evaluation.
type Context struct {
	// User-provided functions.
	Functions map[string]Function
	// Vars defined during evaluation.
	Vars map[string]interface{}
	// defined functions
	Defs map[string]*Fun
}

func (v *Value) Evaluate(ctx *Context) (interface{}, error) {
	switch {
	case v.Number != nil:
		return *v.Number, nil
	case v.String != nil:
		return *v.String, nil
	case v.Variable != nil:
		value, ok := ctx.Vars[*v.Variable]
		if !ok {
			return nil, fmt.Errorf("unknown variable %q", *v.Variable)
		}
		return value, nil
	case v.Subexpression != nil:
		return v.Subexpression.Evaluate(ctx)
	case v.Call != nil:
		return v.Call.Evaluate(ctx)
	}
	panic("unsupported value type" + repr.String(v))
}

func (f *Factor) Evaluate(ctx *Context) (interface{}, error) {
	base, err := f.Base.Evaluate(ctx)
	if err != nil {
		return nil, err
	}
	if f.Exponent == nil {
		return base, nil
	}
	baseNum, exponentNum, err := evaluateFloats(ctx, base, f.Exponent)
	if err != nil {
		return nil, lexer.Errorf(f.Pos, "invalid factor: %s", err)
	}
	return math.Pow(baseNum, exponentNum), nil
}

func (o *OpFactor) Evaluate(ctx *Context, lhs interface{}) (interface{}, error) {
	lhsNumber, rhsNumber, err := evaluateFloats(ctx, lhs, o.Factor)
	if err != nil {
		return nil, lexer.Errorf(o.Pos, "invalid arguments for %s: %s", o.Operator, err)
	}
	switch o.Operator {
	case "*":
		return lhsNumber * rhsNumber, nil
	case "/":
		return lhsNumber / rhsNumber, nil
	}
	panic("unreachable")
}

func (t *Term) Evaluate(ctx *Context) (interface{}, error) {
	lhs, err := t.Left.Evaluate(ctx)
	if err != nil {
		return nil, err
	}
	for _, right := range t.Right {
		rhs, err := right.Evaluate(ctx, lhs)
		if err != nil {
			return nil, err
		}
		lhs = rhs
	}
	return lhs, nil
}

func (o *OpTerm) Evaluate(ctx *Context, lhs interface{}) (interface{}, error) {
	lhsNumber, rhsNumber, err := evaluateFloats(ctx, lhs, o.Term)
	if err != nil {
		return nil, lexer.Errorf(o.Pos, "invalid arguments for %s: %s", o.Operator, err)
	}
	switch o.Operator {
	case "+":
		return lhsNumber + rhsNumber, nil
	case "-":
		return lhsNumber - rhsNumber, nil
	}
	panic("unreachable")
}

func (c *Cmp) Evaluate(ctx *Context) (interface{}, error) {
	lhs, err := c.Left.Evaluate(ctx)
	if err != nil {
		return nil, err
	}
	for _, right := range c.Right {
		rhs, err := right.Evaluate(ctx, lhs)
		if err != nil {
			return nil, err
		}
		lhs = rhs
	}
	return lhs, nil
}

func (o *OpCmp) Evaluate(ctx *Context, lhs interface{}) (interface{}, error) {
	rhs, err := o.Cmp.Evaluate(ctx)
	if err != nil {
		return nil, err
	}
	switch lhs := lhs.(type) {
	case float64:
		rhs, ok := rhs.(float64)
		if !ok {
			return nil, lexer.Errorf(o.Pos, "rhs of %s must be a number", o.Operator)
		}
		switch o.Operator {
		case "=":
			return lhs == rhs, nil
		case "!=":
			return lhs != rhs, nil
		case "<":
			return lhs < rhs, nil
		case ">":
			return lhs > rhs, nil
		case "<=":
			return lhs <= rhs, nil
		case ">=":
			return lhs >= rhs, nil
		}
	case string:
		rhs, ok := rhs.(string)
		if !ok {
			return nil, lexer.Errorf(o.Pos, "rhs of %s must be a string", o.Operator)
		}
		switch o.Operator {
		case "=":
			return lhs == rhs, nil
		case "!=":
			return lhs != rhs, nil
		case "<":
			return lhs < rhs, nil
		case ">":
			return lhs > rhs, nil
		case "<=":
			return lhs <= rhs, nil
		case ">=":
			return lhs >= rhs, nil
		}
	default:
		return nil, lexer.Errorf(o.Pos, "lhs of %s must be a number or string", o.Operator)
	}
	panic("unreachable")
}

func (e *Expression) Evaluate(ctx *Context) (interface{}, error) {
	lhs, err := e.Left.Evaluate(ctx)
	if err != nil {
		return nil, err
	}
	for _, right := range e.Right {
		rhs, err := right.Evaluate(ctx, lhs)
		if err != nil {
			return nil, err
		}
		lhs = rhs
	}
	return lhs, nil
}

func (c *Call) Evaluate(ctx *Context) (interface{}, error) {

	args := []interface{}{}
	for _, arg := range c.Args {
		value, err := arg.Evaluate(ctx)
		if err != nil {
			return nil, err
		}
		args = append(args, value)
	}

	// fmt.Println("Looking for function", c.Name)
	fun, ok := ctx.Defs[c.Name]
	if !ok {
		// fmt.Println("Calling builtin", c.Name)
		// a built-in function
		function, ok := ctx.Functions[c.Name]
		if !ok {
			return nil, lexer.Errorf(c.Pos, "unknown function %q", c.Name)
		}

		value, err := function(args...)
		if err != nil {
			return nil, lexer.Errorf(c.Pos, "call to %s() failed", c.Name)
		}
		return value, nil
	}

	// push the function params
	if len(fun.Params) != len(args) {
		return nil, lexer.Errorf(c.Pos, "Not all function params given in call to %s", c.Name)
	}
	for index := 0; index < len(fun.Params); index++ {
		ctx.Vars[fun.Params[index]] = args[index]
	}

	// a defined function
	// fmt.Println("Calling def", c.Name)
	value, err := fun.Evaluate(ctx)
	if err != nil {
		return nil, lexer.Errorf(c.Pos, "call to %s() failed", c.Name)
	}
	return value, err

}

func (fun *Fun) Evaluate(ctx *Context) (interface{}, error) {
	// fmt.Println("Running function:", fun.Name)
	for index := 0; index < len(fun.Commands); {
		cmd := fun.Commands[index]
		switch {
		case cmd.Remark != nil:

		case cmd.Let != nil:
			cmd := cmd.Let
			value, err := cmd.Value.Evaluate(ctx)
			if err != nil {
				return nil, err
			}
			ctx.Vars[cmd.Variable] = value
		case cmd.Return != nil:
			cmd := cmd.Return
			value, err := cmd.Value.Evaluate(ctx)
			return value, err
		case cmd.Call != nil:
			_, err := cmd.Call.Evaluate(ctx)
			if err != nil {
				return nil, err
			}

		default:
			panic("unsupported command " + repr.String(cmd))
		}

		index++
	}
	return nil, nil
}

func (program *Program) Evaluate() error {
	if len(program.Funs) == 0 {
		return nil
	}

	ctx := &Context{
		Vars: map[string]interface{}{},
		Functions: map[string]Function{
			"print": func(arg ...interface{}) (interface{}, error) {
				fmt.Println(arg[0])
				return nil, nil
			},
		},
		Defs: map[string]*Fun{},
	}

	// find main
	var main *Fun
	for i := 0; i < len(program.Funs); i++ {
		ctx.Defs[program.Funs[i].Name] = program.Funs[i]
		if program.Funs[i].Name == "main" {
			main = program.Funs[i]
		}
	}
	if main == nil {
		panic("No main function found.")
	}

	// fmt.Println("======================================")
	// repr.Println(ctx)
	// fmt.Println("======================================")

	value, err := main.Evaluate(ctx)
	if err != nil {
		return err
	}
	fmt.Println("Final program value: ", value)

	// fmt.Println("======================================")
	// repr.Println(ctx.Vars)
	// fmt.Println("======================================")

	return nil
}

func evaluateFloats(ctx *Context, lhs interface{}, rhsExpr Evaluatable) (float64, float64, error) {
	rhs, err := rhsExpr.Evaluate(ctx)
	if err != nil {
		return 0, 0, err
	}
	lhsNumber, ok := lhs.(float64)
	if !ok {
		return 0, 0, fmt.Errorf("lhs must be a number")
	}
	rhsNumber, ok := rhs.(float64)
	if !ok {
		return 0, 0, fmt.Errorf("rhs must be a number")
	}
	return lhsNumber, rhsNumber, nil
}
