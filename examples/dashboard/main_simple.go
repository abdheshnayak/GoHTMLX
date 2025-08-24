package main

import (
	"log"
	"net/http"

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
	app.Get("/", handleSimpleDashboard)
	app.Get("/simple", handleSimpleDashboard)

	// Start server
	port := "3000"
	log.Printf("🚀 Simple Dashboard server starting on port %s", port)
	log.Printf("📊 Visit http://localhost:%s to view the dashboard", port)
	
	if err := app.Listen(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}

func handleSimpleDashboard(c *fiber.Ctx) error {
	stats := getSimpleStats()

	dashboardPage := gohtmlxc.SimplePageComp(gohtmlxc.SimplePage{
		Title: "GoHTMLX Dashboard",
		Stats: stats,
	}, element.Attrs{})

	return c.Type("html").SendString(dashboardPage.Render())
}

func getSimpleStats() t.DashboardStats {
	return t.DashboardStats{
		TotalUsers:     1250,
		ActiveUsers:    890,
		Revenue:        45678.90,
		Growth:         12.5,
		NewSignups:     45,
		ConversionRate: 3.2,
	}
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
