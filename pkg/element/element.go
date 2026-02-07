// Package element provides the runtime types and rendering used by GoHTMLX-generated code.
// Generated components produce Element values (via E, R, and component constructors) and
// call Render(w) to write HTML. Use this package only from generated code or when
// implementing custom elements; the transpiler lives in pkg/transpiler.
package element

import (
	"fmt"
	"io"
	"strings"

	"github.com/abdheshnayak/gohtmlx/pkg/utils"
)

const (
	nbspChar = " "
)

// Attrs holds attribute key-value pairs for an element (e.g. class, id).
// Used by generated code and by E(tag, attrs, children...).
type Attrs map[string]any

// GetAttr returns a typed attribute value from attrs, or nil if missing or wrong type.
func GetAttr[T any](attrs Attrs, key string) *T {
	if val, ok := attrs[key]; ok {
		resp, ok := val.(T)
		if ok {
			return &resp
		}
		utils.Log.Error("invalid type", "for", key)
		return nil
	}

	utils.Log.Error("key not found", "for", key)
	return nil
}

// Event is a placeholder type for future event handling.
type Event any

// Element is the interface produced by generated components. Render writes HTML to w.
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
		switch item := item.(type) {
		case int:
			buffer.WriteString(fmt.Sprintf("%d", item))
		case float64:
			buffer.WriteString(fmt.Sprintf("%f", item))
		case bool:
			buffer.WriteString(fmt.Sprintf("%t", item))
		case *string:
			buffer.WriteString(*item)
		case *int:
			buffer.WriteString(fmt.Sprintf("%d", *item))
		case *float64:
			buffer.WriteString(fmt.Sprintf("%f", *item))
		case *bool:
			buffer.WriteString(fmt.Sprintf("%t", *item))
		case *Element:
			_, _ = (*item).Render(&buffer)
		case *[]Element:
			for _, child := range *item {
				_, _ = child.Render(&buffer)
			}
		case string:
			buffer.WriteString(strings.ReplaceAll(item, nbspChar, "&nbsp;"))
			// buffer.WriteString(item.(string))

		case Element:
			_, _ = item.Render(&buffer)
		case []Element:
			for _, child := range item {
				_, _ = child.Render(&buffer)
			}
		default:
			utils.Log.Error("error", "for", fmt.Sprintf("%v", item))
			buffer.WriteString(fmt.Sprintf("%v", item))
		}
	}

	return w.Write([]byte(buffer.String()))
}

// R builds an Element from a mix of strings, Elements, and slices of Elements (used by generated code).
func R(items ...interface{}) Element {
	return renderElement{
		items: items,
	}
}

// E builds an HTML element with the given tag, attrs, and children (used by generated code).
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

		switch v := v.(type) {
		case string:
			buffer.WriteString(v)
		case *string:
			buffer.WriteString(*v)
		case Element:
			_, _ = v.Render(&buffer)
		case []Element:
			for _, child := range v {
				_, _ = child.Render(&buffer)
			}
		default:
			utils.Log.Error("unknown type", "for", fmt.Sprintf("%v", v))
			buffer.WriteString(fmt.Sprintf("%v", v))
		}

		buffer.WriteString("\"")
	}

	buffer.WriteString(">")
	for _, child := range e.childrens {
		_, _ = child.Render(&buffer)
	}
	buffer.WriteString("</")
	buffer.WriteString(e.tag)
	buffer.WriteString(">")

	return w.Write([]byte(buffer.String()))
}
