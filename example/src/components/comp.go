package pages

import (
	. "github.com/abdheshnayak/gohtmlx/pkg/element"
)

func Home() Element {
	name := "hi"

	return RenderE("home", name)
}

func InputComp(vals Attrs, _ ...Element) Element {
	name := GetAttr[string](vals, "name")
	placeholder := GetAttr[string](vals, "placeholder")
	vmodel := GetAttr[string](vals, "v-model")

	return RenderE("inputcomp", name, placeholder, vmodel)
}

func Layout(vals Attrs, children ...Element) Element {
	return RenderE("layout", children)
}

func Root(vals Attrs, children ...Element) Element {
	id := vals["id"]

	return RenderE("root", children, id)
}

func Search() Element {
	return RenderE("SearchResult")
}

func SearchCard(vals Attrs, _ ...Element) Element {
	return RenderE("SearchCard")
}

func NoResult(vals Attrs, _ ...Element) Element {
	return RenderE("NoResult")
}

func HomeScript(vals Attrs, _ ...Element) Element {
	return RenderE("HomeScript")
}
