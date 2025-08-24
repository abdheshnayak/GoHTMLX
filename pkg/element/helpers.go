package element

import (
	"fmt"
	"io"
)

// HTML5 semantic elements
func Article(attrs Attrs, children ...Element) Element {
	return E("article", attrs, children...)
}

func Aside(attrs Attrs, children ...Element) Element {
	return E("aside", attrs, children...)
}

func Footer(attrs Attrs, children ...Element) Element {
	return E("footer", attrs, children...)
}

func Header(attrs Attrs, children ...Element) Element {
	return E("header", attrs, children...)
}

func Main(attrs Attrs, children ...Element) Element {
	return E("main", attrs, children...)
}

func Nav(attrs Attrs, children ...Element) Element {
	return E("nav", attrs, children...)
}

func Section(attrs Attrs, children ...Element) Element {
	return E("section", attrs, children...)
}

// Common HTML elements
func Div(attrs Attrs, children ...Element) Element {
	return E("div", attrs, children...)
}

func Span(attrs Attrs, children ...Element) Element {
	return E("span", attrs, children...)
}

func P(attrs Attrs, children ...Element) Element {
	return E("p", attrs, children...)
}

func A(href string, attrs Attrs, children ...Element) Element {
	if attrs == nil {
		attrs = make(Attrs)
	}
	attrs["href"] = href
	return E("a", attrs, children...)
}

func Img(src, alt string, attrs Attrs) Element {
	if attrs == nil {
		attrs = make(Attrs)
	}
	attrs["src"] = src
	attrs["alt"] = alt
	return E("img", attrs)
}

// Heading elements
func H1(attrs Attrs, children ...Element) Element {
	return E("h1", attrs, children...)
}

func H2(attrs Attrs, children ...Element) Element {
	return E("h2", attrs, children...)
}

func H3(attrs Attrs, children ...Element) Element {
	return E("h3", attrs, children...)
}

func H4(attrs Attrs, children ...Element) Element {
	return E("h4", attrs, children...)
}

func H5(attrs Attrs, children ...Element) Element {
	return E("h5", attrs, children...)
}

func H6(attrs Attrs, children ...Element) Element {
	return E("h6", attrs, children...)
}

// List elements
func Ul(attrs Attrs, children ...Element) Element {
	return E("ul", attrs, children...)
}

func Ol(attrs Attrs, children ...Element) Element {
	return E("ol", attrs, children...)
}

func Li(attrs Attrs, children ...Element) Element {
	return E("li", attrs, children...)
}

// Table elements
func Table(attrs Attrs, children ...Element) Element {
	return E("table", attrs, children...)
}

func Thead(attrs Attrs, children ...Element) Element {
	return E("thead", attrs, children...)
}

func Tbody(attrs Attrs, children ...Element) Element {
	return E("tbody", attrs, children...)
}

func Tfoot(attrs Attrs, children ...Element) Element {
	return E("tfoot", attrs, children...)
}

func Tr(attrs Attrs, children ...Element) Element {
	return E("tr", attrs, children...)
}

func Th(attrs Attrs, children ...Element) Element {
	return E("th", attrs, children...)
}

func Td(attrs Attrs, children ...Element) Element {
	return E("td", attrs, children...)
}

// Form elements
func Form(attrs Attrs, children ...Element) Element {
	return E("form", attrs, children...)
}

func Input(inputType string, attrs Attrs) Element {
	if attrs == nil {
		attrs = make(Attrs)
	}
	attrs["type"] = inputType
	return E("input", attrs)
}

func TextInput(attrs Attrs) Element {
	return Input("text", attrs)
}

func EmailInput(attrs Attrs) Element {
	return Input("email", attrs)
}

func PasswordInput(attrs Attrs) Element {
	return Input("password", attrs)
}

func SubmitButton(value string, attrs Attrs) Element {
	if attrs == nil {
		attrs = make(Attrs)
	}
	attrs["type"] = "submit"
	attrs["value"] = value
	return E("input", attrs)
}

func Button(attrs Attrs, children ...Element) Element {
	return E("button", attrs, children...)
}

func Label(forAttr string, attrs Attrs, children ...Element) Element {
	if attrs == nil {
		attrs = make(Attrs)
	}
	if forAttr != "" {
		attrs["for"] = forAttr
	}
	return E("label", attrs, children...)
}

func Select(attrs Attrs, children ...Element) Element {
	return E("select", attrs, children...)
}

func Option(value string, attrs Attrs, children ...Element) Element {
	if attrs == nil {
		attrs = make(Attrs)
	}
	attrs["value"] = value
	return E("option", attrs, children...)
}

func Textarea(attrs Attrs, children ...Element) Element {
	return E("textarea", attrs, children...)
}

// Text formatting
func Strong(attrs Attrs, children ...Element) Element {
	return E("strong", attrs, children...)
}

func Em(attrs Attrs, children ...Element) Element {
	return E("em", attrs, children...)
}

func Small(attrs Attrs, children ...Element) Element {
	return E("small", attrs, children...)
}

func Mark(attrs Attrs, children ...Element) Element {
	return E("mark", attrs, children...)
}

