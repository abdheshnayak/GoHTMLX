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
	vmodel := vals["v-model"]

	return RenderE("inputcomp", name, placeholder, vmodel)
}

func Layout(vals Attr, children ...Element) Element {
	return RenderE("layout", children)
}

func Root(vals Attr, children ...Element) Element {
	id := vals["id"]

	return RenderE("root", children, id)
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

func HomeScript(vals Attr, _ ...Element) Element {
	return RenderE("HomeScript")
}
