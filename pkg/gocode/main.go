package gocode

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"strings"

	"github.com/abdheshnayak/gohtmlx/pkg/utils"
)

// ReplaceRenderE replaces RenderE function calls with the relevant code block
func ReplaceRenderE(input []byte, codes map[string][]byte) ([]byte, error) {
	wrap := false

	// Parse the Go code into an AST
	fs := token.NewFileSet()
	node, err := parser.ParseFile(fs, "", input, parser.ParseComments)
	if err != nil {
		return nil, err
	}

	// Traverse the AST and replace RenderE calls
	ast.Inspect(node, func(n ast.Node) bool {
		// get pkg name
		pkg, ok := n.(*ast.File)
		if ok {
			pkg.Name.Name = fmt.Sprintf("gohtmlx%s", pkg.Name.Name)
		}

		// Check if it's a function call
		if call, ok := n.(*ast.CallExpr); ok {
			// Check if it's a call to RenderE
			if ident, ok := call.Fun.(*ast.Ident); ok && ident.Name == "RenderE" {
				// Replace RenderE function with the relevant code block
				// For simplicity, replace it with a print statement for now
				// You can replace this with your relevant code
				if len(call.Args) == 0 {
					return true
				}

				comp, ok := call.Args[0].(*ast.BasicLit)
				if !ok {
					return true
				}

				if comp.Kind != token.STRING {
					return true
				}

				code, ok := codes[strings.Trim(strings.ToLower(comp.Value), `"`)]
				if !ok {
					utils.Log.Warn("no code found", "for", comp.Value)
					return true
				}

				e, err := parser.ParseExpr(fmt.Sprintf("R(%s)", string(code)))
				if err != nil {
					utils.Log.Error("failed to parse", "expression", err)
					return true
				}
				ce, ok := e.(*ast.CallExpr)
				if !ok {
					return true
				}

				if wrap {
					k := *call
					k.Fun.(*ast.Ident).Name = "Re"
					ce.Args = append(ce.Args, &k)
				}

				*call = *(ce)
			}
		}
		return true
	})

	// Re-generate the Go code from the modified AST (this step can be customized)
	var buf bytes.Buffer
	if err = format.Node(&buf, fs, node); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// Helper function to convert argument list to a string
func getArgsAsString(args []ast.Expr) string {
	var argStrings []string
	for _, arg := range args {
		argStrings = append(argStrings, fmt.Sprintf("%v", arg))
	}
	return strings.Join(argStrings, ", ")
}
