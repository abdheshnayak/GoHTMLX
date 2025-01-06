package utils

import (
	"log/slog"
	"strings"
	"text/template"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/nxtcoder17/fwatcher/pkg/logging"
)

func Capitalize(s string) string {
	if len(s) == 0 {
		return s
	}

	return strings.ToUpper(s[:1]) + s[1:]
}

func GetSections(tmpl *template.Template) []string {
	sectionNames := []string{}
	t := tmpl.Templates()
	for _, v := range t {
		sectionNames = append(sectionNames, v.Name())
	}

	return sectionNames
}

func ParseSections(tmpl *template.Template) (map[string]string, error) {
	sectionNames := GetSections(tmpl)
	t := tmpl.Templates()
	for _, v := range t {
		sectionNames = append(sectionNames, v.Name())
	}

	sections := make(map[string]string)

	for _, section := range sectionNames {
		var buffer strings.Builder
		err := tmpl.ExecuteTemplate(&buffer, section, nil)
		if err != nil {
			return nil, err
		}
		sections[section] = buffer.String()
	}

	return sections, nil
}

var Log = logging.NewSlogLogger(logging.SlogOptions{
	ShowTimestamp: true,
	ShowCaller:    false,
})

func FiberLogger(c *fiber.Ctx) error {
	start := time.Now()

	// Process request
	err := c.Next()

	// Log the request details
	Log.Info("HTTP request",
		slog.String("method", c.Method()),
		slog.String("path", c.Path()),
		slog.Int("status", c.Response().StatusCode()),
		slog.Duration("latency", time.Since(start)),
		slog.String("ip", c.IP()),
	)

	return err
}
