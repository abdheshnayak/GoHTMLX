package element

import (
	"bytes"
	"fmt"
	"regexp"
	"strings"

	"github.com/abdheshnayak/gohtmlx/pkg/utils"
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

type CompInfo struct {
	Name  string
	Props map[string]string
}

type Html interface {
	RenderGolangCode(comps map[string]CompInfo) (string, error)
}

type htmlc struct {
	nodes []*html.Node
}

func processFor(n *html.Node, comps map[string]CompInfo) (string, error) {
	var buffer strings.Builder

	key := ""
	as := ""
	for _, a := range n.Attr {
		if a.Key == "items" {
			key = a.Val
		}
		if a.Key == "as" {
			as = a.Val
		}

		if key != "" && as != "" {
			break
		}
	}

	if as == "" {
		as = "item"
	}

	if key == "" {
		return "", fmt.Errorf("key items not found in 'for' element")

	}
	if strings.HasPrefix(key, "{$attrs.") {
		return "", fmt.Errorf("cannot use $attrs in 'for' element, please use props instead")
	}
	if !strings.HasPrefix(key, "{") || !strings.HasSuffix(key, "}") {
		return "", fmt.Errorf("invalid key %s in 'for' element", key)
	}

	key = processRaws(key)

	// key = strings.Trim(key, "{}")

	buffer.WriteString("R(func() []Element {\n")
	buffer.WriteString("resp := []Element{}\n")

	buffer.WriteString(fmt.Sprintf("for _, %s := range %s {\n", as, key))
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		b, err := render(c, comps)
		if err != nil {
			return "", err
		}

		if len(b) == 0 {
			continue
		}

		buffer.WriteString(fmt.Sprintf("resp = append(resp, %s)\n", string(b)))
	}
	buffer.WriteString("}\n")
	buffer.WriteString("return resp\n")
	buffer.WriteString("}(),)")

	return buffer.String(), nil
}

func generateProps(n *html.Node, comps map[string]CompInfo) (string, error) {
	isStd := isStandard(n.Data)

	var buffer strings.Builder

	if isStd {
		buffer.WriteString("E(`")
		buffer.WriteString(strings.TrimSpace(n.Data))
		buffer.WriteString("`,")
		buffer.WriteString("Attrs{")

		for _, a := range n.Attr {
			if isStd {
				buffer.WriteString(fmt.Sprintf("`%s`:", a.Key))
			} else {
				if prop, ok := comps[n.Data].Props[a.Key]; ok {
					buffer.WriteString(fmt.Sprintf("%s:", utils.Capitalize(prop)))
				} else {
					buffer.WriteString(fmt.Sprintf("%s:", a.Key))
				}
			}
			buffer.WriteString(processRaws(a.Val))
			buffer.WriteString(",")
		}
		buffer.WriteString("},")

		return buffer.String(), nil
	}

	var props strings.Builder
	var attrs strings.Builder
	for _, a := range n.Attr {
		if prop, ok := comps[n.Data].Props[a.Key]; ok {
			props.WriteString(fmt.Sprintf("%s:", utils.Capitalize(prop)))

			props.WriteString(processRaws(a.Val))
			props.WriteString(",")
		} else {
			attrs.WriteString(fmt.Sprintf("`%s`:", a.Key))

			attrs.WriteString(processRaws(a.Val))
			attrs.WriteString(",")
		}

	}

	buffer.WriteString(fmt.Sprintf("%sComp(", comps[strings.TrimSpace(n.Data)].Name))
	buffer.WriteString(fmt.Sprintf("%s{%s},", comps[strings.TrimSpace(n.Data)].Name, props.String()))
	buffer.WriteString(fmt.Sprintf("Attrs{%s},", attrs.String()))

	return buffer.String(), nil
}

func (h htmlc) RenderGolangCode(comps map[string]CompInfo) (string, error) {
	// string writer
	var buffer strings.Builder
	bts := []string{}

	buffer.WriteString("R(")
	for _, n := range h.nodes {
		b, err := render(n, comps)
		if err != nil {
			return "", err
		}

		bts = append(bts, b)
		// buffer.Write(b)
	}

	buffer.WriteString(strings.Join(bts, ","))
	buffer.WriteString(")")

	return buffer.String(), nil
}