func Del(attrs Attrs, children ...Element) Element {
	return E("del", attrs, children...)
}

func Ins(attrs Attrs, children ...Element) Element {
	return E("ins", attrs, children...)
}

func Sub(attrs Attrs, children ...Element) Element {
	return E("sub", attrs, children...)
}

func Sup(attrs Attrs, children ...Element) Element {
	return E("sup", attrs, children...)
}

// Media elements
func Video(attrs Attrs, children ...Element) Element {
	return E("video", attrs, children...)
}

func Audio(attrs Attrs, children ...Element) Element {
	return E("audio", attrs, children...)
}

func Source(src, mediaType string, attrs Attrs) Element {
	if attrs == nil {
		attrs = make(Attrs)
	}
	attrs["src"] = src
	attrs["type"] = mediaType
	return E("source", attrs)
}

// Meta elements
func Meta(name, content string, attrs Attrs) Element {
	if attrs == nil {
		attrs = make(Attrs)
	}
	attrs["name"] = name
	attrs["content"] = content
	return E("meta", attrs)
}

func MetaCharset(charset string) Element {
	return E("meta", Attrs{"charset": charset})
}

func MetaViewport(content string) Element {
	return E("meta", Attrs{"name": "viewport", "content": content})
}

func Link(rel, href string, attrs Attrs) Element {
	if attrs == nil {
		attrs = make(Attrs)
	}
	attrs["rel"] = rel
	attrs["href"] = href
	return E("link", attrs)
}

func StylesheetLink(href string) Element {
	return Link("stylesheet", href, nil)
}

func Script(src string, attrs Attrs, children ...Element) Element {
	if attrs == nil {
		attrs = make(Attrs)
	}
	if src != "" {
		attrs["src"] = src
	}
	return E("script", attrs, children...)
}

func InlineScript(content string) Element {
	return Script("", nil, R(content))
}

func StyleElement(attrs Attrs, children ...Element) Element {
	return E("style", attrs, children...)
}

func InlineStyle(content string) Element {
	return StyleElement(nil, R(content))
}

// Document structure
func HtmlElement(lang string, attrs Attrs, children ...Element) Element {
	if attrs == nil {
		attrs = make(Attrs)
	}
	attrs["lang"] = lang
	return E("html", attrs, children...)
}

func Head(attrs Attrs, children ...Element) Element {
	return E("head", attrs, children...)
}

func Body(attrs Attrs, children ...Element) Element {
	return E("body", attrs, children...)
}

func Title(title string) Element {
	return E("title", nil, R(title))
}

// Utility functions
func Text(content string) Element {
	return R(content)
}

func RawHTML(content string) Element {
	return &rawElement{content: content}
}

// rawElement allows inserting raw HTML without escaping
type rawElement struct {
	content string
}

func (r *rawElement) Render(w io.Writer) (int, error) {
	return w.Write([]byte(r.content))
}

// Fragment groups multiple elements without a wrapper
func Fragment(children ...Element) Element {
	return &fragmentElement{children: children}
}

type fragmentElement struct {
	children []Element
}

func (f *fragmentElement) Render(w io.Writer) (int, error) {
	total := 0
	for _, child := range f.children {
		n, err := child.Render(w)
		if err != nil {
			return total, err
		}
		total += n
	}
	return total, nil
}

// Conditional rendering
func If(condition bool, element Element) Element {
	if condition {
		return element
	}
	return Fragment() // Empty fragment
}

func IfElse(condition bool, ifElement, elseElement Element) Element {
	if condition {
		return ifElement
	}
	return elseElement
}

// For loop helper
func For[T any](items []T, fn func(T, int) Element) Element {
	var children []Element
	for i, item := range items {
		children = append(children, fn(item, i))
	}
	return Fragment(children...)
}

// Map helper for transforming slices
func Map[T any](items []T, fn func(T) Element) Element {
	var children []Element
	for _, item := range items {
		children = append(children, fn(item))
	}
	return Fragment(children...)
}

// Join elements with a separator
func Join(separator Element, elements ...Element) Element {
	if len(elements) == 0 {
		return Fragment()
	}
	if len(elements) == 1 {
		return elements[0]
	}

	var result []Element
	for i, element := range elements {
		if i > 0 {
			result = append(result, separator)
		}
		result = append(result, element)
	}
	return Fragment(result...)
}

// Common attribute builders
func ID(id string) Attrs {
	return Attrs{"id": id}
}

func Class(classes ...string) Attrs {
	return Attrs{"class": ClassNames(classes)}
}

func StyleAttr(styles map[string]string) Attrs {
	var styleStr string
	for prop, value := range styles {
		if styleStr != "" {
			styleStr += "; "
		}
		styleStr += fmt.Sprintf("%s: %s", prop, value)
	}
	return Attrs{"style": styleStr}
}

func DataAttr(key, value string) Attrs {
	return Attrs{fmt.Sprintf("data-%s", key): value}
}

func AriaAttr(key, value string) Attrs {
	return Attrs{fmt.Sprintf("aria-%s", key): value}
}
