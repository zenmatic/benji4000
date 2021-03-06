package bscript

import (
	"fmt"
	"math"
	"os"
	"strings"

	"github.com/alecthomas/participle/lexer"
	"github.com/alecthomas/repr"
	"github.com/uzudil/benji4000/gfx"
)

var ANON_COUNT uint32

const STACK_LIMIT = 1000

type Evaluatable interface {
	Evaluate(ctx *Context) (interface{}, error)
}

type Builtin func(ctx *Context, args ...interface{}) (interface{}, error)

func (c *Closure) String() string {
	return fmt.Sprintf("%s(%s)", c.Function, strings.Join(c.Params, ","))
}

type Closure struct {
	// function defition: params
	Params []string
	// function defition: commands
	Commands []*Command
	// The current function's name
	Function string
	// variables
	Vars map[string]interface{}
	// function definitions
	Defs map[string]*Closure
	// the parent closure
	Parent *Closure
	// the graphics
	Video *gfx.Gfx
}

type Runtime struct {
	Pos      lexer.Position
	Function string
	Vars     map[string]interface{}
}

// Context for evaluation.
type Context struct {
	// built-in functions.
	Builtins map[string]Builtin
	// top level constants
	Consts map[string]interface{}
	// the global closure
	Closure *Closure
	// the runtime stack
	RuntimeStack []Runtime
	// the current position
	Pos lexer.Position
	// the program
	Program *Program
	// the video card
	Video *gfx.Gfx
}

func (v *Value) Evaluate(ctx *Context) (interface{}, error) {
	switch {
	case v.Number != nil:
		if v.Number.Sign != nil && *(v.Number.Sign) == "-" {
			return -v.Number.Number, nil
		}
		return v.Number.Number, nil
	case v.Boolean != nil:
		return *v.Boolean == "true", nil
	case v.Null != nil:
		return nil, nil
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
	case v.AnonFun != nil:
		return v.AnonFun.Evaluate(ctx)
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
	value, ok := ctx.Consts[v.Variable]
	if ok {
		return value, nil
	}
	for closure := ctx.Closure; closure != nil; closure = closure.Parent {
		value, ok = closure.Vars[v.Variable]
		if ok {
			return value, nil
		}
		def, ok := closure.Defs[v.Variable]
		if ok {
			return def, nil
		}
	}
	return nil, lexer.Errorf(v.Pos, "unknown variable %q", v.Variable)
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
	case "%":
		return float64(int(lhsNumber) % int(rhsNumber)), nil
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
	case bool:
		rhs, ok := rhs.(bool)
		if !ok {
			return nil, lexer.Errorf(o.Pos, "rhs of %s must be a string", o.Operator)
		}
		switch o.Operator {
		case "=":
			return lhs == rhs, nil
		case "!=":
			return lhs != rhs, nil
		}
	default:
		return nil, lexer.Errorf(o.Pos, "lhs of %s must be a number or string", o.Operator)
	}
	panic("unreachable")
}

func (b *OpBoolTerm) Evaluate(ctx *Context, lhs interface{}) (interface{}, error) {
	rhs, err := b.Right.Evaluate(ctx)
	if err != nil {
		return nil, err
	}
	switch lhs := lhs.(type) {
	case bool:
		rhs, ok := rhs.(bool)
		if !ok {
			return nil, lexer.Errorf(b.Pos, "rhs of %s must be a boolean", b.Operator)
		}
		switch b.Operator {
		case "&&":
			return rhs && lhs, nil
		case "||":
			return rhs || lhs, nil
		}
	default:
		return nil, lexer.Errorf(b.Pos, "lhs of %s must be a boolean", b.Operator)
	}
	panic("unreachable")
}

func (b *BoolTerm) Evaluate(ctx *Context) (interface{}, error) {
	lhs, err := b.Left.Evaluate(ctx)
	if err != nil {
		return nil, err
	}
	for _, right := range b.Right {
		rhs, err := right.Evaluate(ctx, lhs)
		if err != nil {
			return nil, err
		}
		lhs = rhs
	}
	return lhs, nil
}

