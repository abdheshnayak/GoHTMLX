// Minimal example: one component, no framework. Renders HTML to stdout.
package main

import (
	"os"

	gc "github.com/abdheshnayak/gohtmlx/examples/minimal/dist/gohtmlxc"
)

func main() {
	el := gc.Hello{Name: "GoHTMLX", Attrs: nil}.Get()
	if _, err := el.Render(os.Stdout); err != nil {
		panic(err)
	}
}
