package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/abdheshnayak/gohtmlx/examples/showcase/src/comps"
	"github.com/abdheshnayak/gohtmlx/pkg/element"
	gohtmlxfiber "github.com/abdheshnayak/gohtmlx/pkg/integration/fiber"
	"github.com/abdheshnayak/gohtmlx/pkg/utils"
	"github.com/gofiber/fiber/v2"
)

func main() {
	utils.Log = gohtmlxfiber.Log
	app := fiber.New(fiber.Config{
		DisableStartupMessage: true,
	})

	app.Use(gohtmlxfiber.FiberLogger)

	app.Static("/", "./dist/static", fiber.Static{
		Compress:      true,
		CacheDuration: time.Microsecond,
	})

	// HTMX fragment endpoints — return GoHTMLX-generated HTML for partial updates
	app.Get("/api/time", func(c *fiber.Ctx) error {
		el := comps.ServerTime(time.Now().Format("2006-01-02 15:04:05 MST"), "Server time")
		c.Set("Content-Type", "text/html")
		_, err := el.Render(c)
		return err
	})
	app.Post("/api/feedback", func(c *fiber.Ctx) error {
		name := strings.TrimSpace(c.FormValue("name"))
		msg := strings.TrimSpace(c.FormValue("message"))
		var el element.Element
		if name == "" || msg == "" {
			el = comps.FeedbackErrors("Please fill in both name and message.")
		} else {
			el = comps.FeedbackSuccess("Thanks, " + name + "! We got your message.")
		}
		c.Set("Content-Type", "text/html")
		_, err := el.Render(c)
		return err
	})

	// Route to handle dynamic exports
	app.Get("/", func(c *fiber.Ctx) error {
		response, err := invokeExport("home")
		if err != nil {
			return c.Status(404).SendString(err.Error())
		}

		c.Set("Content-Type", "text/html")
		if _, err := response.Render(c); err != nil {
			return err
		}

		return nil
	})

	app.Post("/:module", func(c *fiber.Ctx) error {
		module := c.Params("module")
		response, err := invokeExport(module)
		if err != nil {
			return c.Status(404).SendString(err.Error())
		}

		c.Set("Content-Type", "text/html")
		if _, err := response.Render(c); err != nil {
			return err
		}

		return nil
	})

	app.Get("/:module", func(c *fiber.Ctx) error {
		module := c.Params("module")
		response, err := invokeExport(module)
		if err != nil {
			return c.Status(404).SendString(err.Error())
		}

		c.Set("Content-Type", "text/html")
		if _, err := response.Render(c); err != nil {
			return err
		}

		return nil
	})

	gohtmlxfiber.Log.Info("Listening on port 3000")
	if err := app.Listen(":3000"); err != nil {
		panic(err)
	}

}

func invokeExport(module string) (element.Element, error) {
	switch module {
	case "home":
		return comps.Home(), nil
	default:
		return nil, fmt.Errorf("module %s not found", module)
	}
}
