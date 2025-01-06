package gocode

import (
	"bytes"
	"fmt"
	"go/format"
	"log/slog"
	"strings"

	"github.com/abdheshnayak/gohtmlx/pkg/utils"
)

func ConstructStruct(props map[string]string, name string) []byte {
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
		slog.Error("Error formatting source", slog.String("source", buffer.String()), slog.String("error", err.Error()))

		return []byte(buffer.String())
	}
	return b
}

func ConstructSource(codes map[string][]byte, structs [][]byte, imports []string) ([]byte, error) {
	var builder strings.Builder

	builder.WriteString("package gohtmlxc\n\n")
	builder.WriteString("import (\n")

	builder.WriteString(". \"github.com/abdheshnayak/gohtmlx/pkg/element\"\n")

	for _, v := range imports {
		builder.WriteString(fmt.Sprintf("%s\n", strings.TrimSpace(v)))
	}

	builder.WriteString(")\n\n")

	structsByte := bytes.Join(structs, []byte("\n\n"))
	builder.WriteString(string(structsByte))

	for k, v := range codes {
		builder.WriteString(fmt.Sprintf("func %s(", k))
		builder.WriteString(fmt.Sprintf("props %sProps, children ...Element", k))
		builder.WriteString(") Element {\n")

		builder.WriteString(fmt.Sprintf("return %s\n", string(v)))
		builder.WriteString("\n}\n\n")
	}

	// b, err := format.Source([]byte(builder.String()))
	// if err != nil {
	// 	fmt.Println(builder.String())
	// 	return nil, err
	// }

	return []byte(builder.String()), nil
}
