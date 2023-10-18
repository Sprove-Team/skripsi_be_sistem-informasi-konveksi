package direktur

import (
	"github.com/gofiber/fiber/v2"

	resGlobal "github.com/be-sistem-informasi-konveksi/common/reponse/global"
	req "github.com/be-sistem-informasi-konveksi/common/request/direktur"
	"github.com/be-sistem-informasi-konveksi/helper"
	usecase "github.com/be-sistem-informasi-konveksi/usecase/direktur/produk"
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
		return c.Status(fiber.StatusBadRequest).JSON(resGlobal.ErrorResWithoutData(fiber.StatusBadRequest))
	}
	errValidate := h.validator.Validate(req)
	if len(errValidate) > 0 {
		return c.Status(fiber.StatusBadRequest).JSON(resGlobal.ErrorResWithData(errValidate, fiber.StatusBadRequest))
	}
	err := h.uc.Create(*req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(resGlobal.ErrorResWithoutData(fiber.StatusInternalServerError))
	}

	return c.Status(fiber.StatusCreated).JSON(resGlobal.SuccessResWithoutData("C"))
}
