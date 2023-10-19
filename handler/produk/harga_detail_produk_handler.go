package produk

import (
	"log"

	"github.com/gofiber/fiber/v2"

	req "github.com/be-sistem-informasi-konveksi/common/request/produk"
	resGlobal "github.com/be-sistem-informasi-konveksi/common/response/global"
	"github.com/be-sistem-informasi-konveksi/helper"
	usecase "github.com/be-sistem-informasi-konveksi/usecase/produk"
)

type HargaDetailProdukHandler interface {
	Create(c *fiber.Ctx) error
}

type hargaDetailProdukHandler struct {
	uc        usecase.HargaDetailProdukUsecase
	validator helper.Validator
}

func NewHargaDetailProdukHandler(uc usecase.HargaDetailProdukUsecase, validator helper.Validator) HargaDetailProdukHandler {
	return &hargaDetailProdukHandler{uc, validator}
}

func (h *hargaDetailProdukHandler) Create(c *fiber.Ctx) error {
	c.Accepts("application/json")
	req := new(req.CreateHargaDetailProduk)
	if err := c.BodyParser(req); err != nil {
		log.Println(err)
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
		return c.Status(fiber.StatusInternalServerError).JSON(resGlobal.ErrorResWithoutData(fiber.StatusInternalServerError))
	}

	return c.Status(fiber.StatusCreated).JSON(resGlobal.SuccessResWithoutData("C"))
}
