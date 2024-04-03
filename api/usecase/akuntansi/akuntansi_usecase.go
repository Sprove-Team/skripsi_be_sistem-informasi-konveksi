package uc_akuntansi

import (
	"bytes"
	"context"
	"fmt"
	"math"
	"strconv"
	"time"
	"unicode/utf8"

	// "github.com/360EntSecGroup-Skylar/excelize"
	repo "github.com/be-sistem-informasi-konveksi/api/repository/akuntansi/mysql/gorm"
	req "github.com/be-sistem-informasi-konveksi/common/request/akuntansi"
	res "github.com/be-sistem-informasi-konveksi/common/response/akuntansi"
	"github.com/be-sistem-informasi-konveksi/helper"
	"github.com/be-sistem-informasi-konveksi/pkg"
)

type AkuntansiUsecase interface {
	GetAllJU(ctx context.Context, reqGetAllJU req.GetAllJU) (res.JurnalUmumRes, error)
	DownloadJU(req req.GetAllJU, JU res.JurnalUmumRes) (*bytes.Buffer, error)
	DownloadBB(req req.GetAllBB, BB []res.BukuBesarRes) (*bytes.Buffer, error)
	GetAllBB(ctx context.Context, reqGetAllBB req.GetAllBB) ([]res.BukuBesarRes, error)
	GetAllNC(ctx context.Context, reqGetAllNC req.GetAllNC) (res.NeracaSaldoRes, error)
	GetAllLBR(ctx context.Context, reqGetAllLBR req.GetAllLBR) ([]res.LabaRugiRes, error)
}

type akuntansiUsecase struct {
	repo     repo.AkuntansiRepo
	excelize pkg.ExcelizePkg
}

