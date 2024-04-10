package middleware_timeout

import (
	"errors"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/timeout"
)

var errorMessage = errors.New("server timeout")

type TimeoutMidleware interface {
	Timeout(t *time.Duration) fiber.Handler
}

type timeoutMidleware struct{}

func NewTimeoutMiddleware() TimeoutMidleware {
	return &timeoutMidleware{}
}

func (t *timeoutMidleware) Timeout(tt *time.Duration) fiber.Handler {
	h := func(c *fiber.Ctx) error {
		return c.Next()
	}
	if tt == nil {
		return timeout.NewWithContext(h, 10*time.Second, errorMessage)
	}
	return timeout.NewWithContext(h, *tt, errorMessage)
}
