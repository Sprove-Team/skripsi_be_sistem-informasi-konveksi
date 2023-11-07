package auth

import (
	"strings"

	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"

	"github.com/be-sistem-informasi-konveksi/common/message"
	"github.com/be-sistem-informasi-konveksi/common/response/global"
	"github.com/be-sistem-informasi-konveksi/helper"
)

type AuthMidleware interface {
	AuthVerifyToken(key string) func(c *fiber.Ctx) error
	AuthNextIfRole(role string) func(c *fiber.Ctx) error
}

type authMidleware struct{}

func NewAuthMiddleware() AuthMidleware {
	return &authMidleware{}
}

func (a *authMidleware) AuthVerifyToken(key string) func(c *fiber.Ctx) error {
	return jwtware.New(jwtware.Config{
		Claims: new(helper.Claims),
		SigningKey: jwtware.SigningKey{
			Key: []byte(key),
		},
		SuccessHandler: func(c *fiber.Ctx) error {
			user := c.Locals("user").(*jwt.Token).Claims.(*helper.Claims)
			c.Locals("user", user)
			c.Locals("user_role", user.Role)
			return c.Next()
		},
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			c.Set("Content-Type", "application/json")
			msgRes := err.Error()
			if strings.Contains(err.Error(), "token is malformed:") {
				msgRes = message.UnauthInvalidFormatToken
			} else if strings.Contains(err.Error(), "expired") {
				msgRes = message.UnauthTokenExpired
			}
			res := global.CustomRes(fiber.StatusUnauthorized, msgRes, nil)
			return c.Status(fiber.StatusUnauthorized).JSON(res)
		},
	})
}

func (a *authMidleware) AuthNextIfRole(role string) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		user := c.Locals("user").(*helper.Claims)
		if user.Role != role {
			return c.Status(fiber.StatusUnauthorized).JSON(global.CustomRes(fiber.StatusUnauthorized, message.UnauthUserNotAllowed, nil))
		}
		return c.Next()
	}
}
