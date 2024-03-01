package sudo

import (
	"context"
	"errors"
)

type CheckCommand struct {
	UserId   string
	DeviceId string
	Token    string
}

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
