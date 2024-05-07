package handler_invoice

import (
	"context"
	"strings"

	ucHutangPiutang "github.com/be-sistem-informasi-konveksi/api/usecase/akuntansi/hutang_piutang"
	ucKontak "github.com/be-sistem-informasi-konveksi/api/usecase/akuntansi/kontak"
	usecase "github.com/be-sistem-informasi-konveksi/api/usecase/invoice"
	ucUser "github.com/be-sistem-informasi-konveksi/api/usecase/user"
	"github.com/be-sistem-informasi-konveksi/common/message"
	reqGlobal "github.com/be-sistem-informasi-konveksi/common/request/global"
	req "github.com/be-sistem-informasi-konveksi/common/request/invoice"
	res_global "github.com/be-sistem-informasi-konveksi/common/response"
	"github.com/be-sistem-informasi-konveksi/entity"
	"github.com/be-sistem-informasi-konveksi/helper"
	"github.com/be-sistem-informasi-konveksi/pkg"
	"github.com/gofiber/fiber/v2"
)

type InvoiceHandler interface {
	GetAll(c *fiber.Ctx) error
	GetById(c *fiber.Ctx) error
	Create(c *fiber.Ctx) error
	Update(c *fiber.Ctx) error
	Delete(c *fiber.Ctx) error
}

type invoiceHandler struct {
	uc              usecase.InvoiceUsecase
	ucUser          ucUser.UserUsecase
	ucKontak        ucKontak.KontakUsecase
	ucHutangPiutang ucHutangPiutang.HutangPiutangUsecase
	validator       pkg.Validator
}

func NewInvoiceHandler(
	uc usecase.InvoiceUsecase,
	ucUser ucUser.UserUsecase,
	ucKontak ucKontak.KontakUsecase,
	ucHutangPiutang ucHutangPiutang.HutangPiutangUsecase,
	validator pkg.Validator,
) InvoiceHandler {
	return &invoiceHandler{uc, ucUser, ucKontak, ucHutangPiutang, validator}
}

func errResponse(c *fiber.Ctx, err error) error {
	// fmt.Println("err -> ", err)
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
	case message.UserNotFound,
		message.KontakNotFound,
		message.SablonNotFound,
		message.ProdukNotFound,
		message.BordirNotFound,
		message.DetailInvoiceNotFound,
		message.BayarMustLessThanTotalHargaInvoice,
		message.AkunNotFound:
		badRequest = append(badRequest, err.Error())
	}

	if strings.Contains(err.Error(), message.UserNotAllowedToModifiedStatusProdusi) {
		badRequest = append(badRequest, err.Error())
	}

	if len(badRequest) > 0 {
		return c.Status(fiber.StatusBadRequest).JSON(res_global.ErrorRes(fiber.ErrBadRequest.Code, fiber.ErrBadRequest.Message, badRequest))
	}

	return c.Status(fiber.StatusInternalServerError).JSON(res_global.ErrorRes(fiber.ErrInternalServerError.Code, fiber.ErrInternalServerError.Message, nil))
}

