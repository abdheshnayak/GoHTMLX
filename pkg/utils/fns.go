package utils

import (
	"strings"
	"text/template"
)

func Capitalize(s string) string {
	if len(s) == 0 {
		return s
	}

	return strings.ToUpper(s[:1]) + s[1:]
}

func GetSections(tmpl *template.Template) []string {
	sectionNames := []string{}
	t := tmpl.Templates()
	for _, v := range t {
		sectionNames = append(sectionNames, v.Name())
	}

	return sectionNames
}

func ParseSections(tmpl *template.Template) (map[string]string, error) {
	sectionNames := GetSections(tmpl)
	t := tmpl.Templates()
	for _, v := range t {
		sectionNames = append(sectionNames, v.Name())
	}

	sections := make(map[string]string)

	for _, section := range sectionNames {
		var buffer strings.Builder
		err := tmpl.ExecuteTemplate(&buffer, section, nil)
		if err != nil {
			return nil, err
		}
		sections[section] = buffer.String()
	}

	return sections, nil
}
