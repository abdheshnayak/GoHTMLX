package element

import (
	"fmt"
	"io"
	"strings"
)

type Attr map[string]any
type Event any

type Element interface {
	Render(io.Writer) (int, error)
}

type element struct {
	tag      string
	attrs    Attr
	children []Element
}

type renderElement struct {
	items []interface{}
}

func (t renderElement) Render(w io.Writer) (int, error) {
	var buffer strings.Builder
	for _, item := range t.items {
		switch item.(type) {
		case string:
			buffer.WriteString(item.(string))
		case Element:
			item.(Element).Render(&buffer)
		case []Element:
			for _, child := range item.([]Element) {
				child.Render(&buffer)
			}
		default:
			fmt.Println(item.(string))
			buffer.WriteString(fmt.Sprintf("%v", item))
		}
	}

	return w.Write([]byte(buffer.String()))
}

func R(items ...interface{}) Element {
	return renderElement{
		items: items,
	}
}

func E(tag string, attrs map[string]any, children ...Element) Element {
	return element{
		tag:      tag,
		attrs:    attrs,
		children: children,
	}
}

func RenderE(element string, props ...interface{}) Element {
	return renderElement{
		items: []interface{}{},
	}
}

func (e element) Render(w io.Writer) (int, error) {
	var buffer strings.Builder
	buffer.WriteString("<")
	buffer.WriteString(e.tag)
	for k, v := range e.attrs {
		buffer.WriteString(" ")
		buffer.WriteString(k)
		buffer.WriteString("=\"")
		buffer.WriteString(fmt.Sprintf("%v", v))
		buffer.WriteString("\"")
	}
	buffer.WriteString(">")
	for _, child := range e.children {
		child.Render(&buffer)
	}
	buffer.WriteString("</")
	buffer.WriteString(e.tag)
	buffer.WriteString(">")

	return w.Write([]byte(buffer.String()))
}
