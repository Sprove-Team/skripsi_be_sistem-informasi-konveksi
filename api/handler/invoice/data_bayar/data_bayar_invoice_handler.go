package handler_invoice_data_bayar

import (
	"context"
	"fmt"
	"time"

	uc_akuntansi_hp "github.com/be-sistem-informasi-konveksi/api/usecase/akuntansi/hutang_piutang"
	uc_invoice_data_bayar "github.com/be-sistem-informasi-konveksi/api/usecase/invoice/data_bayar"
	"github.com/be-sistem-informasi-konveksi/common/message"
	akuntansi "github.com/be-sistem-informasi-konveksi/common/request/akuntansi/hutang_piutang"
	req_global "github.com/be-sistem-informasi-konveksi/common/request/global"
	req "github.com/be-sistem-informasi-konveksi/common/request/invoice/data_bayar"
	res_global "github.com/be-sistem-informasi-konveksi/common/response"
	"github.com/be-sistem-informasi-konveksi/pkg"
	"github.com/gofiber/fiber/v2"
)

type DataBayarInvoiceHandler interface {
	GetAll(c *fiber.Ctx) error
	GetById(c *fiber.Ctx) error
	GetByInvoiceId(c *fiber.Ctx) error
	CreateByInvoiceId(c *fiber.Ctx) error
	Update(c *fiber.Ctx) error
	Delete(c *fiber.Ctx) error
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
	fmt.Println("err -> ", err.Error())
	if err == context.DeadlineExceeded {
		return c.Status(fiber.StatusRequestTimeout).JSON(res_global.ErrorRes(fiber.ErrRequestTimeout.Code, fiber.ErrRequestTimeout.Message, nil))
	}

	if err.Error() == "record not found" {
		return c.Status(fiber.StatusNotFound).JSON(res_global.ErrorRes(fiber.ErrNotFound.Code, fiber.ErrNotFound.Message, nil))
	}

	if err.Error() == "duplicated key not allowed" {
		return c.Status(fiber.StatusConflict).JSON(res_global.ErrorRes(fiber.ErrConflict.Code, fiber.ErrConflict.Message, nil))
	}

	badRequest := make([]string, 0, 1)

	switch err.Error() {
	case
		message.HutangPiutangNotFound,
		message.InvoiceNotFound,
		message.BayarMustLessThanSisaTagihan,
		message.BayarMustLessThanTotalHargaInvoice,
		message.CannotModifiedTerkonfirmasiDataBayar:
		badRequest = append(badRequest, err.Error())
	}

	if len(badRequest) > 0 {
		return c.Status(fiber.StatusBadRequest).JSON(res_global.ErrorRes(fiber.ErrBadRequest.Code, fiber.ErrBadRequest.Message, badRequest))
	}
	return c.Status(fiber.StatusInternalServerError).JSON(res_global.ErrorRes(fiber.ErrInternalServerError.Code, fiber.ErrInternalServerError.Message, nil))
}

func (h *dataBayarInvoiceHandler) GetAll(c *fiber.Ctx) error {
	req := new(req.GetAll)
	
	c.QueryParser(req)

	errValidate := h.validator.Validate(req)
	if errValidate != nil {
		return c.Status(fiber.StatusBadRequest).JSON(errValidate)
	}

	ctx := c.UserContext()
	datas, err := h.uc.GetAll(uc_invoice_data_bayar.ParamGetAll{
		Ctx: ctx,
		Req: *req,
	})
	// Handle errors
	if err != nil {
		return errResponse(c, err)
	}

	// Respond with success status
	return c.Status(fiber.StatusOK).JSON(res_global.SuccessRes(fiber.StatusOK, message.OK, datas))
}
func (h *dataBayarInvoiceHandler) GetById(c *fiber.Ctx) error {
	req := new(req_global.ParamByID)

	c.ParamsParser(req)
	c.BodyParser(req)

	errValidate := h.validator.Validate(req)
	if errValidate != nil {
		return c.Status(fiber.StatusBadRequest).JSON(errValidate)
	}

	ctx := c.UserContext()
	data, err := h.uc.GetByID(uc_invoice_data_bayar.ParamGetByID{
		Ctx: ctx,
		Req: *req,
	})
	// Handle errors
	if err != nil {
		return errResponse(c, err)
	}

	// Respond with success status
	return c.Status(fiber.StatusOK).JSON(res_global.SuccessRes(fiber.StatusOK, message.OK, data))
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
	return c.Status(fiber.StatusOK).JSON(res_global.SuccessRes(fiber.StatusOK, message.OK, datas))
}

func (h *dataBayarInvoiceHandler) CreateByInvoiceId(c *fiber.Ctx) error {
	req := new(req.CreateByInvoiceId)

	c.ParamsParser(req)
	c.BodyParser(req)

	errValidate := h.validator.Validate(req)
	if errValidate != nil {
		return c.Status(fiber.StatusBadRequest).JSON(errValidate)
	}

	ctx := c.UserContext()

	err := h.uc.CreateByInvoiceID(uc_invoice_data_bayar.ParamCreateByInvoiceID{
		Ctx: ctx,
		Req: *req,
	})

	if err != nil {
		return errResponse(c, err)
	}

	// Respond with success status
	return c.Status(fiber.StatusCreated).JSON(res_global.SuccessRes(fiber.StatusCreated, message.Created, nil))
}

func (h *dataBayarInvoiceHandler) Update(c *fiber.Ctx) error {
	req := new(req.Update)

	c.ParamsParser(req)
	c.BodyParser(req)

	errValidate := h.validator.Validate(req)
	if errValidate != nil {
		return c.Status(fiber.StatusBadRequest).JSON(errValidate)
	}

	ctx := c.UserContext()
	claims := c.Locals("user").(*pkg.Claims)
	dataByrInvoice, err := h.uc.UpdateDataBayarInvoice(uc_invoice_data_bayar.ParamUpdateDataBayarInvoice{
		Ctx:    ctx,
		Claims: claims,
		Req:    *req,
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

	return c.Status(fiber.StatusOK).JSON(res_global.SuccessRes(fiber.StatusOK, message.OK, nil))
}

func (h *dataBayarInvoiceHandler) Delete(c *fiber.Ctx) error {
	req := new(req_global.ParamByID)

	c.ParamsParser(req)

	errValidate := h.validator.Validate(req)
	if errValidate != nil {
		return c.Status(fiber.StatusBadRequest).JSON(errValidate)
	}

	ctx := c.UserContext()
	err := h.uc.Delete(uc_invoice_data_bayar.ParamDelete{
		Ctx: ctx,
		Req: *req,
	})

	if err != nil {
		return errResponse(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(res_global.SuccessRes(fiber.StatusOK, message.OK, nil))
}
