package uc_tugas

import (
	"context"
	"errors"
	"time"

	repo_invoice "github.com/be-sistem-informasi-konveksi/api/repository/invoice/mysql/gorm"
	repo_tugas "github.com/be-sistem-informasi-konveksi/api/repository/tugas"
	repo_user "github.com/be-sistem-informasi-konveksi/api/repository/user/mysql/gorm"
	"github.com/be-sistem-informasi-konveksi/common/message"
	req_tugas "github.com/be-sistem-informasi-konveksi/common/request/tugas"
	"github.com/be-sistem-informasi-konveksi/entity"
	"github.com/be-sistem-informasi-konveksi/helper"
	"github.com/be-sistem-informasi-konveksi/pkg"
	"golang.org/x/sync/errgroup"
)

type (
	ParamCreate struct {
		Ctx context.Context
		Req req_tugas.Create
	}
	ParamGetByInvoiceId struct {
		Ctx context.Context
		Req req_tugas.GetByInvoiceId
	}
)

type TugasUsecase interface {
	Create(param ParamCreate) error
	GetByInvoiceId(param ParamGetByInvoiceId) ([]entity.Tugas, error)
}

type tugasUsecase struct {
	repo        repo_tugas.TugasRepo
	userRepo    repo_user.UserRepo
	invoiceRepo repo_invoice.InvoiceRepo
	ulid        pkg.UlidPkg
}

func NewTugasUsecase(repo repo_tugas.TugasRepo,
	userRepo repo_user.UserRepo,
	invoiceRepo repo_invoice.InvoiceRepo,
	ulid pkg.UlidPkg) TugasUsecase {
	return &tugasUsecase{repo, userRepo, invoiceRepo, ulid}
}

func (u *tugasUsecase) Create(param ParamCreate) error {

	g := new(errgroup.Group)
	g.SetLimit(10)
	var users = make([]entity.User, 0, len(param.Req.UserID))
	var count int
	var tanggal_deadline time.Time
	g.Go(func() error {
		var err error
		users, err = u.userRepo.GetUserSpvByIds(repo_user.ParamGetUserSpvByIds{
			Ctx:        param.Ctx,
			JenisSpvID: param.Req.JenisSpvID,
			IDs:        param.Req.UserID,
		})
		if err != nil {
			return err
		}
		return nil
	})

	g.Go(func() error {
		count = helper.CountUniqueElements(param.Req.UserID)
		return nil
	})

	g.Go(func() error {
		var err error
		tanggal_deadline, err = time.Parse(time.DateOnly, param.Req.TanggalDeadline)
		if err != nil {
			helper.LogsError(err)
			return err
		}
		return nil
	})

	if err := g.Wait(); err != nil {
		return err
	}

	if len(users) != count {
		return errors.New(message.UserNotFoundOrNotSpv)
	}

	err := u.repo.Create(repo_tugas.ParamCreate{
		Ctx: param.Ctx,
		Tugas: &entity.Tugas{
			Base: entity.Base{
				ID: u.ulid.MakeUlid().String(),
			},
			InvoiceID:       param.Req.InvoiceID,
			JenisSpvID:      param.Req.JenisSpvID,
			TanggalDeadline: &tanggal_deadline,
			Users:           users,
		},
	})

	if err != nil {
		return err
	}

	return nil
}

func (u *tugasUsecase) GetByInvoiceId(param ParamGetByInvoiceId) ([]entity.Tugas, error) {
	g := new(errgroup.Group)
	g.SetLimit(3)
	var dataTugas []entity.Tugas
	g.Go(func() error {
		var err error
		dataTugas, err = u.repo.GetByInvoiceId(repo_tugas.ParamGetByInvoiceId{
			Ctx:       param.Ctx,
			InvoiceID: param.Req.InvoiceID,
		})
		if err != nil {
			return err
		}
		return nil
	})

	g.Go(func() error {
		if _, err := u.invoiceRepo.GetById(repo_invoice.ParamGetById{
			Ctx: param.Ctx,
			ID:  param.Req.InvoiceID,
		}); err != nil {
			return err
		}
		return nil
	})

	if err := g.Wait(); err != nil {
		return nil, err
	}
	return dataTugas, nil
}
