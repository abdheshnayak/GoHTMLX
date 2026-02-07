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
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of gohtmlx:\n")
		fmt.Fprintf(os.Stderr, "  gohtmlx --src=DIR --dist=DIR\n")
		fmt.Fprintf(os.Stderr, "\n")
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\nExit codes:\n")
		fmt.Fprintf(os.Stderr, "  0  success\n")
		fmt.Fprintf(os.Stderr, "  1  transpilation failed (parse, codegen, or write error)\n")
		fmt.Fprintf(os.Stderr, "  2  invalid arguments or missing flags\n")
	}

	src := flag.String("src", "", "source directory containing .html components")
	dist := flag.String("dist", "", "destination directory for generated Go code")
	flag.Parse()

	if *src == "" || *dist == "" {
		flag.Usage()
		os.Exit(2)
	}

	if err := Run(*src, *dist); err != nil {
		utils.Log.Error("transpiling failed", "err", err)
		os.Exit(1)
	}
}

func Run(src, dist string) error {
	utils.Log.Info("transpiling...")
	t := time.Now()
	defer func(t time.Time) {
		utils.Log.Info(fmt.Sprintf("transpiled in %s", time.Since(t)))
	}(t)

	files, err := utils.WalkAndReadHTMLFiles(src)
	if err != nil {
		return err
	}

	// Merge imports and sections from all files; track source file per component
	imports := []string{}
	sections := make(map[string]string)
	componentSource := make(map[string]string)   // component name -> file path
	componentFileContent := make(map[string][]byte) // component name -> full file content (for line/snippet)

	for _, f := range files {
		tmplGlobal, err := template.New("global").Delims("<!-- *", " -->").Parse(string(f.Content))
		if err != nil {
			return &TranspileError{FilePath: f.Path, Line: 0, Message: err.Error()}
		}
		gsections, err := utils.ParseSections(tmplGlobal)
		if err != nil {
			return &TranspileError{FilePath: f.Path, Line: 0, Message: err.Error()}
		}
		if imps, ok := gsections["imports"]; ok {
			for _, v := range strings.Split(strings.TrimSpace(imps), "\n") {
				s := strings.TrimSpace(v)
				if s != "" {
					imports = append(imports, s)
				}
			}
		}

		tmplSections, err := template.New("sections").Delims("<!-- +", " -->").Parse(string(f.Content))
		if err != nil {
			return &TranspileError{FilePath: f.Path, Line: 0, Message: err.Error()}
		}
		fileSections, err := utils.ParseSections(tmplSections)
		if err != nil {
			return &TranspileError{FilePath: f.Path, Line: 0, Message: err.Error()}
		}
		for name, content := range fileSections {
			if _, ok := sections[name]; ok {
				other := componentSource[name]
				return &TranspileError{
					FilePath: f.Path,
					Message:  fmt.Sprintf("component %q already defined in %s", name, other),
				}
			}
			sections[name] = content
			componentSource[name] = f.Path
			componentFileContent[name] = f.Content
		}
	}
	sort.Strings(imports)

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
		if err != nil {
			return wrapTranspileErr(name, componentSource[name], componentFileContent[name], err)
		}
		m, err := utils.ParseSections(tpl)
		if err != nil {
			return wrapTranspileErr(name, componentSource[name], componentFileContent[name], err)
		}

		var propsMap map[string]string
		if props, ok := m["props"]; ok {
			if err := yaml.Unmarshal([]byte(props), &propsMap); err != nil {
				return wrapTranspileErr(name, componentSource[name], componentFileContent[name], err)
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
		filePath := componentSource[name]
		fileContent := componentFileContent[name]
		hparser := htmltemplate.New("section-data").Delims("<!-- |", "-->")
		tpl, err := hparser.Parse(string(content))
		if err != nil {
			return wrapTranspileErr(name, filePath, fileContent, err)
		}
		m, err := utils.ParseSections(tpl)
		if err != nil {
			return wrapTranspileErr(name, filePath, fileContent, err)
		}

		if html, ok := m["html"]; ok {
			h, err := element.NewHtml([]byte(html))
			if err != nil {
				return wrapTranspileErr(name, filePath, fileContent, err)
			}

			out, err := h.RenderGolangCode(components)
			if err != nil {
				return wrapTranspileErr(name, filePath, fileContent, err)
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
