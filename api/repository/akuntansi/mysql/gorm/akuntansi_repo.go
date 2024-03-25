package repo_akuntansi

import (
	"context"
	"time"

	"github.com/be-sistem-informasi-konveksi/entity"
	"github.com/be-sistem-informasi-konveksi/helper"
	"gorm.io/gorm"
)

type ResultDataJU struct {
	AyatJurnalID string
	TransaksiID  string
	AkunID       string
	NamaAkun     string
	KodeAkun     string
	Tanggal      string
	Keterangan   string
	Debit        float64
	Kredit       float64
}

type ResultDataBB struct {
	TransaksiID string
	AkunID      string
	NamaAkun    string
	KodeAkun    string
	SaldoNormal string
	Tanggal     string
	Keterangan  string
	Debit       float64
	Kredit      float64
	Saldo       float64
}

type ResultSaldoAwalDataBB struct {
	AkunID      string
	KodeAkun    string
	SaldoNormal string
	NamaAkun    string
	Saldo       float64
}

type ResultDataNc struct {
	KodeAkun string
	NamaAkun string
	// SaldoDebit  float64
	// SaldoKredit float64
	SaldoNormal string
	Saldo       float64
}

type ResultDataLBR struct {
	KategoriAkun string
	NamaAkun     string
	KodeAkun     string
	// SaldoDebit   float64 // from ayat jurnal
	// SaldoKredit  float64 // from ayat jurnal
	SaldoNormal string
	Saldo       float64
}

type AkuntansiRepo interface {
	GetDataJU(ctx context.Context, startDate, endDate time.Time) ([]ResultDataJU, error)
	GetDataBB(ctx context.Context, akunIDs []string, startDate, endDate time.Time) ([]ResultDataBB, []ResultSaldoAwalDataBB, error)
	GetDataNC(ctx context.Context, date time.Time) ([]ResultDataNc, error)
	GetDataLBR(ctx context.Context, startDate, endDate time.Time) ([]ResultDataLBR, error)
	// GetNeraca(ctx context.Context)
}

type akuntansiRepo struct {
	DB *gorm.DB
}

func NewAkuntansiRepo(DB *gorm.DB) AkuntansiRepo {
	return &akuntansiRepo{DB}
}

func (r *akuntansiRepo) GetDataJU(ctx context.Context, startDate, endDate time.Time) ([]ResultDataJU, error) {
	resultDatasJU := []ResultDataJU{}

	subQuery := r.DB.Model(&entity.Transaksi{}).
		Select("id").
		Where("DATE(tanggal) >= ? AND DATE(tanggal) <= ?", startDate, endDate)

	err := r.DB.WithContext(ctx).Model(&entity.Akun{}).
		Joins("JOIN ayat_jurnal ON akun.id = ayat_jurnal.akun_id").
		Joins("JOIN transaksi ON ayat_jurnal.transaksi_id = transaksi.id").
		Select("ayat_jurnal.id as AyatJurnalID,transaksi.id as TransaksiID,akun.id as AkunID, akun.kode as KodeAkun, akun.nama as NamaAkun, ayat_jurnal.debit as Debit, ayat_jurnal.kredit as Kredit, transaksi.Tanggal as Tanggal, transaksi.keterangan as Keterangan").
		Where("transaksi.id IN (?)", subQuery).
		Group("ayat_jurnal.id").
		Find(&resultDatasJU).Error
	if err != nil {
		helper.LogsError(err)
		return []ResultDataJU{}, err
	}
	return resultDatasJU, nil
}

