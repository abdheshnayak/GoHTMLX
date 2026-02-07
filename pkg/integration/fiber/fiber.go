// Package fiber provides optional Fiber framework integration for GoHTMLX.
// Use this package only when your app uses Fiber; the core transpiler has no Fiber dependency.
package fiber

import (
	"log/slog"
	"time"

	"github.com/abdheshnayak/gohtmlx/pkg/utils"
	"github.com/gofiber/fiber/v2"
)

// Log is a slog-based logger for use in apps (e.g. set utils.Log = fiber.Log in example main).
var Log utils.Logger = utils.NewSlogLogger(slog.Default())

// FiberLogger is a Fiber middleware that logs each request (method, path, status, latency).
// Use from your Fiber app: app.Use(fiber.FiberLogger).
func FiberLogger(c *fiber.Ctx) error {
	start := time.Now()
	err := c.Next()
	// Log after next so we have status code
	Log.Info("HTTP request",
		slog.String("method", c.Method()),
		slog.String("path", c.Path()),
		slog.Int("status", c.Response().StatusCode()),
		slog.Duration("latency", time.Since(start)),
		slog.String("ip", c.IP()),
	)
	return err
}
