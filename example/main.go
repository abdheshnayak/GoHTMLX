package main

import (
	"fmt"
	"time"

	"github.com/abdheshnayak/gohtmlx/example/src/comps"
	"github.com/abdheshnayak/gohtmlx/pkg/element"
	"github.com/abdheshnayak/gohtmlx/pkg/utils"
	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New(fiber.Config{
		DisableStartupMessage: true,
	})

	app.Use(utils.FiberLogger)

	app.Static("/", "./dist/static", fiber.Static{
		Compress:      true,
		CacheDuration: time.Microsecond,
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

	utils.Log.Info("Listening on port 3000")
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
