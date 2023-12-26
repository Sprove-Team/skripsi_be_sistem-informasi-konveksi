package akuntansi

import (
	"errors"
	"sync"

	"github.com/be-sistem-informasi-konveksi/common/message"
	req "github.com/be-sistem-informasi-konveksi/common/request/akuntansi/transaksi"
)

func updateSaldo(saldoAkun *float64, saldo *float64, amount float64, isDebit bool) {
	if isDebit {
		*saldoAkun -= amount
		*saldo -= amount
	} else {
		*saldoAkun += amount
		*saldo += amount
	}
}

func isSameReqAyJurnals(a, b []req.ReqAyatJurnal) bool {
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

func isKreditEqualToDebit(ayatJurnals []req.ReqAyatJurnal, totalTransaksi *float64, wg *sync.WaitGroup, m *sync.Mutex) error {
	var totalKredit, totalDebit float64
	for _, v := range ayatJurnals {
		wg.Add(1)
		go func(v req.ReqAyatJurnal) {
			defer wg.Done()
			if v.Debit != 0 {
				m.Lock()
				totalDebit += v.Debit
				m.Unlock()
			}
			if v.Kredit != 0 {
				m.Lock()
				totalKredit += v.Kredit
				m.Unlock()
			}
		}(v)
	}
	wg.Wait()

	if totalDebit != totalKredit {
		return errors.New(message.CreditDebitNotSame)
	}

	*totalTransaksi = totalDebit
	return nil
}

func isDuplicateAkun(ayatJurnals []req.ReqAyatJurnal) error {
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