func NewAkuntansiUsecase(repo repo.AkuntansiRepo, excelize pkg.ExcelizePkg) AkuntansiUsecase {
	return &akuntansiUsecase{repo, excelize}
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
			ID:       v.AyatJurnalID,
			AkunID:   v.AkunID,
			KodeAkun: v.KodeAkun,
			NamaAkun: v.NamaAkun,
			Debit:    v.Debit,
			Kredit:   v.Kredit,
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

func (u *akuntansiUsecase) DownloadJU(req req.GetAllJU, JU res.JurnalUmumRes) (*bytes.Buffer, error) {
	f, sheetName := u.excelize.InitExcelize("Jurnal Umum")
	// Create a new Excel file
	startDate, _ := time.Parse(time.DateOnly, req.StartDate)
	endDate, _ := time.Parse(time.DateOnly, req.EndDate)
	title := fmt.Sprintf("Dalam Rupiah (%s - %s)", startDate.Format("01 Jan 2006"), endDate.Format("01 Jan 2006"))
	f.SetCellValue(sheetName, "A1", title)
	f.SetCellStyle(sheetName, "A1", "A1", u.excelize.StyleBold(f))

	// Create headers
	headers := []string{"Tanggal", "Keterangan", "Kode Akun", "Nama Akun", "Debit", "Kredit"}
	// Apply gray styleHeader to header
	maxLength := make([]int, len(headers))
	alpha := "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	styleHeader := u.excelize.StyleHeader(f, "center", "#E0E0E0")
	for i, header := range headers {
		cell := string(alpha[i]) + "3"
		f.SetCellValue(sheetName, cell, header)
		if cellWidth := utf8.RuneCountInString(header) + 2; maxLength[i] < cellWidth {
			maxLength[i] = cellWidth
		}
		f.SetCellStyle(sheetName, cell, cell, styleHeader)
	}

	// Freeze the header row
	// `{"freeze":true,"split":false,"x_split":0,"y_split":3,"top_left_cell":"A4"}`
	// xlsx.SetPanes(sheetName, )
	// Populate data
	row := 4 // Start from row 3
	formatCurrency := u.excelize.StyleCurrencyRpIndo(f, false, "", "", false, "")
	for _, transaksi := range JU.Transaksi {
		for _, ayat := range transaksi.AyatJurnal {
			rowStr := strconv.Itoa(row)
			parse, _ := time.Parse(time.RFC3339, transaksi.Tanggal)
			tanggal := parse.Format(time.DateTime)
			f.SetCellValue(sheetName, "A"+rowStr, tanggal)
			if cellWidth := u.excelize.GetCellWidth(tanggal, 2); maxLength[0] < cellWidth {
				maxLength[0] = cellWidth
			}
			f.SetCellValue(sheetName, "B"+rowStr, transaksi.Keterangan)
			if cellWidth := u.excelize.GetCellWidth(transaksi.Keterangan, 2); maxLength[1] < cellWidth {
				maxLength[1] = cellWidth
			}
			f.SetCellValue(sheetName, "C"+rowStr, string(ayat.KodeAkun))
			if cellWidth := u.excelize.GetCellWidth(ayat.KodeAkun, 2); maxLength[2] < cellWidth {
				maxLength[2] = cellWidth
			}
			f.SetCellValue(sheetName, "D"+rowStr, ayat.NamaAkun)
			if cellWidth := u.excelize.GetCellWidth(ayat.NamaAkun, 2) + 2; maxLength[3] < cellWidth {
				maxLength[3] = cellWidth
			}
			f.SetCellValue(sheetName, "E"+rowStr, ayat.Debit)
			f.SetCellStyle(sheetName, "E"+rowStr, "E"+rowStr, formatCurrency)
			debit, _ := f.GetCellValue(sheetName, "E"+rowStr)
			if cellWidth := u.excelize.GetCellWidth(debit, 2) + 2; maxLength[4] < cellWidth {
				maxLength[4] = cellWidth
			}
			f.SetCellValue(sheetName, "F"+rowStr, ayat.Kredit)
			f.SetCellStyle(sheetName, "F"+rowStr, "F"+rowStr, formatCurrency)
			kredit, _ := f.GetCellValue(sheetName, "F"+rowStr)
			if cellWidth := u.excelize.GetCellWidth(kredit, 2) + 2; maxLength[5] < cellWidth {
				maxLength[5] = cellWidth
			}
			row++
		}
	}
	for i, v := range maxLength {
		cell := string(alpha[i])
		f.SetColWidth(sheetName, cell, cell, float64(v))
	}
	rowStr := strconv.Itoa(row)
	f.MergeCell(sheetName, "A"+rowStr, "D"+rowStr)
	f.SetCellValue(sheetName, "A"+rowStr, "Total")
	f.SetCellStyle(sheetName, "A"+rowStr, "A"+rowStr, styleHeader)

	formatCurrencyWithColor := u.excelize.StyleCurrencyRpIndo(f, true, "", "", true, "#E0E0E0")
	f.SetCellFormula(sheetName, "E"+rowStr, fmt.Sprintf("=SUM(E4:E%d)", row-1))
	f.SetCellStyle(sheetName, "E"+rowStr, "E"+rowStr, formatCurrencyWithColor)

	f.SetCellFormula(sheetName, "F"+rowStr, fmt.Sprintf("=SUM(F4:F%d)", row-1))
	f.SetCellStyle(sheetName, "F"+rowStr, "F"+rowStr, formatCurrencyWithColor)
	buf, err := f.WriteToBuffer()
	if err != nil {
		helper.LogsError(err)
		return nil, err
	}
	return buf, nil
}

func (u *akuntansiUsecase) DownloadBB(req req.GetAllBB, BB []res.BukuBesarRes) (*bytes.Buffer, error) {
	f, sheetName := u.excelize.InitExcelize("Buku Besar")
	// Create a new Excel file
	startDate, _ := time.Parse(time.DateOnly, req.StartDate)
	endDate, _ := time.Parse(time.DateOnly, req.EndDate)
	title := fmt.Sprintf("Dalam IDR (%s - %s)", startDate.Format("01 Jan 2006"), endDate.Format("01 Jan 2006"))
	f.SetCellValue(sheetName, "A1", title)

	// Make the first cell bold
	styleBold := u.excelize.StyleBold(f)

	f.SetCellStyle(sheetName, "A1", "A1", styleBold)

	alpha := "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	// Populate data
	row := 4 // Start from row 3
	headers := []string{"Kode Akun", "Nama Akun", "Saldo Normal"}
	// style
	formatCurrency := u.excelize.StyleCurrencyRpIndo(f, false, "", "", false, "")
	styleHeader := u.excelize.StyleHeader(f, "center", "#E0E0E0")
	formatCurrencyWithColor := u.excelize.StyleCurrencyRpIndo(f, true, "", "", true, "#E0E0E0")
	formatCurrencyWithColorCenter := u.excelize.StyleCurrencyRpIndo(f, true, "", "center", true, "#E0E0E0")
	formatCurrencyWithColorRedCenter := u.excelize.StyleCurrencyRpIndo(f, true, "", "center", true, "#E6340E")

	for _, bb := range BB {
		// rowStr := strconv.Itoa(row)
		for _, header := range headers {
			cell := "A" + strconv.Itoa(row)
			cell2 := "B" + strconv.Itoa(row)
			f.SetCellValue(sheetName, cell, header)
			f.SetCellStyle(sheetName, cell, cell, styleBold)
			switch header {
			case "Kode Akun":
				f.SetCellValue(sheetName, cell2, bb.KodeAkun)
			case "Nama Akun":
				f.SetCellValue(sheetName, cell2, bb.NamaAkun)
			case "Saldo Normal":
				f.SetCellValue(sheetName, cell2, bb.SaldoNormal)
			}
			row++
		}
		row++ // margin

		headers2 := []string{"Tanggal", "Keterangan", "Debit", "Kredit", "Saldo"}
		maxLength := make([]int, len(headers2))
		for j, header := range headers2 {
			cell := string(alpha[j]) + strconv.Itoa(row)
			f.SetCellValue(sheetName, cell, header)
			if cellWidth := u.excelize.GetCellWidth(header, 2); maxLength[j] < cellWidth {
				maxLength[j] = cellWidth
			}
			// Apply gray style to header
			f.SetCellStyle(sheetName, cell, cell, styleHeader)
		}
		row++
		rowStart := row
		for _, content := range bb.AyatJurnal {
			parse, _ := time.Parse(time.RFC3339, content.Tanggal)
			tanggal := parse.Format(time.DateTime)
			rowStr := strconv.Itoa(row)
			f.SetCellValue(sheetName, "A"+rowStr, tanggal)
			if cellWidth := u.excelize.GetCellWidth(tanggal, 4); cellWidth > maxLength[0] {
				maxLength[0] = cellWidth
			}
			f.SetCellValue(sheetName, "B"+rowStr, content.Keterangan)
			if cellWidth := u.excelize.GetCellWidth(content.Keterangan, 4); cellWidth > maxLength[1] {
				maxLength[1] = cellWidth
			}
			f.SetCellValue(sheetName, "C"+rowStr, content.Debit)
			f.SetCellStyle(sheetName, "C"+rowStr, "C"+rowStr, formatCurrency)
			debit, _ := f.GetCellValue(sheetName, "C"+rowStr)
			if cellWidth := u.excelize.GetCellWidth(debit, 4); cellWidth > maxLength[2] {
				maxLength[2] = cellWidth
			}
			f.SetCellValue(sheetName, "D"+rowStr, content.Kredit)
			f.SetCellStyle(sheetName, "D"+rowStr, "D"+rowStr, formatCurrency)
			kredit, _ := f.GetCellValue(sheetName, "D"+rowStr)
			if cellWidth := u.excelize.GetCellWidth(kredit, 4); cellWidth > maxLength[3] {
				maxLength[3] = cellWidth
			}
			f.SetCellValue(sheetName, "E"+rowStr, content.Saldo)
			f.SetCellStyle(sheetName, "E"+rowStr, "E"+rowStr, formatCurrency)
			saldo, _ := f.GetCellValue(sheetName, "E"+rowStr)
			if cellWidth := u.excelize.GetCellWidth(saldo, 4); cellWidth > maxLength[4] {
				maxLength[4] = cellWidth
			}
			row++
		}
		{
			rowStr := strconv.Itoa(row)
			f.MergeCell(sheetName, "A"+rowStr, "B"+rowStr)
			f.SetCellValue(sheetName, "A"+rowStr, "Total")
			f.SetCellStyle(sheetName, "A"+rowStr, "A"+rowStr, styleHeader)
			//? total debit
			cellTotalDebit := "C" + rowStr
			f.SetCellValue(sheetName, cellTotalDebit, bb.TotalDebit)
			f.SetCellStyle(sheetName, cellTotalDebit, cellTotalDebit, formatCurrencyWithColor)
			totalDebit, _ := f.GetCellValue(sheetName, cellTotalDebit)
			if cellWidth := u.excelize.GetCellWidth(totalDebit, 4); cellWidth > maxLength[2] {
				maxLength[2] = cellWidth
			}

			//? total kredit
			cellTotalKredit := "D" + rowStr
			f.SetCellValue(sheetName, cellTotalKredit, bb.TotalKredit)
			f.SetCellStyle(sheetName, cellTotalKredit, cellTotalKredit, formatCurrencyWithColor)
			totalKredit, _ := f.GetCellValue(sheetName, cellTotalKredit)
			if cellWidth := u.excelize.GetCellWidth(totalKredit, 4); cellWidth > maxLength[3] {
				maxLength[3] = cellWidth
			}

			//? total saldo
			cellTotalSaldo := "E" + rowStr
			f.SetCellFormula(sheetName, cellTotalSaldo, fmt.Sprintf("=SUM(E%d:E%d)", rowStart, row-1))
			f.SetCellStyle(sheetName, cellTotalSaldo, cellTotalSaldo, formatCurrencyWithColor)
			totalSaldo, _ := f.GetCellValue(sheetName, cellTotalSaldo)
			if cellWidth := u.excelize.GetCellWidth(totalSaldo, 4); cellWidth > maxLength[4] {
				maxLength[4] = cellWidth
			}
			row++
		}
		{
			//? total saldo yang telah di kurangi debit dan kredit
			rowStr := strconv.Itoa(row)
			f.MergeCell(sheetName, "A"+rowStr, "B"+rowStr)
			f.SetCellValue(sheetName, "A"+rowStr, "Total Saldo Setelah Debit - Kredit (Note: Berdasarkan Saldo Normal)")
			f.SetCellStyle(sheetName, "A"+rowStr, "A"+rowStr, styleHeader)

			f.MergeCell(sheetName, "C"+rowStr, "D"+rowStr)
			f.SetCellValue(sheetName, "C"+rowStr, bb.TotalSaldo)

			if bb.TotalSaldo >= 0 {
				f.SetCellStyle(sheetName, "C"+rowStr, "C"+rowStr, formatCurrencyWithColorCenter)
			} else {
				f.SetCellStyle(sheetName, "C"+rowStr, "C"+rowStr, formatCurrencyWithColorRedCenter)
			}
		}
		for i, v := range maxLength {
			cell := string(alpha[i])
			f.SetColWidth(sheetName, cell, cell, float64(v))
		}
		row += 3 // margin
	}
	buf, err := f.WriteToBuffer()
	if err != nil {
		helper.LogsError(err)
		return nil, err
	}
	return buf, nil
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
		dataAkun.AyatJurnal = append(dataAkun.AyatJurnal, ayatSaldoAwal)
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

		dataAkun.AyatJurnal = append(dataAkun.AyatJurnal, ayatJurnal)

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
		// nsRes.TotalDebit += v.SaldoDebit
		// nsRes.TotalKredit += v.SaldoKredit
		dataSaldoAkuns[i] = res.DataSaldoAkun{
			KodeAkun: v.KodeAkun,
			NamaAkun: v.NamaAkun,
		}
		if v.SaldoNormal == "DEBIT" {
			nsRes.TotalDebit += v.Saldo
			dataSaldoAkuns[i].SaldoDebit = v.Saldo
		} else {
			nsRes.TotalKredit += v.Saldo
			dataSaldoAkuns[i].SaldoKredit = v.Saldo
		}
		i++
	}
	nsRes.DataSaldoAkuns = dataSaldoAkuns

	return nsRes, nil
}

