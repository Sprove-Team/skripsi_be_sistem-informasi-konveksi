package transaksi

import (
	"context"

	usecase "github.com/be-sistem-informasi-konveksi/api/usecase/akuntansi/transaksi"
	"github.com/be-sistem-informasi-konveksi/common/message"
	req "github.com/be-sistem-informasi-konveksi/common/request/akuntansi/transaksi"
	"github.com/be-sistem-informasi-konveksi/common/response"
	resGlobal "github.com/be-sistem-informasi-konveksi/common/response/global"
	"github.com/be-sistem-informasi-konveksi/helper"
	"github.com/be-sistem-informasi-konveksi/pkg"
	"github.com/gofiber/fiber/v2"
)

type TransaksiHandler interface {
	Create(c *fiber.Ctx) error
	GetAll(c *fiber.Ctx) error
}

type transaksiHandler struct {
	uc        usecase.TransaksiUsecase
	validator pkg.Validator
}

func NewTransaksiHandler(uc usecase.TransaksiUsecase, validator pkg.Validator) TransaksiHandler {
	return &transaksiHandler{uc, validator}
}

func (h *transaksiHandler) Create(c *fiber.Ctx) error {
	req := new(req.Create)

	if err := c.BodyParser(req); err != nil {
		helper.LogsError(err)
		return c.Status(fiber.StatusBadRequest).JSON(resGlobal.ErrorResWithoutData(fiber.StatusBadRequest))
	}

	// buktiPembayaran, _ := c.FormFile("bukti_pembayaran")
	//
	// if buktiPembayaran != nil && !h.validate.IsImg(buktiPembayaran) {
	// 	return c.Status(fiber.StatusBadRequest).JSON(resGlobal.ErrorResWithData([]response.BaseFormatError{
	// 		{
	// 			FieldName: "bukti_pembayaran",
	// 			Message:   message.InvalidImageFormat,
	// 		},
	// 	}, fiber.StatusBadRequest))
	// }

	// req.BuktiPembayaran = buktiPembayaran

	errValidate := h.validator.Validate(req)
	if len(errValidate) > 0 {
		return c.Status(fiber.StatusBadRequest).JSON(resGlobal.ErrorResWithData(errValidate, fiber.StatusBadRequest))
	}

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

		if err.Error() == message.AkunCannotBeSame {
			return c.Status(fiber.StatusBadRequest).JSON(resGlobal.ErrorResWithData([]response.BaseFormatError{
				{
					FieldName: "ayat_jurnal",
					Message:   err.Error(),
				},
			}, fiber.StatusBadRequest))
		}
		if err.Error() == message.CreditDebitNotSame {
			return c.Status(fiber.StatusBadRequest).JSON(resGlobal.ErrorResWithData([]response.BaseFormatError{
				{
					FieldName: "debit dan kredit",
					Message:   err.Error(),
				},
			}, fiber.StatusBadRequest))
		}
		if err.Error() == message.AkunIdNotFound {
			return c.Status(fiber.StatusBadRequest).JSON(resGlobal.ErrorResWithData([]response.BaseFormatError{
				{
					FieldName: "akun_id",
					Message:   err.Error(),
				},
			}, fiber.StatusBadRequest))
		}
		helper.LogsError(err)
		return c.Status(fiber.StatusInternalServerError).JSON(resGlobal.ErrorResWithoutData(fiber.StatusInternalServerError))
	}

	// Respond with success status
	return c.Status(fiber.StatusCreated).JSON(resGlobal.SuccessResWithoutData("C"))
}

func (h *transaksiHandler) GetAll(c *fiber.Ctx) error {
	reqU := new(req.GetAll)
	c.QueryParser(reqU)

	errValidate := h.validator.Validate(reqU)
	if len(errValidate) > 0 {
		return c.Status(fiber.StatusBadRequest).JSON(resGlobal.ErrorResWithData(errValidate, fiber.StatusBadRequest))
	}

	ctx := c.UserContext()

	transaksi, err := h.uc.GetAll(ctx, *reqU)
	if err != nil {
		helper.LogsError(err)
		return c.Status(fiber.StatusInternalServerError).JSON(resGlobal.ErrorResWithoutData(fiber.StatusInternalServerError))
	}
	data := fiber.Map{
		"transaksi": transaksi,
	}
	return c.Status(fiber.StatusOK).JSON(resGlobal.SuccessResWithData(data, "R"))
}
