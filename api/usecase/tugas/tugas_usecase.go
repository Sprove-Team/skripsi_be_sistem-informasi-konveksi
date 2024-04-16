package uc_tugas

import (
	"context"
	"errors"
	"time"

	repo_invoice "github.com/be-sistem-informasi-konveksi/api/repository/invoice/mysql/gorm"
	repo_tugas "github.com/be-sistem-informasi-konveksi/api/repository/tugas"
	repo_user "github.com/be-sistem-informasi-konveksi/api/repository/user/mysql/gorm"
	"github.com/be-sistem-informasi-konveksi/common/message"
	req_global "github.com/be-sistem-informasi-konveksi/common/request/global"
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
	ParamUpdate struct {
		Ctx context.Context
		Req req_tugas.Update
	}
	ParamDelete struct {
		Ctx context.Context
		Req req_global.ParamByID
	}
	ParamGetAll struct {
		Ctx context.Context
		Req req_tugas.GetAll
	}
	ParamGetById struct {
		Ctx context.Context
		Req req_global.ParamByID
	}
	ParamGetByInvoiceId struct {
		Ctx context.Context
		Req req_tugas.GetByInvoiceId
	}
)

type TugasUsecase interface {
	Create(param ParamCreate) error
	Update(param ParamUpdate) error
	Delete(param ParamDelete) error
	GetAll(param ParamGetAll) ([]entity.Tugas, error)
	GetById(param ParamGetById) (*entity.Tugas, error)
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

func (u *tugasUsecase) checkAndGetUsers(
	g *errgroup.Group,
	ctx context.Context,
	jenisSpvID string,
	userID []string,
	outUser chan<- []entity.User,
	outCount chan<- int,
) {
	g.Go(func() error {
		users, err := u.userRepo.GetUserSpvByIds(repo_user.ParamGetUserSpvByIds{
			Ctx:        ctx,
			JenisSpvID: jenisSpvID,
			IDs:        userID,
		})
		if err != nil {
			return err
		}
		outUser <- users
		return nil
	})
	g.Go(func() error {
		outCount <- helper.CountUniqueElements(userID)
		return nil
	})
}

func (u *tugasUsecase) Create(param ParamCreate) error {

	g := new(errgroup.Group)
	g.SetLimit(10)

	var tanggal_deadline time.Time
	usersChan := make(chan []entity.User, 1)
	countChan := make(chan int, 1)

	u.checkAndGetUsers(g, param.Ctx, param.Req.JenisSpvID, param.Req.UserID, usersChan, countChan)
	g.Go(func() error {
		var err error
		tanggal_deadline, err = time.Parse(time.DateOnly, param.Req.TanggalDeadline)
		if err != nil {
			helper.LogsError(err)
			return err
		}
		tanggal_deadline = tanggal_deadline.Local().UTC()
		return nil
	})

	if err := g.Wait(); err != nil {
		return err
	}
	users := <-usersChan
	if len(users) != <-countChan {
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

func (u *tugasUsecase) Update(param ParamUpdate) error {
	g := new(errgroup.Group)
	g.SetLimit(10)

	usersChan := make(chan []entity.User, 1)
	countChan := make(chan int, 1)

	u.checkAndGetUsers(g, param.Ctx, param.Req.JenisSpvID, param.Req.UserID, usersChan, countChan)

	var tanggalDeadline time.Time
	if param.Req.TanggalDeadline != "" {
		g.Go(func() error {
			var err error
			tanggalDeadline, err = time.Parse(time.DateOnly, param.Req.TanggalDeadline)
			if err != nil {
				helper.LogsError(err)
				return err
			}
			tanggalDeadline = tanggalDeadline.Local().UTC()
			return nil
		})
	}

	if err := g.Wait(); err != nil {
		return err
	}

	users := <-usersChan

	if len(users) != <-countChan {
		return errors.New(message.UserNotFoundOrNotSpv)
	}

	newTugas := entity.Tugas{
		Base: entity.Base{
			ID: param.Req.ID,
		},
		JenisSpvID: param.Req.JenisSpvID,
		Users:      users,
	}

	if !tanggalDeadline.IsZero() {
		newTugas.TanggalDeadline = &tanggalDeadline
	}
	err := u.repo.Update(repo_tugas.ParamUpdate{
		Ctx:      param.Ctx,
		NewTugas: &newTugas,
	})

	if err != nil {
		return err
	}

	return nil

}

func (u *tugasUsecase) Delete(param ParamDelete) error {
	err := u.repo.Delete(repo_tugas.ParamDelete{
		Ctx: param.Ctx,
		ID:  param.Req.ID,
	})
	if err != nil {
		return err
	}
	return nil
}

func (u *tugasUsecase) GetAll(param ParamGetAll) ([]entity.Tugas, error) {

	dataTugas, err := u.repo.GetAll(repo_tugas.ParamGetAll{
		Ctx:   param.Ctx,
		Tahun: param.Req.Tahun,
		Bulan: param.Req.Bulan,
	})
	if err != nil {
		return nil, err
	}

	return dataTugas, nil

}

func (u *tugasUsecase) GetById(param ParamGetById) (*entity.Tugas, error) {
	data, err := u.repo.GetById(repo_tugas.ParamGetById{
		Ctx: param.Ctx,
		ID:  param.Req.ID,
	})
	if err != nil {
		return nil, err
	}
	return data, nil
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
