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
