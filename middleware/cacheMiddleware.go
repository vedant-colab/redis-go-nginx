package middleware

import (
	"context"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
)

func CacheMiddleware(client *redis.Client, ttl time.Duration) fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx := context.Background()
		key := fmt.Sprintf("cache:%s", c.OriginalURL())

		// check Cache
		val, err := client.Get(ctx, key).Result()
		if err == nil {
			c.Response().Header.Set("Content-type", "application/json")
			c.Response().SetBodyString(val)
			return nil
		}

		if err := c.Next(); err != nil {
			return err
		}

		if c.Response().StatusCode() == fiber.StatusOK {
			client.Set(ctx, key, c.Response().Body(), ttl)
		}
		return nil
	}
}
