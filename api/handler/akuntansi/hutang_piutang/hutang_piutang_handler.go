package akuntansi

import (
	"fmt"

	usecase "github.com/be-sistem-informasi-konveksi/api/usecase/akuntansi/hutang_piutang"
	"github.com/be-sistem-informasi-konveksi/common/message"
	req "github.com/be-sistem-informasi-konveksi/common/request/akuntansi/hutang_piutang"
	"github.com/be-sistem-informasi-konveksi/common/response"
	"github.com/be-sistem-informasi-konveksi/pkg"
	"github.com/gofiber/fiber/v2"
)

type HutangPiutangHandler interface {
	Create(c *fiber.Ctx) error
	CreateBayar(c *fiber.Ctx) error
	GetAll(c *fiber.Ctx) error
}

type hutangPiutangHandler struct {
	uc        usecase.HutangPiutangUsecase
	validator pkg.Validator
	errResponse
}

func NewHutangPiutangHandler(uc usecase.HutangPiutangUsecase, validator pkg.Validator) HutangPiutangHandler {
	return &hutangPiutangHandler{uc, validator, errResponse{}}
}

func (h *hutangPiutangHandler) Create(c *fiber.Ctx) error {

	req := new(req.Create)

	c.BodyParser(req)
	// fmt.Println("ket -> ", req.Keterangan)

	errValidate := h.validator.Validate(req)
	if errValidate != nil {
		return c.Status(fiber.StatusBadRequest).JSON(errValidate)
	}
	ctx := c.UserContext()

	dataHP, err := h.uc.CreateDataHP(usecase.ParamCreateDataHp{
		Ctx: ctx,
		Req: *req,
	})

	if err != nil {
		return h.errHP(c, err)
	}

	err = h.uc.CreateCommitDB(ctx, dataHP)
	// Handle errors
	if err != nil {
		return h.errHP(c, err)
	}

	// Respond with success status
	return c.Status(fiber.StatusCreated).JSON(response.SuccessRes(fiber.StatusCreated, message.Created, nil))
}

func (h *hutangPiutangHandler) CreateBayar(c *fiber.Ctx) error {

	req := new(req.CreateBayar)

	c.BodyParser(req)
	c.ParamsParser(req)

	errValidate := h.validator.Validate(req)
	if errValidate != nil {
		return c.Status(fiber.StatusBadRequest).JSON(errValidate)
	}
	fmt.Println(req.HutangPiutangID)
	ctx := c.UserContext()
	err := h.uc.CreateBayar(ctx, *req)
	// Handle errors
	if err != nil {
		return h.errBayar(c, err)
	}

	// Respond with success status
	return c.Status(fiber.StatusCreated).JSON(response.SuccessRes(fiber.StatusCreated, message.Created, nil))
}
func (h *hutangPiutangHandler) GetAll(c *fiber.Ctx) error {
	req := new(req.GetAll)

	c.QueryParser(req)

	errValidate := h.validator.Validate(req)
	if errValidate != nil {
		return c.Status(fiber.StatusBadRequest).JSON(errValidate)
	}
	ctx := c.UserContext()
	data, err := h.uc.GetAll(ctx, *req)
	// Handle errors
	if err != nil {
		return h.errHP(c, err)
	}

	// Respond with success status
	return c.Status(fiber.StatusOK).JSON(response.SuccessRes(fiber.StatusOK, message.OK, data))
}
