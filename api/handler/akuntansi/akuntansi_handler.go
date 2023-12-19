package akuntansi

import "github.com/gofiber/fiber/v2"

type AkuntansiHandler interface {
	JurnalUmum(c *fiber.Ctx) error
	NeracaSaldo(c *fiber.Ctx) error
	BukuBesar(c *fiber.Ctx) error
}
