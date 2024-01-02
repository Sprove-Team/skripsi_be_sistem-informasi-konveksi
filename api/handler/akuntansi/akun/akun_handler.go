package akuntansi

import (
	"context"

	"github.com/be-sistem-informasi-konveksi/common/message"
	req "github.com/be-sistem-informasi-konveksi/common/request/akuntansi/akun"
	"github.com/be-sistem-informasi-konveksi/common/request/global"
	"github.com/be-sistem-informasi-konveksi/common/response"

	usecase "github.com/be-sistem-informasi-konveksi/api/usecase/akuntansi/akun"
	"github.com/be-sistem-informasi-konveksi/pkg"
	"github.com/gofiber/fiber/v2"
)

type AkunHandler interface {
	Create(c *fiber.Ctx) error
	GetAll(c *fiber.Ctx) error
	Update(c *fiber.Ctx) error
	Delete(c *fiber.Ctx) error
	// GetById(c *fiber.Ctx) error
	// Add other methods as needed
}

type akunHandler struct {
	uc        usecase.AkunUsecase
	validator pkg.Validator
}

func NewAkunHandler(uc usecase.AkunUsecase, validator pkg.Validator) AkunHandler {
	return &akunHandler{uc, validator}
}

func (h *akunHandler) Create(c *fiber.Ctx) error {
	// Parse request body
	req := new(req.Create)
	c.BodyParser(req)
	// if err := c.BodyParser(req); err != nil {
	// 	return c.Status(fiber.StatusBadRequest).JSON(resGlobal.ErrorResWithoutData(fiber.StatusBadRequest))
	// }

	// Validate request
	errValidate := h.validator.Validate(req)
	if errValidate != nil {
		return c.Status(fiber.StatusBadRequest).JSON(errValidate)
	}

	// Create context
	ctx := c.UserContext()

	// Call usecase to create Akun
	err := h.uc.Create(ctx, *req)

	// Handle context timeout
	if ctx.Err() == context.DeadlineExceeded {
		return c.Status(fiber.StatusRequestTimeout).JSON(response.ErrorRes(fiber.ErrRequestTimeout.Code, fiber.ErrRequestTimeout.Message, nil))
	}

	// Handle errors
	if err != nil {

		if err.Error() == message.KelompokAkunIdNotFound {
			return c.Status(fiber.StatusBadRequest).JSON(response.ErrorRes(fiber.ErrBadRequest.Code, fiber.ErrBadRequest.Message, map[string][]string{
				"kelompok_akun_id": {message.KelompokAkunIdNotFound},
			}))
		}

		if err.Error() == "duplicated key not allowed" {
			return c.Status(fiber.StatusConflict).JSON(response.ErrorRes(fiber.ErrConflict.Code, fiber.ErrConflict.Message, nil))
		}

		return c.Status(fiber.StatusInternalServerError).JSON(response.ErrorRes(fiber.ErrInternalServerError.Code, fiber.ErrInternalServerError.Message, nil))
	}

	// Respond with success status
	return c.Status(fiber.StatusCreated).JSON(response.SuccessRes(fiber.StatusCreated, message.Created, nil))
}

func (h *akunHandler) Update(c *fiber.Ctx) error {
	// Parse request body
	req := new(req.Update)
	c.ParamsParser(req)
	c.BodyParser(req)

	// Validate request
	errValidate := h.validator.Validate(req)
	if errValidate != nil {
		return c.Status(fiber.StatusBadRequest).JSON(errValidate)
	}

	// Create context
	ctx := c.UserContext()

	// Call usecase to create Akun
	err := h.uc.Update(ctx, *req)

	// Handle context timeout
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

func (h *akunHandler) Delete(c *fiber.Ctx) error {
	// Parse request body
	req := new(global.ParamByID)
	c.ParamsParser(req)

	// Validate request
	errValidate := h.validator.Validate(req)
	if errValidate != nil {
		return c.Status(fiber.StatusBadRequest).JSON(errValidate)
	}

	// Create context
	ctx := c.UserContext()

	// Call usecase to create Akun
	err := h.uc.Delete(ctx, req.ID)

	// Handle context timeout
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

func (h *akunHandler) GetAll(c *fiber.Ctx) error {
	req := new(req.GetAll)
	c.QueryParser(req)

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
