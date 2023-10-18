package direktur

import (
	"log"

	"github.com/gofiber/fiber/v2"

	resGlobal "github.com/be-sistem-informasi-konveksi/common/reponse/global"
	req "github.com/be-sistem-informasi-konveksi/common/request/direktur"
	helper "github.com/be-sistem-informasi-konveksi/helper"
	usecase "github.com/be-sistem-informasi-konveksi/usecase/direktur/produk"
)

type ProdukHandler interface {
	Create(c *fiber.Ctx) error
}

type produkHandler struct {
	uc        usecase.ProdukUsecase
	validator helper.Validator
}

func NewProdukHandler(uc usecase.ProdukUsecase, validator helper.Validator) ProdukHandler {
	return &produkHandler{uc, validator}
}

func (h *produkHandler) Create(c *fiber.Ctx) error {
	c.Accepts("application/json")
	req := new(req.CreateProduk)
	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(resGlobal.ErrorResWithoutData(fiber.StatusBadRequest))
	}
	errValidate := h.validator.Validate(req)

	if len(errValidate) > 0 {
		return c.Status(fiber.StatusBadRequest).JSON(resGlobal.ErrorResWithData(errValidate, fiber.StatusBadRequest))
	}
	err := h.uc.Create(*req)
	if err != nil {
		if err.Error() == "duplicated key not allowed" {
			return c.Status(fiber.StatusConflict).JSON(resGlobal.ErrorResWithoutData(fiber.StatusConflict))
		}
		log.Println(err)
		return c.Status(fiber.StatusInternalServerError).JSON(resGlobal.ErrorResWithoutData(fiber.StatusInternalServerError))
	}
	return c.Status(fiber.StatusCreated).JSON(resGlobal.SuccessResWithoutData("C"))
}
