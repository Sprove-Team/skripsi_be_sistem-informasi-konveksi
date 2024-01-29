package akuntansi

import (
	"errors"
	"strings"

	"github.com/be-sistem-informasi-konveksi/common/message"
	"github.com/be-sistem-informasi-konveksi/entity"
)

func UpdateSaldo(saldo *float64, ayKredit, ayDebit float64, saldoNormal string) {
	if saldoNormal == "DEBIT" {
		*saldo = ayDebit - ayKredit
	}
	if saldoNormal == "KREDIT" {
		*saldo = ayKredit - ayDebit
	}
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

func IsValidAkunHutangPiutang(nama string) bool {
	validNames := []string{"piutang", "hutang", "kas & bank"}
	if strings.HasPrefix(nama, "pendapatan") {
		return true
	}
	if strings.HasPrefix(nama, "beban") {
		return true
	}
	for _, validName := range validNames {
		if nama == validName {
			return true
		}
	}
	return false
}
