package sudo_fiber

import (
	"github.com/9ssi7/sudo"
	"github.com/cilloparch/cillop/result"
	"github.com/gofiber/fiber/v2"
)

// Config holds the configuration for the sudo middleware.
type Config struct {
	// Service is the sudo service.
	Service sudo.Service

	// UserParser is the function used to parse the current user.
	UserParser CurrentUserParser

	// LocaleParser is the function used to parse the locale.
	LocaleParser LocaleParser

	// AccessKey is the name of the header used to store the access key.
	AccessKey string

	// VerifyKey is the name of the header used to store the verify key.
	VerifyKey string

	// CodeKey is the name of the header used to store the code.
	CodeKey string
}

// New creates a new sudo middleware.
func New(cnf Config) fiber.Handler {
	if cnf.AccessKey == "" {
		cnf.AccessKey = "X-Sudo-Access-KEY"
	}
	if cnf.VerifyKey == "" {
		cnf.VerifyKey = "X-Sudo-Verify-KEY"
	}
	if cnf.CodeKey == "" {
		cnf.CodeKey = "X-Sudo-Code"
	}
	return func(ctx *fiber.Ctx) error {
		user := cnf.UserParser(ctx)
		device_id := DeviceIdParse(ctx)
		if !user.TwoFactorEnabled {
			return ctx.Next()
		}
		l := cnf.LocaleParser(ctx)
		accessKey := ctx.Get(cnf.AccessKey)
		if accessKey != "" {
			err := cnf.Service.Check(ctx.UserContext(), sudo.CheckCommand{
				UserId:   user.UUID,
				DeviceId: device_id,
				Token:    accessKey,
			})
			if err != nil {
				return result.Error(err.Error(), fiber.StatusForbidden)
			}
			return ctx.Next()
		}
		verifyKey := ctx.Get(cnf.VerifyKey)
		code := ctx.Get(cnf.CodeKey)
		if verifyKey == "" || code == "" {
			verifyToken, err := cnf.Service.Start(ctx.UserContext(), sudo.StartCommand{
				UserId:   user.UUID,
				DeviceId: device_id,
				Phone:    user.Phone,
				Email:    user.Email,
				Locale:   l,
			})
			if err != nil {
				return result.Error(err.Error(), fiber.StatusForbidden)
			}
			ctx.Set(cnf.VerifyKey, *verifyToken)
			return result.ErrorDetail(sudo.Messages.VerifyStarted, map[string]interface{}{"verify_required": true}, fiber.StatusForbidden)
		}
		if err := cnf.Service.VerifyToken(verifyKey); err != nil {
			return result.Error(err.Error(), fiber.StatusForbidden)
		}
		tkn, err := cnf.Service.Verify(ctx.UserContext(), sudo.VerifyCommand{
			UserId:      user.UUID,
			DeviceId:    device_id,
			VerifyToken: verifyKey,
			Code:        code,
		})
		if err != nil {
			return result.Error(err.Error(), fiber.StatusForbidden)
		}
		ctx.Set(cnf.AccessKey, *tkn)
		return ctx.Next()
	}
}
