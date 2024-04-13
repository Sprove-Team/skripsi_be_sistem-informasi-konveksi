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
		Ctx         context.Context
		PreloadAkun bool
		ID          string
	}
	ParamGetAll struct {
		Ctx      context.Context
		KontakID string
		Status   string
		Next     string
		Limit    int
	}
)

type DataBayarInvoiceRepo interface {
	Create(param ParamCreate) error
	Update(param ParamUpdate) error
	Delete(param ParamDelete) error
	GetAll(param ParamGetAll) ([]entity.DataBayarInvoice, error)
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

func (r *dataBayarInvoiceRepo) GetAll(param ParamGetAll) ([]entity.DataBayarInvoice, error) {
	datas := make([]entity.DataBayarInvoice, 0, param.Limit)
	tx := r.DB.WithContext(param.Ctx).Model(datas).Order("id ASC")
	if param.Next != "" {
		tx = tx.Where("id > ?", param.Next)
	}

	if param.Status != "" {
		tx = tx.Where("status = ?", param.Status)
	}

	if param.KontakID != "" {
		// Preload the Invoice relationship conditionally
		tx = tx.Joins("JOIN invoice ON data_bayar_invoice.invoice_id = invoice.id AND invoice.kontak_id = ?", param.KontakID)
	}

	tx = tx.Preload("Invoice", func(db *gorm.DB) *gorm.DB {
		return db.Select("id", "kontak_id", "nomor_referensi", "created_at")
	}).Preload("Invoice.Kontak", func(db *gorm.DB) *gorm.DB {
		return db.Select("id", "nama")
	})

	if err := tx.Limit(param.Limit).Find(&datas).Error; err != nil {
		return nil, err
	}
	return datas, nil
}

func (r *dataBayarInvoiceRepo) GetByID(param ParamGetById) (*entity.DataBayarInvoice, error) {
	data := new(entity.DataBayarInvoice)
	tx := r.DB.WithContext(param.Ctx).Where("id = ?", param.ID)
	if param.PreloadAkun {
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
