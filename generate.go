package substate

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

// TemplateData contains all data which is required by the template.
type TemplateData struct {
	// The source of the go:generate directive.
	Source string

	// The package in which the output file will be written.
	Package string

	Imports []string
	Fields  []Field
}

// Import contains metadata for a package import.
type Import struct {
	Alias string
	Path  string
}

// Field contains template data about a method in the Substate interface.
type Field struct {
	// The name of the attribute within the generated struct.
	Name string

	// The name of the method in the Substate interface which should be
	// implemented on the generated struct.
	Method string

	// The Go type of the field.
	Type string

	// String of arguments the method expects.
	Params string

	// String of arguments the method returns.
	Results string
}

// Generate reads the source file and generates a Substate implementation on a
// concrete type in order to be used to testing.
func Generate(src, dest, typeName string) error {
	pwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("get working dir: %w", err)
	}

	// Grab the source code - we want to be able to grab certain snippets.
	srcData, err := ioutil.ReadFile(filepath.Join(pwd, src))
	if err != nil {
		log.Fatal(err)
	}

	fset := token.NewFileSet()
	tree, err := parser.ParseFile(fset, src, nil, 0)
	if err != nil {
		log.Fatal(err)
	}

	data := TemplateData{
		Source: src,
		// "testing" won't be imported by the source file, so we've got to add
		// it here by default.
		Imports: []string{"\"testing\""},
	}

	var (
		errs  = make([]error, 0)
		found = false
	)

	ast.Inspect(tree, func(n ast.Node) bool {
		switch t := n.(type) {
		case *ast.File:
			data.Package = t.Name.Name
			return true
		case *ast.TypeSpec:
			// We only care about the one interface.
			if strings.EqualFold(t.Name.Name, typeName) {
				found = true
				return true
			}

			return false
		case *ast.ImportSpec:
			var alias string
			if t.Name != nil {
				alias = t.Name.Name
			}

			importString := strings.TrimSpace(fmt.Sprintf("%s %s", alias,
				t.Path.Value))
			data.Imports = append(data.Imports, importString)

			return true
		case *ast.InterfaceType:
			for _, method := range t.Methods.List {
				var params []string
				var returns []string

				fn, ok := method.Type.(*ast.FuncType)
				if !ok {
					return false
				}

				if fn.Params != nil {
					for _, param := range fn.Params.List {
						typ := srcData[param.Type.Pos()-1 : param.Type.End()-1]
						params = append(params, string(typ))
					}
				}

				if fn.Results != nil {
					if len(fn.Results.List) != 1 {
						// TODO: Support multiple return args.
						errs = append(errs, ErrMultipleReturnArguments)
						break
					}

					arg := fn.Results.List[0]

					if _, ok := arg.Type.(*ast.FuncType); ok {
						errs = append(errs, ErrFunctionReturnArgument)
						break
					}

					argType := srcData[arg.Type.Pos()-1 : arg.Type.End()-1]
					returns = append(returns, string(argType))
				}

				f := Field{
					Method:  method.Names[0].Name,
					Name:    camelCase(method.Names[0].Name),
					Params:  paramString(params),
					Results: returnString(returns),
					Type:    returnString(returns),
				}

				data.Fields = append(data.Fields, f)
			}

			return false
		default:
			return true
		}
	})

	// Only generate if a Substate interface was found.
	if !found {
		return ErrSubstateNotFound
	}

	// If any errors occurred during the inspection phase, return the earliest
	// one.
	if len(errs) > 0 {
		return errs[0]
	}

	t, err := template.New("").Parse(tpl)
	if err != nil {
		return fmt.Errorf("parse template: %w", err)
	}

	var b bytes.Buffer
	if err = t.Execute(&b, data); err != nil {
		return fmt.Errorf("execute template: %w", err)
	}

	// Format the source code before writing.
	output, err := format.Source(b.Bytes())
	if err != nil {
		return fmt.Errorf("format output bytes: %w", err)
	}

	if err = os.WriteFile(dest, output, 0644); err != nil {
		return fmt.Errorf("write output file: %w", err)
	}

	return nil
}

func paramString(params []string) string {
	var builder strings.Builder

	builder.WriteString("(")
	for i, paramType := range params {
		builder.WriteString(fmt.Sprintf("arg%d %s", i, paramType))
	}

	builder.WriteString(")")

	return builder.String()
}

func returnString(returns []string) string {
	if len(returns) == 0 {
		return ""
	}

	if len(returns) == 1 {
		return returns[0]
	}

	return fmt.Sprintf("(%s)", strings.Join(returns, ", "))
}

func camelCase(s string) string {
	return strings.ToLower(s[:1]) + s[1:]
}
