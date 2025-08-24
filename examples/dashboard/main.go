package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/abdheshnayak/gohtmlx/examples/dashboard/dist/gohtmlxc"
	t "github.com/abdheshnayak/gohtmlx/examples/dashboard/src/types"
	"github.com/abdheshnayak/gohtmlx/pkg/element"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func main() {
	app := fiber.New(fiber.Config{
		Views:        nil,
		ErrorHandler: errorHandler,
	})

	// Middleware
	app.Use(recover.New())
	app.Use(logger.New())
	app.Use(cors.New())

	// Static files
	app.Static("/static", "./static")

	// Routes
	setupRoutes(app)

	// Start server
	port := getEnv("PORT", "3000")
	log.Printf("🚀 Dashboard server starting on port %s", port)
	log.Printf("📊 Visit http://localhost:%s to view the dashboard", port)
	
	if err := app.Listen(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}

func setupRoutes(app *fiber.App) {
	// Dashboard routes
	app.Get("/", handleDashboard)
	app.Get("/dashboard", handleDashboard)
	app.Get("/users", handleUsers)
	app.Get("/settings", handleSettings)
	app.Get("/profile", handleProfile)
	
	// API routes
	api := app.Group("/api")
	api.Get("/stats", handleAPIStats)
	api.Get("/users", handleAPIUsers)
	api.Get("/notifications", handleAPINotifications)
	api.Post("/users", handleAPICreateUser)
	api.Put("/users/:id", handleAPIUpdateUser)
	api.Delete("/users/:id", handleAPIDeleteUser)
}

func handleDashboard(c *fiber.Ctx) error {
	stats := getDashboardStats()
	users := getRecentUsers(5)
	activities := getRecentActivity(10)
	chartData := getChartData()

	dashboardPage := gohtmlxc.DashboardPageComp(gohtmlxc.DashboardPage{
		Stats:          stats,
		Users:          users,
		RecentActivity: activities,
		ChartData:      chartData,
	}, element.Attrs{})

	return c.Type("html").SendString(dashboardPage.Render())
}

func handleUsers(c *fiber.Ctx) error {
	page := c.QueryInt("page", 1)
	limit := c.QueryInt("limit", 10)
	search := c.Query("search", "")

	users := getUsers(page, limit, search)
	totalUsers := getTotalUsers()
	totalPages := (totalUsers + limit - 1) / limit

	columns := getUserTableColumns(true)
	userData := formatUsersForTable(users)

	tableComponent := gohtmlxc.TableComp(gohtmlxc.Table{
		Columns:      columns,
		Data:         userData,
		Striped:      true,
		Hoverable:    true,
		Loading:      false,
		EmptyMessage: "No users found",
	}, element.Attrs{})

	paginationComponent := gohtmlxc.PaginationComp(gohtmlxc.Pagination{
		CurrentPage:  page,
		TotalPages:   totalPages,
		TotalItems:   totalUsers,
		ItemsPerPage: limit,
		ShowInfo:     true,
	}, element.Attrs{})

	layoutComponent := gohtmlxc.LayoutComp(gohtmlxc.Layout{
		Title:        "Users",
		User:         getCurrentUser(),
		Navigation:   getNavigation(),
		Notifications: getNotifications(),
		ShowSearch:   true,
		Breadcrumbs:  []map[string]string{
			{"label": "Dashboard", "url": "/"},
			{"label": "Users", "url": ""},
		},
		PageActions: []map[string]string{
			{
				"label":   "Add User",
				"variant": "primary",
				"icon":    "fas fa-plus",
				"onclick": "openAddUserModal()",
			},
		},
	}, element.Attrs{}, 
		element.Div(element.Attrs{"class": "space-y-6"},
			tableComponent,
			paginationComponent,
		),
	)

	return c.Type("html").SendString(layoutComponent.Render())
}

func handleSettings(c *fiber.Ctx) error {
	layoutComponent := gohtmlxc.LayoutComp(gohtmlxc.Layout{
		Title:        "Settings",
		User:         getCurrentUser(),
		Navigation:   getNavigation(),
		Notifications: getNotifications(),
		ShowSearch:   false,
		Breadcrumbs:  []map[string]string{
			{"label": "Dashboard", "url": "/"},
			{"label": "Settings", "url": ""},
		},
	}, element.Attrs{}, 
		gohtmlxc.CardComp(gohtmlxc.Card{
			Title: "Application Settings",
			Icon:  "fas fa-cog",
		}, element.Attrs{}, 
			element.P(element.Attrs{"class": "text-gray-600"}, 
				element.Text("Settings page is under construction."),
			),
		),
	)

	return c.Type("html").SendString(layoutComponent.Render())
}

func handleProfile(c *fiber.Ctx) error {
	user := getCurrentUser()
	
	layoutComponent := gohtmlxc.LayoutComp(gohtmlxc.Layout{
		Title:        "Profile",
		User:         user,
		Navigation:   getNavigation(),
		Notifications: getNotifications(),
		ShowSearch:   false,
		Breadcrumbs:  []map[string]string{
			{"label": "Dashboard", "url": "/"},
			{"label": "Profile", "url": ""},
		},
	}, element.Attrs{}, 
		element.Div(element.Attrs{"class": "grid grid-cols-1 lg:grid-cols-3 gap-8"},
			element.Div(element.Attrs{"class": "lg:col-span-1"},
				gohtmlxc.CardComp(gohtmlxc.Card{
					Title: "Profile Information",
					Icon:  "fas fa-user",
				}, element.Attrs{}, 
					element.Div(element.Attrs{"class": "text-center"},
						gohtmlxc.AvatarComp(gohtmlxc.Avatar{
							Src:  user.Avatar,
							Name: user.Name,
							Size: "lg",
						}, element.Attrs{"class": "mx-auto mb-4"}),
						element.H3(element.Attrs{"class": "text-lg font-medium text-gray-900"}, 
							element.Text(user.Name),
						),
						element.P(element.Attrs{"class": "text-gray-500"}, 
							element.Text(user.Email),
						),
						gohtmlxc.BadgeComp(gohtmlxc.Badge{
							Text:    user.Role,
							Variant: "primary",
						}, element.Attrs{"class": "mt-2"}),
					),
				),
			),
			element.Div(element.Attrs{"class": "lg:col-span-2"},
				gohtmlxc.CardComp(gohtmlxc.Card{
					Title: "Account Details",
					Icon:  "fas fa-info-circle",
				}, element.Attrs{}, 
					element.P(element.Attrs{"class": "text-gray-600"}, 
						element.Text("Profile editing functionality is under construction."),
					),
				),
			),
		),
	)

	return c.Type("html").SendString(layoutComponent.Render())
}

// API Handlers
func handleAPIStats(c *fiber.Ctx) error {
	stats := getDashboardStats()
	return c.JSON(stats)
}

func handleAPIUsers(c *fiber.Ctx) error {
	page := c.QueryInt("page", 1)
	limit := c.QueryInt("limit", 10)
	search := c.Query("search", "")

	users := getUsers(page, limit, search)
	return c.JSON(fiber.Map{
		"users": users,
		"total": getTotalUsers(),
		"page":  page,
		"limit": limit,
	})
}

func handleAPINotifications(c *fiber.Ctx) error {
	notifications := getNotifications()
	return c.JSON(notifications)
}

func handleAPICreateUser(c *fiber.Ctx) error {
	var user t.User
	if err := c.BodyParser(&user); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}
	
	// Simulate user creation
	user.ID = rand.Intn(10000)
	user.Status = "active"
	
	return c.Status(201).JSON(user)
}

