package element

import (
	"bytes"
	"fmt"
	"regexp"
	"strings"

	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

func capitalize(s string) string {
	if len(s) == 0 {
		return s
	}
	return strings.ToUpper(string(s[0])) + s[1:]
}

type Html interface {
	RenderGolangCode(comps map[string]string) ([]byte, error)
}

type htmlc struct {
	nodes []*html.Node
}

func (h htmlc) RenderGolangCode(comps map[string]string) ([]byte, error) {
	// string writer
	var buffer strings.Builder
	bts := [][]byte{}

	for _, n := range h.nodes {
		b, err := render(n, comps)
		if err != nil {
			return nil, err
		}

		bts = append(bts, b)
		// buffer.Write(b)
	}

	buffer.Write(bytes.Join(bts, []byte(",")))

	return []byte(buffer.String()), nil
}

func NewHtml(htmlCode []byte) (Html, error) {
	context := &html.Node{
		Type:     html.ElementNode,
		DataAtom: atom.Div,
		Data:     "div",
	}

	if bytes.Contains(htmlCode, []byte("<html>")) && bytes.Contains(htmlCode, []byte("</html>")) {
		context = nil
	}

	n, err := html.ParseFragment(bytes.NewReader(bytes.Trim(bytes.TrimSpace(htmlCode), "\n")), context)
	if err != nil {
		return nil, err
	}

	return htmlc{
		nodes: n,
	}, nil
}

func trimString(s string) string {
	return strings.Trim(strings.TrimSpace(s), "\n")
}
func trimBytes(b []byte) []byte {
	return bytes.Trim(bytes.TrimSpace(b), "\n")
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
			result = append(result, fmt.Sprintf("%s", trimString(matches[i-1][1])))
		}

		if trimString(part) == "" {
			continue
		}

		// Add the literal part (e.g., "hello ")
		result = append(result, fmt.Sprintf("`%s`", trimString(part)))
	}

	if len(result) == 0 {
		return ""
	}

	// Join the parts with commas
	return fmt.Sprintf("R(%s)", trimString(strings.Join(result, ", ")))
}

func render(n *html.Node, comps map[string]string) ([]byte, error) {
	var buffer strings.Builder

	switch n.Type {
	case html.TextNode:
		buffer.WriteString(transformString(strings.TrimSpace(n.Data)))

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
			buffer.WriteString(fmt.Sprintf("%s(", comps[trimString(n.Data)]))
		} else {
			buffer.WriteString("E(`")
			buffer.WriteString(strings.TrimSpace(n.Data))
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
		if n.Data == "script" {
			// childs = append(childs, "E(``)")

      buffer.WriteString("R(`")
			for c := n.FirstChild; c != nil; c = c.NextSibling {
				buffer.WriteString(c.Data)
			}

			buffer.WriteString("`))")
		} else {
			for c := n.FirstChild; c != nil; c = c.NextSibling {
				b, err := render(c, comps)
				if err != nil {
					return nil, err
				}

				if len(b) == 0 {
					continue
				}

				childs = append(childs, string(b))
			}

			buffer.WriteString(strings.Join(childs, ","))
			buffer.WriteString(")")
		}

	case html.DocumentNode:

		childs := []string{}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			b, err := render(c, comps)
			if err != nil {
				return nil, err
			}

			if len(b) == 0 {
				continue
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
