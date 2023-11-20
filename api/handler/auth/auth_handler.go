package auth

import "github.com/gofiber/fiber/v2"

type AuthHandler interface {
	Login(c *fiber.Ctx) error
	GetNewAccessToken(c *fiber.Ctx) error
}

type authHandler struct{}

// func (h *authHandler) Login(c *fiber.Ctx) error {
// }
//
// func (h *authHandler) GetNewAccessToken(c *fiber.Ctx) error {
// }
