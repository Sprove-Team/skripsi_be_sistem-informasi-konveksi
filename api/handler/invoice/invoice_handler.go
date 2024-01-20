package invoice

import (
	"context"
	"fmt"

	usecase "github.com/be-sistem-informasi-konveksi/api/usecase/invoice"
	"github.com/be-sistem-informasi-konveksi/common/message"
	req "github.com/be-sistem-informasi-konveksi/common/request/invoice"
	"github.com/be-sistem-informasi-konveksi/common/response"
	"github.com/be-sistem-informasi-konveksi/pkg"
	"github.com/gofiber/fiber/v2"
)

type InvoiceHandler interface {
	GetAll(c *fiber.Ctx) error
	Create(c *fiber.Ctx) error
}

type invoiceHandler struct {
	uc        usecase.InvoiceUsecase
	validator pkg.Validator
}

func NewInvoiceHandler(uc usecase.InvoiceUsecase, validator pkg.Validator) InvoiceHandler {
	return &invoiceHandler{uc, validator}
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

	badRequest := map[string][]string{}

	fmt.Println(err.Error())
	if err.Error() == message.UserNotFound {
		badRequest["user_id"] = []string{err.Error()}
	}
	switch err.Error()[2:] {
	case message.BordirNotFound:
		badRequest[fmt.Sprintf("detail_invoice[%c].bordir_id", err.Error()[0])] = []string{err.Error()[2:]}
	case message.ProdukNotFound:
		badRequest[fmt.Sprintf("detail_invoice[%c].produk_id", err.Error()[0])] = []string{err.Error()[2:]}
	case message.SablonNotFound:
		badRequest[fmt.Sprintf("detail_invoice[%c].sablon_id", err.Error()[0])] = []string{err.Error()[2:]}
	case message.HargaDetailProdukNotFound:
		badRequest[fmt.Sprintf("detail_invoice[%c].harga_produk_id", err.Error()[0])] = []string{err.Error()[2:]}
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
	// Call usecase to create KelompokAkun
	datas, err := h.uc.GetAll(ctx, *req)
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
	// Call usecase to create KelompokAkun
	err := h.uc.Create(ctx, *req)
	// Handle errors
	if err != nil {
		return errResponse(c, err)
	}

	// Respond with success status
	return c.Status(fiber.StatusCreated).JSON(response.SuccessRes(fiber.StatusCreated, message.Created, nil))
}
