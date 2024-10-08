package pkg

import (
	"strings"

	"github.com/google/uuid"
)

type UuidGenerator interface {
	GenerateUUID() (string, error)
}

type googleUUID struct{}

func NewGoogleUUID() googleUUID {
	return googleUUID{}
}

func (g googleUUID) GenerateUUID() (string, error) {
	uuid, err := uuid.NewUUID()
	if err != nil {
		return "", err
	}
	uuidString := strings.ReplaceAll(uuid.String(), "-", "")
	return uuidString, nil
}

func (g googleUUID) CheckValidUUID(uids string) bool {
	_, err := uuid.Parse(uids)
	return err == nil
}
