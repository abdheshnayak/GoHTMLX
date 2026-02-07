package types

type Table struct {
	Id string
}

type Form struct{}

// Showcase types
type NavLink struct {
	Label string
	Href  string
	IsCta string // optional CSS class, e.g. "nav-cta"
}

type Feature struct {
	Title       string
	Description string
	Icon        string
	Code        string // optional; HTML-escaped snippet shown below description when ShowCode is true
	ShowCode    bool
	Language    string // optional; for syntax highlighting, e.g. "language-go", "language-bash", "language-html"
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
