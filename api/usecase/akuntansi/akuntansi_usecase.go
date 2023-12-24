package akuntansi

import (
	"context"
	"math"
	"time"

	repo "github.com/be-sistem-informasi-konveksi/api/repository/akuntansi/mysql/gorm"
	req "github.com/be-sistem-informasi-konveksi/common/request/akuntansi"
	res "github.com/be-sistem-informasi-konveksi/common/response/akuntansi"
	"github.com/be-sistem-informasi-konveksi/helper"
)

type AkuntansiUsecase interface {
	GetAllJU(ctx context.Context, reqGetAllJU req.GetAllJU) (res.JurnalUmumRes, error)
	GetAllBB(ctx context.Context, reqGetAllBB req.GetAllBB) ([]res.BukuBesarRes, error)
	GetAllNC(ctx context.Context, reqGetAllNC req.GetAllNC) (res.NeracaSaldoRes, error)
}

type akuntansiUsecase struct {
	repo repo.AkuntansiRepo
}

func NewAkuntansiUsecase(repo repo.AkuntansiRepo) AkuntansiUsecase {
	return &akuntansiUsecase{repo}
}

func getLastDateOfPreviousMonth(startDate string) (string, error) {
	// Parse the start date string
	startDateObj, err := time.Parse(time.DateOnly, startDate)
	if err != nil {
		return "", err
	}

	// Calculate the first day of the current month
	startDateObj = startDateObj.AddDate(0, 0, 1-startDateObj.Day())

	// Calculate the last day of the previous month
	lastDayOfPreviousMonth := startDateObj.Add(-1)

	// Format the result as a string
	result := lastDayOfPreviousMonth.Format("2006-01-02")

	return result, nil
}

func (u *akuntansiUsecase) GetAllJU(ctx context.Context, reqGetAllJU req.GetAllJU) (res.JurnalUmumRes, error) {
	startDate, err := time.Parse(time.DateOnly, reqGetAllJU.StartDate)
	if err != nil {
		helper.LogsError(err)
		return res.JurnalUmumRes{}, err
	}

	endDate, err := time.Parse(time.DateOnly, reqGetAllJU.EndDate)
	if err != nil {
		helper.LogsError(err)
		return res.JurnalUmumRes{}, err
	}

	dataJU, err := u.repo.GetDataJU(ctx, startDate, endDate)
	if err != nil {
		return res.JurnalUmumRes{}, err
	}

	dataTransaksisMap := map[string]res.DataTransaksiJU{}
	var totalDebit, totalKredit float64
	for _, v := range dataJU {
		dataTr, exist := dataTransaksisMap[v.TransaksiID]
		ayatJurnalJU := res.DataAyatJurnalJU{
			AyatJurnalID: v.AyatJurnalID,
			AkunID:       v.AkunID,
			KodeAkun:     v.KodeAkun,
			NamaAkun:     v.NamaAkun,
			Debit:        v.Debit,
			Kredit:       v.Kredit,
		}
		if !exist {
			dataTr = res.DataTransaksiJU{
				TransaksiID: v.TransaksiID,
				Tanggal:     v.Tanggal,
				Keterangan:  v.Keterangan,
			}
		}
		totalKredit += math.Abs(v.Kredit)
		totalDebit += math.Abs(v.Debit)

		dataTr.AyatJurnal = append(dataTr.AyatJurnal, ayatJurnalJU)

		dataTransaksisMap[v.TransaksiID] = dataTr
	}

	sliceDataTrJU := make([]res.DataTransaksiJU, len(dataTransaksisMap))
	i := 0
	for _, v := range dataTransaksisMap {
		sliceDataTrJU[i] = v
		i++
	}

	var res res.JurnalUmumRes
	res.TotalDebit = totalDebit
	res.TotalKredit = totalKredit
	res.Transaksi = sliceDataTrJU

	return res, nil
}