func handleAPIUpdateUser(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid user ID"})
	}
	
	var user t.User
	if err := c.BodyParser(&user); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}
	
	user.ID = id
	return c.JSON(user)
}

func handleAPIDeleteUser(c *fiber.Ctx) error {
	id := c.Params("id")
	return c.JSON(fiber.Map{"message": fmt.Sprintf("User %s deleted successfully", id)})
}

// Data functions
func getDashboardStats() t.DashboardStats {
	return t.DashboardStats{
		TotalUsers:     1250,
		ActiveUsers:    890,
		Revenue:        45678.90,
		Growth:         12.5,
		NewSignups:     45,
		ConversionRate: 3.2,
	}
}

func getCurrentUser() t.User {
	return t.User{
		ID:       1,
		Name:     "John Doe",
		Email:    "john.doe@example.com",
		Avatar:   "https://images.unsplash.com/photo-1472099645785-5658abf4ff4e?ixlib=rb-1.2.1&ixid=eyJhcHBfaWQiOjEyMDd9&auto=format&fit=facepad&facepad=2&w=256&h=256&q=80",
		Role:     "Administrator",
		Status:   "online",
		LastSeen: time.Now(),
	}
}

func getNavigation() []t.NavigationItem {
	return []t.NavigationItem{
		{
			ID:     "dashboard",
			Label:  "Dashboard",
			Icon:   "fas fa-home",
			URL:    "/",
			Active: true,
		},
		{
			ID:    "users",
			Label: "Users",
			Icon:  "fas fa-users",
			URL:   "/users",
			Badge: "1.2k",
		},
		{
			ID:    "analytics",
			Label: "Analytics",
			Icon:  "fas fa-chart-bar",
			Children: []t.NavigationItem{
				{ID: "reports", Label: "Reports", URL: "/analytics/reports"},
				{ID: "insights", Label: "Insights", URL: "/analytics/insights"},
			},
		},
		{
			ID:    "settings",
			Label: "Settings",
			Icon:  "fas fa-cog",
			URL:   "/settings",
		},
	}
}

