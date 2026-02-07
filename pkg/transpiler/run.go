package transpiler

import (
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

// RunOptions configures Run. Nil means default: one file per component, package "gohtmlxc".
type RunOptions struct {
	// SingleFile emits one comp_generated.go instead of one file per component (legacy behavior).
	SingleFile bool
	// Pkg is the generated package name (default "gohtmlxc").
	Pkg string
}

func defaultOptions(opts *RunOptions) RunOptions {
	if opts == nil {
		return RunOptions{Pkg: "gohtmlxc"}
	}
	if opts.Pkg == "" {
		opts.Pkg = "gohtmlxc"
	}
	return *opts
}

// importsUsedInComponent returns only imports whose package alias appears in structStr or codeStr (e.g. "t.").
// Avoids "imported and not used" in per-component files.
func importsUsedInComponent(imports []string, structStr string, codeStr string) []string {
	var out []string
	for _, imp := range imports {
		alias := importAlias(imp)
		if alias == "" {
			out = append(out, imp)
			continue
		}
		needle := alias + "."
		if strings.Contains(structStr, needle) || strings.Contains(codeStr, needle) {
			out = append(out, imp)
		}
	}
	return out
}

func importAlias(imp string) string {
	imp = strings.TrimSpace(imp)
	// Format: alias "path" or "path"
	parts := strings.Fields(imp)
	if len(parts) >= 2 && !strings.HasPrefix(parts[0], `"`) {
		return parts[0]
	}
	if len(parts) >= 1 && strings.HasPrefix(parts[0], `"`) {
		path := strings.Trim(parts[0], `"`)
		if i := strings.LastIndex(path, "/"); i >= 0 {
			return path[i+1:]
		}
		return path
	}
	return ""
}

// componentFileName returns a safe filename for the component (e.g. "SampleTable" -> "SampleTable.go").
func componentFileName(name string) string {
	var b strings.Builder
	for _, r := range name {
		if (r >= 'A' && r <= 'Z') || (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == '_' {
			b.WriteRune(r)
		} else if r == ' ' || r == '-' {
			b.WriteRune('_')
		}
	}
	s := b.String()
	if s == "" {
		s = "component"
	}
	return s + ".go"
}

// Run transpiles HTML components from src to Go code in dist.
// It uses utils.Log for progress when set (default is no-op).
// If opts is nil, one file per component is emitted with package "gohtmlxc".
func Run(src, dist string, opts *RunOptions) error {
	opt := defaultOptions(opts)
	if utils.Log != nil {
		utils.Log.Info("transpiling...")
	}
	t := time.Now()
	defer func(t time.Time) {
		if utils.Log != nil {
			utils.Log.Info("transpiled", "duration", time.Since(t))
		}
	}(t)

	files, err := utils.WalkAndReadHTMLFiles(src)
	if err != nil {
		return err
	}

	imports := []string{}
	sections := make(map[string]string)
	componentSource := make(map[string]string)
	componentFileContent := make(map[string][]byte)

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
	structMap := make(map[string]string)

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

			s := gocode.ConstructStruct(propsMap, name)
			structs = append(structs, s)
			structMap[name] = s
		} else {
			s := gocode.ConstructStruct(map[string]string{}, name)
			structs = append(structs, s)
			structMap[name] = s
		}
	}

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

	outDir := path.Join(dist, opt.Pkg)
	if err := os.MkdirAll(outDir, 0755); err != nil {
		return err
	}
	// Remove existing .go files so we never mix single-file and multi-file output
	matches, _ := filepath.Glob(filepath.Join(outDir, "*.go"))
	for _, m := range matches {
		_ = os.Remove(m)
	}

	if opt.SingleFile {
		b, err := gocode.ConstructSourceWithPkg(goCodes, structs, imports, opt.Pkg)
		if err != nil {
			return err
		}
		outPath := path.Join(outDir, "comp_generated.go")
		if err := os.WriteFile(outPath, []byte(b), 0644); err != nil {
			return err
		}
		return nil
	}

	// One file per component; each file gets package + only the imports it uses (so no "imported and not used")
	for _, name := range sectionNames {
		codeStr, ok := goCodes[name]
		if !ok {
			continue
		}
		structStr := structMap[name]
		usedImports := importsUsedInComponent(imports, structStr, codeStr)
		compContent, err := gocode.ConstructComponentFile(opt.Pkg, usedImports, name, structStr, codeStr)
		if err != nil {
			return err
		}
		filename := componentFileName(name)
		if err := os.WriteFile(path.Join(outDir, filename), []byte(compContent), 0644); err != nil {
			return err
		}
	}

	return nil
}