func (h *invoiceHandler) GetAll(c *fiber.Ctx) error {
	req := new(req.GetAll)

	c.QueryParser(req)

	errValidate := h.validator.Validate(req)
	if errValidate != nil {
		return c.Status(fiber.StatusBadRequest).JSON(errValidate)
	}

	ctx := c.UserContext()
	datas, err := h.uc.GetAll(usecase.ParamGetAll{
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

func (h *invoiceHandler) GetById(c *fiber.Ctx) error {
	req := new(reqGlobal.ParamByID)

	c.ParamsParser(req)

	errValidate := h.validator.Validate(req)
	if errValidate != nil {
		return c.Status(fiber.StatusBadRequest).JSON(errValidate)
	}

	ctx := c.UserContext()
	datas, err := h.uc.GetById(usecase.ParamGetById{
		Ctx: ctx,
		ID:  req.ID,
	})
	// Handle errors
	if err != nil {
		return errResponse(c, err)
	}

	// Respond with success status
	return c.Status(fiber.StatusOK).JSON(res_global.SuccessRes(fiber.StatusOK, message.OK, datas))
}

func (h *invoiceHandler) Create(c *fiber.Ctx) error {
	req := new(req.Create)

	files, err := helper.GetFiles("bukti_pembayaran", c.MultipartForm)
	if err != nil {
		helper.LogsError(err)
		return errResponse(c, err)
	}

	filesGambarDesign, err := helper.GetFiles("gambar_design", c.MultipartForm)
	if err != nil {
		helper.LogsError(err)
		return errResponse(c, err)
	}
	
	req.BuktiPembayaran = files
	req.GambarDesign = filesGambarDesign

	if err := helper.DataParser(c.FormValue("data"), req); err != nil {
		return errResponse(c, err)
	}

	errValidate := h.validator.Validate(req)
	if errValidate != nil {
		return c.Status(fiber.StatusBadRequest).JSON(errValidate)
	}

	ctx := c.UserContext()
	claims := c.Locals("user").(*pkg.Claims)

	dataInvoice, dataReqHp, err := h.uc.CreateDataInvoice(usecase.ParamCreateDataInvoice{
		Ctx:    ctx,
		Req:    *req,
		Claims: claims,
	})

	if err != nil {
		return errResponse(c, err)
	}

	dataHp, err := h.ucHutangPiutang.CreateDataHP(ucHutangPiutang.ParamCreateDataHp{
		Ctx: ctx,
		Req: *dataReqHp,
	})

	if err != nil {
		return errResponse(c, err)
	}

	// set data hp into invoice
	dataInvoice.HutangPiutang = dataHp

	err = h.uc.CreateCommitDB(usecase.ParamCommitDB{
		Ctx:     ctx,
		Invoice: dataInvoice,
	})

	if err != nil {
		return errResponse(c, err)
	}

	// Respond with success status
	return c.Status(fiber.StatusCreated).JSON(res_global.SuccessRes(fiber.StatusCreated, message.Created, nil))
}

func (h *invoiceHandler) Update(c *fiber.Ctx) error {
	req := new(req.Update)

	c.ParamsParser(req)
	c.BodyParser(req)

	errValidate := h.validator.Validate(req)
	if errValidate != nil {
		return c.Status(fiber.StatusBadRequest).JSON(errValidate)
	}

	ctx := c.UserContext()
	claims := c.Locals("user").(*pkg.Claims)

	dataInvoice, err := h.uc.UpdateDataInvoice(usecase.ParamUpdateDataInvoice{
		Ctx:    ctx,
		Req:    *req,
		Claims: claims,
	})

	if err != nil {
		return errResponse(c, err)
	}

	if claims.Role == entity.RolesById[4] {
		err = h.uc.UpdateCommitDB(usecase.ParamCommitDB{
			Ctx:     ctx,
			Invoice: dataInvoice,
		})
	}else {
		err = h.uc.SaveCommitDB(usecase.ParamCommitDB{
			Ctx:     ctx,
			Invoice: dataInvoice,
		})
	}


	if err != nil {
		return errResponse(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(res_global.SuccessRes(fiber.StatusOK, message.OK, nil))
}

func (h *invoiceHandler) Delete(c *fiber.Ctx) error {
	req := new(reqGlobal.ParamByID)

	c.ParamsParser(req)
	c.BodyParser(req)

	errValidate := h.validator.Validate(req)
	if errValidate != nil {
		return c.Status(fiber.StatusBadRequest).JSON(errValidate)
	}

	ctx := c.UserContext()
	err := h.uc.Delete(usecase.ParamDelete{Ctx: ctx, ID: req.ID})

	if err != nil {
		return errResponse(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(res_global.SuccessRes(fiber.StatusOK, message.OK, nil))
}
