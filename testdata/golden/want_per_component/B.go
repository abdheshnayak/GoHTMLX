package gohtmlxc

import (
	. "github.com/abdheshnayak/gohtmlx/pkg/element"
)

type B struct {
	Attrs Attrs
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
