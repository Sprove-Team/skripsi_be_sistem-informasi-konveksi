package akuntansi

import (
	"context"
	"errors"

	"github.com/be-sistem-informasi-konveksi/common/message"
	"github.com/be-sistem-informasi-konveksi/entity"
	"github.com/be-sistem-informasi-konveksi/helper"
	"github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
)

type AkunTransactionDetails struct {
	ID           string
	SaldoNormal  string
	Saldo        float64
	TotalSaldoTr float64
	KelID        string
	Nama         string
	Kode         string
}

type AkunRepo interface {
	Create(ctx context.Context, akun *entity.Akun) error
	Update(ctx context.Context, akun *entity.Akun) error
	Delete(ctx context.Context, id string) error
	GetAll(ctx context.Context, search SearchAkun) ([]entity.Akun, error)
	GetById(ctx context.Context, id string) (entity.Akun, error)
	GetByIds(ctx context.Context, ids []string) ([]entity.Akun, error)
	GetAkunDetailsByTransactionID(ctx context.Context, id string) ([]AkunTransactionDetails, error)
	// GetAllWithouFilterPreload(ctx context.Context) ([]entity.Akun, error)
}

type akunRepo struct {
	DB *gorm.DB
}

func NewAkunRepo(DB *gorm.DB) AkunRepo {
	return &akunRepo{DB}
}

func (r *akunRepo) GetById(ctx context.Context, id string) (entity.Akun, error) {
	data := entity.Akun{}
	err := r.DB.WithContext(ctx).Omit("created_at", "updated_at", "deleted_at").Preload("KelompokAkun", func(db *gorm.DB) *gorm.DB {
		return db.Select("id", "nama", "kode")
	}).First(&data, "id = ?", id).Error
	if err != nil {
		helper.LogsError(err)
		return data, err
	}
	return data, nil
}

func (r *akunRepo) GetByIds(ctx context.Context, ids []string) ([]entity.Akun, error) {
	datas := []entity.Akun{}

	err := r.DB.WithContext(ctx).Where("id IN ?", ids).Find(&datas).Error
	if err != nil {
		helper.LogsError(err)
		return datas, err
	}
	return datas, nil
}

func (r *akunRepo) GetAkunDetailsByTransactionID(ctx context.Context, id string) ([]AkunTransactionDetails, error) {
	var datas []AkunTransactionDetails
	err := r.DB.WithContext(ctx).Model(&entity.Akun{}).
		Joins("JOIN ayat_jurnal ON ayat_jurnal.akun_id = akun.id").
		Where("ayat_jurnal.transaksi_id", id).
		Select("akun.id as ID, akun.saldo_normal as SaldoNormal, SUM(ayat_jurnal.saldo) as TotalSaldoTr, akun.saldo as Saldo, akun.kelompok_akun_id as KelID, akun.nama as Nama, akun.kode as Kode").
		Group("akun.id").
		Find(&datas).Error
	if err != nil {
		helper.LogsError(err)
		return nil, err
	}
	return datas, nil
}

func (r *akunRepo) Create(ctx context.Context, akun *entity.Akun) error {
	err := r.DB.WithContext(ctx).Create(akun).Error
	if err != nil {
		helper.LogsError(err)
		return err
	}
	return nil
}

func (r *akunRepo) Update(ctx context.Context, akun *entity.Akun) error {
	err := r.DB.WithContext(ctx).Omit("id, created_at, updated_at, deleted_at").Updates(akun).Error
	if err != nil {
		helper.LogsError(err)
		return err
	}
	return nil
}

func (r *akunRepo) Delete(ctx context.Context, id string) error {
	err := r.DB.WithContext(ctx).Where("id = ?", id).Delete(&entity.Akun{}).Error
	if err != nil {
		if mysqlErr, ok := err.(*mysql.MySQLError); ok {
			switch mysqlErr.Number {
			case 1451: // MySQL code for duplicate entry
				return errors.New(message.AkunCannotDeleted)
			default:
				return err
			}
		}
		helper.LogsError(err)
		return err
	}
	return nil
}

func (r *akunRepo) DeleteById(ctx context.Context, id string) error {
	err := r.DB.WithContext(ctx).Delete(&entity.Akun{}, "id = ?", id).Error
	if err != nil {
		helper.LogsError(err)
		return err
	}
	return nil
}

type SearchAkun struct {
	Nama  string
	Kode  string
	Next  string
	Limit int
}

func (r *akunRepo) GetAll(ctx context.Context, searchAkun SearchAkun) ([]entity.Akun, error) {
	datas := []entity.Akun{}

	tx := r.DB.WithContext(ctx).Model(&entity.Akun{}).Order("id ASC").Omit("created_at", "deleted_at", "updated_at")

	conditions := map[string]interface{}{
		"id > ?":      searchAkun.Next,
		"nama LIKE ?": "%" + searchAkun.Nama + "%",
		"kode = ?":    searchAkun.Kode,
	}

	for condition, value := range conditions {
		if value != "" {
			tx = tx.Where(condition, value)
		}
	}

	err := tx.Limit(searchAkun.Limit).Preload("KelompokAkun", func(db *gorm.DB) *gorm.DB {
		return db.Select("id", "nama", "kode")
	}).Find(&datas).Error
	if err != nil {
		helper.LogsError(err)
		return datas, err
	}
	return datas, nil
}
