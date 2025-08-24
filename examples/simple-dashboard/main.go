package main

import (
	"bufio"
	"log"

	"github.com/abdheshnayak/gohtmlx/examples/simple-dashboard/dist/gohtmlxc"
	"github.com/abdheshnayak/gohtmlx/pkg/element"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func main() {
	app := fiber.New(fiber.Config{
		Views: nil,
	})

	// Middleware
	app.Use(recover.New())
	app.Use(logger.New())
	app.Use(cors.New())

	// Static files
	app.Static("/static", "./static")

	// Routes
	app.Get("/", handleDashboard)

	// Start server
	port := "3000"
	log.Printf("🚀 GoHTMLX Simple Dashboard starting on port %s", port)
	log.Printf("📊 Visit http://localhost:%s to view the dashboard", port)

	if err := app.Listen(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}

func handleDashboard(c *fiber.Ctx) error {
	dashboard := gohtmlxc.DashboardComp(gohtmlxc.Dashboard{
		Title:       "GoHTMLX Dashboard",
		TotalUsers:  1250,
		ActiveUsers: 890,
		Revenue:     45678.90,
		Growth:      12.5,
	}, element.Attrs{})

	w := bufio.NewWriter(c)
	i, err := dashboard.Render(w)
	if err != nil {
		return err
	}

	return c.Type("html").SendString(i)
}
