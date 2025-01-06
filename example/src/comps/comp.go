package comps

import (
	gc "github.com/abdheshnayak/gohtmlx/example/dist/gohtmlxc"
	. "github.com/abdheshnayak/gohtmlx/pkg/element"
)

func Home() Element {
	name := "hi"

	return gc.Home(gc.HomeProps{
		Id:   "app",
		Name: name,
	})
}

func Search() Element {
	return RenderE("SearchResult")
}

func NoResult() Element {
	return RenderE("NoResult")
}
