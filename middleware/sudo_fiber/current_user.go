package sudo_fiber

import "github.com/gofiber/fiber/v2"

// CurrentUser represents the current user.
type CurrentUser struct {
	// TwoFactorEnabled indicates whether the two-factor authentication is enabled for the user.
	TwoFactorEnabled bool

	// UUID is the unique identifier of the user.
	UUID string

	// Phone is the phone number of the user.
	Phone string

	// Email is the email of the user.
	Email string
}

// CurrentUserParser represents the current user parser.
type CurrentUserParser func(ctx *fiber.Ctx) CurrentUser
