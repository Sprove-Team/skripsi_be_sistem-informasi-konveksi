package invoice

import (
	"context"
	"time"

	"github.com/be-sistem-informasi-konveksi/entity"
	"github.com/be-sistem-informasi-konveksi/helper"
	"gorm.io/gorm"
)

type CreateParam struct {
	Invoice       *entity.Invoice
	DetailInvoice []entity.DetailInvoice
}

type InvoiceRepo interface {
	Create(ctx context.Context, param CreateParam) error
	GetLastInvoiceCrrYear(ctx context.Context) (entity.Invoice, error)
}

type invoiceRepo struct {
	DB *gorm.DB
}

func NewInvoiceRepo(DB *gorm.DB) InvoiceRepo {
	return &invoiceRepo{DB}
}

func (r *invoiceRepo) Create(ctx context.Context, param CreateParam) error {
	err := r.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(param.Invoice).Error; err != nil {
			return err
		}
		if err := tx.Create(&param.DetailInvoice).Error; err != nil {
			return err
		}
		return nil
	})
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
