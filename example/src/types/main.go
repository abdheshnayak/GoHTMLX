package types

type Table struct {
	Id string
}

type Form struct{}

// Showcase types
type NavLink struct {
	Label string
	Href  string
}

type Feature struct {
	Title       string
	Description string
	Icon        string
}

type Stat struct {
	Value string
	Label string
}

type Testimonial struct {
	Quote  string
	Author string
	Role   string
}

// DocSection is used for the feature-showcase explanation blocks.
type DocSection struct {
	Title     string
	Body      string
	Badge     string
	ShowBadge bool
}

// CodeExample holds a teaching snippet (title, optional description, code shown as-is; code should be HTML-escaped).
type CodeExample struct {
	Title           string
	Description     string
	ShowDescription bool
	Code            string
}
