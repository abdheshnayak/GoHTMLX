package types

import "time"

// User represents a user in the system
type User struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Avatar   string `json:"avatar"`
	Role     string `json:"role"`
	Status   string `json:"status"`
	LastSeen time.Time `json:"last_seen"`
}

// DashboardStats represents key metrics for the dashboard
type DashboardStats struct {
	TotalUsers     int     `json:"total_users"`
	ActiveUsers    int     `json:"active_users"`
	Revenue        float64 `json:"revenue"`
	Growth         float64 `json:"growth"`
	NewSignups     int     `json:"new_signups"`
	ConversionRate float64 `json:"conversion_rate"`
}

// ChartData represents data for charts
type ChartData struct {
	Label string  `json:"label"`
	Value float64 `json:"value"`
	Color string  `json:"color"`
}

// TableColumn represents a table column configuration
type TableColumn struct {
	Key    string `json:"key"`
	Label  string `json:"label"`
	Width  string `json:"width"`
	Align  string `json:"align"`
	Type   string `json:"type"` // text, number, badge, avatar, date
}

// NavigationItem represents a navigation menu item
type NavigationItem struct {
	ID       string `json:"id"`
	Label    string `json:"label"`
	Icon     string `json:"icon"`
	URL      string `json:"url"`
	Active   bool   `json:"active"`
	Badge    string `json:"badge"`
	Children []NavigationItem `json:"children"`
}

// Notification represents a system notification
type Notification struct {
	ID        string    `json:"id"`
	Type      string    `json:"type"` // success, warning, error, info
	Title     string    `json:"title"`
	Message   string    `json:"message"`
	Timestamp time.Time `json:"timestamp"`
	Read      bool      `json:"read"`
}

// ButtonVariant represents button styling variants
type ButtonVariant string

const (
	ButtonPrimary   ButtonVariant = "primary"
	ButtonSecondary ButtonVariant = "secondary"
	ButtonSuccess   ButtonVariant = "success"
	ButtonWarning   ButtonVariant = "warning"
	ButtonDanger    ButtonVariant = "danger"
	ButtonGhost     ButtonVariant = "ghost"
)

// ButtonSize represents button sizes
type ButtonSize string

const (
	ButtonSmall  ButtonSize = "sm"
	ButtonMedium ButtonSize = "md"
	ButtonLarge  ButtonSize = "lg"
)
