package sudo

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/google/uuid"
)

type VerifyCommand struct {
	UserId      string
	DeviceId    string
	VerifyToken string
	Code        string
}

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
