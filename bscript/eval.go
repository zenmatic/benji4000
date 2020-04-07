package bscript

import (
	"fmt"
	"math"

	"github.com/alecthomas/participle/lexer"
	"github.com/alecthomas/repr"
)

const STACK_LIMIT = 1000

type Evaluatable interface {
	Evaluate(ctx *Context) (interface{}, error)
}

type Function func(args ...interface{}) (interface{}, error)

type ContextLevel struct {
	Function string
	Vars     map[string]interface{}
}

// Context for evaluation.
type Context struct {
	Stack []ContextLevel
	// built-in functions.
	Builtins map[string]Function
	// The current function's name
	Function string
	// Vars defined during evaluation.
	Vars map[string]interface{}
	// Global variables
	Globals map[string]interface{}
	// top level constants
	Consts map[string]interface{}
	// defined functions
	Defs map[string]*Fun
}

func (v *Value) Evaluate(ctx *Context) (interface{}, error) {
	switch {
	case v.Number != nil:
		return *v.Number, nil
	case v.Map != nil:
		m := make(map[string]interface{})
		if v.Map.LeftNameValuePair != nil {
			value, err := v.Map.LeftNameValuePair.Value.Evaluate(ctx)
			if err != nil {
				return value, err
			}
			m[v.Map.LeftNameValuePair.Name] = value
			for i := 0; i < len(v.Map.RightNameValuePairs); i++ {
				value, err := v.Map.RightNameValuePairs[i].Value.Evaluate(ctx)
				if err != nil {
					return value, err
				}
				m[v.Map.RightNameValuePairs[i].Name] = value
			}
		}
		return m, nil
	case v.Array != nil:
		length := 0
		if v.Array.LeftValue != nil {
			length = 1 + len(v.Array.RightValues)
		}
		a := make([]interface{}, length)
		if v.Array.LeftValue != nil {
			value, err := v.Array.LeftValue.Evaluate(ctx)
			if err != nil {
				return value, err
			}
			a[0] = value
			for i := 0; i < len(v.Array.RightValues); i++ {
				value, err := v.Array.RightValues[i].Evaluate(ctx)
				if err != nil {
					return value, err
				}
				a[1+i] = value
			}
		}
		return &a, nil
	case v.ArrayElement != nil:
		currentValue, err := v.ArrayElement.Variable.Evaluate(ctx)
		if err != nil {
			return nil, err
		}
		for _, arrayIndex := range v.ArrayElement.Indexes {
			ivalue, err := arrayIndex.Index.Evaluate(ctx)
			if err != nil {
				return nil, err
			}

			a, ok := currentValue.(*[]interface{})
			if ok {
				// it's an array
				index := (int)(ivalue.(float64))
				if index < 0 || index >= len(*a) {
					return nil, lexer.Errorf(v.Pos, "Index out of bounds")
				}
				currentValue = (*a)[index]
			} else {
				// map?
				m, ok := currentValue.(map[string]interface{})
				if !ok {
					return nil, lexer.Errorf(v.Pos, "Array element should refer to array or map")
				}
				currentValue = m[ivalue.(string)]
			}
		}
		return currentValue, nil
	case v.String != nil:
		return *v.String, nil
	case v.Variable != nil:
		return v.Variable.Evaluate(ctx)
	case v.Subexpression != nil:
		return v.Subexpression.Evaluate(ctx)
	case v.Call != nil:
		return v.Call.Evaluate(ctx)
	}
	panic("unsupported value type" + repr.String(v))
}

