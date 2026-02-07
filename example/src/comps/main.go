package comps

import (
	"strings"

	gc "github.com/abdheshnayak/gohtmlx/example/dist/gohtmlxc"
	"github.com/abdheshnayak/gohtmlx/pkg/element"
	t "github.com/abdheshnayak/gohtmlx/example/src/types"
)

// htmlEscape escapes code so the browser shows it as text instead of rendering tags/comments.
func htmlEscape(s string) string {
	s = strings.ReplaceAll(s, "&", "&amp;")
	s = strings.ReplaceAll(s, "<", "&lt;")
	s = strings.ReplaceAll(s, ">", "&gt;")
	s = strings.ReplaceAll(s, "\"", "&quot;")
	return s
}

func Home() element.Element {
	return gc.Home{
		NavLinks: []t.NavLink{
			{Label: "Features", Href: "#features"},
			{Label: "Docs", Href: "https://github.com/abdheshnayak/gohtmlx#readme"},
			{Label: "GitHub", Href: "https://github.com/abdheshnayak/gohtmlx", IsCta: "nav-cta"},
		},
		Features: []t.Feature{
			{
				Title: "Quick start", Description: "Install the CLI, point it at your HTML, and use the generated package. Works with any HTTP framework.",
				Code: htmlEscape("go install github.com/abdheshnayak/gohtmlx@latest\ngohtmlx --src=./src --dist=./dist\n# In your app: import the generated package and call ComponentName{...}.Get().Render(w)"),
				ShowCode: true,
			},
			{
				Title: "Define a component",
				Description: htmlEscape("Wrap the component in <!-- + define \"Name\" --> ... <!-- + end -->. Use <!-- | define \"props\" --> for YAML props and <!-- | define \"html\" --> for the template."),
				Code: htmlEscape("<!-- + define \"Greet\" -->\n<!-- | define \"props\" -->\nname: string\n<!-- | end -->\n<!-- | define \"html\" -->\n<div>Hello, {props.Name}!</div>\n<!-- | end -->\n<!-- + end -->"),
				ShowCode: true,
			},
			{
				Title: "Expressions and props", Description: "Use {props.Field} in the HTML and attr={value} for attributes. Pass props when using the component.",
				Code: htmlEscape("<Greet name={props.UserName}></Greet>\n<p>{props.A} — {props.B}</p>"),
				ShowCode: true,
			},
			{
				Title: "Loops",
				Description: htmlEscape("Use <for items={props.Items} as=\"item\">. The body is repeated for each element."),
				Code: htmlEscape("<for items={props.Links} as=\"link\">\n  <li><a href={link.Href}>{link.Label}</a></li>\n</for>"),
				ShowCode: true,
			},
			{
				Title: "Conditionals",
				Description: htmlEscape("Use <if condition={bool}>, optional <elseif>, and <else>. Condition must be a boolean expression."),
				Code: htmlEscape("<if condition={props.ShowHero}>\n  <Hero title={props.Title}></Hero>\n</if>\n<else>\n  <p>Default</p>\n</else>"),
				ShowCode: true,
			},
			{
				Title: "Slots",
				Description: htmlEscape("Layouts declare <slot name=\"header\"/>; callers pass <slot name=\"header\">content</slot> as direct children."),
				Code: htmlEscape("// In layout:\n<div><header><slot name=\"header\"/></header><main><slot name=\"body\"/></main></div>\n\n// At call site:\n<Card><slot name=\"header\">Title</slot><slot name=\"body\">Body</slot></Card>"),
				ShowCode: true,
			},
		},
		Stats: []t.Stat{
			{Value: "0", Label: "JS runtime"},
			{Value: "1", Label: "CLI command"},
			{Value: "100%", Label: "Go"},
		},
		ShowHero:          true,
		HeroTitle:         "HTML-first server components for Go",
		HeroSubtitle:      "Write components in HTML, transpile to Go. No JavaScript, no framework lock-in — just server-rendered views that fit any stack.",
		HeroBadge:         "No JavaScript • Server-only • MIT",
		ShowHeroBadge:     true,
		ShowCtaPrimary:    true,
		CtaPrimaryText:    "Get started",
		CtaPrimaryHref:    "#features",
		ShowCtaSecondary:  true,
		CtaSecondaryText:  "View on GitHub",
		CtaSecondaryHref:  "https://github.com/abdheshnayak/gohtmlx",
		ShowAlert:         true,
		AlertMessage:      "This page is built 100% with GoHTMLX. Examples below are in the Features section only.",
		Attrs:             element.Attrs{},
	}.Get()
}
