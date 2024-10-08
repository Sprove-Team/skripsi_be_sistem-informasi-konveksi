package uc_invoice_data_bayar

import (
	"context"
	"errors"
	"sort"

	repo_akuntansi_hp "github.com/be-sistem-informasi-konveksi/api/repository/akuntansi/mysql/gorm/hutang_piutang"
	repoInvoice "github.com/be-sistem-informasi-konveksi/api/repository/invoice/mysql/gorm"
	repo "github.com/be-sistem-informasi-konveksi/api/repository/invoice/mysql/gorm/data_bayar"
	"github.com/be-sistem-informasi-konveksi/common/message"
	reqGlobal "github.com/be-sistem-informasi-konveksi/common/request/global"
	req_global "github.com/be-sistem-informasi-konveksi/common/request/global"
	req "github.com/be-sistem-informasi-konveksi/common/request/invoice/data_bayar"
	"github.com/be-sistem-informasi-konveksi/entity"
	"github.com/be-sistem-informasi-konveksi/pkg"
)

type (
	ParamCreateByInvoiceID struct {
		Ctx context.Context
		Req req.CreateByInvoiceId
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
	ParamGetByID struct {
		Ctx context.Context
		Req req_global.ParamByID
	}
	ParamGetAll struct {
		Ctx context.Context
		Req req.GetAll
	}
)

type DataBayarInvoice interface {
	CreateByInvoiceID(param ParamCreateByInvoiceID) error
	UpdateDataBayarInvoice(param ParamUpdateDataBayarInvoice) (*entity.DataBayarInvoice, error)
	UpdateCommitDB(param ParamUpdateCommitDB) error
	Delete(param ParamDelete) error
	GetByInvoiceID(param ParamGetByInvoiceID) ([]entity.DataBayarInvoice, error)
	GetByID(param ParamGetByID) (*entity.DataBayarInvoice, error)
	GetAll(param ParamGetAll) ([]entity.DataBayarInvoice, error)
}

type dataBayarInvoice struct {
	repo        repo.DataBayarInvoiceRepo
	repoInvoice repoInvoice.InvoiceRepo
	repoHP      repo_akuntansi_hp.HutangPiutangRepo
	ulid        pkg.UlidPkg
}

func NewDataInvoice(repo repo.DataBayarInvoiceRepo, repoInvoice repoInvoice.InvoiceRepo, repoHP repo_akuntansi_hp.HutangPiutangRepo, ulid pkg.UlidPkg) DataBayarInvoice {
	return &dataBayarInvoice{repo, repoInvoice, repoHP, ulid}
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
	invoice, err := u.repoInvoice.GetByIdWithoutPreload(repoInvoice.ParamGetByIdWithoutPreload{
		Ctx: param.Ctx,
		ID:  param.Req.InvoiceID,
	})
	if err != nil {
		if err.Error() == "record not found" {
			return errors.New(message.InvoiceNotFound)
		}
		return err
	}

	hp, err := u.repoHP.GetByInvoiceId(repo_akuntansi_hp.ParamGetByInvoiceId{
		Ctx: param.Ctx,
		ID:  param.Req.InvoiceID,
	})

	if err != nil {
		if err.Error() != "record not found" {
			return err
		}
	}

	var sisa float64
	var errValidate error
	if hp.ID == "" {
		sisa = invoice.TotalHarga - param.Req.Total
		if sisa < 0 {
			errValidate = errors.New(message.BayarMustLessThanTotalHargaInvoice)
		}
	} else {
		sisa = hp.Sisa - param.Req.Total
		if sisa < 0 {
			errValidate = errors.New(message.BayarMustLessThanSisaTagihan)
		}
	}

	if errValidate != nil {
		return errValidate
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
	case entity.RolesById[2], entity.RolesById[1]:
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
func (u *dataBayarInvoice) GetByID(param ParamGetByID) (*entity.DataBayarInvoice, error) {
	data, err := u.repo.GetByID(repo.ParamGetById{
		Ctx:         param.Ctx,
		ID:          param.Req.ID,
		PreloadAkun: true,
	})
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (u *dataBayarInvoice) GetAll(param ParamGetAll) ([]entity.DataBayarInvoice, error) {
	if param.Req.Limit <= 0 {
		param.Req.Limit = 10
	}
	datas, err := u.repo.GetAll(repo.ParamGetAll{
		Ctx:      param.Ctx,
		KontakID: param.Req.KontakID,
		Status:   param.Req.Status,
		Next:     param.Req.Next,
		Limit:    param.Req.Limit,
	})
	if err != nil {
		return nil, err
	}
	
	sort.Slice(datas, func(i, j int) bool {
		switch param.Req.Sort {
		case "DESC" :
			return datas[i].CreatedAt.After(*datas[j].CreatedAt)
		default:
			return datas[i].CreatedAt.Before(*datas[j].CreatedAt)
		}
	})

	return datas, nil
}