func getNotifications() []t.Notification {
	return []t.Notification{
		{
			ID:        "1",
			Type:      "info",
			Title:     "New user registered",
			Message:   "Alice Johnson just signed up for an account",
			Timestamp: time.Now().Add(-5 * time.Minute),
			Read:      false,
		},
		{
			ID:        "2",
			Type:      "success",
			Title:     "Payment received",
			Message:   "$299.00 payment from Acme Corp",
			Timestamp: time.Now().Add(-1 * time.Hour),
			Read:      false,
		},
		{
			ID:        "3",
			Type:      "warning",
			Title:     "Server maintenance",
			Message:   "Scheduled maintenance in 2 hours",
			Timestamp: time.Now().Add(-2 * time.Hour),
			Read:      true,
		},
	}
}

func getRecentUsers(limit int) []t.User {
	users := []t.User{
		{ID: 1, Name: "Alice Johnson", Email: "alice@example.com", Avatar: "https://images.unsplash.com/photo-1494790108755-2616b612b786?ixlib=rb-1.2.1&ixid=eyJhcHBfaWQiOjEyMDd9&auto=format&fit=facepad&facepad=2&w=256&h=256&q=80", Role: "User", Status: "online", LastSeen: time.Now().Add(-5 * time.Minute)},
		{ID: 2, Name: "Bob Smith", Email: "bob@example.com", Avatar: "https://images.unsplash.com/photo-1519244703995-f4e0f30006d5?ixlib=rb-1.2.1&ixid=eyJhcHBfaWQiOjEyMDd9&auto=format&fit=facepad&facepad=2&w=256&h=256&q=80", Role: "User", Status: "offline", LastSeen: time.Now().Add(-2 * time.Hour)},
		{ID: 3, Name: "Carol Davis", Email: "carol@example.com", Avatar: "https://images.unsplash.com/photo-1517841905240-472988babdf9?ixlib=rb-1.2.1&ixid=eyJhcHBfaWQiOjEyMDd9&auto=format&fit=facepad&facepad=2&w=256&h=256&q=80", Role: "Manager", Status: "away", LastSeen: time.Now().Add(-30 * time.Minute)},
		{ID: 4, Name: "David Wilson", Email: "david@example.com", Avatar: "https://images.unsplash.com/photo-1506794778202-cad84cf45f1d?ixlib=rb-1.2.1&ixid=eyJhcHBfaWQiOjEyMDd9&auto=format&fit=facepad&facepad=2&w=256&h=256&q=80", Role: "User", Status: "online", LastSeen: time.Now().Add(-1 * time.Minute)},
		{ID: 5, Name: "Emma Brown", Email: "emma@example.com", Avatar: "https://images.unsplash.com/photo-1438761681033-6461ffad8d80?ixlib=rb-1.2.1&ixid=eyJhcHBfaWQiOjEyMDd9&auto=format&fit=facepad&facepad=2&w=256&h=256&q=80", Role: "Admin", Status: "busy", LastSeen: time.Now().Add(-10 * time.Minute)},
	}
	
	if limit > 0 && limit < len(users) {
		return users[:limit]
	}
	return users
}

