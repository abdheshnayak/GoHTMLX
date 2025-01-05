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
	tag      string
	children []Element

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

func E(tag string, attrs map[string]any, children ...Element) Element {
	return element{
		tag:      tag,
		attrs:    attrs,
		children: children,
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
		case []byte:
			buffer.WriteString(string(v.([]byte)))
		case int:
			buffer.WriteString(fmt.Sprintf("%d", v.(int)))
		case float64:
			buffer.WriteString(fmt.Sprintf("%f", v.(float64)))
		case bool:
			buffer.WriteString(fmt.Sprintf("%t", v.(bool)))
		case []string:
			for _, child := range v.([]string) {
				buffer.WriteString(child)
			}
		case [][]byte:
			for _, child := range v.([][]byte) {
				buffer.WriteString(string(child))
			}
		case []int:
			for _, child := range v.([]int) {
				buffer.WriteString(fmt.Sprintf("%d", child))
			}
		case []float64:
			for _, child := range v.([]float64) {
				buffer.WriteString(fmt.Sprintf("%f", child))
			}
		case []bool:
			for _, child := range v.([]bool) {
				buffer.WriteString(fmt.Sprintf("%t", child))
			}
		case Element:
			v.(Element).Render(&buffer)
		case []Element:
			for _, child := range v.([]Element) {
				child.Render(&buffer)
			}
		case *string:
			buffer.WriteString(fmt.Sprintf("%s", *v.(*string)))
		case *[]byte:
			buffer.WriteString(fmt.Sprintf("%s", *v.(*[]byte)))
		case *int:
			buffer.WriteString(fmt.Sprintf("%d", *v.(*int)))
		case *float64:
			buffer.WriteString(fmt.Sprintf("%f", *v.(*float64)))
		case *bool:
			buffer.WriteString(fmt.Sprintf("%t", *v.(*bool)))
		case *Element:
			(*v.(*Element)).Render(&buffer)
		case *[]Element:
			for _, child := range *v.(*[]Element) {
				child.Render(&buffer)
			}
		case *[]string:
			for _, child := range *v.(*[]string) {
				buffer.WriteString(child)
			}
		case *[][]byte:
			for _, child := range *v.(*[][]byte) {
				buffer.WriteString(string(child))
			}
		case *[]int:
			for _, child := range *v.(*[]int) {
				buffer.WriteString(fmt.Sprintf("%d", child))
			}
		case *[]float64:
			for _, child := range *v.(*[]float64) {
				buffer.WriteString(fmt.Sprintf("%f", child))
			}
		case *[]bool:
			for _, child := range *v.(*[]bool) {
				buffer.WriteString(fmt.Sprintf("%t", child))
			}

		default:
			fmt.Printf("Unknown type: %T\n", v)
			buffer.WriteString(fmt.Sprintf("%v", v))
		}

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
