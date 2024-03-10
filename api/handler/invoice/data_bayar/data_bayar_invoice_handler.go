package handler_invoice_data_bayar

import (
	"context"
	"time"

	uc_akuntansi_hp "github.com/be-sistem-informasi-konveksi/api/usecase/akuntansi/hutang_piutang"
	uc_invoice_data_bayar "github.com/be-sistem-informasi-konveksi/api/usecase/invoice/data_bayar"
	"github.com/be-sistem-informasi-konveksi/common/message"
	akuntansi "github.com/be-sistem-informasi-konveksi/common/request/akuntansi/hutang_piutang"
	req "github.com/be-sistem-informasi-konveksi/common/request/invoice/data_bayar"
	"github.com/be-sistem-informasi-konveksi/common/response"
	"github.com/be-sistem-informasi-konveksi/entity"
	"github.com/be-sistem-informasi-konveksi/pkg"
	"github.com/gofiber/fiber/v2"
)

type DataBayarInvoiceHandler interface {
	GetByInvoiceId(c *fiber.Ctx) error
	// Create(c *fiber.Ctx) error
	Update(c *fiber.Ctx) error
	// Delete(c *fiber.Ctx) error
}

type dataBayarInvoiceHandler struct {
	uc        uc_invoice_data_bayar.DataBayarInvoice
	uc_hp     uc_akuntansi_hp.HutangPiutangUsecase
	validator pkg.Validator
}

func NewDataBayarInvoiceHandler(
	uc uc_invoice_data_bayar.DataBayarInvoice,
	uc_hp uc_akuntansi_hp.HutangPiutangUsecase,
	validator pkg.Validator,
) DataBayarInvoiceHandler {
	return &dataBayarInvoiceHandler{uc, uc_hp, validator}
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

	badRequest := make([]string, 0, 1)

	switch err.Error() {
	case
		message.CannotModifiedTerkonfirmasiDataBayar:
		badRequest = append(badRequest, err.Error())
	}

	if len(badRequest) > 0 {
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrorRes(fiber.ErrBadRequest.Code, fiber.ErrBadRequest.Message, badRequest))
	}
	return c.Status(fiber.StatusInternalServerError).JSON(response.ErrorRes(fiber.ErrInternalServerError.Code, fiber.ErrInternalServerError.Message, nil))
}

func (h *dataBayarInvoiceHandler) GetByInvoiceId(c *fiber.Ctx) error {
	req := new(req.GetByInvoiceID)

	c.ParamsParser(req)
	c.BodyParser(req)

	errValidate := h.validator.Validate(req)
	if errValidate != nil {
		return c.Status(fiber.StatusBadRequest).JSON(errValidate)
	}

	ctx := c.UserContext()
	datas, err := h.uc.GetByInvoiceID(uc_invoice_data_bayar.ParamGetByInvoiceID{
		Ctx: ctx,
		Req: *req,
	})
	// Handle errors
	if err != nil {
		return errResponse(c, err)
	}

	// Respond with success status
	return c.Status(fiber.StatusOK).JSON(response.SuccessRes(fiber.StatusOK, message.OK, datas))
}

// func (h *dataBayarInvoiceHandler) GetById(c *fiber.Ctx) error {
// 	req := new(reqGlobal.ParamByID)

// 	c.ParamsParser(req)

// 	errValidate := h.validator.Validate(req)
// 	if errValidate != nil {
// 		return c.Status(fiber.StatusBadRequest).JSON(errValidate)
// 	}

// 	ctx := c.UserContext()
// 	datas, err := h.uc.GetById(usecase.ParamGetById{
// 		Ctx: ctx,
// 		ID:  req.ID,
// 	})
// 	// Handle errors
// 	if err != nil {
// 		return errResponse(c, err)
// 	}

// 	// Respond with success status
// 	return c.Status(fiber.StatusOK).JSON(response.SuccessRes(fiber.StatusOK, message.OK, datas))
// }

// func (h *dataBayarInvoiceHandler) Create(c *fiber.Ctx) error {
// 	req := new(req.Create)

// 	c.BodyParser(req)

// 	errValidate := h.validator.Validate(req)
// 	if errValidate != nil {
// 		return c.Status(fiber.StatusBadRequest).JSON(errValidate)
// 	}

// 	ctx := c.UserContext()

// 	dataInvoice, dataReqHp, err := h.uc.CreateDataInvoice(usecase.ParamCreateDataInvoice{
// 		Ctx: ctx,
// 		Req: *req,
// 	})

