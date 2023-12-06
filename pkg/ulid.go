package pkg

import "github.com/oklog/ulid/v2"

type UlidPkg interface {
	MakeUlid() ulid.ULID
	ParseStrict(ulidStr string) (ulid.ULID, error)
}

type ulidPkg struct{}

func NewUlidPkg() UlidPkg {
	return &ulidPkg{}
}

func (u *ulidPkg) MakeUlid() ulid.ULID {
	return ulid.Make()
}

func (u *ulidPkg) ParseStrict(ulidStr string) (ulid.ULID, error) {
	return ulid.ParseStrict(ulidStr)
}
