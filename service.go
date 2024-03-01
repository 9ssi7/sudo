package sudo

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"time"

	"github.com/cilloparch/cillop/db/mredis"
	"github.com/google/uuid"
)

type Config struct {
	Redis        mredis.Service
	NotifySender NotifySender

	Expire time.Duration
}

type Service interface {
	Check(ctx context.Context, cmd CheckCommand) error
	Start(ctx context.Context, cmd StartCommand) (*string, error)
	Verify(ctx context.Context, cmd VerifyCommand) (*string, error)
	VerifyToken(token string) error
}

type client struct {
	redis        mredis.Service
	notifySender NotifySender
	expire       time.Duration
}

type entity struct {
	DeviceId    string
	AccessToken *string
	VerifyToken *string
	Code        string
	IsVerified  bool
	TryCount    int
	ExpiresAt   int64
}

func New(cnf Config) Service {
	if cnf.Expire == 0 {
		cnf.Expire = 5 * time.Minute
	}
	return &client{
		redis:        cnf.Redis,
		notifySender: cnf.NotifySender,
		expire:       cnf.Expire,
	}
}

func (c *client) VerifyToken(token string) error {
	_, err := uuid.Parse(token)
	if err != nil {
		return errors.New(Messages.InvalidToken)
	}
	return nil
}

func (c *client) calcKey(deviceId string, userId string) string {
	return "sudo" + "__" + userId + "__" + deviceId
}

func (c *client) getByKey(ctx context.Context, key string) (*entity, bool, error) {
	res, err := c.redis.Get(ctx, key)
	if err != nil {
		return nil, true, errors.New(Messages.FailedRedisFetch)
	}
	var e entity
	if err := json.Unmarshal([]byte(res), &e); err != nil {
		return nil, false, errors.New(Messages.FailedMarshal)
	}
	return &e, false, nil
}

func (c *client) makeCode() string {
	num := rand.Intn(9999)
	return fmt.Sprintf("%04d", num)
}