func getUsers(page, limit int, search string) []t.User {
	allUsers := getRecentUsers(0)
	// Simple pagination simulation
	start := (page - 1) * limit
	end := start + limit
	if end > len(allUsers) {
		end = len(allUsers)
	}
	if start > len(allUsers) {
		return []t.User{}
	}
	return allUsers[start:end]
}

func getTotalUsers() int {
	return len(getRecentUsers(0))
}

func getRecentActivity(limit int) []map[string]interface{} {
	activities := []map[string]interface{}{
		{"type": "user", "user": "Alice Johnson", "action": "signed up", "target": "", "timestamp": time.Now().Add(-5 * time.Minute)},
		{"type": "payment", "user": "Bob Smith", "action": "made a payment of", "target": "$299.00", "timestamp": time.Now().Add(-15 * time.Minute)},
		{"type": "edit", "user": "Carol Davis", "action": "updated profile for", "target": "Emma Brown", "timestamp": time.Now().Add(-30 * time.Minute)},
		{"type": "delete", "user": "David Wilson", "action": "deleted", "target": "Project Alpha", "timestamp": time.Now().Add(-1 * time.Hour)},
		{"type": "create", "user": "Emma Brown", "action": "created", "target": "New Campaign", "timestamp": time.Now().Add(-2 * time.Hour)},
	}
	
	if limit > 0 && limit < len(activities) {
		return activities[:limit]
	}
	return activities
}

func getChartData() []t.ChartData {
	return []t.ChartData{
		{Label: "Jan", Value: 12000, Color: "#3B82F6"},
		{Label: "Feb", Value: 19000, Color: "#10B981"},
		{Label: "Mar", Value: 15000, Color: "#F59E0B"},
		{Label: "Apr", Value: 25000, Color: "#EF4444"},
		{Label: "May", Value: 22000, Color: "#8B5CF6"},
		{Label: "Jun", Value: 30000, Color: "#06B6D4"},
	}
}

func getUserTableColumns(showActions bool) []t.TableColumn {
	columns := []t.TableColumn{
		{Key: "name", Label: "Name", Type: "avatar", Align: "left"},
		{Key: "email", Label: "Email", Type: "text", Align: "left"},
		{Key: "role", Label: "Role", Type: "badge", Align: "left"},
		{Key: "status", Label: "Status", Type: "badge", Align: "center"},
		{Key: "last_seen", Label: "Last Seen", Type: "date", Align: "right"},
	}
	
	if showActions {
		columns = append(columns, t.TableColumn{
			Key: "actions", Label: "Actions", Type: "actions", Align: "center", Width: "20",
		})
	}
	
	return columns
}

