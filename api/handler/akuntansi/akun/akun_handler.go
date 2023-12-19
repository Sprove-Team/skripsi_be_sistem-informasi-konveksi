package akuntansi

import (
	"context"
	"log"

	req "github.com/be-sistem-informasi-konveksi/common/request/akuntansi/akun"
	resGlobal "github.com/be-sistem-informasi-konveksi/common/response/global"

	usecase "github.com/be-sistem-informasi-konveksi/api/usecase/akuntansi/akun"
	"github.com/be-sistem-informasi-konveksi/pkg"
	"github.com/gofiber/fiber/v2"
)

type AkunHandler interface {
	Create(c *fiber.Ctx) error
	GetAll(c *fiber.Ctx) error
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
	c.Accepts("application/json")

	// Parse request body
	req := new(req.Create)
	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(resGlobal.ErrorResWithoutData(fiber.StatusBadRequest))
	}

	// Validate request
	errValidate := h.validator.Validate(req)
	if len(errValidate) > 0 {
		return c.Status(fiber.StatusBadRequest).JSON(resGlobal.ErrorResWithData(errValidate, fiber.StatusBadRequest))
	}

	// Create context
	ctx := c.UserContext()

	// Call usecase to create Akun
	err := h.uc.Create(ctx, *req)

	// Handle context timeout
	if ctx.Err() == context.DeadlineExceeded {
		return c.Status(fiber.StatusRequestTimeout).JSON(resGlobal.ErrorResWithoutData(fiber.StatusRequestTimeout))
	}

	// Handle errors
	if err != nil {
		if err.Error() == "duplicated key not allowed" {
			return c.Status(fiber.StatusConflict).JSON(resGlobal.ErrorResWithoutData(fiber.StatusConflict))
		}

		log.Println(err)
		return c.Status(fiber.StatusInternalServerError).JSON(resGlobal.ErrorResWithoutData(fiber.StatusInternalServerError))
	}

	// Respond with success status
	return c.Status(fiber.StatusCreated).JSON(resGlobal.SuccessResWithoutData("C"))
}

func (h *akunHandler) GetAll(c *fiber.Ctx) error {
	req := new(req.GetAll)
	c.QueryParser(req)

	ctx := c.UserContext()

	data, err := h.uc.GetAll(ctx, *req)

	c.Accepts("application/json")
	if ctx.Err() == context.DeadlineExceeded {
		return c.Status(fiber.StatusRequestTimeout).JSON(resGlobal.ErrorResWithoutData(fiber.StatusRequestTimeout))
	}
	if err != nil {
		if err.Error() == "record not found" {
			return c.Status(fiber.StatusNotFound).JSON(resGlobal.ErrorResWithoutData(fiber.StatusNotFound))
		}
		log.Println(err)
		return c.Status(fiber.StatusInternalServerError).JSON(resGlobal.ErrorResWithoutData(fiber.StatusInternalServerError))
	}

	dataRes := fiber.Map{
		"akun": data,
	}
	return c.Status(fiber.StatusOK).JSON(resGlobal.SuccessResWithData(dataRes, "R"))
}
