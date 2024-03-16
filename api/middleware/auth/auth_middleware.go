package middleware_auth

import (
	"context"
	"os"
	"strings"

	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"

	repo_user "github.com/be-sistem-informasi-konveksi/api/repository/user/mysql/gorm"
	"github.com/be-sistem-informasi-konveksi/common/message"
	res_global "github.com/be-sistem-informasi-konveksi/common/response"
	"github.com/be-sistem-informasi-konveksi/helper"
	"github.com/be-sistem-informasi-konveksi/pkg"
)

type AuthMidleware interface {
	Authorization(roles []string) fiber.Handler
}

type authMidleware struct {
	repoUser repo_user.UserRepo
}

func NewAuthMiddleware(repoUser repo_user.UserRepo) AuthMidleware {
	return &authMidleware{repoUser}
}

func (a *authMidleware) Authorization(roles []string) fiber.Handler {
	if len(roles) <= 0 {
		roles = []string{"DIREKTUR", "BENDAHARA", "ADMIN", "MANAJER_PRODUKSI", "SUPERVISOR"}
	}

	return jwtware.New(jwtware.Config{
		Claims: new(pkg.Claims),
		SigningKey: jwtware.SigningKey{
			Key: []byte(os.Getenv("JWT_TOKEN")),
		},
		SuccessHandler: func(c *fiber.Ctx) error {
			user, ok := c.Locals("user").(*jwt.Token).Claims.(*pkg.Claims)
			if !ok {
				return c.Status(fiber.StatusUnauthorized).JSON(res_global.ErrorRes(fiber.ErrUnauthorized.Code, fiber.ErrUnauthorized.Message, []string{message.UnauthInvalidToken}))
			}
			var pass bool
			for _, role := range roles {
				if strings.EqualFold(user.Role, role) {
					pass = true
					break
				}
			}

			if !pass {
				return c.Status(fiber.StatusUnauthorized).JSON(res_global.ErrorRes(fiber.ErrUnauthorized.Code, fiber.ErrUnauthorized.Message, nil))
			}

			ctx := c.UserContext()
			dataUser, err := a.repoUser.GetById(repo_user.ParamGetById{
				Ctx: ctx,
				ID:  user.ID,
			})

			if err != nil {
				if err == context.DeadlineExceeded {
					return c.Status(fiber.StatusRequestTimeout).JSON(res_global.ErrorRes(fiber.ErrRequestTimeout.Code, fiber.ErrRequestTimeout.Message, nil))
				}
				if err.Error() == "record not found" {
					return c.Status(fiber.StatusUnauthorized).JSON(res_global.ErrorRes(fiber.ErrUnauthorized.Code, fiber.ErrUnauthorized.Message, []string{message.UnauthInvalidToken}))
				}
				return c.Status(fiber.StatusInternalServerError).JSON(res_global.ErrorRes(fiber.ErrInternalServerError.Code, fiber.ErrInternalServerError.Message, nil))
			}

			if user.Nama != dataUser.Nama ||
				user.Username != dataUser.Username ||
				user.Role != dataUser.Role {
				return c.Status(fiber.StatusUnauthorized).JSON(res_global.ErrorRes(fiber.ErrUnauthorized.Code, fiber.ErrUnauthorized.Message, []string{message.UnauthInvalidToken}))
			}

			c.Locals("user", user)
			return c.Next()
		},
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			helper.LogsError(err)
			if strings.Contains(err.Error(), "token is expired") {
				return c.Status(fiber.StatusUnauthorized).JSON(res_global.ErrorRes(fiber.ErrUnauthorized.Code, fiber.ErrUnauthorized.Message, []string{message.UnauthTokenExpired}))
			}
			return c.Status(fiber.StatusUnauthorized).JSON(res_global.ErrorRes(fiber.ErrUnauthorized.Code, fiber.ErrUnauthorized.Message, []string{message.UnauthInvalidToken}))
		},
	})
}