func (u *akuntansiUsecase) GetAllLBR(ctx context.Context, reqGetAllLBR req.GetAllLBR) ([]res.LabaRugiRes, error) {
	startDate, err := time.Parse(time.DateOnly, reqGetAllLBR.StartDate)
	if err != nil {
		helper.LogsError(err)
		return nil, err
	}

	endDate, err := time.Parse(time.DateOnly, reqGetAllLBR.EndDate)
	if err != nil {
		helper.LogsError(err)
		return nil, err
	}

	lbrRes, err := u.repo.GetDataLBR(ctx, startDate, endDate)
	if err != nil {
		return nil, err
	}

	labaRugiMap := make(map[string]res.LabaRugiRes)
	var saldoKreditDebit float64
	for _, v := range lbrRes {
		labaRugi, ok := labaRugiMap[v.KategoriAkun]
		if !ok {
			labaRugi = res.LabaRugiRes{
				NamaKategori: v.KategoriAkun,
			}
		}
		labaRugi.Total += v.Saldo
		dataAkunLBR := res.DataAkunLBR{
			KodeAkun: v.KodeAkun,
			NamaAkun: v.NamaAkun,
			Saldo:    v.Saldo,
		}

		saldoKreditDebit = math.Abs(v.Saldo)

		if v.SaldoNormal == "DEBIT" {
			dataAkunLBR.SaldoDebit = saldoKreditDebit
		} else {
			dataAkunLBR.SaldoKredit = saldoKreditDebit
		}

		labaRugi.DataAkunLBR = append(labaRugi.DataAkunLBR, dataAkunLBR)

		labaRugiMap[v.KategoriAkun] = labaRugi
	}

	labaRugiRes := make([]res.LabaRugiRes, len(labaRugiMap))

	i := 0
	for _, v := range labaRugiMap {
		labaRugiRes[i] = v
		i++
	}

	return labaRugiRes, nil
}
