package gocode

import (
	"fmt"
	"go/format"
	"sort"
	"strings"

	"github.com/abdheshnayak/gohtmlx/pkg/utils"
)

func ConstructStruct(props map[string]string, name string) string {
	var buffer strings.Builder
	buffer.WriteString("type ")
	buffer.WriteString(fmt.Sprintf("%s", name))
	buffer.WriteString(" struct {\n")

	keys := make([]string, 0, len(props))
	for k := range props {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		v := props[k]
		buffer.WriteString(utils.Capitalize(k))
		buffer.WriteString(" ")
		buffer.WriteString(v)
		buffer.WriteString("\n")
	}

	buffer.WriteString("Attrs Attrs\n")

	buffer.WriteString("}\n")

	b, err := format.Source([]byte(buffer.String()))
	if err != nil {
		utils.Log.Error("failed to format source", buffer.String(), err)

		return buffer.String()
	}
	return string(b)
}

func ConstructSource(codes map[string]string, structs []string, imports []string) (string, error) {
	var builder strings.Builder

	builder.WriteString("package gohtmlxc\n\n")
	builder.WriteString("import (\n")

	builder.WriteString(". \"github.com/abdheshnayak/gohtmlx/pkg/element\"\n")

	importList := make([]string, len(imports))
	copy(importList, imports)
	sort.Strings(importList)
	for _, v := range importList {
		s := strings.TrimSpace(v)
		if s != "" {
			builder.WriteString(fmt.Sprintf("%s\n", s))
		}
	}

	builder.WriteString(")\n\n")

	structsByte := strings.Join(structs, "\n\n")
	builder.WriteString(string(structsByte))

	codeKeys := make([]string, 0, len(codes))
	for k := range codes {
		codeKeys = append(codeKeys, k)
	}
	sort.Strings(codeKeys)
	for _, k := range codeKeys {
		v := codes[k]
		builder.WriteString(fmt.Sprintf("func %sComp(", k))
		builder.WriteString(fmt.Sprintf("props %s, attrs Attrs, children ...Element", k))
		builder.WriteString(") Element {\n")

		builder.WriteString(`
        props.Attrs = attrs
        if props.Attrs == nil{
            props.Attrs = Attrs{}
        }
    `)

		builder.WriteString(fmt.Sprintf("\nreturn %s\n", v))
		builder.WriteString("\n}\n\n")

		builder.WriteString(fmt.Sprintf("func (c %s) Get(children ...Element) Element {\n", k))

		builder.WriteString(fmt.Sprintf("return %sComp(c, c.Attrs, children...)\n", k))
		builder.WriteString("}\n\n")
	}

	b, err := format.Source([]byte(builder.String()))
	if err != nil {
		return "", err
	}

	return string(b), nil
}
