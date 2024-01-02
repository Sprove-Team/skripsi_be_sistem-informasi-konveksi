package bordir

import (
	"context"

	"github.com/gofiber/fiber/v2"

	usecase "github.com/be-sistem-informasi-konveksi/api/usecase/bordir"
	"github.com/be-sistem-informasi-konveksi/common/message"
	req "github.com/be-sistem-informasi-konveksi/common/request/bordir"
	reqGlobal "github.com/be-sistem-informasi-konveksi/common/request/global"
	"github.com/be-sistem-informasi-konveksi/common/response"
	"github.com/be-sistem-informasi-konveksi/pkg"
)

type BordirHandler interface {
	Delete(c *fiber.Ctx) error
	Create(c *fiber.Ctx) error
	Update(c *fiber.Ctx) error
	GetById(c *fiber.Ctx) error
	GetAll(c *fiber.Ctx) error
}

type bordirHandler struct {
	uc        usecase.BordirUsecase
	validator pkg.Validator
}

func NewBordirHandler(uc usecase.BordirUsecase, validator pkg.Validator) BordirHandler {
	return &bordirHandler{uc, validator}
}

func (h *bordirHandler) Create(c *fiber.Ctx) error {
	req := new(req.Create)
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

		return c.Status(fiber.StatusInternalServerError).JSON(response.ErrorRes(fiber.ErrInternalServerError.Code, fiber.ErrInternalServerError.Message, nil))
	}
	return c.Status(fiber.StatusCreated).JSON(response.SuccessRes(fiber.StatusCreated, message.Created, nil))
}

func (h *bordirHandler) Update(c *fiber.Ctx) error {
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
		if err.Error() == "record not found" {
			return c.Status(fiber.StatusNotFound).JSON(response.ErrorRes(fiber.ErrNotFound.Code, fiber.ErrNotFound.Message, nil))
		}
		return c.Status(fiber.StatusInternalServerError).JSON(response.ErrorRes(fiber.ErrInternalServerError.Code, fiber.ErrInternalServerError.Message, nil))
	}

	return c.Status(fiber.StatusOK).JSON(response.SuccessRes(fiber.StatusOK, message.OK, nil))
}

func (h *bordirHandler) Delete(c *fiber.Ctx) error {
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

func (h *bordirHandler) GetById(c *fiber.Ctx) error {
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

func (h *bordirHandler) GetAll(c *fiber.Ctx) error {
	req := new(req.GetAll)
	c.BodyParser(req)
	c.QueryParser(req)

	errValidate := h.validator.Validate(req)
	if errValidate != nil {
		return c.Status(fiber.StatusBadRequest).JSON(errValidate)
	}

	ctx := c.UserContext()
	data, err := h.uc.GetAll(ctx, *req)

	if ctx.Err() == context.DeadlineExceeded {
		return c.Status(fiber.StatusRequestTimeout).JSON(response.ErrorRes(fiber.ErrRequestTimeout.Code, fiber.ErrRequestTimeout.Message, nil))
	}
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response.ErrorRes(fiber.ErrInternalServerError.Code, fiber.ErrInternalServerError.Message, nil))
	}

	return c.Status(fiber.StatusOK).JSON(response.SuccessRes(fiber.StatusOK, message.OK, data))
}
