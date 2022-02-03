package main

import (
	"fmt"
	"go/format"
	"io/ioutil"
	"os"
	"strings"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Error, invalid args.")
		return
	}
	outDir := os.Args[1]
	var types []string

	types = []string{
		"Binary   : Expr left, Token operator, Expr right",
		"Grouping : Expr expression",
		"Literal  : Object value",
		"Unary    : Token operator, Expr right",
	}

	defineAst(outDir, "expr", types)
}

func defineAst(out, baseName string, types []string) {
	var src string
	src += fmt.Sprintln("package main")
	src += fmt.Sprintln("")
	src += fmt.Sprintln("")
	src += defineVisitor(baseName, types)
	src += "\n"

	// define base interface
	src += fmt.Sprintf("type %v interface{\n", strings.Title(baseName))
	src += fmt.Sprintf("    Accept(v %v) interface{}\n", strings.Title(baseName)+"Visitor")
	src += "}\n"
	src += "\n"

	for _, t := range types {
		className := strings.TrimRight(strings.Split(t, ":")[0], " ")
		fields := strings.TrimLeft(strings.Split(t, ":")[1], " ")
		fields = strings.TrimRight(fields, " ")
		src += defineType(baseName, className, fields)
		src += fmt.Sprintln("")
	}

	err := writeFile(fmt.Sprintf("%v/%v.go", out, baseName), src)
	if err != nil {
		fmt.Println("Error", err)
	}
}

func defineVisitor(baseName string, types []string) string {
	var src string
	src += fmt.Sprintf("type %v interface {", strings.Title(baseName)+"Visitor")
	src += fmt.Sprintln("")

	for _, t := range types {
		ftype, _ := parseField(t)
		src += fmt.Sprintf("    Visit%v%v(%v *%v) interface{}\n",
			strings.Title(ftype), strings.Title(baseName), strings.ToLower(baseName), strings.Title(ftype))
	}

	src += "}"
	src += "\n"
	return src
}

func defineType(baseName, className, fields string) string {
	var src string

	// define struct
	src += fmt.Sprintf("type %v struct {", className)
	src += fmt.Sprintln("")
	for _, f := range strings.Split(fields, ",") {
		ftype, fname := parseField(f)
		if ftype == "Object" {
			ftype = "interface{}"
		}
		if ftype == "Token" {
			ftype = "*Token"
		}
		src += fmt.Sprintf("    %v %v", fname, ftype)
		src += fmt.Sprintln("")
	}
	src += "}"
	src += fmt.Sprintln("")

	var params []string
	for _, f := range strings.Split(fields, ",") {
		ftype, fname := parseField(f)
		if ftype == "Object" {
			ftype = "interface{}"
		}
		if ftype == "Token" {
			ftype = "*Token"
		}
		params = append(params, fmt.Sprintf("%v %v", fname, ftype))
	}

	// define NewXXX
	src += fmt.Sprintf("func New%v(%v) %v {", className, strings.Join(params, ","), strings.Title(baseName))
	src += fmt.Sprintln("")
	src += fmt.Sprintln("")
	src += fmt.Sprintf("    %v := &%v{\n", baseName, className)
	for _, f := range strings.Split(fields, ",") {
		_, fname := parseField(f)
		src += fmt.Sprintf("        %v: %v,\n", fname, fname)
	}
	src += "    }"
	src += fmt.Sprintln("")
	src += fmt.Sprintf("    return %v", baseName)
	src += fmt.Sprintln("")
	src += "}"
	src += fmt.Sprintln("")
	src += fmt.Sprintln("")
	// define accept method.
	src += fmt.Sprintf("func (%v *%v) Accept(v %vVisitor) interface{} {", strings.ToLower(baseName), className, strings.Title(baseName))
	src += fmt.Sprintln("")
	src += fmt.Sprintf("    return v.Visit%v%v(%v)", strings.Title(className), strings.Title(baseName), strings.ToLower(baseName))
	src += fmt.Sprintln("")
	src += "}"

	return src
}

func parseField(field string) (typ, name string) {
	field = strings.Trim(field, " ")
	typ, name = strings.Split(field, " ")[0], strings.Split(field, " ")[1]
	return
}

func writeFile(path, str string) error {
	buf, err := format.Source([]byte(str))
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(path, []byte(buf), 0644)
	return err
}