// 	if err != nil {
// 		return errResponse(c, err)
// 	}

// 	dataHp, err := h.ucHutangPiutang.CreateDataHP(ucHutangPiutang.ParamCreateDataHp{
// 		Ctx: ctx,
// 		Req: *dataReqHp,
// 	})

// 	if err != nil {
// 		return errResponse(c, err)
// 	}

// 	// set data hp into invoice
// 	dataInvoice.HutangPiutang = *dataHp

// 	err = h.uc.CreateCommitDB(usecase.ParamCommitDB{
// 		Ctx:     ctx,
// 		Invoice: dataInvoice,
// 	})

// 	if err != nil {
// 		return errResponse(c, err)
// 	}

// 	// Respond with success status
// 	return c.Status(fiber.StatusCreated).JSON(response.SuccessRes(fiber.StatusCreated, message.Created, nil))
// }

func (h *dataBayarInvoiceHandler) Update(c *fiber.Ctx) error {
	req := new(req.Update)

	c.ParamsParser(req)
	c.BodyParser(req)

	errValidate := h.validator.Validate(req)
	if errValidate != nil {
		return c.Status(fiber.StatusBadRequest).JSON(errValidate)
	}

	ctx := c.UserContext()
	userData := c.Locals("user").(*entity.User)
	dataByrInvoice, err := h.uc.UpdateDataBayarInvoice(uc_invoice_data_bayar.ParamUpdateDataBayarInvoice{
		Ctx:  ctx,
		User: userData,
		Req:  *req,
	})

	if err != nil {
		return errResponse(c, err)
	}
	var paramUpdate = uc_invoice_data_bayar.ParamUpdateCommitDB{
		Ctx:         ctx,
		DataBayar:   dataByrInvoice,
		DataBayarHP: nil,
	}

	if dataByrInvoice.Status == "TERKONFIRMASI" {
		dataHp, err := h.uc_hp.GetHPByInvoiceID(ctx, dataByrInvoice.InvoiceID)
		if err != nil {
			return errResponse(c, err)
		}
		dataByrHP, err := h.uc_hp.CreateDataBayar(ctx, akuntansi.CreateBayar{
			HutangPiutangID: dataHp.ID,
			ReqBayar: akuntansi.ReqBayar{
				Tanggal:         time.Now().Format(time.RFC3339),
				BuktiPembayaran: dataByrInvoice.BuktiPembayaran,
				Keterangan:      dataByrInvoice.Keterangan,
				AkunBayarID:     dataByrInvoice.AkunID,
				Total:           dataByrInvoice.Total,
			},
		})
		if err != nil {
			return errResponse(c, err)
		}
		paramUpdate.DataBayarHP = dataByrHP
	}

	err = h.uc.UpdateCommitDB(paramUpdate)

	if err != nil {
		return errResponse(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(response.SuccessRes(fiber.StatusOK, message.OK, nil))
}

// func (h *dataBayarInvoiceHandler) UpdateStatusProduksi(c *fiber.Ctx) error {
// 	req := new(req.UpdateStatusProduksi)

// 	c.ParamsParser(req)
// 	c.BodyParser(req)

// 	errValidate := h.validator.Validate(req)
// 	if errValidate != nil {
// 		return c.Status(fiber.StatusBadRequest).JSON(errValidate)
// 	}

// 	ctx := c.UserContext()
// 	err := h.uc.UpdateStatusProduksi(usecase.ParamUpdateStatusProduksi{
// 		Ctx: ctx,
// 		Req: *req,
// 	})

// 	if err != nil {
// 		return errResponse(c, err)
// 	}

// 	return c.Status(fiber.StatusOK).JSON(response.SuccessRes(fiber.StatusOK, message.OK, nil))
// }

// func (h *dataBayarInvoiceHandler) Delete(c *fiber.Ctx) error {
// 	req := new(reqGlobal.ParamByID)

// 	c.ParamsParser(req)
// 	c.BodyParser(req)

// 	errValidate := h.validator.Validate(req)
// 	if errValidate != nil {
// 		return c.Status(fiber.StatusBadRequest).JSON(errValidate)
// 	}

// 	ctx := c.UserContext()
// 	err := h.uc.Delete(usecase.ParamDelete{Ctx: ctx, ID: req.ID})

// 	if err != nil {
// 		return errResponse(c, err)
// 	}

// 	return c.Status(fiber.StatusOK).JSON(response.SuccessRes(fiber.StatusOK, message.OK, nil))
// }
