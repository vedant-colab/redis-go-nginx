package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"fmt"
	"redisgo/db"
	"redisgo/middleware"
	"redisgo/service"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
)

func Logger(ip string, logger *log.Logger) {
	logLine := fmt.Sprintf("BLOCKED - %v - %s\n", time.Now(), ip)
	fmt.Print(logLine)
	logger.Print(logLine)
}

func main() {
	app := fiber.New()

	app.Use(logger.New())

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file", err)
	}

	logFile, _ := os.OpenFile("blocked_ips.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	logger := log.New(logFile, "", log.LstdFlags)

	client := redis.NewClient(&redis.Options{
		Addr:     "redis:6379",
		Password: "",
		DB:       0,
		Protocol: 2,
	})

	ctx := context.Background()

	app.Use(limiter.New(limiter.Config{
		Next: func(c *fiber.Ctx) bool {
			return c.Path() == "/health"
		},
		Expiration: 20 * time.Second,
	}))

	app.Use(func(c *fiber.Ctx) error {
		ip := c.IP()
		key := fmt.Sprintf("rate<%s>", ip)

		count, err := client.Incr(ctx, key).Result()
		if err != nil {
			return c.Status(500).SendString("Internal Server Error")
		}
		if count == 1 {
			client.Expire(ctx, key, 30*time.Second)
		}
		if count > 5 {
			Logger(ip, logger)
			return c.Status(429).SendString("Too Many Requests")
		}
		return c.Next()
	})

	err = client.Set(ctx, "foo", "bar", 0).Err()
	if err != nil {
		panic(err)
	}

	val, err := client.Get(ctx, "foo").Result()
	if err != nil {
		panic(err)
	}
	fmt.Println("foo: ", val)

	app.Get("/home", func(c *fiber.Ctx) error {
		return c.SendString("This is golang")
	})

	app.Get("/test", func(c *fiber.Ctx) error {
		sum := 0
		for i := range 500_000 {
			sum += i
		}
		fmt.Println("Handled by: ", os.Getenv("HOSTNAME"))
		return c.SendString(fmt.Sprintf("Done by %s", os.Getenv("HOSTNAME")))
	})

	app.Get("/health", func(c *fiber.Ctx) error {
		return c.SendString("This api monitors health of server")
	})

	app.Get("/:id", middleware.CacheMiddleware(client, 30*time.Second), func(c *fiber.Ctx) error {
		strId := c.Params("id")

		ctx := c.UserContext()
		if ctx == nil {
			ctx = context.Background()
		}

		id, convErr := strconv.Atoi(strId)
		if convErr != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "invalid ID",
			})
		}

		data, dbErr := service.FetchStudentById(ctx, id)
		if dbErr != nil {
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
				"error": dbErr.Error(),
			})
		}

		return c.Status(http.StatusOK).JSON(data)
	})

	db.ConnectDB()
	app.Listen("0.0.0.0:8083")

}
