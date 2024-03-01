package sudo

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/google/uuid"
)

type StartCommand struct {
	UserId   string
	DeviceId string
	Phone    string
	Email    string
	Locale   string
}

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
