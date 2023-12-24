package akuntansi

import (
	"context"
	"log"

	"github.com/be-sistem-informasi-konveksi/common/message"
	req "github.com/be-sistem-informasi-konveksi/common/request/akuntansi/golongan_akun"
	"github.com/be-sistem-informasi-konveksi/common/response"
	resGlobal "github.com/be-sistem-informasi-konveksi/common/response/global"

	usecase "github.com/be-sistem-informasi-konveksi/api/usecase/akuntansi/golongan_akun"
	"github.com/be-sistem-informasi-konveksi/pkg"
	"github.com/gofiber/fiber/v2"
)

type GolonganAkunHandler interface {
	Create(c *fiber.Ctx) error
}

type golonganAkunHandler struct {
	uc        usecase.GolonganAkunUsecase
	validator pkg.Validator
}

func NewGolonganAkunHandler(uc usecase.GolonganAkunUsecase, validator pkg.Validator) GolonganAkunHandler {
	return &golonganAkunHandler{uc, validator}
}

func (h *golonganAkunHandler) Create(c *fiber.Ctx) error {
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

	// Call usecase to create GolonganAkun
	err := h.uc.Create(ctx, *req)

	// Handle context timeout
	if ctx.Err() == context.DeadlineExceeded {
		return c.Status(fiber.StatusRequestTimeout).JSON(resGlobal.ErrorResWithoutData(fiber.StatusRequestTimeout))
	}

	// Handle errors
	if err != nil {
		if err.Error() == message.KelompokAkunIdNotFound {
			return c.Status(fiber.StatusBadRequest).JSON(resGlobal.ErrorResWithData([]response.BaseFormatError{
				{
					FieldName: "kelompok_akun_id",
					Message:   err.Error(),
				},
			}, fiber.StatusBadRequest))
		}
		if err.Error() == "duplicated key not allowed" {
			return c.Status(fiber.StatusConflict).JSON(resGlobal.ErrorResWithoutData(fiber.StatusConflict))
		}

		log.Println(err)
		return c.Status(fiber.StatusInternalServerError).JSON(resGlobal.ErrorResWithoutData(fiber.StatusInternalServerError))
	}

	// Respond with success status
	return c.Status(fiber.StatusCreated).JSON(resGlobal.SuccessResWithoutData("C"))
}
