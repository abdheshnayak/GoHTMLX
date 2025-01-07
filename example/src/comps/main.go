package comps

import (
	gc "github.com/abdheshnayak/gohtmlx/example/dist/gohtmlxc"
	. "github.com/abdheshnayak/gohtmlx/pkg/element"
)

func Home() Element {
	name := "hi"

	return gc.Home{
		Id:   "app",
		Name: name,
	}.Get()
}

func Search() Element {
	return gc.SearchResult{}.Get()
}

func NoResult() Element {
	return gc.NoResult{}.Get()
}