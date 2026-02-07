package comps

import (
	gc "github.com/abdheshnayak/gohtmlx/example/dist/gohtmlxc"
	"github.com/abdheshnayak/gohtmlx/pkg/element"
	t "github.com/abdheshnayak/gohtmlx/example/src/types"
)

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
		ShowHero:     true,
		HeroTitle:    "GoHTMLX Showcase",
		HeroSubtitle: "Server-side HTML components with conditionals, loops, and slots — all in Go.",
		ShowAlert:    true,
		AlertMessage: "This page is built entirely with GoHTMLX: layout slots, conditionals, and for-loops.",
		Attrs:        element.Attrs{},
	}.Get()
}
