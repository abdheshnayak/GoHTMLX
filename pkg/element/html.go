package element

import (
	"bytes"
	"fmt"
	"regexp"
	"strings"

	"golang.org/x/net/html"
)

func capitalize(s string) string {
	if len(s) == 0 {
		return s
	}
	return strings.ToUpper(string(s[0])) + s[1:]
}

type Html interface {
	RenderGolangCode() ([]byte, error)
}

type htmlc struct {
	node *html.Node
}

func (h htmlc) RenderGolangCode() ([]byte, error) {
	// string writer
	var buffer strings.Builder
	b, err := render(h.node)
	if err != nil {
		return nil, err
	}

	buffer.Write(b)

	return []byte(buffer.String()), nil
}

func NewHtml(htmlCode []byte) (Html, error) {
	n, err := html.Parse(bytes.NewReader(htmlCode))
	if err != nil {
		return nil, err
	}

	return htmlc{
		node: n,
	}, nil
}

func transformString(template string) string {
	// Compile the regular expression to match {variable}
	// re := regexp.MustCompile(`\{([^\{\}]+(?:\{[^\{\}]*\}[^\{\}]*)*)\}`)

	re := regexp.MustCompile(`\{(.*)\}`)

	// Find all matches and replace them with the appropriate format
	parts := re.Split(template, -1)
	matches := re.FindAllStringSubmatch(template, -1)

	var result []string
	for i, part := range parts {
		if i > 0 {
			// For each match, add the variable name (e.g., "name")
			result = append(result, fmt.Sprintf("%s", matches[i-1][1]))
		}
		// Add the literal part (e.g., "hello ")
		result = append(result, fmt.Sprintf("`%s`", part))
	}

	// Join the parts with commas
	return fmt.Sprintf("R(%s)", strings.Join(result, ", "))
}

func render(n *html.Node) ([]byte, error) {
	var buffer strings.Builder

	switch n.Type {
	case html.TextNode:
		buffer.WriteString(transformString(n.Data))

	// case html.CommentNode:
	// 	buffer.WriteString("<!--")
	// 	buffer.WriteString(n.Data)
	// 	buffer.WriteString("-->")
	// case html.DoctypeNode:
	// 	buffer.WriteString("<!DOCTYPE ")
	// 	buffer.WriteString(n.Data)
	// 	buffer.WriteString(">")
	case html.ElementNode:
		if !isStandard(n.Data) {
			buffer.WriteString(fmt.Sprintf("%s(", capitalize(n.Data)))
		} else {
			buffer.WriteString("E(`")
			buffer.WriteString(n.Data)
			buffer.WriteString("`,")
		}

		buffer.WriteString("Attr{")
		for _, a := range n.Attr {
			buffer.WriteString(fmt.Sprintf("`%s`:", a.Key))
			// like {data} then remove {}
			if regexp.MustCompile(`\{(.*)\}`).MatchString(a.Val) {
				buffer.WriteString(regexp.MustCompile(`\{(.*)\}`).FindStringSubmatch(a.Val)[1])
			} else {
				buffer.WriteString(fmt.Sprintf("`%s`", a.Val))
			}
			buffer.WriteString(",")
		}
		buffer.WriteString("},")
		childs := []string{}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			b, err := render(c)
			if err != nil {
				return nil, err
			}
			childs = append(childs, string(b))
		}
		buffer.WriteString(strings.Join(childs, ","))
		buffer.WriteString(")")

	case html.DocumentNode:
		childs := []string{}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			b, err := render(c)
			if err != nil {
				return nil, err
			}
			childs = append(childs, string(b))
			// buffer.Write(b)
		}
		buffer.WriteString(strings.Join(childs, ","))
	default:
		return nil, nil
	}

	return []byte(buffer.String()), nil
}
