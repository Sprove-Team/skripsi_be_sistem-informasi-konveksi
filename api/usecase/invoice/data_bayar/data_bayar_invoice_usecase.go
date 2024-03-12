package uc_invoice_data_bayar

import (
	"context"
	"errors"

	repoInvoice "github.com/be-sistem-informasi-konveksi/api/repository/invoice/mysql/gorm"
	repo "github.com/be-sistem-informasi-konveksi/api/repository/invoice/mysql/gorm/data_bayar"
	"github.com/be-sistem-informasi-konveksi/common/message"
	reqGlobal "github.com/be-sistem-informasi-konveksi/common/request/global"
	req "github.com/be-sistem-informasi-konveksi/common/request/invoice/data_bayar"
	"github.com/be-sistem-informasi-konveksi/entity"
	"github.com/be-sistem-informasi-konveksi/pkg"
)

type (
	ParamCreateByInvoiceID struct {
		Ctx context.Context
		Req req.Create
	}
	ParamUpdateDataBayarInvoice struct {
		Ctx    context.Context
		Claims *pkg.Claims
		Req    req.Update
	}
	ParamUpdateCommitDB struct {
		Ctx         context.Context
		DataBayar   *entity.DataBayarInvoice
		DataBayarHP *entity.DataBayarHutangPiutang
	}
	ParamDelete struct {
		Ctx context.Context
		Req reqGlobal.ParamByID
	}
	ParamGetByInvoiceID struct {
		Ctx context.Context
		Req req.GetByInvoiceID
	}
)

type DataBayarInvoice interface {
	CreateByInvoiceID(param ParamCreateByInvoiceID) error
	UpdateDataBayarInvoice(param ParamUpdateDataBayarInvoice) (*entity.DataBayarInvoice, error)
	UpdateCommitDB(param ParamUpdateCommitDB) error
	Delete(param ParamDelete) error
	GetByInvoiceID(param ParamGetByInvoiceID) ([]entity.DataBayarInvoice, error)
}

type dataBayarInvoice struct {
	repo        repo.DataBayarInvoiceRepo
	repoInvoice repoInvoice.InvoiceRepo
	ulid        pkg.UlidPkg
}

func NewDataInvoice(repo repo.DataBayarInvoiceRepo, repoInvoice repoInvoice.InvoiceRepo, ulid pkg.UlidPkg) DataBayarInvoice {
	return &dataBayarInvoice{repo, repoInvoice, ulid}
}

func (u *dataBayarInvoice) IsInvoiceExist(ctx context.Context, id string) error {
	err := u.repoInvoice.CheckInvoice(repoInvoice.ParamGetById{Ctx: ctx, ID: id})
	if err != nil {
		if err.Error() == "record not found" {
			return errors.New(message.InvoiceNotFound)
		}
		return err
	}
	return nil
}

func (u *dataBayarInvoice) IsStatusTerkonfirmasi(ctx context.Context, id string) error {
	oldData, err := u.repo.GetByID(repo.ParamGetById{
		Ctx: ctx,
		ID:  id,
	})

	if err != nil {
		return err
	}

	if oldData.Status == "TERKONFIRMASI" {
		return errors.New(message.CannotModifiedTerkonfirmasiDataBayar)
	}

	return nil
}

func (u *dataBayarInvoice) CreateByInvoiceID(param ParamCreateByInvoiceID) error {
	if err := u.IsInvoiceExist(param.Ctx, param.Req.InvoiceID); err != nil {
		return err
	}
	dataBayar := entity.DataBayarInvoice{
		Base: entity.Base{
			ID: u.ulid.MakeUlid().String(),
		},
		InvoiceID:       param.Req.InvoiceID,
		AkunID:          param.Req.AkunID,
		Keterangan:      param.Req.Keterangan,
		BuktiPembayaran: param.Req.BuktiPembayaran,
		Total:           param.Req.Total,
	}
	if err := u.repo.Create(repo.ParamCreate{
		Ctx:       param.Ctx,
		DataBayar: &dataBayar,
	}); err != nil {
		return err
	}
	return nil
}

func (u *dataBayarInvoice) UpdateDataBayarInvoice(param ParamUpdateDataBayarInvoice) (*entity.DataBayarInvoice, error) {

	if param.Claims == nil {
		return nil, errors.New(message.InternalServerError)
	}

	if err := u.IsStatusTerkonfirmasi(param.Ctx, param.Req.ID); err != nil {
		return nil, err
	}

	oldDataByrInvoice, err := u.repo.GetByID(repo.ParamGetById{
		Ctx: param.Ctx,
		ID:  param.Req.ID,
	})

	if err != nil {
		return nil, err
	}

	dataBayarInvoice := &entity.DataBayarInvoice{
		Base: entity.Base{
			ID: param.Req.ID,
		},
		InvoiceID:       oldDataByrInvoice.InvoiceID,
		Keterangan:      oldDataByrInvoice.Keterangan,
		AkunID:          oldDataByrInvoice.AkunID,
		BuktiPembayaran: oldDataByrInvoice.BuktiPembayaran,
		Total:           oldDataByrInvoice.Total,
	}

	switch param.Claims.Role {
	case "BENDAHARA", "DIREKTUR":
		dataBayarInvoice.Status = param.Req.Status
	}

	if param.Req.Keterangan != "" {
		dataBayarInvoice.Keterangan = param.Req.Keterangan
	}
	if param.Req.AkunID != "" {
		dataBayarInvoice.AkunID = param.Req.AkunID
	}
	if len(param.Req.BuktiPembayaran) > 0 {
		dataBayarInvoice.BuktiPembayaran = param.Req.BuktiPembayaran
	}
	if param.Req.Total != 0 {
		dataBayarInvoice.Total = param.Req.Total
	}

	return dataBayarInvoice, nil

}

func (u *dataBayarInvoice) UpdateCommitDB(param ParamUpdateCommitDB) error {
	err := u.repo.Update(repo.ParamUpdate(param))

	if err != nil {
		return err
	}

	return nil
}

func (u *dataBayarInvoice) Delete(param ParamDelete) error {
	if err := u.IsStatusTerkonfirmasi(param.Ctx, param.Req.ID); err != nil {
		return err
	}
	err := u.repo.Delete(repo.ParamDelete{Ctx: param.Ctx, ID: param.Req.ID})
	if err != nil {
		return err
	}
	return nil
}

func (u *dataBayarInvoice) GetByInvoiceID(param ParamGetByInvoiceID) ([]entity.DataBayarInvoice, error) {
	if err := u.IsInvoiceExist(param.Ctx, param.Req.InvoiceID); err != nil {
		return nil, err
	}

	datas, err := u.repo.GetByInvoiceID(repo.ParamGetByInvoiceID{
		Ctx:       param.Ctx,
		InvoiceID: param.Req.InvoiceID,
		Status:    param.Req.Status,
		AkunID:    param.Req.AkunID,
	})

	if err != nil {
		return nil, err
	}
	return datas, nil
}
