package sudo_fiber

import "github.com/gofiber/fiber/v2"

type LocaleParser func(ctx *fiber.Ctx) string
