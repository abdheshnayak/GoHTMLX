package gocode

import (
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"strings"
)

// ReplaceRenderE replaces RenderE function calls with the relevant code block
func ReplaceRenderE(input string, codes map[string][]byte) (string, error) {
	// Parse the Go code into an AST
	fs := token.NewFileSet()
	node, err := parser.ParseFile(fs, "", input, parser.ParseComments)
	if err != nil {
		return "", err
	}

	// Traverse the AST and replace RenderE calls
	ast.Inspect(node, func(n ast.Node) bool {
		// get pkg name
		pkg, ok := n.(*ast.File)
		if ok {
			pkg.Name.Name = fmt.Sprintf("gox%s", pkg.Name.Name)
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

				code := codes[strings.Trim(comp.Value, `"`)]

				e, err := parser.ParseExpr(string(code))
				if err != nil {
					fmt.Println("Error:", err)
					return true
				}
				ce, ok := e.(*ast.CallExpr)
				if !ok {
					return true
				}

				*call = *(ce)
			}
		}
		return true
	})

	// Re-generate the Go code from the modified AST (this step can be customized)
	var builder strings.Builder
	if err = format.Node(&builder, fs, node); err != nil {
		return "", err
	}
	return builder.String(), nil
}

// Helper function to convert argument list to a string
func getArgsAsString(args []ast.Expr) string {
	var argStrings []string
	for _, arg := range args {
		argStrings = append(argStrings, fmt.Sprintf("%v", arg))
	}
	return strings.Join(argStrings, ", ")
}
