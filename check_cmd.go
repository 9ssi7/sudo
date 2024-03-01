// Package sudo provides a service for handling sudo operations.
package sudo

import (
	"context"
	"errors"
)

// CheckCommand represents the command structure for the Check operation.
type CheckCommand struct {
	// UserId is the unique identifier of the user.
	UserId string

	// DeviceId is the unique identifier of the device.
	DeviceId string

	// Token is the access token to be checked.
	Token string
}

// Check validates the provided token for a given user and device.
func (c *client) Check(ctx context.Context, cmd CheckCommand) error {
	e, ok, err := c.getByKey(ctx, c.calcKey(cmd.DeviceId, cmd.UserId))
	if err != nil {
		return err
	}
	if ok {
		return errors.New(Messages.NotFound)
	}
	if e.AccessToken != nil && *e.AccessToken != cmd.Token {
		return errors.New(Messages.InvalidToken)
	}
	return nil
}
