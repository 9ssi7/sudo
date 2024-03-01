// Package sudo provides a service for handling sudo operations.
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

// Config holds the configuration for the sudo service.
type Config struct {
	// Redis is the Redis service.
	Redis mredis.Service

	// NotifySender is the function used to send notifications.
	NotifySender NotifySender

	// Expire is the expiration time for sudo operations.
	Expire time.Duration
}

// Service defines the interface for the sudo service.
type Service interface {
	// Check validates the provided token for a given user and device.
	Check(ctx context.Context, cmd CheckCommand) error

	// Start initiates the sudo operation for a given user and device.
	Start(ctx context.Context, cmd StartCommand) (*string, error)

	// Verify validates the provided code and verify token for a given user and device.
	Verify(ctx context.Context, cmd VerifyCommand) (*string, error)

	// VerifyToken validates the format of a given token.
	VerifyToken(token string) error
}

// client is the implementation of the Service interface.
type client struct {
	redis        mredis.Service
	notifySender NotifySender
	expire       time.Duration
}

// entity represents the data structure stored in Redis for sudo operations.
type entity struct {
	DeviceId    string
	AccessToken *string
	VerifyToken *string
	Code        string
	IsVerified  bool
	TryCount    int
	ExpiresAt   int64
}

// New creates a new instance of the sudo service.
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

// VerifyToken validates the format of a given token.
func (c *client) VerifyToken(token string) error {
	_, err := uuid.Parse(token)
	if err != nil {
		return errors.New(Messages.InvalidToken)
	}
	return nil
}

// calcKey generates a unique key for storing sudo data in Redis.
func (c *client) calcKey(deviceId string, userId string) string {
	return "sudo" + "__" + userId + "__" + deviceId
}

// getByKey retrieves sudo data from Redis using the specified key.
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

// makeCode generates a random verification code.
func (c *client) makeCode() string {
	num := rand.Intn(9999)
	return fmt.Sprintf("%04d", num)
}
