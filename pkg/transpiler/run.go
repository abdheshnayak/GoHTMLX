// Package transpiler implements the GoHTMLX pipeline: read .html component files,
// parse sections (define, props, html), discover slots and props, generate Go code
// via pkg/element and pkg/gocode, and write to --dist. It does not depend on
// Fiber or file watchers; the CLI (main) calls Run and handles exit codes.
package transpiler

import (
	"fmt"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"regexp"
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
	// ValidateTypes runs go build on the generated package after codegen and returns a TranspileError
	// with file/line when the build fails (e.g. invalid prop types). Run from module root so go build can resolve the package.
	ValidateTypes bool
	// Incremental skips transpilation when no .html under src is newer than the generated .go files under dist.
	// Useful in watch scripts to avoid work when nothing changed. Best-effort; a full run is always correct.
	Incremental bool
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

// importPath returns the quoted path from an import line (e.g. `t "pkg/path"` -> "pkg/path").
// Used to deduplicate imports by path.
func importPath(imp string) string {
	imp = strings.TrimSpace(imp)
	parts := strings.Fields(imp)
	for _, p := range parts {
		if strings.HasPrefix(p, `"`) && strings.HasSuffix(p, `"`) {
			return strings.Trim(p, `"`)
		}
	}
	return ""
}

// deduplicateImports keeps one import per path (same path in multiple files → single import).
// Prefers the form that has an explicit alias when present. Deterministic: sorted by path.
func deduplicateImports(imports []string) []string {
	byPath := make(map[string]string)
	for _, imp := range imports {
		path := importPath(imp)
		if path == "" {
			continue
		}
		existing, ok := byPath[path]
		// Prefer form with alias (two or more tokens before the quoted path)
		if !ok || (importAlias(imp) != "" && importAlias(existing) == "") {
			byPath[path] = imp
		}
	}
	out := make([]string, 0, len(byPath))
	for _, imp := range byPath {
		out = append(out, imp)
	}
	return out
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

// findModuleRoot walks up from dir until it finds a directory containing go.mod.
func findModuleRoot(dir string) (string, error) {
	dir, err := filepath.Abs(dir)
	if err != nil {
		return "", err
	}
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir, nil
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return "", fmt.Errorf("no go.mod found in %s or any parent directory", dir)
		}
		dir = parent
	}
}

// goBuildErrorRe matches "path/file.go:line:col: message" or "path/file.go:line: message" from go build stderr.
var goBuildErrorRe = regexp.MustCompile(`([^:]+\.go):(\d+)(?::(\d+))?:\s*(.+)`)

// validateGeneratedPackage runs go build on the generated package and returns a TranspileError on failure.
func validateGeneratedPackage(outDir, dist, pkg string, componentSource map[string]string, sectionNames []string, singleFile bool) error {
	wd, err := os.Getwd()
	if err != nil {
		return &TranspileError{Message: "validate-types: cannot get working directory: " + err.Error()}
	}
	root, err := findModuleRoot(wd)
	if err != nil {
		return &TranspileError{Message: "validate-types: " + err.Error() + " (run from module root)"}
	}
	outDirAbs, err := filepath.Abs(outDir)
	if err != nil {
		return &TranspileError{Message: "validate-types: cannot resolve output path: " + err.Error()}
	}
	rel, err := filepath.Rel(root, outDirAbs)
	if err != nil {
		return &TranspileError{Message: "validate-types: generated output is outside module: " + err.Error()}
	}
	pkgPath := "./" + filepath.ToSlash(rel)
	cmd := exec.Command("go", "build", "-o", os.DevNull, pkgPath)
	cmd.Dir = root
	out, err := cmd.CombinedOutput()
	if err == nil {
		return nil
	}
	// Parse first error line from stderr (go build prints file:line: message)
	basenameToComponent := make(map[string]string)
	if singleFile {
		basenameToComponent["comp_generated.go"] = ""
	} else {
		for _, name := range sectionNames {
			basenameToComponent[componentFileName(name)] = name
		}
	}
	lines := strings.Split(string(out), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		sub := goBuildErrorRe.FindStringSubmatch(line)
		if len(sub) < 4 {
			continue
		}
		filePart, lineStr, msg := sub[1], sub[2], sub[len(sub)-1]
		base := filepath.Base(filePart)
		component := basenameToComponent[base]
		var filePath string
		if component != "" && componentSource[component] != "" {
			filePath = componentSource[component]
		} else {
			filePath = filepath.Join(outDir, base)
		}
		var lineNum int
		_, _ = fmt.Sscanf(lineStr, "%d", &lineNum)
		return &TranspileError{
			Component: component,
			FilePath:  filePath,
			Line:      lineNum,
			Message:   "go build failed: " + msg,
			Snippet:   strings.TrimSpace(line),
		}
	}
	return &TranspileError{Message: "go build failed", Snippet: string(out)}
}

