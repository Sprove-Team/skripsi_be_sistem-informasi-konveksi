package auth

import (
	"strings"

	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"

	userRepo "github.com/be-sistem-informasi-konveksi/api/repository/user/mysql/gorm"
	"github.com/be-sistem-informasi-konveksi/common/response"
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
				return c.Status(fiber.StatusUnauthorized).JSON(response.ErrorRes(fiber.ErrUnauthorized.Code, fiber.ErrUnauthorized.Message, nil))
			}
			ctx := c.UserContext()
			_, err := a.userRepo.GetById(ctx, user.ID)
			if err != nil {
				if err.Error() == "record not found" {
					return c.Status(fiber.StatusUnauthorized).JSON(response.ErrorRes(fiber.ErrUnauthorized.Code, fiber.ErrUnauthorized.Message, nil))
				}
				return c.Status(fiber.StatusInternalServerError).JSON(response.ErrorRes(fiber.ErrInternalServerError.Code, fiber.ErrInternalServerError.Message, nil))
			}

			c.Locals("user", user)
			return c.Next()
		},
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			return c.Status(fiber.StatusUnauthorized).JSON(response.ErrorRes(fiber.ErrUnauthorized.Code, fiber.ErrUnauthorized.Message, nil))
		},
	})
}
