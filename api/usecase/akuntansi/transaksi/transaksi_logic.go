package akuntansi

// import (
// 	"context"
// 	"errors"

// 	repoAkun "github.com/be-sistem-informasi-konveksi/api/repository/akuntansi/mysql/gorm/akun"
// 	"github.com/be-sistem-informasi-konveksi/common/message"
// 	req "github.com/be-sistem-informasi-konveksi/common/request/akuntansi/transaksi"
// 	"github.com/be-sistem-informasi-konveksi/entity"
// 	"github.com/be-sistem-informasi-konveksi/helper"
// )

// func updateSaldo(saldo *float64, amount float64, isDebit bool) {
// 	if isDebit {
// 		*saldo -= amount
// 	} else {
// 		*saldo += amount
// 	}
// }

// if v.Kredit != 0 {
// 				updateSaldo(&akun.Saldo, &saldo, v.Kredit, akun.SaldoNormal == "DEBIT")
// 			}
// 			if v.Debit != 0 {
// 				updateSaldo(&akun.Saldo, &saldo, v.Debit, akun.SaldoNormal == "KREDIT")
// 			}

// func isSameReqAyJurnals(a, b []req.ReqAyatJurnal) bool {
// 	if len(a) != len(b) {
// 		return false
// 	}
// 	for i := range a {
// 		if a[i] != b[i] {
// 			return false
// 		}
// 	}
// 	return true
// }

// func isKreditEqualToDebit(ayatJurnals []req.ReqAyatJurnal, totalTransaksi *float64) error {
// 	var totalKredit, totalDebit float64
// 	for _, v := range ayatJurnals {
// 		if v.Debit != 0 {
// 			totalDebit += v.Debit
// 		}
// 		if v.Kredit != 0 {
// 			totalKredit += v.Kredit
// 		}
// 	}

// 	if totalDebit != totalKredit {
// 		return errors.New(message.CreditDebitNotSame)
// 	}

// 	*totalTransaksi = totalDebit
// 	return nil
// }

// func isDuplicateAkun(ayatJurnals []req.ReqAyatJurnal) error {
// 	akunCount := make(map[string]int)
// 	var debit, kredit float64

// 	for _, v := range ayatJurnals {
// 		akunCount[v.AkunID]++

// 		if akunCount[v.AkunID] > 1 {
// 			return errors.New(message.AkunCannotBeSame)
// 		}

// 		debit += v.Debit
// 		kredit += v.Kredit

// 		if debit == kredit {
// 			// Reset counters and map for the next pair
// 			debit = 0
// 			kredit = 0
// 			akunCount = make(map[string]int)
// 		}
// 	}

// 	return nil
// }

// // logic transaksi hutang piutang

// func setAyatJurnal(ctx context.Context, i int, id string, repoAkun repoAkun.AkunRepo, reqAy req.ReqAyatJurnal, ay *entity.AyatJurnal) error {
// 	akun, err := repoAkun.GetById(ctx, reqAy.AkunID)
// 	if err != nil {
// 		if err.Error() == "record not found" {
// 			return errors.New(message.AkunNotFound)
// 		}
// 		return err
// 	}

// 	var saldo float64
// 	if reqAy.Kredit != 0 {
// 		helper.UpdateSaldo(&saldo, reqAy.Kredit, akun.SaldoNormal == "DEBIT")
// 	}
// 	if reqAy.Debit != 0 {
// 		helper.UpdateSaldo(&saldo, reqAy.Debit, akun.SaldoNormal == "KREDIT")
// 	}
// 	*ay = entity.AyatJurnal{
// 		Base: entity.Base{
// 			ID: id,
// 		},
// 		AkunID: reqAy.AkunID,
// 		Debit:  reqAy.Debit,
// 		Kredit: reqAy.Kredit,
// 		Saldo:  saldo,
// 	}
// 	return nil
// }

// // func createHutangPiutang(ctx context.Context, reqTransaksi req.Create) error {

// // }
