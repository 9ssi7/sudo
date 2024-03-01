package sudo_fiber

import "github.com/gofiber/fiber/v2"

type CurrentUser struct {
	TwoFactorEnabled bool
	UUID             string
	Phone            string
	Email            string
}

type CurrentUserParser func(ctx *fiber.Ctx) CurrentUser