// Run transpiles HTML components from src to Go code in dist.
// It walks src for .html files, parses component sections, merges imports, discovers
// slots from HTML, generates structs and component code, and writes to dist/<pkg>/*.go
// (or a single comp_generated.go when opts.SingleFile is true). Uses utils.Log for
// progress when set. If opts is nil, one file per component is emitted with package "gohtmlxc".
//
// On failure, Run returns an error that is always a *TranspileError (use errors.As to extract it).
// TranspileError provides FilePath, Line, Message, Snippet, and Component so callers can show
// file:line and context in UIs or logs.
// newestModTime returns the newest modification time of files under dir with the given extension (e.g. ".html").
// Returns zero time if no matching files exist.
func newestModTime(dir, ext string) (time.Time, error) {
	var newest time.Time
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		if strings.HasSuffix(strings.ToLower(info.Name()), ext) {
			if info.ModTime().After(newest) {
				newest = info.ModTime()
			}
		}
		return nil
	})
	return newest, err
}

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

	if opt.Incremental {
		srcNewest, err := newestModTime(src, ".html")
		if err != nil {
			return &TranspileError{FilePath: src, Message: "incremental check: " + err.Error()}
		}
		outDir := path.Join(dist, opt.Pkg)
		distNewest, err := newestModTime(outDir, ".go")
		if err != nil {
			// outDir may not exist yet; treat as "no output" and run full transpile
			distNewest = time.Time{}
		}
		if !srcNewest.IsZero() && !distNewest.IsZero() && !srcNewest.After(distNewest) {
			if utils.Log != nil {
				utils.Log.Info("incremental skip (no .html changes)")
			}
			return nil
		}
	}

	files, err := utils.WalkAndReadHTMLFiles(src)
	if err != nil {
		return &TranspileError{FilePath: src, Message: err.Error()}
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
			// Skip the root template name (e.g. "sections") which ParseSections includes from template.New("sections")
			if name == "sections" {
				continue
			}
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
	imports = deduplicateImports(imports)
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
		}
		if propsMap == nil {
			propsMap = make(map[string]string)
		}
		if html, ok := m["html"]; ok {
			slotNames, err := element.SlotNamesFromHTML([]byte(html))
			if err != nil {
				return wrapTranspileErr(name, componentSource[name], componentFileContent[name], err)
			}
			for _, slotName := range slotNames {
				key := "slot" + utils.Capitalize(slotName)
				propsMap[key] = "Element"
			}
		}
		for k := range propsMap {
			components[strings.ToLower(name)].Props[strings.ToLower(k)] = k
		}
		s := gocode.ConstructStruct(propsMap, name)
		structs = append(structs, s)
		structMap[name] = s
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
		return &TranspileError{FilePath: outDir, Message: "failed to create output directory: " + err.Error()}
	}
	// Remove existing .go files so we never mix single-file and multi-file output
	matches, _ := filepath.Glob(filepath.Join(outDir, "*.go"))
	for _, m := range matches {
		_ = os.Remove(m)
	}

	if opt.SingleFile {
		b, err := gocode.ConstructSourceWithPkg(goCodes, structs, imports, opt.Pkg)
		if err != nil {
			return &TranspileError{Message: "codegen: " + err.Error()}
		}
		outPath := path.Join(outDir, "comp_generated.go")
		if err := os.WriteFile(outPath, []byte(b), 0644); err != nil {
			return &TranspileError{FilePath: outPath, Message: err.Error()}
		}
		if !opt.ValidateTypes {
			return nil
		}
	} else {
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
				return &TranspileError{Component: name, FilePath: componentSource[name], Message: "codegen: " + err.Error()}
			}
			filename := componentFileName(name)
			if err := os.WriteFile(path.Join(outDir, filename), []byte(compContent), 0644); err != nil {
				return &TranspileError{FilePath: path.Join(outDir, filename), Message: err.Error()}
			}
		}
	}

	if opt.ValidateTypes {
		if err := validateGeneratedPackage(outDir, dist, opt.Pkg, componentSource, sectionNames, opt.SingleFile); err != nil {
			return err
		}
	}
	return nil
}
