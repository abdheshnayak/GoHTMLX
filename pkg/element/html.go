package element

import (
	"bytes"
	"fmt"
	"regexp"
	"sort"
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

// processIfChain handles <if condition={expr}>...</if> and optional <elseif condition={}>...</elseif>, <else>...</else>.
// Returns generated code and the last node consumed (so caller can skip to last.NextSibling).
func processIfChain(ifNode *html.Node, comps map[string]CompInfo) (string, *html.Node, error) {
	cond, err := getConditionAttr(ifNode)
	if err != nil {
		return "", nil, err
	}
	condGo := processRaws(cond)

	var parts []string // "if cond { return []Element{...} }" etc.
	thenCode, err := renderChildren(ifNode, comps)
	if err != nil {
		return "", nil, err
	}
	parts = append(parts, fmt.Sprintf("if %s {\nreturn []Element{%s}\n}", condGo, thenCode))

	last := ifNode
	hadElse := false
	for sib := ifNode.NextSibling; sib != nil; sib = sib.NextSibling {
		switch sib.Data {
		case "elseif":
			c, err := getConditionAttr(sib)
			if err != nil {
				return "", nil, err
			}
			cGo := processRaws(c)
			body, err := renderChildren(sib, comps)
			if err != nil {
				return "", nil, err
			}
			parts = append(parts, fmt.Sprintf("if %s {\nreturn []Element{%s}\n}", cGo, body))
			last = sib
		case "else":
			body, err := renderChildren(sib, comps)
			if err != nil {
				return "", nil, err
			}
			parts = append(parts, fmt.Sprintf("return []Element{%s}", body))
			last = sib
			hadElse = true
			goto done
		default:
			goto done
		}
	}
done:
	if !hadElse {
		parts = append(parts, "return []Element{}")
	}
	inner := strings.Join(parts, "\n")
	return fmt.Sprintf("R(func() []Element {\n%s\n}(),)", inner), last, nil
}

func getConditionAttr(n *html.Node) (string, error) {
	for _, a := range n.Attr {
		if a.Key == "condition" {
			return a.Val, nil
		}
	}
	return "", fmt.Errorf("'if' or 'elseif' element requires condition attribute")
}

func renderChildren(n *html.Node, comps map[string]CompInfo) (string, error) {
	var parts []string
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		b, err := render(c, comps)
		if err != nil {
			return "", err
		}
		if len(b) == 0 {
			continue
		}
		parts = append(parts, b)
	}
	return strings.Join(parts, ","), nil
}

// collectChildNodes returns direct children of n, excluding elseif/else (those are handled by processIfChain).
func collectChildNodes(n *html.Node) []*html.Node {
	var out []*html.Node
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if c.Type == html.ElementNode && (c.Data == "elseif" || c.Data == "else") {
			continue
		}
		out = append(out, c)
	}
	return out
}