func (r *akuntansiRepo) GetDataBB(ctx context.Context, akunIDs []string, startDate, endDate time.Time) ([]ResultDataBB, []ResultSaldoAwalDataBB, error) {
	resultDatasBB := []ResultDataBB{}
	resultSaldoAwalDatasBB := []ResultSaldoAwalDataBB{}

	isAkunIdExist := len(akunIDs) == 0
	err := r.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		subQuery := tx.Model(&entity.Transaksi{}).
			Select("id").
			Where("DATE(tanggal) >= ? AND DATE(tanggal) <= ?", startDate, endDate)

		tx2 := tx.Model(&entity.Akun{})

		tx2 = tx2.Joins("JOIN ayat_jurnal ON akun.id = ayat_jurnal.akun_id").
			Joins("JOIN transaksi ON ayat_jurnal.transaksi_id = transaksi.id").
			Select("transaksi.id as TransaksiID,akun.id as AkunID, akun.kode as KodeAkun, akun.nama as NamaAkun, akun.saldo_normal as SaldoNormal, ayat_jurnal.debit as Debit, ayat_jurnal.kredit as Kredit, transaksi.Tanggal as Tanggal, transaksi.keterangan as Keterangan, ayat_jurnal.saldo as Saldo")

		if !isAkunIdExist {
			tx2 = tx2.Where("akun.id IN (?)", akunIDs)
		}

		err2 := tx2.Where("transaksi.id IN (?)", subQuery).
			Group("ayat_jurnal.id").
			Find(&resultDatasBB).Error

		if err2 != nil {
			return err2
		}

		// saldo awal
		subQuery2 := tx.Model(&entity.Transaksi{}).Select("id").Where("DATE(tanggal) < ?", startDate)

		tx3 := tx.Model(&entity.Akun{}).
			Joins("JOIN ayat_jurnal ON akun.id = ayat_jurnal.akun_id").
			Select("akun.id as AkunID,akun.Kode as KodeAkun, akun.nama as NamaAkun, akun.saldo_normal as SaldoNormal, COALESCE(SUM(ayat_jurnal.saldo),0) as Saldo")

		if !isAkunIdExist {
			tx3 = tx3.Where("akun.id IN (?)", akunIDs)
		}

		err3 := tx3.Where("ayat_jurnal.transaksi_id IN (?)", subQuery2).Group("akun.id").
			Find(&resultSaldoAwalDatasBB).Error

		if err3 != nil {
			return err3
		}

		return nil
	})
	if err != nil {
		helper.LogsError(err)
		return []ResultDataBB{}, []ResultSaldoAwalDataBB{}, err
	}

	return resultDatasBB, resultSaldoAwalDatasBB, nil
}

func (r *akuntansiRepo) GetDataNC(ctx context.Context, date time.Time) ([]ResultDataNc, error) {
	resultDatasNc := []ResultDataNc{}
	err := r.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {

		subQueryTr := tx.Model(&entity.Transaksi{}).Select("id").Where("MONTH(tanggal) = ? AND YEAR(tanggal) = ?", date.Month(), date.Year())

		// Select("akun.kode as KodeAkun, akun.nama as NamaAkun, COALESCE(ABS(sum(ayat_jurnal.saldo)), 0) AS Saldo, akun.saldo_normal AS SaldoNormal"). -> hitung manual dari tabel ayatjurnals
		// Select("akun.kode as KodeAkun, akun.nama as NamaAkun, ABS(akun.saldo) AS Saldo, akun.saldo_normal AS SaldoNormal").
		err2 := tx.Model(&entity.Akun{}).
			Joins("JOIN ayat_jurnal ON ayat_jurnal.akun_id = akun.id").
			Joins("JOIN transaksi ON transaksi.id = ayat_jurnal.transaksi_id").
			Select("akun.kode as KodeAkun, akun.nama as NamaAkun, COALESCE(ABS(sum(ayat_jurnal.saldo)), 0) AS Saldo, akun.saldo_normal AS SaldoNormal").
			Where("transaksi.id IN (?)", subQueryTr).
			Group("akun.id").
			Find(&resultDatasNc).Error
		if err2 != nil {
			return err2
		}
		return nil
	})
	if err != nil {
		helper.LogsError(err)
		return []ResultDataNc{}, err
	}

	return resultDatasNc, nil
}

func (r *akuntansiRepo) GetDataLBR(ctx context.Context, startDate, endDate time.Time) ([]ResultDataLBR, error) {
	datas := []ResultDataLBR{}

	err := r.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		subQueryTr := tx.Model(&entity.Transaksi{}).
			Select("id").
			Where("DATE(tanggal) >= ? AND DATE(tanggal) <= ?", startDate, endDate)

			// Select("k.kategori_akun as KategoriAkun, akun.Nama as NamaAkun, akun.saldo as Saldo, COALESCE(sum(ay.debit), 0) as SaldoDebit, COALESCE(sum(ay.kredit), 0) as SaldoKredit, COALESCE(sum(ay.saldo), 0) as Saldo, akun.kode as KodeAkun").
			// Select("akun.kode AS KodeAkun,k.kategori_akun AS KategoriAkun, akun.Nama AS NamaAkun, akun.saldo AS Saldo").
		err := r.DB.WithContext(ctx).Model(&entity.Akun{}).
			Joins("JOIN kelompok_akun k ON k.id = akun.kelompok_akun_id AND k.kategori_akun IN (?)", []string{"PENDAPATAN", "BEBAN"}).
			Joins("JOIN ayat_jurnal ay ON ay.akun_id = akun.id").
			Select("k.kategori_akun as KategoriAkun, akun.Nama as NamaAkun, COALESCE(sum(ay.saldo), 0) as Saldo, akun.kode as KodeAkun").
			Where("ay.transaksi_id IN (?)", subQueryTr).
			Group("akun.id").
			Find(&datas).Error
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		helper.LogsError(err)
		return nil, err
	}
	return datas, nil
}
