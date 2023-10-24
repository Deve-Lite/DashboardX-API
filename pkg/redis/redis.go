package redis

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/Deve-Lite/DashboardX-API/config"
	"github.com/redis/go-redis/v9"
)

func NewDB(c *config.Config) *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", c.Redis.Host, c.Redis.Port),
		Password: c.Redis.Password,
		DB:       int(c.Redis.Database),
	})

	ctx := context.Background()
	defer ctx.Done()

	key := "healthcheck"
	val := "some-value"
	client.Set(ctx, key, val, time.Minute*5)

	r := client.Get(ctx, key)
	if r.Val() != val {
		log.Panic("Can not connect to Redis.")
	}

	client.Del(ctx, key)

	log.Print("Redis connected successfully")
	return client
}

func FlushDB(c *config.Config) error {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", c.Redis.Host, c.Redis.Port),
		Password: c.Redis.Password,
		DB:       int(c.Redis.Database),
	})

	ctx := context.Background()
	defer ctx.Done()

	if err := client.FlushDB(ctx).Err(); err != nil {
		return err
	}

	return nil
}
