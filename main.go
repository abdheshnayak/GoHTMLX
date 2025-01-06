package main

import (
	"flag"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"
	"text/template"
	htmltemplate "text/template"

	"github.com/abdheshnayak/gohtmlx/pkg/element"
	"github.com/abdheshnayak/gohtmlx/pkg/gocode"
	"github.com/abdheshnayak/gohtmlx/pkg/utils"
	"sigs.k8s.io/yaml"
)

func main() {
	src := flag.String("src", "", "source directory")
	dist := flag.String("dist", "", "destination directory")
	flag.Parse()

	if src == nil || dist == nil {
		flag.PrintDefaults()
		return
	}

	if *src == "" || *dist == "" {
		flag.PrintDefaults()
		return
	}

	if err := Run(*src, *dist); err != nil {
		fmt.Println("Error:", err)
	}
}

func Run(src, dist string) error {
	fmt.Println("transpiling...")
	defer fmt.Println("transpiled")

	input, err := utils.WalkAndConcatenateHTML(src)
	if err != nil {
		return err
	}

	parser := template.New("sections").Delims("<!-- {{-", "}} -->")

	// Parse the template
	tmpl, err := parser.Parse(string(input))
	if err != nil {
		return err
	}

	sections, err := utils.ParseSections(tmpl)
	if err != nil {
		return err
	}

	goCodes := map[string][]byte{}

	components := map[string]element.CompInfo{}
	for k := range sections {
		components[strings.ToLower(k)] = element.CompInfo{
			Name:  k,
			Props: map[string]string{},
		}
	}

	structs := [][]byte{}

	for name, content := range sections {
		hparser := htmltemplate.New("sections").Delims("<!-- {{+", "}} -->")
		tpl, err := hparser.Parse(string(content))
		m, err := utils.ParseSections(tpl)
		if err != nil {
			return err
		}

		var propsMap map[string]string
		if props, ok := m["props"]; ok {
			if err := yaml.Unmarshal([]byte(props), &propsMap); err != nil {
				return err
			}

			if _, ok := components[name]; !ok {
				for k := range propsMap {
					components[strings.ToLower(name)].Props[strings.ToLower(k)] = k
				}
			}

			structs = append(structs, gocode.ConstructStruct(propsMap, name))
		} else {
			structs = append(structs, gocode.ConstructStruct(map[string]string{}, name))
		}
	}

	// Output the parsed map
	for name, content := range sections {
		hparser := htmltemplate.New("sections").Delims("<!-- {{+", "}} -->")
		tpl, err := hparser.Parse(string(content))
		m, err := utils.ParseSections(tpl)
		if err != nil {
			return err
		}

		if html, ok := m["html"]; ok {
			h, err := element.NewHtml([]byte(html))
			if err != nil {
				return err
			}

			out, err := h.RenderGolangCode(components)
			if err != nil {
				return err
			}

			goCodes[name] = out
		}
	}

	b, err := gocode.ConstructSource(goCodes, structs, nil)
	if err != nil {
		return err
	}

	outPath := path.Join(dist, "gohtmlxc", "comp_generated.go")
	if err := os.MkdirAll(filepath.Dir(outPath), 0755); err != nil {
		return err
	}
	if err := os.WriteFile(outPath, b, 0644); err != nil {
		return err
	}

	return nil
}
