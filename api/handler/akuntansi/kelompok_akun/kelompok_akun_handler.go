package akuntansi

import (
	"context"
	"log"

	req "github.com/be-sistem-informasi-konveksi/common/request/akuntansi/kelompok_akun"
	resGlobal "github.com/be-sistem-informasi-konveksi/common/response/global"

	usecase "github.com/be-sistem-informasi-konveksi/api/usecase/akuntansi/kelompok_akun"
	"github.com/be-sistem-informasi-konveksi/pkg"
	"github.com/gofiber/fiber/v2"
)

type KelompokAkunHandler interface {
	Create(c *fiber.Ctx) error
}

type kelompokAkunHandler struct {
	uc        usecase.KelompokAkunUsecase
	validator pkg.Validator
}

func NewKelompokAkunHandler(uc usecase.KelompokAkunUsecase, validator pkg.Validator) KelompokAkunHandler {
	return &kelompokAkunHandler{uc, validator}
}

func (h *kelompokAkunHandler) Create(c *fiber.Ctx) error {
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

	// Call usecase to create KelompokAkun
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
