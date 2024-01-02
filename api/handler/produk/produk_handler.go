package produk

import (
	"context"
	"strings"

	"github.com/gofiber/fiber/v2"

	usecase "github.com/be-sistem-informasi-konveksi/api/usecase/produk"
	"github.com/be-sistem-informasi-konveksi/common/message"
	reqGlobal "github.com/be-sistem-informasi-konveksi/common/request/global"
	req "github.com/be-sistem-informasi-konveksi/common/request/produk"
	"github.com/be-sistem-informasi-konveksi/common/response"
	"github.com/be-sistem-informasi-konveksi/pkg"
)

type ProdukHandler interface {
	Delete(c *fiber.Ctx) error
	Create(c *fiber.Ctx) error
	Update(c *fiber.Ctx) error
	GetById(c *fiber.Ctx) error
	GetAll(c *fiber.Ctx) error
}

type produkHandler struct {
	uc        usecase.ProdukUsecase
	validator pkg.Validator
}

func NewProdukHandler(uc usecase.ProdukUsecase, validator pkg.Validator) ProdukHandler {
	return &produkHandler{uc, validator}
}

func (h *produkHandler) Create(c *fiber.Ctx) error {
	req := new(req.Create)
	c.BodyParser(req)

	errValidate := h.validator.Validate(req)
	if errValidate != nil {
		return c.Status(fiber.StatusBadRequest).JSON(errValidate)
	}

	ctx := c.UserContext()
	err := h.uc.Create(ctx, *req)

	if ctx.Err() == context.DeadlineExceeded {
		return c.Status(fiber.StatusRequestTimeout).JSON(response.ErrorRes(fiber.ErrRequestTimeout.Code, fiber.ErrRequestTimeout.Message, nil))
	}

	if err != nil {
		if err.Error() == "duplicated key not allowed" {
			return c.Status(fiber.StatusConflict).JSON(response.ErrorRes(fiber.ErrConflict.Code, fiber.ErrConflict.Message, nil))
		}

		if err.Error() == message.KategoriNotFound {
			return c.Status(fiber.StatusBadRequest).JSON(response.ErrorRes(fiber.ErrBadRequest.Code, fiber.ErrBadRequest.Message, nil))
		}

		return c.Status(fiber.StatusInternalServerError).JSON(response.ErrorRes(fiber.ErrInternalServerError.Code, fiber.ErrInternalServerError.Message, nil))
	}
	return c.Status(fiber.StatusCreated).JSON(response.SuccessRes(fiber.StatusCreated, message.Created, nil))
}

func (h *produkHandler) Update(c *fiber.Ctx) error {
	req := new(req.Update)

	c.ParamsParser(req)
	c.BodyParser(req)

	errValidate := h.validator.Validate(req)

	if errValidate != nil {
		return c.Status(fiber.StatusBadRequest).JSON(errValidate)
	}
	ctx := c.UserContext()
	err := h.uc.Update(ctx, *req)
	if ctx.Err() == context.DeadlineExceeded {
		return c.Status(fiber.StatusRequestTimeout).JSON(response.ErrorRes(fiber.ErrRequestTimeout.Code, fiber.ErrRequestTimeout.Message, nil))
	}

	if err != nil {
		if err.Error() == message.KategoriNotFound {
			return c.Status(fiber.StatusBadRequest).JSON(response.ErrorRes(fiber.ErrBadRequest.Code, fiber.ErrBadRequest.Message, map[string][]string{
				"kategori_id": {message.KategoriNotFound},
			}))
		}

		if err.Error() == "record not found" {
			return c.Status(fiber.StatusNotFound).JSON(response.ErrorRes(fiber.ErrNotFound.Code, fiber.ErrNotFound.Message, nil))
		}
		return c.Status(fiber.StatusInternalServerError).JSON(response.ErrorRes(fiber.ErrInternalServerError.Code, fiber.ErrInternalServerError.Message, nil))
	}

	return c.Status(fiber.StatusOK).JSON(response.SuccessRes(fiber.StatusOK, message.OK, nil))
}

func (h *produkHandler) Delete(c *fiber.Ctx) error {
	req := new(reqGlobal.ParamByID)
	c.ParamsParser(req)

	errValidate := h.validator.Validate(req)
	if errValidate != nil {
		return c.Status(fiber.StatusBadRequest).JSON(errValidate)
	}

	ctx := c.UserContext()
	err := h.uc.Delete(ctx, req.ID)
	if ctx.Err() == context.DeadlineExceeded {
		return c.Status(fiber.StatusRequestTimeout).JSON(response.ErrorRes(fiber.ErrRequestTimeout.Code, fiber.ErrRequestTimeout.Message, nil))
	}

	if err != nil {
		if err.Error() == "record not found" {
			return c.Status(fiber.StatusNotFound).JSON(response.ErrorRes(fiber.ErrNotFound.Code, fiber.ErrNotFound.Message, nil))
		}
		return c.Status(fiber.StatusInternalServerError).JSON(response.ErrorRes(fiber.ErrInternalServerError.Code, fiber.ErrInternalServerError.Message, nil))
	}

	return c.Status(fiber.StatusOK).JSON(response.SuccessRes(fiber.StatusOK, message.OK, nil))
}

func (h *produkHandler) GetById(c *fiber.Ctx) error {
	req := new(reqGlobal.ParamByID)
	c.ParamsParser(req)
	errValidate := h.validator.Validate(req)
	if errValidate != nil {
		return c.Status(fiber.StatusBadRequest).JSON(errValidate)
	}
	ctx := c.UserContext()
	data, err := h.uc.GetById(ctx, req.ID)
	if ctx.Err() == context.DeadlineExceeded {
		return c.Status(fiber.StatusRequestTimeout).JSON(response.ErrorRes(fiber.ErrRequestTimeout.Code, fiber.ErrRequestTimeout.Message, nil))
	}

	if err != nil {
		if err.Error() == "record not found" {
			return c.Status(fiber.StatusNotFound).JSON(response.ErrorRes(fiber.ErrNotFound.Code, fiber.ErrNotFound.Message, nil))
		}
		return c.Status(fiber.StatusInternalServerError).JSON(response.ErrorRes(fiber.ErrInternalServerError.Code, fiber.ErrInternalServerError.Message, nil))
	}

	return c.Status(fiber.StatusOK).JSON(response.SuccessRes(fiber.StatusOK, message.OK, data))
}

func (h *produkHandler) GetAll(c *fiber.Ctx) error {
	reqU := new(req.GetAll)
	c.BodyParser(reqU)
	c.QueryParser(reqU)

	reqU.HargaDetail = strings.ToUpper(reqU.HargaDetail)

	errValidate := h.validator.Validate(reqU)
	if errValidate != nil {
		return c.Status(fiber.StatusBadRequest).JSON(errValidate)
	}

	ctx := c.UserContext()
	datas, err := h.uc.GetAll(ctx, *reqU)
	if ctx.Err() == context.DeadlineExceeded {
		return c.Status(fiber.StatusRequestTimeout).JSON(response.ErrorRes(fiber.ErrRequestTimeout.Code, fiber.ErrRequestTimeout.Message, nil))
	}

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response.ErrorRes(fiber.ErrInternalServerError.Code, fiber.ErrInternalServerError.Message, nil))
	}

	return c.Status(fiber.StatusOK).JSON(response.SuccessRes(fiber.StatusOK, message.OK, datas))
}