func formatUsersForTable(users []t.User) []map[string]interface{} {
	var data []map[string]interface{}
	for _, user := range users {
		data = append(data, map[string]interface{}{
			"id":        fmt.Sprintf("%d", user.ID),
			"name":      user.Name,
			"email":     user.Email,
			"avatar":    user.Avatar,
			"role":      user.Role,
			"status":    user.Status,
			"last_seen": user.LastSeen,
		})
	}
	return data
}

// Helper functions (these would normally be in separate utility files)
func getEnv(key, defaultValue string) string {
	// In a real application, you'd use os.Getenv
	return defaultValue
}

func errorHandler(c *fiber.Ctx, err error) error {
	code := fiber.StatusInternalServerError
	if e, ok := err.(*fiber.Error); ok {
		code = e.Code
	}
	
	return c.Status(code).JSON(fiber.Map{
		"error": err.Error(),
	})
}

// Template helper functions (these would be available in templates)
func getBadgeVariant(value string) string {
	switch value {
	case "online", "active", "Admin":
		return "success"
	case "offline", "inactive":
		return "secondary"
	case "away", "Manager":
		return "warning"
	case "busy":
		return "danger"
	default:
		return "primary"
	}
}

func formatDate(value interface{}) string {
	if t, ok := value.(time.Time); ok {
		return t.Format("Jan 2, 2006")
	}
	return fmt.Sprintf("%v", value)
}

func formatNumber(value interface{}) string {
	if f, ok := value.(float64); ok {
		return fmt.Sprintf("%.0f", f)
	}
	return fmt.Sprintf("%v", value)
}

func formatCurrency(value interface{}) string {
	if f, ok := value.(float64); ok {
		return fmt.Sprintf("$%.2f", f)
	}
	return fmt.Sprintf("$%v", value)
}

func formatTimeAgo(timestamp interface{}) string {
	if t, ok := timestamp.(time.Time); ok {
		duration := time.Since(t)
		if duration < time.Minute {
			return "just now"
		} else if duration < time.Hour {
			return fmt.Sprintf("%d minutes ago", int(duration.Minutes()))
		} else if duration < 24*time.Hour {
			return fmt.Sprintf("%d hours ago", int(duration.Hours()))
		} else {
			return fmt.Sprintf("%d days ago", int(duration.Hours()/24))
		}
	}
	return "unknown"
}

func getActivityColor(activityType interface{}) string {
	switch fmt.Sprintf("%v", activityType) {
	case "user":
		return "bg-blue-500"
	case "payment":
		return "bg-green-500"
	case "edit":
		return "bg-yellow-500"
	case "delete":
		return "bg-red-500"
	case "create":
		return "bg-purple-500"
	default:
		return "bg-gray-500"
	}
}

func getActivityIcon(activityType interface{}) string {
	switch fmt.Sprintf("%v", activityType) {
	case "user":
		return "fas fa-user"
	case "payment":
		return "fas fa-dollar-sign"
	case "edit":
		return "fas fa-edit"
	case "delete":
		return "fas fa-trash"
	case "create":
		return "fas fa-plus"
	default:
		return "fas fa-info"
	}
}

func getUnreadNotificationCount(notifications []t.Notification) int {
	count := 0
	for _, n := range notifications {
		if !n.Read {
			count++
		}
	}
	return count
}

func getStatusFromUser(user t.User) string {
	return user.Status
}

func getNotificationColor(notificationType string) string {
	switch notificationType {
	case "success":
		return "bg-green-400"
	case "warning":
		return "bg-yellow-400"
	case "error":
		return "bg-red-400"
	default:
		return "bg-blue-400"
	}
}

func getPaginationRange(current, total int) []int {
	// Simple pagination range - in a real app this would be more sophisticated
	var pages []int
	for i := 1; i <= total; i++ {
		pages = append(pages, i)
	}
	return pages
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
