package pages

import (
	. "github.com/abdheshnayak/govelte/pkg/element"
)

func Home() Element {
	name := "hello"

	return RenderE("home", name)
}

func Inputcomp(vals Attr, _ ...Element) Element {
	name := vals["name"]
	placeholder := vals["placeholder"]

	return RenderE("inputcomp", name, placeholder)
}

func Layout(vals Attr, children ...Element) Element {
	return RenderE("layout", children)
}

func Root(vals Attr, children ...Element) Element {
	return RenderE("root", children)
}

func Search() Element {
	return RenderE("SearchResult")
}

func Searchcard(vals Attr, _ ...Element) Element {
	return RenderE("SearchCard")
}

func Noresult(vals Attr, _ ...Element) Element {
	return RenderE("NoResult")
}