func (u *akuntansiUsecase) GetAllBB(ctx context.Context, reqGetAllBB req.GetAllBB) ([]res.BukuBesarRes, error) {
	startDate, err := time.Parse(time.DateOnly, reqGetAllBB.StartDate)
	if err != nil {
		helper.LogsError(err)
		return nil, err
	}

	endDate, err := time.Parse(time.DateOnly, reqGetAllBB.EndDate)
	if err != nil {
		helper.LogsError(err)
		return nil, err
	}
	dataAyatJurnals, dataSaldoAwalJurnals, err := u.repo.GetDataBB(ctx, reqGetAllBB.AkunID, startDate, endDate)
	if err != nil {
		return nil, err
	}

	akunDataMap := map[string]res.BukuBesarRes{}
	tanggalSaldoAwal, _ := getLastDateOfPreviousMonth(reqGetAllBB.StartDate)
	for _, v := range dataSaldoAwalJurnals {
		dataAkun, exist := akunDataMap[v.AkunID]
		ayatSaldoAwal := res.DataAyatJurnalBB{
			Tanggal:    tanggalSaldoAwal,
			Keterangan: "saldo awal",
			Saldo:      v.Saldo,
		}
		if !exist {
			dataAkun = res.BukuBesarRes{
				NamaAkun:    v.NamaAkun,
				KodeAkun:    v.KodeAkun,
				SaldoNormal: v.SaldoNormal,
			}
		}
		absSaldo := math.Abs(v.Saldo)
		if v.SaldoNormal == "DEBIT" {
			if v.Saldo < 0 {
				dataAkun.TotalKredit += absSaldo
				ayatSaldoAwal.Kredit = absSaldo
			} else {
				dataAkun.TotalDebit += absSaldo
				ayatSaldoAwal.Debit = absSaldo
			}
		} else {
			if v.Saldo > 0 {
				dataAkun.TotalKredit += absSaldo
				ayatSaldoAwal.Kredit = absSaldo
			} else {
				dataAkun.TotalDebit += absSaldo
				ayatSaldoAwal.Debit = absSaldo
			}
		}
		dataAkun.TotalSaldo += v.Saldo
		dataAkun.AyatJurnals = append(dataAkun.AyatJurnals, ayatSaldoAwal)
		akunDataMap[v.AkunID] = dataAkun
	}

	// fmt.Println(ayatJurnals)
	for _, ay := range dataAyatJurnals {
		dataAkun, exist := akunDataMap[ay.AkunID]
		ayatJurnal := res.DataAyatJurnalBB{
			TransaksiID: ay.TransaksiID,
			Tanggal:     ay.Tanggal,
			Keterangan:  ay.Keterangan,
			Kredit:      ay.Kredit,
			Debit:       ay.Debit,
		}

		if !exist {
			dataAkun = res.BukuBesarRes{
				NamaAkun:    ay.NamaAkun,
				SaldoNormal: ay.SaldoNormal,
				KodeAkun:    ay.KodeAkun,
			}
		}

		dataAkun.TotalKredit += ay.Kredit
		dataAkun.TotalDebit += ay.Debit
		dataAkun.TotalSaldo += ay.Saldo
		ayatJurnal.Saldo = dataAkun.TotalSaldo

		dataAkun.AyatJurnals = append(dataAkun.AyatJurnals, ayatJurnal)

		akunDataMap[ay.AkunID] = dataAkun

	}

	bukuBesarRes := make([]res.BukuBesarRes, len(akunDataMap))
	i := 0
	for _, value := range akunDataMap {
		bukuBesarRes[i] = value
		i++
	}
	return bukuBesarRes, nil
}

func (u *akuntansiUsecase) GetAllNC(ctx context.Context, reqGetAllNC req.GetAllNC) (res.NeracaSaldoRes, error) {
	date, err := time.Parse("2006-01", reqGetAllNC.Date)
	if err != nil {
		helper.LogsError(err)
		return res.NeracaSaldoRes{}, err
	}
	dataAkunsSaldo, err := u.repo.GetDataNC(ctx, date)
	if err != nil {
		return res.NeracaSaldoRes{}, err
	}
	nsRes := res.NeracaSaldoRes{}

	dataSaldoAkuns := make([]res.DataSaldoAkun, len(dataAkunsSaldo))
	i := 0
	for _, v := range dataAkunsSaldo {
		nsRes.TotalDebit += v.SaldoDebit
		nsRes.TotalKredit += v.SaldoKredit
		dataSaldoAkuns[i] = res.DataSaldoAkun{
			KodeAkun:    v.KodeAkun,
			NamaAkun:    v.NamaAkun,
			SaldoDebit:  v.SaldoDebit,
			SaldoKredit: v.SaldoKredit,
		}
		i++
	}
	nsRes.DataSaldoAkuns = dataSaldoAkuns

	return nsRes, nil
}
