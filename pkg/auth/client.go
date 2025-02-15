// Copyright (c) Mainflux
// SPDX-License-Identifier: Apache-2.0

package auth

import (
	"context"

	"github.com/go-redis/redis/v8"
	"github.com/mainflux/mainflux"
)

// Client represents Auth cache.
type Client interface {
	Authorize(ctx context.Context, chanID, thingID string) error
	Identify(ctx context.Context, thingKey string) (string, error)
}

const (
	chanPrefix = "channel"
	keyPrefix  = "thing_key"
)

type client struct {
	redisClient  *redis.Client
	thingsClient mainflux.ThingsServiceClient
}

// New returns redis channel cache implementation.
func New(redisClient *redis.Client, thingsClient mainflux.ThingsServiceClient) Client {
	return client{
		redisClient:  redisClient,
		thingsClient: thingsClient,
	}
}

func (c client) Identify(ctx context.Context, thingKey string) (string, error) {
	tkey := keyPrefix + ":" + thingKey
	thingID, err := c.redisClient.Get(ctx, tkey).Result()
	if err != nil {
		t := &mainflux.Token{
			Value: thingKey,
		}

		thid, err := c.thingsClient.Identify(context.TODO(), t)
		if err != nil {
			return "", err
		}
		return thid.GetValue(), nil
	}
	return thingID, nil
}

func (c client) Authorize(ctx context.Context, chanID, thingID string) error {
	if c.redisClient.SIsMember(ctx, chanPrefix+":"+chanID, thingID).Val() {
		return nil
	}

	ar := &mainflux.AccessByIDReq{
		ThingID: thingID,
		ChanID:  chanID,
	}
	_, err := c.thingsClient.CanAccessByID(ctx, ar)
	return err
}
