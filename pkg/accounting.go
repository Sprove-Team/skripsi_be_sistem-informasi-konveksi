package pkg

import "github.com/leekchan/accounting"

type AccountingPkg interface {
	New() accounting.Accounting
}

type accountingPkg struct{}

func NewAccounting() AccountingPkg {
	return &accountingPkg{}
}

func (a *accountingPkg) New() accounting.Accounting {
	ac := accounting.Accounting{Symbol: "Rp", Precision: 2}
	return ac
}