func (e *Expression) Evaluate(ctx *Context) (interface{}, error) {
	lhs, err := e.BoolTerm.Evaluate(ctx)
	if err != nil {
		return nil, err
	}
	for _, right := range e.OpBoolTerm {
		rhs, err := right.Evaluate(ctx, lhs)
		if err != nil {
			return nil, err
		}
		lhs = rhs
	}
	return lhs, nil
}

func (ctx *Context) debug(message string) {
	ctx.Builtins["print"](ctx, message)
	indent := "  "
	ctx.Builtins["print"](ctx, "Constants:")
	for k, v := range ctx.Consts {
		ctx.Builtins["print"](ctx, fmt.Sprintf("  %s=%v", k, v))
	}
	ctx.Builtins["print"](ctx, "Closures:")
	for closure := ctx.Closure; closure != nil; closure = closure.Parent {
		ctx.Builtins["print"](ctx, "-----------------")
		ctx.Builtins["print"](ctx, fmt.Sprintf("%sFunction: %s\n", indent, closure.Function))
		ctx.Builtins["print"](ctx, fmt.Sprintf("%sVars: %v\n", indent, closure.Vars))
		ctx.Builtins["print"](ctx, fmt.Sprintf("%sDefs: %v\n", indent, closure.Defs))
		indent = indent + "  "
	}
	ctx.Builtins["print"](ctx, "------------------------------------")
	ctx.Builtins["print"](ctx, "Runtime Call Stack:")
	indent = "  "
	for _, runtime := range ctx.RuntimeStack {
		ctx.Builtins["print"](ctx, fmt.Sprintf("%s%s at %s Vars=%s\n", indent, runtime.Function, runtime.Pos, runtime.Vars))
		indent = indent + "  "
	}
	ctx.Builtins["print"](ctx, "------------------------------------")
	ctx.Builtins["print"](ctx, fmt.Sprintf("Currently: %s\n", ctx.Pos))
}

func evalBuiltinCall(ctx *Context, c *Call, builtin Builtin, args []interface{}) (interface{}, error) {
	// fmt.Println("Calling builtin", c.Name)
	// a built-in function

	value, err := builtin(ctx, args...)
	if err != nil {
		return nil, err
	}
	return value, nil
}

func evalFunctionCall(ctx *Context, c *Call, closure *Closure, args []interface{}) (interface{}, error) {
	if len(ctx.RuntimeStack) > STACK_LIMIT {
		panic("Stack limit exceeded")
	}

	// save local variables (needed when a recursive call modifies the closure's variables)
	saved := make(map[string]interface{}, len(ctx.Closure.Vars))
	for k, v := range ctx.Closure.Vars {
		saved[k] = v
	}
	ctx.RuntimeStack = append(ctx.RuntimeStack, Runtime{
		Pos:      c.Pos,
		Function: c.Name,
		Vars:     saved,
	})
	savedClosure := ctx.Closure
	ctx.Closure = closure

	// create function call param variables
	if len(closure.Params) != len(args) {
		return nil, lexer.Errorf(c.Pos, "Not all function params given in call to %s", c.Name)
	}
	for index := 0; index < len(closure.Params); index++ {
		closure.Vars[closure.Params[index]] = args[index]
	}

	// make the call (evaluate the function's code)
	value, err := evalBlock(ctx, closure.Commands)
	if err != nil {
		return nil, err
	}

	// restore local vars and environment
	ctx.Closure = savedClosure
	for k, v := range saved {
		ctx.Closure.Vars[k] = v
	}
	// drop the last frame of the stack
	ctx.RuntimeStack = ctx.RuntimeStack[:len(ctx.RuntimeStack)-1]

	return value, err
}

