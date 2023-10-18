package direktur

import (
	"log"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"

	resGlobal "github.com/be-sistem-informasi-konveksi/common/reponse/global"
	req "github.com/be-sistem-informasi-konveksi/common/request/direktur"
	"github.com/be-sistem-informasi-konveksi/helper"
	usecase "github.com/be-sistem-informasi-konveksi/usecase/direktur/produk"
)

type KategoriProdukHandler interface {
	Create(c *fiber.Ctx) error
	Delete(c *fiber.Ctx) error
}

type kategoriProdukHandler struct {
	uc        usecase.KategoriProdukUsecase
	validator helper.Validator
}

func NewKategoriProdukHandler(uc usecase.KategoriProdukUsecase, validator helper.Validator) KategoriProdukHandler {
	return &kategoriProdukHandler{uc, validator}
}

func (h *kategoriProdukHandler) Create(c *fiber.Ctx) error {
	c.Accepts("application/json")
	req := new(req.CreateKategoriProduk)
	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(resGlobal.ErrorResWithoutData(fiber.StatusBadRequest))
	}
	errValidate := h.validator.Validate(req)
	if len(errValidate) > 0 {
		return c.Status(fiber.StatusBadRequest).JSON(resGlobal.ErrorResWithData(errValidate, fiber.StatusBadRequest))
	}
	req.Nama = strings.ToLower(req.Nama)
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

func (h *kategoriProdukHandler) Delete(c *fiber.Ctx) error {
	id := c.Params("id", "0")
	id64, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(resGlobal.ErrorResWithoutData(fiber.StatusInternalServerError))
	}
	err = h.uc.Delete(id64)
	if err != nil && err.Error() == "record not found" {
		return c.Status(fiber.StatusNotFound).JSON(resGlobal.ErrorResWithoutData(fiber.StatusNotFound))
	}
	return c.Status(fiber.StatusOK).JSON(resGlobal.SuccessResWithoutData("D"))
}
