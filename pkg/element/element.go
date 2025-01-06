package element

import (
	"fmt"
	"io"
	"log/slog"
	"strings"
)

type Attrs map[string]any

func GetAttr[T any](attrs Attrs, key string) *T {
	if val, ok := attrs[key]; ok {
		resp, ok := val.(T)
		if ok {
			return &resp
		}
		slog.Error(fmt.Sprintf("value for %s is not of type %T", key, val))
	}

	slog.Error(fmt.Sprintf("key %s not found", key))
	return nil
}

type Event any

type Element interface {
	Render(io.Writer) (int, error)
}

type element struct {
	tag       string
	childrens []Element

	attrs Attrs
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

type renderComp struct {
	name      string
	props     interface{}
	childrens []Element
}

func E(tag string, attrs Attrs, childrens ...Element) Element {
	return element{
		tag:       tag,
		attrs:     attrs,
		childrens: childrens,
	}
}

func Re(element string, props ...interface{}) Element {
	return renderElement{
		items: []interface{}{},
	}
}

func RenderE(element string, props ...interface{}) Element {
	return Re(element, props...)
}

func (e element) Render(w io.Writer) (int, error) {
	var buffer strings.Builder
	buffer.WriteString("<")
	buffer.WriteString(e.tag)
	for k, v := range e.attrs {
		buffer.WriteString(" ")
		buffer.WriteString(k)
		buffer.WriteString("=\"")

		switch v.(type) {
		case string:
			buffer.WriteString(v.(string))
		case *string:
			buffer.WriteString(fmt.Sprintf("%s", *v.(*string)))
		case Element:
			v.(Element).Render(&buffer)
		case []Element:
			for _, child := range v.([]Element) {
				child.Render(&buffer)
			}
		default:
			fmt.Printf("Unknown type: %T\n", v)
			buffer.WriteString(fmt.Sprintf("%v", v))
		}

		buffer.WriteString("\"")
	}

	buffer.WriteString(">")
	for _, child := range e.childrens {
		child.Render(&buffer)
	}
	buffer.WriteString("</")
	buffer.WriteString(e.tag)
	buffer.WriteString(">")

	return w.Write([]byte(buffer.String()))
}
