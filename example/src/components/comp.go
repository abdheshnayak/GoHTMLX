package pages

import (
	. "github.com/abdheshnayak/govelte/pkg/element"
)

func Home() Element {
	name := "hi"

	return RenderE("home", name)
}

func InputComp(vals Attr, _ ...Element) Element {
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

func SearchCard(vals Attr, _ ...Element) Element {
	return RenderE("SearchCard")
}

func NoResult(vals Attr, _ ...Element) Element {
	return RenderE("NoResult")
}