func (v *Variable) Evaluate(ctx *Context) (interface{}, error) {
	value, ok := ctx.Vars[v.Variable]
	if !ok {
		value, ok = ctx.Consts[v.Variable]
		if !ok {
			value, ok = ctx.Globals[v.Variable]
			if !ok {
				return nil, lexer.Errorf(v.Pos, "unknown variable %q", v.Variable)
			}
		}
	}
	return value, nil
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
		if o.Operator == "+" {
			// special handling for string concat
			lhsStr, rhsStr, err := evaluateStrings(ctx, lhs, o.Term)
			if err != nil {
				return nil, lexer.Errorf(o.Pos, "invalid arguments for %s: %s", o.Operator, err)
			}
			return lhsStr + rhsStr, nil
		}
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

func (ctx *Context) debug(message string) {
	fmt.Println("====================================")
	fmt.Println(message)
	fmt.Printf("Stack=%v\n", ctx.Stack)
	fmt.Println("------------------------------------")
	fmt.Printf("Function=%s", ctx.Function)
	fmt.Printf("Globals=%v\n", ctx.Globals)
	fmt.Printf("Vars=%v\n", ctx.Vars)
	fmt.Println("====================================")
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
		function, ok := ctx.Builtins[c.Name]
		if !ok {
			return nil, lexer.Errorf(c.Pos, "unknown function %q", c.Name)
		}

		value, err := function(args...)
		if err != nil {
			fmt.Println(err)
			ctx.debug("Failed to call function")
			panic(lexer.Errorf(c.Pos, "call to %s() failed", c.Name))
		}
		return value, nil
	}

	if len(ctx.Stack) > STACK_LIMIT {
		panic("Stack limit exceeded")
	}

	// update globals
	for k := range ctx.Globals {
		ctx.Globals[k] = ctx.Vars[k]
	}
	// push a variable context
	newVars := map[string]interface{}{}
	for k, v := range ctx.Vars {
		newVars[k] = v
	}
	ctx.Stack = append(ctx.Stack, ContextLevel{
		Function: ctx.Function,
		Vars:     newVars,
	})
	// clear the variables
	ctx.Vars = map[string]interface{}{}
	// add the global vars
	for k, v := range ctx.Globals {
		ctx.Vars[k] = v
	}
	ctx.Function = c.Name

	// ctx.debug(fmt.Sprintf("BEFORE call to %s", c.Name))

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
		fmt.Println(err)
		ctx.debug("Failed to call function")
		panic(lexer.Errorf(c.Pos, "call to %s() failed", c.Name))
	}

	// update globals
	for k := range ctx.Globals {
		ctx.Globals[k] = ctx.Vars[k]
	}
	// pop  the stack and restore Vars
	ctx.Function = ctx.Stack[len(ctx.Stack)-1].Function
	ctx.Vars = map[string]interface{}{}
	for k, v := range ctx.Stack[len(ctx.Stack)-1].Vars {
		ctx.Vars[k] = v
	}
	// update vars to globals
	for k := range ctx.Globals {
		ctx.Vars[k] = ctx.Globals[k]
	}
	ctx.Stack = ctx.Stack[:len(ctx.Stack)-1]
	// ctx.debug(fmt.Sprintf("AFTER call to %s", c.Name))

	return value, err

}

func (cmd *Let) Evaluate(ctx *Context) (interface{}, error) {
	value, err := cmd.Value.Evaluate(ctx)

	if err != nil {
		return nil, err
	}
	if cmd.Variable != nil {
		ctx.Vars[*cmd.Variable] = value
	} else if cmd.ArrayElement != nil {
		currentValue, err := cmd.ArrayElement.Variable.Evaluate(ctx)
		if err != nil {
			return nil, err
		}

		for arrayIndexIndex, arrayIndex := range cmd.ArrayElement.Indexes {
			ivalue, err := arrayIndex.Index.Evaluate(ctx)
			if err != nil {
				return nil, err
			}

			lastElement := arrayIndexIndex == len(cmd.ArrayElement.Indexes)-1

			a, ok := currentValue.(*[]interface{})
			if ok {
				// it's an array
				index := (int)(ivalue.(float64))
				if lastElement {
					if index < 0 || index > len(*a) {
						return nil, lexer.Errorf(cmd.Pos, "Index out of bounds")
					}
					if index < len(*a) {
						(*a)[index] = value
					} else {
						*a = append(*a, value)
					}
				} else {
					if index < 0 || index >= len(*a) {
						return nil, lexer.Errorf(cmd.Pos, "Index out of bounds")
					}
					currentValue = (*a)[index]
				}
			} else {
				m, ok := currentValue.(map[string]interface{})
				if ok {
					key := ivalue.(string)
					if lastElement {
						m[key] = value
					} else {
						currentValue = m[key]
					}
				} else {
					return nil, lexer.Errorf(cmd.Pos, "Invalid array element: should be an array or a map")
				}
			}
		}
	} else {
		return nil, lexer.Errorf(cmd.Pos, "Let needs a variable or array element on the LHS.")
	}
	return nil, nil
}

