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

// ConstructSharedFile returns the package and import block for the generated package.
// Used when emitting one file per component; write this to imports.go (or similar).
func ConstructSharedFile(pkg string, imports []string) (string, error) {
	var builder strings.Builder
	builder.WriteString("package " + pkg + "\n\n")
	builder.WriteString("import (\n")
	builder.WriteString("\t. \"github.com/abdheshnayak/gohtmlx/pkg/element\"\n")
	importList := make([]string, len(imports))
	copy(importList, imports)
	sort.Strings(importList)
	for _, v := range importList {
		s := strings.TrimSpace(v)
		if s != "" {
			builder.WriteString("\t" + s + "\n")
		}
	}
	builder.WriteString(")\n")
	b, err := format.Source([]byte(builder.String()))
	if err != nil {
		return "", err
	}
	return string(b), nil
}

// ConstructComponentFile returns the Go code for a single component (package, imports, type, Comp, Get).
// Each file needs its own import block so Attrs, Element, and user types (e.g. t) are in scope.
func ConstructComponentFile(pkg string, imports []string, name string, structStr string, codeStr string) (string, error) {
	var builder strings.Builder
	builder.WriteString("package " + pkg + "\n\n")
	builder.WriteString("import (\n")
	builder.WriteString("\t. \"github.com/abdheshnayak/gohtmlx/pkg/element\"\n")
	importList := make([]string, len(imports))
	copy(importList, imports)
	sort.Strings(importList)
	for _, v := range importList {
		s := strings.TrimSpace(v)
		if s != "" {
			builder.WriteString("\t" + s + "\n")
		}
	}
	builder.WriteString(")\n\n")
	builder.WriteString(structStr)
	builder.WriteString("\n")
	builder.WriteString(fmt.Sprintf("func %sComp(", name))
	builder.WriteString(fmt.Sprintf("props %s, attrs Attrs, children ...Element", name))
	builder.WriteString(") Element {\n")
	builder.WriteString("\tprops.Attrs = attrs\n")
	builder.WriteString("\tif props.Attrs == nil {\n")
	builder.WriteString("\t\tprops.Attrs = Attrs{}\n")
	builder.WriteString("\t}\n")
	builder.WriteString(fmt.Sprintf("\treturn %s\n", codeStr))
	builder.WriteString("}\n\n")
	builder.WriteString(fmt.Sprintf("func (c %s) Get(children ...Element) Element {\n", name))
	builder.WriteString(fmt.Sprintf("\treturn %sComp(c, c.Attrs, children...)\n", name))
	builder.WriteString("}\n")
	b, err := format.Source([]byte(builder.String()))
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func ConstructSource(codes map[string]string, structs []string, imports []string) (string, error) {
	return ConstructSourceWithPkg(codes, structs, imports, "gohtmlxc")
}

// ConstructSourceWithPkg generates a single-file output with the given package name.
func ConstructSourceWithPkg(codes map[string]string, structs []string, imports []string, pkg string) (string, error) {
	var builder strings.Builder

	builder.WriteString("package " + pkg + "\n\n")
	builder.WriteString("import (\n")

	builder.WriteString("\t. \"github.com/abdheshnayak/gohtmlx/pkg/element\"\n")

	importList := make([]string, len(imports))
	copy(importList, imports)
	sort.Strings(importList)
	for _, v := range importList {
		s := strings.TrimSpace(v)
		if s != "" {
			builder.WriteString(fmt.Sprintf("\t%s\n", s))
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
