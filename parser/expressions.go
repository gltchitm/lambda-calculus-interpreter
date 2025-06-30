package parser

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/gltchitm/lambda-calculus-interpreter/tokenizer"
)

type Expression interface {
	Identity() string
	String() string
}

type Application struct {
	Left  Expression
	Right Expression
}

func NewApplication(left, right Expression) *Application {
	return &Application{left, right}
}

func (a Application) Identity() string {
	return "Application"
}

func (a Application) String() string {
	return fmt.Sprintf("(%v %v)", a.Left, a.Right)
}

type Assignment struct {
	To    *Variable
	Value Expression
}

func NewAssignment(to *Variable, value Expression) *Assignment {
	return &Assignment{to, value}
}

func (a Assignment) Identity() string {
	return "Assignment"
}

func (a Assignment) String() string {
	return fmt.Sprintf("%v = %v", a.To, a.Value)
}

type Function struct {
	Parameter *Variable
	Body      Expression
}

func NewFunction(parameter *Variable, body Expression) *Function {
	return &Function{parameter, body}
}

func (f Function) Identity() string {
	return "Function"
}

func (f Function) String() string {
	return fmt.Sprintf(
		"(%s%v.%v)",
		color.HiGreenString(tokenizer.LambdaSymbol),
		f.Parameter,
		f.Body,
	)
}

type Populate struct {
	From int
	To   int
}

func NewPopulate(from, to int) *Populate {
	return &Populate{from, to}
}

func (p Populate) Identity() string {
	return "Populate"
}

func (p Populate) String() string {
	return fmt.Sprintf(
		"(%s %d %s %d)",
		color.New(color.FgHiYellow, color.Italic).Sprint("populate"),
		p.From,
		color.HiYellowString("to"),
		p.To,
	)
}

type Run struct {
	Body Expression
}

func NewRun(body Expression) *Run {
	return &Run{body}
}

func (r Run) String() string {
	return fmt.Sprintf(
		"(%s %v)",
		color.New(color.FgHiYellow, color.Italic).Sprint("run"),
		r.Body,
	)
}

func (r Run) Identity() string {
	return "Run"
}

type Variable struct {
	Name    string
	IsBound bool
}

func NewVariable(name string, isBound bool) *Variable {
	return &Variable{name, isBound}
}

func (v Variable) Identity() string {
	return "Variable"
}

func (v Variable) String() string {
	highlightColor := color.FgHiCyan
	if v.IsBound {
		highlightColor = color.FgHiMagenta
	}

	return color.New(highlightColor).Sprint(v.Name)
}
