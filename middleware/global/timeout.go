package global

import (
	"errors"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/timeout"
)

var errorMessage = errors.New("server timeout")

func TimeoutMid(t *time.Duration) func(c *fiber.Ctx) error {
	h := func(c *fiber.Ctx) error {
		return c.Next()
	}
	if t == nil {
		return timeout.NewWithContext(h, 10*time.Second, errorMessage)
	}
	return timeout.NewWithContext(h, *t, errorMessage)
}
