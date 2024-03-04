package repo_invoice

import (
	"context"
	"reflect"
	"strings"
	"time"

	"github.com/be-sistem-informasi-konveksi/entity"
	"github.com/be-sistem-informasi-konveksi/helper"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type (
	ParamCreate struct {
		Ctx     context.Context
		Invoice *entity.Invoice
	}
	ParamGetAll struct {
		Ctx             context.Context
		StatusProduksi  string
		KontakID        string
		TanggalDeadline time.Time
		TanggalKirim    time.Time
		Order           string
		Next            string
		Limit           int
	}
	ParamGetById struct {
		Ctx context.Context
		ID  string
	}
	ParamUpdateFullAssoc struct {
		Ctx     context.Context
		Invoice *entity.Invoice
	}
	ParamUpdateStatusProduksi struct {
		Ctx    context.Context
		ID     string
		Status string
	}
	ParamDelete struct {
		Ctx     context.Context
		Invoice *entity.Invoice
	}
)

type InvoiceRepo interface {
	Create(param ParamCreate) error
	UpdateFullAssoc(param ParamUpdateFullAssoc) error
	UpdateStatusProduksi(param ParamUpdateStatusProduksi) error
	GetLastInvoiceCrrYear(ctx context.Context) (entity.Invoice, error)
	GetById(param ParamGetById) (*entity.Invoice, error)
	CheckInvoice(param ParamGetById) error
	GetByIdFullAssoc(param ParamGetById) (*entity.Invoice, error)
	GetAll(param ParamGetAll) ([]entity.Invoice, error)
	Delete(param ParamDelete) error
}

type invoiceRepo struct {
	DB *gorm.DB
}

func NewInvoiceRepo(DB *gorm.DB) InvoiceRepo {
	return &invoiceRepo{DB}
}

func (r *invoiceRepo) Create(param ParamCreate) error {
	err := r.DB.WithContext(param.Ctx).Create(param.Invoice).Error
	if err != nil {
		helper.LogsError(err)
		return err
	}
	return nil
}

func (r *invoiceRepo) UpdateFullAssoc(param ParamUpdateFullAssoc) error {
	err := r.DB.WithContext(param.Ctx).Session(&gorm.Session{FullSaveAssociations: true}).Omit("status_produksi").Updates(param.Invoice).Error
	if err != nil {
		helper.LogsError(err)
		return err
	}
	return nil
}

func (r *invoiceRepo) UpdateStatusProduksi(param ParamUpdateStatusProduksi) error {
	err := r.DB.WithContext(param.Ctx).Model(&entity.Invoice{}).Where("id = ?", param.ID).Update("status_produksi", param.Status).Error
	if err != nil {
		helper.LogsError(err)
		return err
	}
	return nil
}

func (r *invoiceRepo) GetLastInvoiceCrrYear(ctx context.Context) (entity.Invoice, error) {
	currentYear := time.Now().Year()
	invoice := new(entity.Invoice)
	err := r.DB.WithContext(ctx).Model(invoice).Select("nomor_referensi").Order("nomor_referensi DESC").First(invoice, "YEAR(created_at) = ?", currentYear).Error
	if err != nil {
		helper.LogsError(err)
		return *invoice, err
	}
	return *invoice, err
}

func (r *invoiceRepo) CheckInvoice(param ParamGetById) error {
	err := r.DB.WithContext(param.Ctx).First(&entity.Invoice{}, "id = ?", param.ID).Error
	if err != nil {
		helper.LogsError(err)
		return err
	}
	return nil
}

func (r *invoiceRepo) GetById(param ParamGetById) (*entity.Invoice, error) {
	data := new(entity.Invoice)
	tx := r.DB.WithContext(param.Ctx).
		Preload("DataBayarInvoice").
		Preload("DataBayarInvoice.Akun").
		Preload("DetailInvoice").
		Preload("Kontak").
		Preload("User").
		First(data, "id = ?", param.ID)

	if err := tx.Error; err != nil {
		helper.LogsError(err)
		return nil, err
	}
	return data, nil
}

func (r *invoiceRepo) GetByIdFullAssoc(param ParamGetById) (*entity.Invoice, error) {
	data := new(entity.Invoice)
	tx := r.DB.WithContext(param.Ctx).
		Preload("HutangPiutang").
		Preload("DetailInvoice").
		Preload("HutangPiutang.Transaksi").
		Preload("HutangPiutang.Transaksi.AyatJurnals").
		Preload("Kontak").
		Preload("User").
		First(data, "id = ?", param.ID)

	if err := tx.Error; err != nil {
		helper.LogsError(err)
		return nil, err
	}
	return data, nil
}

func (r *invoiceRepo) GetAll(param ParamGetAll) ([]entity.Invoice, error) {
	tx := r.DB.WithContext(param.Ctx).Model(&entity.Invoice{})

	if param.Order == "" {
		tx = tx.Order("id ASC")
		param.Order = "id ASC"
	}

	s := reflect.ValueOf(param)
	orderBy := strings.Split(param.Order, " ")[1]
	for i := 0; i < s.NumField(); i++ {
		key := s.Type().Field(i).Name
		value := s.Field(i).Interface()

		switch value.(type) {
		case string:
			if value != "" {
				switch key {
				case "Next":
					if orderBy == "DESC" {
						tx = tx.Where("id < ?", value)
					} else {
						tx = tx.Where("id > ?", value)
					}
				case "Order":
					tx = tx.Order(value)
				case "KontakID":
					tx = tx.Where("kontak_id = ?", value)
				case "StatusProduksi":
					tx = tx.Where("status_produksi = ?", value)
				}
			}
		case time.Time:
			if !value.(time.Time).IsZero() {
				switch key {
				case "TanggalDeadline":
					tx = tx.Where("DATE(tanggal_deadline) = ?", value)
				case "TanggalKirim":
					tx = tx.Where("DATE(tanggal_kirim) = ?", value)
				}
			}
		}
	}

	invoices := make([]entity.Invoice, param.Limit)

	if err := tx.Limit(param.Limit).Preload("Kontak").Preload("User").Find(&invoices).Error; err != nil {
		return nil, err
	}

	return invoices, nil
}

func (r *invoiceRepo) Delete(param ParamDelete) error {
	err := r.DB.WithContext(param.Ctx).Select(clause.Associations).Delete(param.Invoice).Error
	if err != nil {
		helper.LogsError(err)
		return err
	}
	return nil
}
