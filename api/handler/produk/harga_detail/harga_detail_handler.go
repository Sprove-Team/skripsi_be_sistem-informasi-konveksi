package produk

import (
	"context"

	"github.com/gofiber/fiber/v2"

	usecase "github.com/be-sistem-informasi-konveksi/api/usecase/produk/harga_detail"
	"github.com/be-sistem-informasi-konveksi/common/message"
	reqGlobal "github.com/be-sistem-informasi-konveksi/common/request/global"
	req "github.com/be-sistem-informasi-konveksi/common/request/produk/harga_detail"
	"github.com/be-sistem-informasi-konveksi/common/response"
	"github.com/be-sistem-informasi-konveksi/pkg"
)

type HargaDetailProdukHandler interface {
	CreateByProdukId(c *fiber.Ctx) error
	UpdateByProdukId(c *fiber.Ctx) error
	DeleteByProdukId(c *fiber.Ctx) error
	Delete(c *fiber.Ctx) error
	GetByProdukId(c *fiber.Ctx) error
}

type hargaDetailProdukHandler struct {
	uc        usecase.HargaDetailProdukUsecase
	validator pkg.Validator
}

func NewHargaDetailProdukHandler(uc usecase.HargaDetailProdukUsecase, validator pkg.Validator) HargaDetailProdukHandler {
	return &hargaDetailProdukHandler{uc, validator}
}

func errResponse(c *fiber.Ctx, err error) error {
	if err == context.DeadlineExceeded {
		return c.Status(fiber.StatusRequestTimeout).JSON(response.ErrorRes(fiber.ErrRequestTimeout.Code, fiber.ErrRequestTimeout.Message, nil))
	}

	if err.Error() == "record not found" {
		return c.Status(fiber.StatusNotFound).JSON(response.ErrorRes(fiber.ErrNotFound.Code, fiber.ErrNotFound.Message, nil))
	}

	if err.Error() == "duplicated key not allowed" {
		return c.Status(fiber.StatusConflict).JSON(response.ErrorRes(fiber.ErrConflict.Code, fiber.ErrConflict.Message, nil))
	}

	// badRequest := map[string][]string{}

	// if err.Error() == message.ProdukNotFound {
	// 	badRequest["produk_id"] = []string{err.Error()}
	// }

	// if len(badRequest) > 0 {
	// 	return c.Status(fiber.StatusBadRequest).JSON(response.ErrorRes(fiber.ErrBadRequest.Code, fiber.ErrBadRequest.Message, badRequest))
	// }

	return c.Status(fiber.StatusInternalServerError).JSON(response.ErrorRes(fiber.ErrInternalServerError.Code, fiber.ErrInternalServerError.Message, nil))
}

func (h *hargaDetailProdukHandler) CreateByProdukId(c *fiber.Ctx) error {
	// c.Accepts("application/json")
	req := new(req.CreateByProdukId)
	c.BodyParser(req)
	c.ParamsParser(req)

	errValidate := h.validator.Validate(req)
	if errValidate != nil {
		return c.Status(fiber.StatusBadRequest).JSON(errValidate)
	}
	ctx := c.UserContext()
	err := h.uc.CreateByProdukId(ctx, *req)

	if err != nil {
		return errResponse(c, err)
	}
	return c.Status(fiber.StatusCreated).JSON(response.SuccessRes(fiber.StatusCreated, message.Created, nil))
}

func (h *hargaDetailProdukHandler) UpdateByProdukId(c *fiber.Ctx) error {
	req := new(req.UpdateByProdukId)
	c.ParamsParser(req)
	c.BodyParser(req)
	errValidate := h.validator.Validate(req)
	if errValidate != nil {
		return c.Status(fiber.StatusBadRequest).JSON(errValidate)
	}
	ctx := c.UserContext()
	err := h.uc.UpdateByProdukId(ctx, *req)
	if err != nil {
		return errResponse(c, err)
	}
	return c.Status(fiber.StatusOK).JSON(response.SuccessRes(fiber.StatusOK, message.OK, nil))
}

func (h *hargaDetailProdukHandler) DeleteByProdukId(c *fiber.Ctx) error {
	req := new(req.DeleteByProdukId)
	c.ParamsParser(req)
	errValidate := h.validator.Validate(req)
	if errValidate != nil {
		return c.Status(fiber.StatusBadRequest).JSON(errValidate)
	}
	ctx := c.UserContext()
	err := h.uc.DeleteByProdukId(ctx, req.ProdukId)
	if err != nil {
		return errResponse(c, err)
	}
	return c.Status(fiber.StatusOK).JSON(response.SuccessRes(fiber.StatusOK, message.OK, nil))
}

func (h *hargaDetailProdukHandler) Delete(c *fiber.Ctx) error {
	req := new(reqGlobal.ParamByID)
	c.ParamsParser(req)
	errValidate := h.validator.Validate(req)
	if errValidate != nil {
		return c.Status(fiber.StatusBadRequest).JSON(errValidate)
	}
	ctx := c.UserContext()
	err := h.uc.Delete(ctx, req.ID)
	if err != nil {
		return errResponse(c, err)
	}
	return c.Status(fiber.StatusOK).JSON(response.SuccessRes(fiber.StatusOK, message.OK, nil))
}

func (h *hargaDetailProdukHandler) GetByProdukId(c *fiber.Ctx) error {
	req := new(req.GetByProdukId)
	c.ParamsParser(req)
	errValidate := h.validator.Validate(req)
	if errValidate != nil {
		return c.Status(fiber.StatusBadRequest).JSON(errValidate)
	}
	ctx := c.UserContext()
	data, err := h.uc.GetByProdukId(ctx, *req)
	if err != nil {
		return errResponse(c, err)
	}
	return c.Status(fiber.StatusOK).JSON(response.SuccessRes(fiber.StatusOK, message.OK, data))
}
