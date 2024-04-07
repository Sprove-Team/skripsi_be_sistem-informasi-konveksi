package repo_invoice_data_bayar

import (
	"context"

	"github.com/be-sistem-informasi-konveksi/entity"
	"github.com/be-sistem-informasi-konveksi/helper"
	"gorm.io/gorm"
)

type (
	ParamCreate struct {
		Ctx       context.Context
		DataBayar *entity.DataBayarInvoice
	}
	ParamUpdate struct {
		Ctx         context.Context
		DataBayar   *entity.DataBayarInvoice
		DataBayarHP *entity.DataBayarHutangPiutang
	}
	ParamDelete struct {
		Ctx context.Context
		ID  string
	}
	ParamGetByInvoiceID struct {
		Ctx       context.Context
		InvoiceID string
		Status    string
		AkunID    string
	}
	ParamGetById struct {
		Ctx context.Context
		PreloadAkun bool
		ID  string
	}
)

type DataBayarInvoiceRepo interface {
	Create(param ParamCreate) error
	Update(param ParamUpdate) error
	Delete(param ParamDelete) error
	GetByID(param ParamGetById) (*entity.DataBayarInvoice, error)
	GetByInvoiceID(param ParamGetByInvoiceID) ([]entity.DataBayarInvoice, error)
}

type dataBayarInvoiceRepo struct {
	DB *gorm.DB
}

func NewDataBayarInvoiceRepo(DB *gorm.DB) DataBayarInvoiceRepo {
	return &dataBayarInvoiceRepo{DB}
}

func (r *dataBayarInvoiceRepo) Create(param ParamCreate) error {
	err := r.DB.WithContext(param.Ctx).Create(param.DataBayar).Error
	if err != nil {
		helper.LogsError(err)
		return err
	}
	return nil
}

func (r *dataBayarInvoiceRepo) Update(param ParamUpdate) error {
	err := r.DB.WithContext(param.Ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Updates(param.DataBayar).Error; err != nil {
			return nil
		}
		if param.DataBayarHP != nil {
			if err := tx.Create(param.DataBayarHP).Error; err != nil {
				return err
			}

			if err := tx.Select("sisa", "status").Updates(&param.DataBayarHP.HutangPiutang).Error; err != nil {
				return err
			}
		}

		return nil
	})
	if err != nil {
		helper.LogsError(err)
		return err
	}
	return nil
}

func (r *dataBayarInvoiceRepo) Delete(param ParamDelete) error {
	err := r.DB.WithContext(param.Ctx).Delete(&entity.DataBayarInvoice{}, "id = ?", param.ID).Error
	if err != nil {
		helper.LogsError(err)
		return err
	}
	return nil
}

func (r *dataBayarInvoiceRepo) GetByID(param ParamGetById) (*entity.DataBayarInvoice, error) {
	data := new(entity.DataBayarInvoice)
	tx := r.DB.WithContext(param.Ctx).Where("id = ?", param.ID)
	if param.PreloadAkun{
		tx = tx.Preload("Akun")
	}
	if err := tx.First(data).Error; err != nil {
		return nil, err
	}
	return data, nil
}

func (r *dataBayarInvoiceRepo) GetByInvoiceID(param ParamGetByInvoiceID) ([]entity.DataBayarInvoice, error) {

	datas := []entity.DataBayarInvoice{}

	tx := r.DB.WithContext(param.Ctx).Model(&entity.DataBayarInvoice{}).Where("invoice_id = ?", param.InvoiceID)

	if param.AkunID != "" {
		tx = tx.Where("akun_id = ?", param.AkunID)
	}

	if param.Status != "" {
		tx = tx.Where("status = ?", param.Status)
	}

	if err := tx.Preload("Akun").Find(&datas).Error; err != nil {
		return nil, err
	}

	return datas, nil
}
