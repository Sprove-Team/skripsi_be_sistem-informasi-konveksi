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
	Auth(key, role string) fiber.Handler
}

type authMidleware struct{}

func NewAuthMiddleware() AuthMidleware {
	return &authMidleware{}
}

func (a *authMidleware) Auth(keyToken, role string) fiber.Handler {
	return jwtware.New(jwtware.Config{
		Claims: new(helper.Claims),
		SigningKey: jwtware.SigningKey{
			Key: []byte(keyToken),
		},
		SuccessHandler: func(c *fiber.Ctx) error {
			user := c.Locals("user").(*jwt.Token).Claims.(*helper.Claims)
			if !strings.EqualFold(user.Role, role) {
				return c.Status(fiber.StatusUnauthorized).JSON(global.CustomRes(fiber.StatusUnauthorized, message.UnauthUserNotAllowed, nil))
			}
			c.Locals("user", user)
			return c.Next()
		},
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			c.Set("Content-Type", "application/json")
			msgRes := message.UnauthInvalidToken
			if strings.Contains(err.Error(), "expired") {
				msgRes = message.UnauthTokenExpired
			}

			res := global.CustomRes(fiber.StatusUnauthorized, msgRes, nil)
			return c.Status(fiber.StatusUnauthorized).JSON(res)
		},
	})
}
