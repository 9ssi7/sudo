// Package sudo provides a service for handling sudo operations.
package sudo

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/google/uuid"
)

// VerifyCommand represents the command structure for the Verify operation.
type VerifyCommand struct {
	// UserId is the unique identifier of the user.
	UserId string

	// DeviceId is the unique identifier of the device.
	DeviceId string

	// VerifyToken is the token to be verified.
	VerifyToken string

	// Code is the verification code to be verified.
	Code string
}

// Verify validates the provided code and verify token for a given user and device.
func (c *client) Verify(ctx context.Context, cmd VerifyCommand) (*string, error) {
	e, ok, err := c.getByKey(ctx, c.calcKey(cmd.DeviceId, cmd.UserId))
	if err != nil {
		return nil, err
	}
	if ok {
		return nil, errors.New(Messages.NotFound)
	}
	if e.VerifyToken == nil || *e.VerifyToken != cmd.VerifyToken {
		return nil, errors.New(Messages.InvalidToken)
	}
	if e.ExpiresAt < time.Now().Unix() {
		_ = c.redis.Del(ctx, c.calcKey(cmd.DeviceId, cmd.UserId))
		return nil, errors.New(Messages.ExpiredCode)
	}
	if e.TryCount > 3 {
		_ = c.redis.Del(ctx, c.calcKey(cmd.DeviceId, cmd.UserId))
		return nil, errors.New(Messages.ExceedTryCount)
	}
	if e.Code != cmd.Code {
		e.TryCount++
		b, err := json.Marshal(e)
		if err != nil {
			return nil, errors.New(Messages.FailedMarshal)
		}
		if _err := c.redis.Set(ctx, c.calcKey(cmd.DeviceId, cmd.UserId), b); _err != nil {
			return nil, errors.New(Messages.FailedRedisSet)
		}
		return nil, errors.New(Messages.InvalidCode)
	}
	tkn := uuid.New().String()
	e.AccessToken = &tkn
	e.IsVerified = true
	b, _err := json.Marshal(e)
	if _err != nil {
		return nil, errors.New(Messages.FailedMarshal)
	}
	if _err = c.redis.Set(ctx, c.calcKey(cmd.DeviceId, cmd.UserId), b); _err != nil {
		return nil, errors.New(Messages.FailedRedisSet)
	}
	return &tkn, nil
}