// generateProps returns (code, complete, err). When complete is true, code is the full component call
// (including slot props and default children). When false, caller must append children and ")".
func generateProps(n *html.Node, comps map[string]CompInfo, children []*html.Node) (string, bool, error) {
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

		return buffer.String(), false, nil
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

	// Slot content: partition children into <slot name="..."> and rest
	if len(children) > 0 {
		comp, ok := comps[strings.TrimSpace(n.Data)]
		if ok {
			slotRendered := make(map[string]string) // slot name -> R(...) code
			var defaultRendered []string
			for _, c := range children {
				if c.Type == html.ElementNode && c.Data == "slot" {
					name := getAttr(c, "name")
					if name == "" {
						continue
					}
					propName := "Slot" + utils.Capitalize(name)
					if _, exists := comp.Props[strings.ToLower("slot"+utils.Capitalize(name))]; !exists {
						continue
					}
					var slotContent []string
					for ch := c.FirstChild; ch != nil; ch = ch.NextSibling {
						s, err := render(ch, comps)
						if err != nil {
							return "", false, err
						}
						if len(s) > 0 {
							slotContent = append(slotContent, s)
						}
					}
					if len(slotContent) > 0 {
						slotRendered[propName] = strings.Join(slotContent, ",")
					} else {
						slotRendered[propName] = ""
					}
				} else {
					s, err := render(c, comps)
					if err != nil {
						return "", false, err
					}
					if len(s) > 0 {
						defaultRendered = append(defaultRendered, s)
					}
				}
			}
			// Append slot props in deterministic order (by prop name)
			slotPropNames := make([]string, 0, len(slotRendered))
			for k := range slotRendered {
				slotPropNames = append(slotPropNames, k)
			}
			sort.Strings(slotPropNames)
			for _, propName := range slotPropNames {
				val := slotRendered[propName]
				if val == "" {
					props.WriteString(fmt.Sprintf("%s:R(),", propName))
				} else {
					props.WriteString(fmt.Sprintf("%s:R(%s),", propName, val))
				}
			}
			buffer.WriteString(fmt.Sprintf("%sComp(", comps[strings.TrimSpace(n.Data)].Name))
			buffer.WriteString(fmt.Sprintf("%s{%s},", comps[strings.TrimSpace(n.Data)].Name, props.String()))
			buffer.WriteString(fmt.Sprintf("Attrs{%s},", attrs.String()))
			buffer.WriteString(strings.Join(defaultRendered, ","))
			buffer.WriteString(")")
			return buffer.String(), true, nil
		}
	}

	buffer.WriteString(fmt.Sprintf("%sComp(", comps[strings.TrimSpace(n.Data)].Name))
	buffer.WriteString(fmt.Sprintf("%s{%s},", comps[strings.TrimSpace(n.Data)].Name, props.String()))
	buffer.WriteString(fmt.Sprintf("Attrs{%s},", attrs.String()))

	return buffer.String(), false, nil
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

// matchOneBraced matches a single {expr} (no nested braces) so multiple {a} {b} in one string work.
var reOneBraced = regexp.MustCompile(`\{([^{}]*)\}`)

func processRaws(input string) string {
	matches := reOneBraced.FindAllStringSubmatchIndex(input, -1)
	if len(matches) == 0 {
		return fmt.Sprintf("`%s`", input)
	}

	var tokens []string
	lastEnd := 0
	for _, m := range matches {
		// m[0], m[1] = full match; m[2], m[3] = capture group (content inside {})
		if lit := input[lastEnd:m[0]]; lit != "" {
			tokens = append(tokens, fmt.Sprintf("`%s`", lit))
		}
		val := input[m[2]:m[3]]
		if val != "" {
			if strings.HasPrefix(val, "$") {
				f := strings.Split(val, ".")
				tokens = append(tokens, fmt.Sprintf("%s[\"%s\"]", strings.Replace(f[0], "$", "", 1), f[1]))
			} else {
				if strings.HasPrefix(val, "props.") && len(val) > 6 {
					val = val[:6] + strings.ToUpper(val[6:7]) + val[7:]
				}
				tokens = append(tokens, val)
			}
		}
		lastEnd = m[1]
	}
	if lastEnd < len(input) {
		if lit := input[lastEnd:]; lit != "" {
			tokens = append(tokens, fmt.Sprintf("`%s`", lit))
		}
	}

	if len(tokens) > 1 {
		return fmt.Sprintf("R(%s)", strings.Join(tokens, ","))
	}
	if len(tokens) == 1 {
		return tokens[0]
	}
	return fmt.Sprintf("`%s`", input)
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
		} else if n.Data == "if" {
			s, _, err := processIfChain(n, comps)
			if err != nil {
				return "", err
			}
			buffer.WriteString(s)
		} else if n.Data == "elseif" || n.Data == "else" {
			// Consumed by a preceding <if>; skip (processIfChain already emitted code)
			return "", nil
		} else if n.Data == "slot" {
			s, err := processSlot(n)
			if err != nil {
				return "", err
			}
			buffer.WriteString(s)
		} else {

			childNodes := collectChildNodes(n)
			s, complete, err := generateProps(n, comps, childNodes)
			if err != nil {
				return "", err
			}

			buffer.WriteString(s)

			if complete {
				// Custom component with slots/default children already included
			} else if n.Data == "script" || n.Data == "style" {
				buffer.WriteString("R(`")
				for c := n.FirstChild; c != nil; c = c.NextSibling {
					buffer.WriteString(c.Data)
				}
				buffer.WriteString("`))")
			} else {
				childs := []string{}
				for _, c := range childNodes {
					b, err := render(c, comps)
					if err != nil {
						return "", err
					}
					if len(b) > 0 {
						childs = append(childs, b)
					}
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

// processSlot returns Go code that renders the slot content: R(props.SlotName).
// The node must be <slot name="..."/> or <slot name="...">; name is required.
func processSlot(n *html.Node) (string, error) {
	name := getAttr(n, "name")
	if name == "" {
		return "", fmt.Errorf("<slot> requires a name attribute")
	}
	propName := "Slot" + utils.Capitalize(name)
	return fmt.Sprintf("R(props.%s)", propName), nil
}

func getAttr(n *html.Node, key string) string {
	for _, a := range n.Attr {
		if a.Key == key {
			return a.Val
		}
	}
	return ""
}

// SlotNamesFromHTML parses htmlContent and returns unique slot names from <slot name="..."> elements.
// Used by the transpiler to add slot fields to component structs.
func SlotNamesFromHTML(htmlContent []byte) ([]string, error) {
	ctx := &html.Node{Type: html.ElementNode, Data: "div", DataAtom: atom.Div}
	nodes, err := html.ParseFragment(bytes.NewReader(bytes.TrimSpace(htmlContent)), ctx)
	if err != nil {
		return nil, err
	}
	seen := make(map[string]bool)
	var names []string
	var walk func(*html.Node)
	walk = func(n *html.Node) {
		if n == nil {
			return
		}
		if n.Type == html.ElementNode && n.Data == "slot" {
			if name := getAttr(n, "name"); name != "" && !seen[name] {
				seen[name] = true
				names = append(names, name)
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			walk(c)
		}
	}
	for _, n := range nodes {
		walk(n)
	}
	return names, nil
}
