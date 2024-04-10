package pkg

import (
	"github.com/go-playground/locales/id"
	ut "github.com/go-playground/universal-translator"
)

type translator struct{}

func (t *translator) Translator() ut.Translator {
	id := id.New()
	uni := ut.New(id, id)
	trans, _ := uni.GetTranslator("id")
	return trans
}
