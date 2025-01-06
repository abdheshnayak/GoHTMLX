package gocode

import (
	"fmt"
	"go/format"
	"strings"

	"github.com/abdheshnayak/gohtmlx/pkg/utils"
)

func ConstructStruct(props map[string]string, name string) string {
	var buffer strings.Builder
	buffer.WriteString("type ")
	buffer.WriteString(fmt.Sprintf("%sProps", name))
	buffer.WriteString(" struct {\n")

	for k, v := range props {
		buffer.WriteString(utils.Capitalize(k))
		buffer.WriteString(" ")
		buffer.WriteString(v)
		buffer.WriteString("\n")
	}

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

	for _, v := range imports {
		builder.WriteString(fmt.Sprintf("%s\n", strings.TrimSpace(v)))
	}

	builder.WriteString(")\n\n")

	structsByte := strings.Join(structs, "\n\n")
	builder.WriteString(string(structsByte))

	for k, v := range codes {
		builder.WriteString(fmt.Sprintf("func %s(", k))
		builder.WriteString(fmt.Sprintf("props %sProps, children ...Element", k))
		builder.WriteString(") Element {\n")

		builder.WriteString(fmt.Sprintf("return %s\n", v))
		builder.WriteString("\n}\n\n")
	}

	b, err := format.Source([]byte(builder.String()))
	if err != nil {
		return "", err
	}

	return string(b), nil
}