func (cmd *Command) Evaluate(ctx *Context) (interface{}, error) {
	switch {
	case cmd.Remark != nil:

	case cmd.Let != nil:
		_, err := cmd.Let.Evaluate(ctx)
		if err != nil {
			return nil, err
		}
	case cmd.Del != nil:
		cmd := cmd.Del
		if cmd.ArrayElement != nil {
			currentValue, err := cmd.ArrayElement.Variable.Evaluate(ctx)
			if err != nil {
				return nil, err
			}

			for arrayIndexIndex, arrayIndex := range cmd.ArrayElement.Indexes {
				ivalue, err := arrayIndex.Index.Evaluate(ctx)
				if err != nil {
					return nil, err
				}

				lastElement := arrayIndexIndex == len(cmd.ArrayElement.Indexes)-1

				a, ok := currentValue.(*[]interface{})
				if ok {
					// it's an array
					index := (int)(ivalue.(float64))
					if index < 0 || index >= len(*a) {
						return nil, lexer.Errorf(cmd.Pos, "Index out of bounds")
					}
					if lastElement {
						*a = append((*a)[:index], (*a)[index+1:]...)
					} else {
						currentValue = (*a)[index]
					}
				} else {
					m, ok := currentValue.(map[string]interface{})
					if ok {
						key := ivalue.(string)
						if lastElement {
							delete(m, key)
						} else {
							currentValue = m[key]
						}
					} else {
						return nil, lexer.Errorf(cmd.Pos, "Invalid array element: should be an array or a map")
					}
				}
			}
		} else {
			// in the future, del can take other types (map, maybe struct, etc)
			return nil, lexer.Errorf(cmd.Pos, "can't delete this type of expression")
		}
	case cmd.Return != nil:
		cmd := cmd.Return
		value, err := cmd.Value.Evaluate(ctx)
		return value, err
	case cmd.Call != nil:
		_, err := cmd.Call.Evaluate(ctx)
		if err != nil {
			return nil, err
		}
	case cmd.If != nil:
		value, err := cmd.If.Evaluate(ctx)
		return value, err
	case cmd.While != nil:
		value, err := cmd.While.Evaluate(ctx)
		return value, err
	default:
		panic("unsupported command " + repr.String(cmd))
	}
	return nil, nil
}

func evalBlock(ctx *Context, commands []*Command) (interface{}, error) {
	for index := 0; index < len(commands); {
		cmd := commands[index]
		value, err := cmd.Evaluate(ctx)
		if err != nil {
			return nil, err
		}
		if value != nil {
			return value, nil
		}
		// ctx.debug("debug")
		index++
	}
	return nil, nil
}

func (whilecommand *While) Evaluate(ctx *Context) (interface{}, error) {
	for {
		value, err := whilecommand.Condition.Evaluate(ctx)
		if err != nil {
			return nil, err
		}

		if value != true {
			return nil, nil
		}

		value, err = evalBlock(ctx, whilecommand.Commands)
		if err != nil {
			return nil, err
		}
		if value != nil {
			return value, err
		}
	}
}

func (ifcommand *If) Evaluate(ctx *Context) (interface{}, error) {
	value, err := ifcommand.Condition.Evaluate(ctx)
	if err != nil {
		return nil, err
	}

	if value == true {
		return evalBlock(ctx, ifcommand.Commands)
	}
	return evalBlock(ctx, ifcommand.ElseCommands)
}

func (fun *Fun) Evaluate(ctx *Context) (interface{}, error) {
	// fmt.Println("Running function:", fun.Name)
	return evalBlock(ctx, fun.Commands)
}

func (program *Program) Evaluate() error {
	if len(program.TopLevel) == 0 {
		return nil
	}

	ctx := &Context{
		Stack:    []ContextLevel{},
		Consts:   map[string]interface{}{},
		Vars:     map[string]interface{}{},
		Globals:  map[string]interface{}{},
		Function: "",
		Builtins: Builtins(),
		Defs:     map[string]*Fun{},
	}

	// find main
	var main *Fun
	for i := 0; i < len(program.TopLevel); i++ {
		if program.TopLevel[i].Fun != nil {
			ctx.Defs[program.TopLevel[i].Fun.Name] = program.TopLevel[i].Fun
			if program.TopLevel[i].Fun.Name == "main" {
				main = program.TopLevel[i].Fun
			}
		} else if program.TopLevel[i].Const != nil {
			value, err := program.TopLevel[i].Const.Value.Evaluate(ctx)
			if err != nil {
				return err
			}
			ctx.Consts[program.TopLevel[i].Const.Name] = value
		} else if program.TopLevel[i].Let != nil {
			_, err := program.TopLevel[i].Let.Evaluate(ctx)
			if err != nil {
				return err
			}
			// now copy Vars to Globals
			for k, v := range ctx.Vars {
				ctx.Globals[k] = v
			}
		}
	}
	if main == nil {
		panic("No main function found.")
	}

	// fmt.Println("======================================")
	// repr.Println(ctx)
	// fmt.Println("======================================")

	ctx.Function = "main"
	_, err := main.Evaluate(ctx)
	if err != nil {
		ctx.debug("Main error")
		return err
	}
	// fmt.Println("Final program value: ", value)

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

func EvalString(value interface{}) string {
	avalue, ok := value.(*[]interface{})
	if ok {
		a := make([]string, len(*avalue))
		for idx, aa := range *avalue {
			a[idx] = EvalString(aa)
		}
		return fmt.Sprintf("%v", a)
	}
	return fmt.Sprintf("%v", value)
}

func evaluateStrings(ctx *Context, lhs interface{}, rhsExpr Evaluatable) (string, string, error) {
	rhs, err := rhsExpr.Evaluate(ctx)
	if err != nil {
		return "", "", err
	}
	return EvalString(lhs), EvalString(rhs), nil
}
