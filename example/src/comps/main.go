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
			{Label: "Home", Href: "/"},
			{Label: "Features", Href: "#features"},
			{Label: "GitHub", Href: "https://github.com/abdheshnayak/gohtmlx"},
		},
		Features: []t.Feature{
			{Title: "Components", Description: "Define reusable HTML components with props and compose them in Go.", Icon: ""},
			{Title: "Conditionals", Description: "Use <if>, <elseif>, and <else> for conditional rendering without leaving the template.", Icon: ""},
			{Title: "Loops", Description: "Render lists with <for items={...} as=\"item\"> and keep logic in one place.", Icon: ""},
			{Title: "Slots", Description: "Layout components with named slots; fill them at the call site for flexible layouts.", Icon: ""},
			{Title: "Server-side", Description: "Pure server-side rendering. No JS framework; output plain HTML from Go.", Icon: ""},
		},
		Stats: []t.Stat{
			{Value: "0", Label: "Runtime deps"},
			{Value: "1", Label: "Generated package"},
			{Value: "100%", Label: "Go"},
		},
		Testimonials: []t.Testimonial{
			{Quote: "Finally, HTML components that feel like Go. Slots and conditionals made our layout code much cleaner.", Author: "Server Dev", Role: "Backend team"},
			{Quote: "We use it for admin dashboards and internal tools. One template, one build, no npm.", Author: "Ops", Role: "Platform"},
		},
		DocSections: []t.DocSection{
			{Title: "Components & props", Body: "Components are defined with <!-- + define \"Name\" --> and optional props in YAML (e.g. title: string). Use {props.Title} in HTML and pass props at the call site: <Component title={value}>.", Badge: "props", ShowBadge: true},
			{Title: "Loops", Body: "Use <for items={props.Items} as=\"item\"> to iterate; the body is emitted inside a Go range loop. The features grid and stats above use this.", Badge: "for", ShowBadge: true},
			{Title: "Conditionals", Body: "Use <if condition={expr}>, <elseif condition={...}>, and <else> for conditional rendering. The hero and alert at the top are shown when ShowHero and ShowAlert are true.", Badge: "if", ShowBadge: true},
			{Title: "Slots", Body: "Layout components declare placeholders with <slot name=\"header\"/>; callers fill them with <SlotLayout><slot name=\"header\">content</slot></SlotLayout>. See the live demo below.", Badge: "slots", ShowBadge: true},
		},
		CodeExamples: []t.CodeExample{
			{
				Title:           "1. File structure",
				Description:     "Put each component in a .html file under your --src directory. The CLI walks all .html files and merges them; component names must be unique.",
				ShowDescription: true,
				Code:            htmlEscape("src/comps/\n  header.html    # AppHeader\n  hero.html      # Hero, Alert\n  cards.html     # FeatureCard, StatCard\n  home.html      # imports + Home"),
			},
			{
				Title:           "2. Define a component",
				Description:     "Wrap the component in <!-- + define \"Name\" --> ... <!-- + end -->. Use <!-- | define \"props\" --> for YAML props and <!-- | define \"html\" --> for the template.",
				ShowDescription: true,
				Code:            htmlEscape("<!-- + define \"Greet\" -->\n<!-- | define \"props\" -->\nname: string\n<!-- | end -->\n<!-- | define \"html\" -->\n<div>Hello, {props.Name}!</div>\n<!-- | end -->\n<!-- + end -->"),
			},
			{
				Title:           "3. Use expressions and pass props",
				Description:     "In the HTML block use {props.Field} for output and attr={value} for attributes. When using the component, pass props as attributes.",
				ShowDescription: true,
				Code:            htmlEscape("<Greet name={props.UserName}></Greet>\n<p>Two expressions: {props.A} — {props.B}</p>"),
			},
			{
				Title:           "4. Loop over a slice",
				Description:     "Use <for items={props.Items} as=\"item\">. The body is repeated for each element; use the loop variable (e.g. item) inside.",
				ShowDescription: true,
				Code:            htmlEscape("<for items={props.Links} as=\"link\">\n  <li><a href={link.Href}>{link.Label}</a></li>\n</for>"),
			},
			{
				Title:           "5. Conditionals",
				Description:     "Use <if condition={bool}>, optional <elseif condition={...}>, and optional <else>. The condition must be a boolean expression.",
				ShowDescription: true,
				Code:            htmlEscape("<if condition={props.ShowHero}>\n  <Hero title={props.Title}></Hero>\n</if>\n<elseif condition={props.ShowAlt}>\n  <p>Alternative</p>\n</elseif>\n<else>\n  <p>Default</p>\n</else>"),
			},
			{
				Title:           "6. Slots (layout placeholders)",
				Description:     "In a layout component put <slot name=\"header\"/> where content should go. At the call site pass content with <slot name=\"header\">...</slot> as direct children.",
				ShowDescription: true,
				Code:            htmlEscape("// In layout (e.g. Card.html):\n<div class=\"card\">\n  <header><slot name=\"header\"/></header>\n  <main><slot name=\"body\"/></main>\n</div>\n\n// At call site:\n<Card>\n  <slot name=\"header\">Title</slot>\n  <slot name=\"body\">Body text</slot>\n</Card>"),
			},
			{
				Title:           "7. Transpile and use in Go",
				Description:     "Run gohtmlx --src=./src --dist=./dist. Import the generated package and call ComponentName{...}.Get() to get an element; call .Render(w) to write HTML.",
				ShowDescription: true,
				Code:            htmlEscape("// Terminal:\ngohtmlx --src=./src --dist=./dist\n\n// Go:\nimport gc \"yourmod/dist/gohtmlxc\"\n\nel := gc.Greet{Name: \"World\", Attrs: nil}.Get()\nel.Render(w)  // or el.Render(os.Stdout)"),
			},
		},
		ShowHero:     true,
		HeroTitle:    "GoHTMLX Showcase",
		HeroSubtitle: "Server-side HTML components with conditionals, loops, and slots — all in Go.",
		ShowAlert:    true,
		AlertMessage: "This page is built entirely with GoHTMLX: layout slots, conditionals, and for-loops.",
		Attrs:        element.Attrs{},
	}.Get()
}
