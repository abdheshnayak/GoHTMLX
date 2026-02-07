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

type forElement[T any] struct {
	items    []T
	children Element
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
		case int:
			buffer.WriteString(fmt.Sprintf("%d", item.(int)))
		case float64:
			buffer.WriteString(fmt.Sprintf("%f", item.(float64)))
		case bool:
			buffer.WriteString(fmt.Sprintf("%t", item.(bool)))
		case *string:
			buffer.WriteString(fmt.Sprintf("%s", *item.(*string)))
		case *int:
			buffer.WriteString(fmt.Sprintf("%d", *item.(*int)))
		case *float64:
			buffer.WriteString(fmt.Sprintf("%f", *item.(*float64)))
		case *bool:
			buffer.WriteString(fmt.Sprintf("%t", *item.(*bool)))
		case *Element:
			(*item.(*Element)).Render(&buffer)
		case *[]Element:
			for _, child := range *item.(*[]Element) {
				child.Render(&buffer)
			}
		case string:
			buffer.WriteString(strings.ReplaceAll(item.(string), nbspChar, "&nbsp;"))
			// buffer.WriteString(item.(string))

		case Element:
			item.(Element).Render(&buffer)
		case []Element:
			for _, child := range item.([]Element) {
				child.Render(&buffer)
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

type renderComp struct {
	name      string
	props     interface{}
	childrens []Element
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
			utils.Log.Error("unknown type", "for", fmt.Sprintf("%v", v))
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
