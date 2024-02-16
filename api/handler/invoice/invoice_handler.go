package invoice

import (
	"context"

	ucHutangPiutang "github.com/be-sistem-informasi-konveksi/api/usecase/akuntansi/hutang_piutang"
	ucKontak "github.com/be-sistem-informasi-konveksi/api/usecase/akuntansi/kontak"
	usecase "github.com/be-sistem-informasi-konveksi/api/usecase/invoice"
	ucUser "github.com/be-sistem-informasi-konveksi/api/usecase/user"
	"github.com/be-sistem-informasi-konveksi/common/message"
	reqHP "github.com/be-sistem-informasi-konveksi/common/request/akuntansi/hutang_piutang"
	reqKontak "github.com/be-sistem-informasi-konveksi/common/request/akuntansi/kontak"
	req "github.com/be-sistem-informasi-konveksi/common/request/invoice"
	"github.com/be-sistem-informasi-konveksi/common/response"
	"github.com/be-sistem-informasi-konveksi/entity"
	"github.com/be-sistem-informasi-konveksi/pkg"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/sync/errgroup"
)

type InvoiceHandler interface {
	GetAll(c *fiber.Ctx) error
	Create(c *fiber.Ctx) error
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
	// fmt.Println(err)
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
	case message.UserNotFound,
		message.KontakNotFound,
		message.AkunNotFound:
		badRequest = append(badRequest, err.Error())
	}

	if len(badRequest) > 0 {
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrorRes(fiber.ErrBadRequest.Code, fiber.ErrBadRequest.Message, badRequest))
	}
	return c.Status(fiber.StatusInternalServerError).JSON(response.ErrorRes(fiber.ErrInternalServerError.Code, fiber.ErrInternalServerError.Message, nil))
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
	return c.Status(fiber.StatusOK).JSON(response.SuccessRes(fiber.StatusOK, message.OK, datas))
}

func (h *invoiceHandler) Create(c *fiber.Ctx) error {
	req := new(req.Create)

	c.BodyParser(req)

	errValidate := h.validator.Validate(req)
	if errValidate != nil {
		return c.Status(fiber.StatusBadRequest).JSON(errValidate)
	}

	ctx := c.UserContext()

	g := new(errgroup.Group)

	var dataKontak *entity.Kontak

	g.Go(func() error {
		var err error
		if req.KontakID == "" {
			dataKontak, err = h.ucKontak.CreateDataKontak(ucKontak.ParamCreateDataKontak{
				Ctx: ctx,
				Req: reqKontak.Create{
					Nama:       req.NewKontak.Nama,
					NoTelp:     req.NewKontak.NoTelp,
					Alamat:     req.NewKontak.Alamat,
					Keterangan: req.Keterangan,
					Email:      req.NewKontak.Email,
				},
			})

			if err != nil {
				return errResponse(c, err)
			}
		}
		return nil
	})

	var dataInvoice *entity.Invoice
	var dataReqHp *reqHP.Create
	g.Go(func() error {
		var err error
		dataInvoice, dataReqHp, err = h.uc.CreateDataInvoice(usecase.ParamCreateDataInvoice{
			Ctx: ctx,
			Req: *req,
		})

		if err != nil {
			return errResponse(c, err)
		}
		return nil
	})

	if err := g.Wait(); err != nil {
		return err
	}

	dataHp, err := h.ucHutangPiutang.CreateDataHP(ucHutangPiutang.ParamCreateDataHp{
		Ctx: ctx,
		Req: *dataReqHp,
	})
	if err != nil {
		return errResponse(c, err)
	}

	// set data hp into invoice
	dataInvoice.HutangPiutang = *dataHp

	// set kontak into invoice
	if req.KontakID == "" {
		dataInvoice.Kontak = dataKontak
	} else {
		dataInvoice.KontakID = req.KontakID
	}

	err = h.uc.CreateCommitDB(usecase.ParamCreateCommitDB{
		Ctx:     ctx,
		Invoice: dataInvoice,
	})

	if err != nil {
		return errResponse(c, err)
	}

	// Respond with success status
	return c.Status(fiber.StatusCreated).JSON(response.SuccessRes(fiber.StatusCreated, message.Created, nil))
}