func (closure *Closure) findClosure(name string) (*Closure, bool) {
	// a defined function
	fx, ok := closure.Defs[name]
	if ok {
		return fx, ok
	}

	// variable pointing to a function
	variable, ok := closure.Vars[name]
	if ok {
		fx, ok = variable.(*Closure)
		if ok {
			return fx, true
		}
	}

	return nil, false
}

func (callParams *CallParams) evalParams(ctx *Context) ([]interface{}, error) {
	args := []interface{}{}
	for _, arg := range callParams.Args {
		value, err := arg.Evaluate(ctx)
		if err != nil {
			return nil, err
		}
		args = append(args, value)
	}
	return args, nil
}

func (c *Call) Evaluate(ctx *Context) (interface{}, error) {
	args, err := c.CallParams[0].evalParams(ctx)
	if err != nil {
		return nil, err
	}

	// call builtin function
	builtin, ok := ctx.Builtins[c.Name]
	if ok {
		return evalBuiltinCall(ctx, c, builtin, args)
	}

	// call function the first time
	var result interface{}
	var fx *Closure
	for closure := ctx.Closure; closure != nil; closure = closure.Parent {
		// a defined function
		fx, ok = closure.findClosure(c.Name)
		if ok {
			result, err = evalFunctionCall(ctx, c, fx, args)
			if err != nil {
				return nil, err
			}
			break
		}
	}
	if fx == nil {
		return nil, lexer.Errorf(c.Pos, "Unknown function %s()", c.Name)
	}

	// subsequent function calls
	for i := 1; i < len(c.CallParams); i++ {
		args, err := c.CallParams[i].evalParams(ctx)
		if err != nil {
			return nil, err
		}
		fx, ok := result.(*Closure)
		if !ok {
			return nil, lexer.Errorf(c.Pos, "Function call references non-function variable: %v", result)
		}
		result, err = evalFunctionCall(ctx, c, fx, args)
		if err != nil {
			return nil, err
		}
	}
	return result, nil
}

