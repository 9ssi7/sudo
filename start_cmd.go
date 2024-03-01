// Package sudo provides a service for handling sudo operations.
package sudo

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/google/uuid"
)

// StartCommand represents the command structure for the Start operation.
type StartCommand struct {
	// UserId is the unique identifier of the user.
	UserId string

	// DeviceId is the unique identifier of the device.
	DeviceId string

	// Phone is the phone number to send the code to.
	Phone string

	// Email is the email address to send the code to.
	Email string

	// Locale is the language and region to use for the notification.
	Locale string
}

// Start initiates the sudo operation for a given user and device.
func (c *client) Start(ctx context.Context, cmd StartCommand) (*string, error) {
	verifyToken := uuid.New().String()
	e := entity{
		DeviceId:    cmd.DeviceId,
		AccessToken: nil,
		VerifyToken: &verifyToken,
		Code:        c.makeCode(),
		IsVerified:  false,
		ExpiresAt:   time.Now().Add(c.expire).Unix(),
	}
	b, err := json.Marshal(e)
	if err != nil {
		return nil, errors.New(Messages.FailedMarshal)
	}
	if err = c.redis.Set(ctx, c.calcKey(cmd.DeviceId, cmd.UserId), b); err != nil {
		return nil, errors.New(Messages.FailedRedisSet)
	}
	c.notifySender(NotifyCommand{
		DeviceId: cmd.DeviceId,
		Code:     e.Code,
		Phone:    cmd.Phone,
		Email:    cmd.Email,
		Locale:   cmd.Locale,
	})
	return &verifyToken, nil
}
