package sudo_fiber

import (
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type DeviceIdConfig struct {
	CookieName string
}

func NewDeviceId(cnf DeviceIdConfig) fiber.Handler {
	if cnf.CookieName == "" {
		cnf.CookieName = "device_id"
	}
	httpReplacer := func(host string) string {
		return strings.Replace(strings.Replace(host, "http://", "", 1), "https://", "", 1)
	}
	return func(c *fiber.Ctx) error {
		reqDomain := httpReplacer(c.Get("host"))
		deviceId := c.Cookies(cnf.CookieName)
		_, err := uuid.Parse(deviceId)
		if deviceId == "" || err != nil {
			deviceId = uuid.New().String()
			c.Cookie(&fiber.Cookie{
				Name:     cnf.CookieName,
				Value:    deviceId,
				Expires:  time.Now().Add(24 * time.Hour * 365),
				HTTPOnly: true,
				Secure:   true,
				Domain:   reqDomain,
				SameSite: "Strict",
			})
		}
		c.Locals("device_id", deviceId)
		return c.Next()
	}
}

func DeviceIdParse(c *fiber.Ctx) string {
	return c.Locals("device_id").(string)
}
