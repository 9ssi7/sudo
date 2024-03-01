package sudo_fiber

import "github.com/gofiber/fiber/v2"

// LocaleConfig holds the configuration for the locale middleware.
type LocaleParser func(ctx *fiber.Ctx) string
