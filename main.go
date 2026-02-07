package main

import (
	"flag"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"
	"text/template"
	htmltemplate "text/template"
	"time"

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
		utils.Log.Error("transpiling failed", "err", err)
	}
}

func Run(src, dist string) error {
	utils.Log.Info("transpiling...")
	t := time.Now()
	defer func(t time.Time) {
		utils.Log.Info(fmt.Sprintf("transpiled in %s", time.Since(t)))
	}(t)

	input, err := utils.WalkAndConcatenateHTML(src)
	if err != nil {
		return err
	}

	tmpl, err := template.New("global").Delims("<!-- *", " -->").Parse(string(input))
	if err != nil {
		return err
	}

	gsections, err := utils.ParseSections(tmpl)
	if err != nil {
		return err
	}
	imports := []string{}
	if imps, ok := gsections["imports"]; ok {
		for _, v := range strings.Split(strings.TrimSpace(imps), "\n") {
			s := strings.TrimSpace(v)
			if s != "" {
				imports = append(imports, s)
			}
		}
		sort.Strings(imports)
	}

	// Parse the template
	tmpl, err = template.New("sections").Delims("<!-- +", " -->").Parse(string(input))
	if err != nil {
		return err
	}

	sections, err := utils.ParseSections(tmpl)
	if err != nil {
		return err
	}

	goCodes := map[string]string{}

	components := map[string]element.CompInfo{}
	sectionNames := make([]string, 0, len(sections))
	for k := range sections {
		sectionNames = append(sectionNames, k)
		components[strings.ToLower(k)] = element.CompInfo{
			Name:  k,
			Props: map[string]string{},
		}
	}
	sort.Strings(sectionNames)

	structs := []string{}

	for _, name := range sectionNames {
		content := sections[name]
		hparser := htmltemplate.New("sections").Delims("<!-- |", " -->")
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

	// Output the parsed map (same order for deterministic output)
	for _, name := range sectionNames {
		content := sections[name]
		hparser := htmltemplate.New("section-data").Delims("<!-- |", "-->")
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

	b, err := gocode.ConstructSource(goCodes, structs, imports)
	if err != nil {
		return err
	}

	outPath := path.Join(dist, "gohtmlxc", "comp_generated.go")
	if err := os.MkdirAll(filepath.Dir(outPath), 0755); err != nil {
		return err
	}
	if err := os.WriteFile(outPath, []byte(b), 0644); err != nil {
		return err
	}

	return nil
}