func NewHtml(htmlCode []byte) (Html, error) {
	context := &html.Node{
		Type:     html.ElementNode,
		DataAtom: atom.Div,
		Data:     "div",
	}

	if bytes.Contains(htmlCode, []byte("<html>")) || bytes.Contains(htmlCode, []byte("</html>")) {
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

func processNode(input string) string {
	// Regular expression to match {item} or {{item}}
	varPattern := regexp.MustCompile(`\{{2,}(.*)\}{2,}`)
	tokens := []string{}

	// Split the input string into parts based on variable patterns
	splitParts := varPattern.Split(input, -1)
	matches := varPattern.FindAllStringSubmatch(input, -1)

	// Iterate over the split parts and matches to construct the result
	for i, part := range splitParts {
		tokens = append(tokens, processRaws(part))

		if i < len(matches) {
			match := matches[i]
			tokens = append(tokens, fmt.Sprintf("`%s`", match[0]))
		}
	}

	inners := strings.Join(tokens, ", ")
	if len(inners) == 0 {
		return ""
	}
	// Join tokens to form the final R(...) string
	result := fmt.Sprintf("R(%s)", strings.Join(tokens, ", "))
	return result
}

func processRaws(input string) string {
	// TODO: update this regex to handle cases where multiple expressions can present in single line
	// eg: {item} {item2}
	// TODO: update this regex to handle cases witch single expression present in multiple lines
	/*
			   eg:
		        {func(inp []string) string {
		            return strings.Join(inp)
		        }(items)}
	*/

	re := regexp.MustCompile(`\{(.*)\}`)

	if re.MatchString(input) {
		splitParts := re.Split(input, -1)
		matches := re.FindAllStringSubmatch(input, -1)

		tokens := []string{}

		for i, v := range splitParts {
			if v != "" {
				tokens = append(tokens, fmt.Sprintf("`%s`", v))
			}

			if i < len(matches) {
				val := matches[i][1]
				if val == "" {
					continue
				}

				if strings.HasPrefix(val, "$") {
					f := strings.Split(val, ".")
					tokens = append(tokens, fmt.Sprintf("%s[\"%s\"]", strings.Replace(f[0], "$", "", 1), f[1]))
				} else {
					if strings.HasPrefix(val, "props.") {
						val = val[:6] + strings.ToUpper(val[6:7]) + val[7:]
					}
					tokens = append(tokens, val)
				}
			}
		}

		if len(tokens) > 1 {
			return fmt.Sprintf("R(%s)", strings.Join(tokens, ","))
		}

		return strings.Join(tokens, ",")

	} else {
		return fmt.Sprintf("`%s`", input)
	}
}

func render(n *html.Node, comps map[string]CompInfo) (string, error) {
	var buffer strings.Builder

	switch n.Type {
	case html.TextNode:
		buffer.WriteString(processNode(n.Data))

	// case html.CommentNode:
	// 	buffer.WriteString("<!--")
	// 	buffer.WriteString(n.Data)
	// 	buffer.WriteString("-->")
	// case html.DoctypeNode:
	// 	buffer.WriteString("<!DOCTYPE ")
	// 	buffer.WriteString(n.Data)
	// 	buffer.WriteString(">")
	case html.ElementNode:
		if n.Data == "for" {
			s, err := processFor(n, comps)
			if err != nil {
				return "", err
			}
			buffer.WriteString(s)
		} else {

			s, err := generateProps(n, comps)
			if err != nil {
				return "", err
			}

			buffer.WriteString(s)

			childs := []string{}
			if n.Data == "script" || n.Data == "style" {
				buffer.WriteString("R(`")
				for c := n.FirstChild; c != nil; c = c.NextSibling {
					buffer.WriteString(c.Data)
				}

				buffer.WriteString("`))")
			} else {

				for c := n.FirstChild; c != nil; c = c.NextSibling {
					b, err := render(c, comps)
					if err != nil {
						return "", err
					}

					if len(b) == 0 {
						continue
					}

					childs = append(childs, string(b))
				}

				buffer.WriteString(strings.Join(childs, ","))
				buffer.WriteString(")")
			}

		}

	case html.DocumentNode:

		childs := []string{}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			b, err := render(c, comps)
			if err != nil {
				return "", err
			}

			if len(b) == 0 {
				continue
			}

			childs = append(childs, string(b))
		}

		buffer.WriteString(strings.Join(childs, ","))
	default:
		return "", nil
	}

	return buffer.String(), nil
}
