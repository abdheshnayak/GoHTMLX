package main

import (
	"fmt"

	goxpages "github.com/abdheshnayak/gox/example/dist/components"
	"github.com/abdheshnayak/gox/pkg/element"
	"github.com/gofiber/fiber/v2"
)

func main() {

	app := fiber.New()

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

	if err := app.Listen(":3000"); err != nil {
		panic(err)
	}

}

func invokeExport(module string) (element.Element, error) {
	switch module {
	case "home":
		return goxpages.Home(), nil
	default:
		return nil, fmt.Errorf("module %s not found", module)
	}
}
