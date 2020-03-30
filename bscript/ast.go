package bscript

import (
	"github.com/alecthomas/participle"
	"github.com/alecthomas/participle/lexer"
	"github.com/alecthomas/participle/lexer/ebnf"

	"strings"
)

type Program struct {
	Pos lexer.Position

	Funs []*Fun `( @@ )*`
}

type Fun struct {
	Pos lexer.Position

	Name     string     `"def" @Ident "("`
	Params   []string   `( @Ident ( "," @Ident )* )*`
	Commands []*Command `")" ( @@ )* "end"`
}

type Command struct {
	Pos lexer.Position

	Remark *Remark `(   @@ `
	Let    *Let    `  | @@ ";" `
	Return *Return `  | @@ ";" `
	If     *If     `  | @@ `
	While  *While  `  | @@ `
	Call   *Call   `  | @@ ";" )`
}

type While struct {
	Pos lexer.Position

	Condition *Expression `"while" "(" @@ ")"`
	Commands  []*Command  `( @@ )* "end"`
}

type If struct {
	Pos lexer.Position

	Condition *Expression `"if" "(" @@ ")"`
	Commands  []*Command  `( @@ )* "end"`
}

type Remark struct {
	Pos lexer.Position

	Comment string `@Comment`
}

type Call struct {
	Pos lexer.Position

	Name string        `@Ident`
	Args []*Expression `"(" [ @@ { "," @@ } ] ")"`
}

type Let struct {
	Pos lexer.Position

	Variable string      `"let" @Ident`
	Value    *Expression `"=" @@`
}

type Return struct {
	Pos lexer.Position

	Value *Expression `"return" @@`
}

type Operator string

func (o *Operator) Capture(s []string) error {
	*o = Operator(strings.Join(s, ""))
	return nil
}

type Value struct {
	Pos lexer.Position

	Number        *float64    `  @Number`
	Call          *Call       `| @@`
	Variable      *string     `| @Ident`
	String        *string     `| @String`
	Subexpression *Expression `| "(" @@ ")"`
}

type Factor struct {
	Pos lexer.Position

	Base     *Value `@@`
	Exponent *Value `[ "^" @@ ]`
}

type OpFactor struct {
	Pos lexer.Position

	Operator Operator `@("*" | "/")`
	Factor   *Factor  `@@`
}

type Term struct {
	Pos lexer.Position

	Left  *Factor     `@@`
	Right []*OpFactor `{ @@ }`
}

type OpTerm struct {
	Pos lexer.Position

	Operator Operator `@("+" | "-")`
	Term     *Term    `@@`
}

type Cmp struct {
	Pos lexer.Position

	Left  *Term     `@@`
	Right []*OpTerm `{ @@ }`
}

type OpCmp struct {
	Pos lexer.Position

	Operator Operator `@("=" | "<" "=" | ">" "=" | "<" | ">" | "!" "=")`
	Cmp      *Cmp     `@@`
}

type Expression struct {
	Pos lexer.Position

	Left  *Cmp     `@@`
	Right []*OpCmp `{ @@ }`
}

var (
	benjiLexer = lexer.Must(ebnf.New(`
		Comment = "#" { "\u0000"…"\uffff"-"\n"-"\r" } .
		Ident = (alpha | "_") { "_" | alpha | digit } .
		String = "\"" { "\u0000"…"\uffff"-"\""-"\\" | "\\" any } "\"" .
		Number = [ "-" | "+" ] ("." | digit) { "." | digit } .
		Punct = "!"…"/" | ":"…"@" | "["…` + "\"`\"" + ` | "{"…"~" .
		Whitespace = ( " " | "\t" | "\n" | "\r" ) { " " | "\t" | "\n" | "\r" } .

		alpha = "a"…"z" | "A"…"Z" .
		digit = "0"…"9" .
		any = "\u0000"…"\uffff" .
	`))

	Parser = participle.MustBuild(&Program{},
		participle.Lexer(benjiLexer),
		participle.CaseInsensitive("Ident"),
		participle.Unquote("String"),
		participle.UseLookahead(2),
		participle.Elide("Whitespace"),
	)
)
