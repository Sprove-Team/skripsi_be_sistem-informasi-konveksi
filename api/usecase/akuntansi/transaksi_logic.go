package uc_akuntansi

import (
	"errors"

	"github.com/be-sistem-informasi-konveksi/common/message"

	"github.com/be-sistem-informasi-konveksi/entity"
)

// transaksi
func UpdateSaldo(ayKredit, ayDebit float64, saldoNormal string) (saldo float64) {
	if saldoNormal == "DEBIT" {
		saldo = ayDebit - ayKredit
	}
	if saldoNormal == "KREDIT" {
		saldo = ayKredit - ayDebit
	}
	return
}

func IsSameReqAyJurnals(a, b []entity.AyatJurnal) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func IsKreditEqualToDebit(ayatJurnals []entity.AyatJurnal, totalTransaksi *float64) error {
	var totalKredit, totalDebit float64
	for _, v := range ayatJurnals {
		if v.Debit != 0 {
			totalDebit += v.Debit
		}
		if v.Kredit != 0 {
			totalKredit += v.Kredit
		}
	}

	if totalDebit != totalKredit {
		return errors.New(message.CreditDebitNotSame)
	}

	*totalTransaksi = totalDebit
	return nil
}

func IsDuplicateAkun(ayatJurnals []entity.AyatJurnal) error {
	akunCount := make(map[string]int)
	var debit, kredit float64

	for _, v := range ayatJurnals {
		akunCount[v.AkunID]++

		if akunCount[v.AkunID] > 1 {
			return errors.New(message.AkunCannotBeSame)
		}

		debit += v.Debit
		kredit += v.Kredit

		if debit == kredit {
			// Reset counters and map for the next pair
			debit = 0
			kredit = 0
			akunCount = make(map[string]int)
		}
	}

	return nil
}
