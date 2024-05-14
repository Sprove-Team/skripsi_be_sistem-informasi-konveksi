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
		TanggalDipesan  time.Time
		Order           string
		Next            string
		Limit           int
	}
	ParamGetById struct {
		Ctx context.Context
		ID  string
	}
	ParamGetByIdWithoutPreload struct {
		Ctx context.Context
		ID  string
	}
	ParamUpdateFullAssoc struct {
		Ctx     context.Context
		Invoice *entity.Invoice
	}
	ParamUpdate struct {
		Ctx     context.Context
		Invoice *entity.Invoice
	}

	ParamDelete struct {
		Ctx     context.Context
		Invoice *entity.Invoice
	}
)

type InvoiceRepo interface {
	Create(param ParamCreate) error
	UpdateFullAssoc(param ParamUpdateFullAssoc) error
	Update(param ParamUpdate) error
	GetLastInvoiceCrrYear(ctx context.Context) (entity.Invoice, error)
	GetById(param ParamGetById) (*entity.Invoice, error)
	CheckInvoice(param ParamGetById) error
	GetByIdWithoutPreload(param ParamGetByIdWithoutPreload) (*entity.Invoice, error)
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
	err := r.DB.WithContext(param.Ctx).Session(&gorm.Session{FullSaveAssociations: true}).Updates(param.Invoice).Error
	if err != nil {
		helper.LogsError(err)
		return err
	}
	return nil
}

func (r *invoiceRepo) Update(param ParamUpdate) error {
	err := r.DB.WithContext(param.Ctx).Updates(param.Invoice).Error
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

func (r *invoiceRepo) GetByIdWithoutPreload(param ParamGetByIdWithoutPreload) (*entity.Invoice, error) {
	data := new(entity.Invoice)
	err := r.DB.WithContext(param.Ctx).First(data, "id = ?", param.ID).Error
	if err != nil {
		helper.LogsError(err)
		return nil, err
	}
	return data, nil
}

func (r *invoiceRepo) GetById(param ParamGetById) (*entity.Invoice, error) {
	data := new(entity.Invoice)
	tx := r.DB.WithContext(param.Ctx).
		Preload("HutangPiutang", func(db *gorm.DB) *gorm.DB {
			return db.Select("invoice_id", "sisa")
		}).
		Preload("DataBayarInvoice").
		Preload("DataBayarInvoice.Akun").
		Preload("DetailInvoice", func(db *gorm.DB) *gorm.DB {
			return db.Order("qty ASC")
		}).
		Preload("DetailInvoice.Produk", func(db *gorm.DB) *gorm.DB {
			return db.Omit("created_at")
		}).
		Preload("DetailInvoice.Produk.HargaDetails", func(db *gorm.DB) *gorm.DB {
			return db.Omit("created_at").Order("qty ASC")
		}).
		Preload("DetailInvoice.Bordir", func(db *gorm.DB) *gorm.DB {
			return db.Omit("created_at")
		}).
		Preload("DetailInvoice.Sablon", func(db *gorm.DB) *gorm.DB {
			return db.Omit("created_at")
		}).
		Preload("Kontak", func(db *gorm.DB) *gorm.DB {
			return db.Unscoped().Omit("created_at")
		}).Preload("User", func(db *gorm.DB) *gorm.DB {
		return db.Omit("no_telp", "alamat", "created_at")
	}).First(data, "id = ?", param.ID)

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
		Preload("Kontak", func(db *gorm.DB) *gorm.DB {
			return db.Unscoped()
		}).
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

	order := strings.Split(param.Order, " ")
	if param.Order == "" {
		tx = tx.Order("id ASC")
	}

	s := reflect.ValueOf(param)

	for i := 0; i < s.NumField(); i++ {
		key := s.Type().Field(i).Name
		value := s.Field(i).Interface()

		switch value.(type) {
		case string:
			if value != "" {
				switch key {
				case "Next":
					if param.Order != "" {
						invoice := new(entity.Invoice)
						if err := r.DB.Model(&entity.Invoice{}).Where("id = ?", value).First(invoice).Error; err != nil {
							if err != gorm.ErrRecordNotFound {
								return nil, err
							}
						}
						if invoice.ID != "" {
							switch order[0] {
							case "tanggal_kirim":
								if order[1] == "DESC" {
									tx = tx.Where("tanggal_kirim < ?", invoice.TanggalKirim)
								} else {
									tx = tx.Where("tanggal_kirim > ?", invoice.TanggalKirim)
								}
							case "tanggal_deadline":
								if order[1] == "DESC" {
									tx = tx.Where("tanggal_deadline < ?", invoice.TanggalDeadline)
								} else {
									tx = tx.Where("tanggal_deadline > ?", invoice.TanggalDeadline)
								}
							case "tanggal_dipesan":
								if order[1] == "DESC" {
									tx = tx.Where("created_at < ?", invoice.CreatedAt)
								} else {
									tx = tx.Where("created_at > ?", invoice.CreatedAt)
								}
							}
						}
					} else {
						tx = tx.Where("id > ?", value)
					}
				case "Order":
					if order[0] == "tanggal_dipesan" {
						tx = tx.Order("created_at " + order[1])	
					}else {
						tx = tx.Order(value)
					}
				case "KontakID":
					tx = tx.Where("kontak_id = ?", value)
				case "StatusProduksi":
					tx = tx.Where("status_produksi = ?", value)
				}
			}
		case time.Time:
			if !value.(time.Time).IsZero() {
				t := value.(time.Time).Format(time.DateOnly)
				switch key {
				case "TanggalDeadline":
					tx = tx.Where("DATE(tanggal_deadline) = ?", t)
				case "TanggalKirim":
					tx = tx.Where("DATE(tanggal_kirim) = ?", t)
				case "TanggalDipesan":
					tx = tx.Where("DATE(created_at) = ?", t)
				}
			}
		}
	}

	invoices := make([]entity.Invoice, param.Limit)

	if err := tx.Limit(param.Limit).Preload("Kontak", func(db *gorm.DB) *gorm.DB {
		return db.Unscoped().Omit("created_at")
	}).Preload("HutangPiutang", func(db *gorm.DB) *gorm.DB {
		return db.Select("invoice_id", "sisa")
	}).
		Preload("User", func(db *gorm.DB) *gorm.DB {
			return db.Omit("no_telp", "alamat", "created_at")
		}).Find(&invoices).Error; err != nil {
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
