package pages

import (
	"fmt"

	. "github.com/abdheshnayak/gox/pkg/element"
)

func Home() Element {
	list := "list"
	name := "hello"
	onClick := func(e Event) {
		fmt.Println("Clicked")
	}

	return RenderE("home", list, name, onClick)
}

func Inputcomp(vals Attr, _ ...Element) Element {
	name := vals["name"]
	value := vals["value"]

	return RenderE("inputcomp", name, value)
}
