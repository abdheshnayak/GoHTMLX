package gohtmlxc

import (
	. "github.com/abdheshnayak/gohtmlx/pkg/element"
)

type Hello struct {
	Name  string
	Attrs Attrs
}

func HelloComp(props Hello, attrs Attrs, children ...Element) Element {
	props.Attrs = attrs
	if props.Attrs == nil {
		props.Attrs = Attrs{}
	}
	return R(E(`div`, Attrs{`class`: `greeting`}, R(`
  `), E(`p`, Attrs{}, R(R(`Hello, `, props.Name, `!`))), R(`
`)))
}

func (c Hello) Get(children ...Element) Element {
	return HelloComp(c, c.Attrs, children...)
}
