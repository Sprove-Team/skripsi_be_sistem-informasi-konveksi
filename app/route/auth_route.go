package route

import (
	"github.com/gofiber/fiber/v2"

	"github.com/be-sistem-informasi-konveksi/app/handler_init"
)

type AuthRoute interface {
	Auth(router fiber.Router)
}

type authRoute struct {
	h handler_init.AuthHandlerInit
}

func NewAuthRoute(h handler_init.AuthHandlerInit) AuthRoute {
	return &authRoute{h}
}

func (ro *authRoute) Auth(router fiber.Router) {
	router.Post("/login", ro.h.Auth().Login)
	router.Post("/refresh_token", ro.h.Auth().RefreshToken)
}
