package element

import (
	"fmt"
	"html"
	"io"
	"strings"
)

const (
	nbspChar = " "
)

// Attrs represents HTML attributes as a map
type Attrs map[string]any

// GetAttr safely retrieves a typed attribute value
func GetAttr[T any](attrs Attrs, key string) *T {
	if val, ok := attrs[key]; ok {
		if resp, ok := val.(T); ok {
			return &resp
		}
	}
	return nil
}

// GetAttrWithDefault retrieves a typed attribute value with a default fallback
func GetAttrWithDefault[T any](attrs Attrs, key string, defaultValue T) T {
	if val := GetAttr[T](attrs, key); val != nil {
		return *val
	}
	return defaultValue
}

// HasAttr checks if an attribute exists
func HasAttr(attrs Attrs, key string) bool {
	_, ok := attrs[key]
	return ok
}

// SetAttr sets an attribute value
func SetAttr(attrs Attrs, key string, value any) {
	if attrs == nil {
		attrs = make(Attrs)
	}
	attrs[key] = value
}

// MergeAttrs merges multiple attribute maps, with later maps overriding earlier ones
func MergeAttrs(attrMaps ...Attrs) Attrs {
	result := make(Attrs)
	for _, attrs := range attrMaps {
		for k, v := range attrs {
			result[k] = v
		}
	}
	return result
}

// ClassNames builds a class string from multiple class values
func ClassNames(classes ...any) string {
	var result []string
	for _, class := range classes {
		switch v := class.(type) {
		case string:
			if v != "" {
				result = append(result, v)
			}
		case []string:
			for _, c := range v {
				if c != "" {
					result = append(result, c)
				}
			}
		case map[string]bool:
			for c, include := range v {
				if include && c != "" {
					result = append(result, c)
				}
			}
		}
	}
	return strings.Join(result, " ")
}

type Event any

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
			buffer.WriteString(html.EscapeString(item.(string)))

		case Element:
			item.(Element).Render(&buffer)
		case []Element:
			for _, child := range item.([]Element) {
				child.Render(&buffer)
			}
		default:
			// Log error if needed - for now just render as string
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

	// Self-closing tags
	selfClosing := map[string]bool{
		"area": true, "base": true, "br": true, "col": true, "embed": true,
		"hr": true, "img": true, "input": true, "link": true, "meta": true,
		"param": true, "source": true, "track": true, "wbr": true,
	}

	buffer.WriteString("<")
	buffer.WriteString(e.tag)

	// Render attributes
	for k, v := range e.attrs {
		// Handle boolean attributes
		if boolVal, ok := v.(bool); ok {
			if boolVal {
				// True boolean attributes are rendered without values
				buffer.WriteString(" ")
				buffer.WriteString(html.EscapeString(k))
			}
			// False boolean attributes are omitted entirely
			continue
		}

		buffer.WriteString(" ")
		buffer.WriteString(html.EscapeString(k))
		buffer.WriteString("=\"")

		switch val := v.(type) {
		case string:
			buffer.WriteString(html.EscapeString(val))
		case *string:
			if val != nil {
				buffer.WriteString(html.EscapeString(*val))
			}
		case int:
			buffer.WriteString(fmt.Sprintf("%d", val))
		case float64:
			buffer.WriteString(fmt.Sprintf("%g", val))
		case Element:
			val.Render(&buffer)
		case []Element:
			for _, child := range val {
				child.Render(&buffer)
			}
		default:
			buffer.WriteString(html.EscapeString(fmt.Sprintf("%v", val)))
		}

		buffer.WriteString("\"")
	}

	if selfClosing[e.tag] && len(e.childrens) == 0 {
		buffer.WriteString(" />")
	} else {
		buffer.WriteString(">")
		for _, child := range e.childrens {
			child.Render(&buffer)
		}
		buffer.WriteString("</")
		buffer.WriteString(e.tag)
		buffer.WriteString(">")
	}

	return w.Write([]byte(buffer.String()))
}
