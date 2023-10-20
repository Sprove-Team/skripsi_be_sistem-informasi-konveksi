package global

import (
	"errors"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/timeout"
)

var errorMessage = errors.New("server timeout")

func TimeoutMid(h func(c *fiber.Ctx) error, t *time.Duration) func(c *fiber.Ctx) error {
	if t == nil {
		return timeout.NewWithContext(h, 5*time.Second, errorMessage)
	}
	return timeout.NewWithContext(h, *t, errorMessage)
}
