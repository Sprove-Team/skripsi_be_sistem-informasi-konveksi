package auth

import (
	"strings"

	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"

	userRepo "github.com/be-sistem-informasi-konveksi/api/repository/user/mysql/gorm"
	"github.com/be-sistem-informasi-konveksi/common/message"
	"github.com/be-sistem-informasi-konveksi/common/response/global"
	resGlobal "github.com/be-sistem-informasi-konveksi/common/response/global"
	"github.com/be-sistem-informasi-konveksi/pkg"
)

type AuthMidleware interface {
	Authorization(key, role string) fiber.Handler
}

type authMidleware struct {
	userRepo userRepo.UserRepo
}

func NewAuthMiddleware(userRepo userRepo.UserRepo) AuthMidleware {
	return &authMidleware{userRepo}
}

func (a *authMidleware) Authorization(keyToken, role string) fiber.Handler {
	return jwtware.New(jwtware.Config{
		Claims: new(pkg.Claims),
		SigningKey: jwtware.SigningKey{
			Key: []byte(keyToken),
		},
		SuccessHandler: func(c *fiber.Ctx) error {
			user := c.Locals("user").(*jwt.Token).Claims.(*pkg.Claims)
			if !strings.EqualFold(user.Role, role) {
				return c.Status(fiber.StatusUnauthorized).JSON(global.CustomRes(fiber.StatusUnauthorized, message.UnauthUserNotAllowed, nil))
			}
			ctx := c.UserContext()
			_, err := a.userRepo.GetById(ctx, user.ID)
			if err != nil {
				if err.Error() == "record not found" {
					return c.Status(fiber.StatusUnauthorized).JSON(global.CustomRes(fiber.StatusUnauthorized, message.UnauthUserNotFound, nil))
				}
				return c.Status(fiber.StatusInternalServerError).JSON(resGlobal.ErrorResWithoutData(fiber.StatusInternalServerError))
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
