package gohtmlxc

import (
	. "github.com/abdheshnayak/gohtmlx/pkg/element"
)

type A struct {
	Name  string
	Attrs Attrs
}

type B struct {
	Attrs Attrs
}

type sections struct {
	Attrs Attrs
}

func AComp(props A, attrs Attrs, children ...Element) Element {

	props.Attrs = attrs
	if props.Attrs == nil {
		props.Attrs = Attrs{}
	}

	return R(E(`div`, Attrs{`class`: `x`}, R(props.Name)))

}

func (c A) Get(children ...Element) Element {
	return AComp(c, c.Attrs, children...)
}

func BComp(props B, attrs Attrs, children ...Element) Element {

	props.Attrs = attrs
	if props.Attrs == nil {
		props.Attrs = Attrs{}
	}

	return R(E(`span`, Attrs{}, R(`static`)))

}

func (c B) Get(children ...Element) Element {
	return BComp(c, c.Attrs, children...)
}
