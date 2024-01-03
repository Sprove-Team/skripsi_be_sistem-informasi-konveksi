package akuntansi

import (
	"context"

	"github.com/be-sistem-informasi-konveksi/common/message"
	req "github.com/be-sistem-informasi-konveksi/common/request/akuntansi/kelompok_akun"
	"github.com/be-sistem-informasi-konveksi/common/response"

	usecase "github.com/be-sistem-informasi-konveksi/api/usecase/akuntansi/kelompok_akun"
	"github.com/be-sistem-informasi-konveksi/pkg"
	"github.com/gofiber/fiber/v2"
)

type KelompokAkunHandler interface {
	// GetAll(c *fiber.Ctx) error
	Create(c *fiber.Ctx) error
	Update(c *fiber.Ctx) error
}

type kelompokAkunHandler struct {
	uc        usecase.KelompokAkunUsecase
	validator pkg.Validator
}

func NewKelompokAkunHandler(uc usecase.KelompokAkunUsecase, validator pkg.Validator) KelompokAkunHandler {
	return &kelompokAkunHandler{uc, validator}
}

func (h *kelompokAkunHandler) Create(c *fiber.Ctx) error {
	// Parse request body
	req := new(req.Create)
	c.BodyParser(req)

	// Validate request
	errValidate := h.validator.Validate(req)
	if errValidate != nil {
		return c.Status(fiber.StatusBadRequest).JSON(errValidate)
	}

	// Create context
	ctx := c.UserContext()

	// Call usecase to create KelompokAkun
	err := h.uc.Create(ctx, *req)

	// Handle context timeout
	if ctx.Err() == context.DeadlineExceeded {
		return c.Status(fiber.StatusRequestTimeout).JSON(response.ErrorRes(fiber.ErrRequestTimeout.Code, fiber.ErrRequestTimeout.Message, nil))
	}

	// Handle errors
	if err != nil {
		if err.Error() == "duplicated key not allowed" {
			return c.Status(fiber.StatusConflict).JSON(response.ErrorRes(fiber.ErrConflict.Code, fiber.ErrConflict.Message, nil))
		}

		return c.Status(fiber.StatusInternalServerError).JSON(response.ErrorRes(fiber.ErrInternalServerError.Code, fiber.ErrInternalServerError.Message, nil))
	}

	// Respond with success status
	return c.Status(fiber.StatusCreated).JSON(response.SuccessRes(fiber.StatusCreated, message.Created, nil))
}

func (h *kelompokAkunHandler) Update(c *fiber.Ctx) error {
	// Parse request body
	req := new(req.Update)
	c.BodyParser(req)
	c.ParamsParser(req)

	// Validate request
	errValidate := h.validator.Validate(req)
	if errValidate != nil {
		return c.Status(fiber.StatusBadRequest).JSON(errValidate)
	}

	// Create context
	ctx := c.UserContext()

	// Call usecase to create KelompokAkun
	err := h.uc.Update(ctx, *req)

	// Handle context timeout
	if ctx.Err() == context.DeadlineExceeded {
		return c.Status(fiber.StatusRequestTimeout).JSON(response.ErrorRes(fiber.ErrRequestTimeout.Code, fiber.ErrRequestTimeout.Message, nil))
	}

	// Handle errors
	if err != nil {
		if err.Error() == "record not found" {
			return c.Status(fiber.StatusNotFound).JSON(response.ErrorRes(fiber.ErrNotFound.Code, fiber.ErrNotFound.Message, nil))
		}

		return c.Status(fiber.StatusInternalServerError).JSON(response.ErrorRes(fiber.ErrInternalServerError.Code, fiber.ErrInternalServerError.Message, nil))
	}

	// Respond with success status
	return c.Status(fiber.StatusCreated).JSON(response.SuccessRes(fiber.StatusCreated, message.Created, nil))
}