func (cmd *Let) Evaluate(ctx *Context) (interface{}, error) {
	value, err := cmd.Value.Evaluate(ctx)

	if err != nil {
		return nil, err
	}
	if cmd.Variable != nil {
		for c := ctx.Closure; c != nil; c = c.Parent {
			_, ok := c.Vars[*cmd.Variable]
			if ok {
				// existing var
				c.Vars[*cmd.Variable] = value
				return nil, nil
			}
		}
		// new var
		ctx.Closure.Vars[*cmd.Variable] = value
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

// Evaluate a Command.
// some commands return a value which causes the exection of a block to stop (eg. return, while, if)
func (cmd *Command) Evaluate(ctx *Context) (interface{}, error) {
	ctx.Pos = cmd.Pos

	switch {
	case cmd.Remark != nil:
		return nil, nil
	case cmd.Let != nil:
		_, err := cmd.Let.Evaluate(ctx)
		return nil, err
	case cmd.Fun != nil:
		_, err := cmd.Fun.Evaluate(ctx)
		return nil, err
	case cmd.Del != nil:
		_, err := cmd.Del.Evaluate(ctx)
		return nil, err
	case cmd.Return != nil:
		return cmd.Return.Value.Evaluate(ctx)
	case cmd.Call != nil:
		_, err := cmd.Call.Evaluate(ctx)
		return nil, err
	case cmd.If != nil:
		return cmd.If.Evaluate(ctx)
	case cmd.While != nil:
		return cmd.While.Evaluate(ctx)
	default:
		panic("unsupported command " + repr.String(cmd))
	}
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

func (cmd *Del) Evaluate(ctx *Context) (interface{}, error) {
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

func makeClosure(ctx *Context, name string, params []string, commands []*Command) *Closure {
	return &Closure{
		Params:   params,
		Commands: commands,
		Vars:     map[string]interface{}{},
		Defs:     map[string]*Closure{},
		Function: name,
		Parent:   ctx.Closure,
	}
}

func (fun *Fun) Evaluate(ctx *Context) (interface{}, error) {
	ctx.Closure.Defs[fun.Name] = makeClosure(ctx, fun.Name, fun.Params, fun.Commands)
	return nil, nil
}

func (anonFun *AnonFun) Evaluate(ctx *Context) (interface{}, error) {
	name := fmt.Sprintf("_anon_%d", ANON_COUNT)
	ANON_COUNT++
	var params []string
	if anonFun.SingleParam != nil {
		params = []string{*anonFun.SingleParam}
	} else {
		params = anonFun.Params
	}
	var commands []*Command
	if anonFun.SingleCommand != nil {
		commands = []*Command{
			&Command{
				Pos: anonFun.Pos,
				Return: &Return{
					Pos:   anonFun.Pos,
					Value: anonFun.SingleCommand,
				},
			},
		}
	} else {
		commands = anonFun.Commands
	}
	return makeClosure(ctx, name, params, commands), nil
}

func CreateContext(program *Program) *Context {
	global := &Closure{
		Function: "global",
		Params:   []string{},
		Commands: []*Command{},
		Vars:     map[string]interface{}{},
		Defs:     map[string]*Closure{},
		Parent:   nil,
	}
	return &Context{
		Consts:       Constants(),
		Builtins:     Builtins(),
		Closure:      global,
		RuntimeStack: []Runtime{},
		Pos:          lexer.Position{},
		Program:      program,
		Video:        nil,
	}
}

func load(source string, showAst *bool) (*Program, error) {
	r, err := os.Open(source)
	if err != nil {
		return nil, err
	}
	defer r.Close()

	ast := &Program{}
	Parser.Parse(r, ast)
	if showAst != nil && *showAst {
		// print the ast
		repr.Println(ast)
		os.Exit(0)
	}
	return ast, nil
}

func Load(source string, showAst *bool, ctx *Context) (interface{}, error) {
	ast, err := load(source, showAst)
	if err != nil {
		return nil, err
	}
	return ast.init(ctx)
}

func Run(source string, showAst *bool, ctx *Context, video *gfx.Gfx) (interface{}, error) {
	// run it
	ast, err := load(source, showAst)
	if err != nil {
		return nil, err
	}

	ctx, err = ast.init(ctx)
	if err != nil {
		return nil, err
	}
	ctx.Video = video

	return ast.Evaluate(ctx)
}

func (program *Program) init(ctx *Context) (*Context, error) {
	if ctx == nil {
		ctx = CreateContext(program)
	}

	ctx.Program = program

	if len(program.TopLevel) == 0 {
		return ctx, nil
	}

	// define constants and globals
	for i := 0; i < len(program.TopLevel); i++ {
		if program.TopLevel[i].Const != nil {
			value, err := program.TopLevel[i].Const.Value.Evaluate(ctx)
			if err != nil {
				return ctx, err
			}
			ctx.Consts[program.TopLevel[i].Const.Name] = value
		} else if program.TopLevel[i].Let != nil {
			_, err := program.TopLevel[i].Let.Evaluate(ctx)
			if err != nil {
				return ctx, err
			}
		}
	}

	// define functions
	for i := 0; i < len(program.TopLevel); i++ {
		if program.TopLevel[i].Fun != nil {
			_, err := program.TopLevel[i].Fun.Evaluate(ctx)
			if err != nil {
				return ctx, err
			}
		}
	}

	_, ok := ctx.Closure.Defs["main"]
	if !ok {
		return ctx, fmt.Errorf("no main function found")
	}

	return ctx, nil
}

func (program *Program) Evaluate(ctx *Context) (interface{}, error) {

	// Call main()
	call := &Call{
		Name: "main",
		CallParams: []*CallParams{
			&CallParams{
				Args: []*Expression{},
			},
		},
	}
	return call.Evaluate(ctx)
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
