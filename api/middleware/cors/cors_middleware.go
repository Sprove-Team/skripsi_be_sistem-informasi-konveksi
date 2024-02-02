package cors

import (
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

type CorsMidleware interface {
	Cors() fiber.Handler
}

type corsMidleware struct{}

func NewCorsMiddleware() CorsMidleware {
	return &corsMidleware{}
}

func (c *corsMidleware) Cors() fiber.Handler {
	config := cors.Config{
		AllowHeaders:     "Origin,Content-Type,Accept,Content-Length,Accept-Language,Accept-Encoding,Connection,Access-Control-Allow-Origin,Authorization",
		AllowOrigins:     "*",
		AllowCredentials: true,
		AllowMethods:     "GET,POST,HEAD,PUT,DELETE,PATCH,OPTIONS",
	}
	if os.Getenv("ENVIRONMENT") == "PRODUCTION" {
		config.AllowOrigins = os.Getenv("ALLOW_ORIGINS")
		return cors.New(config)
	}
	return cors.New(config)
}
