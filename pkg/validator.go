package pkg

import (
	"bytes"
	"strings"
	"sync"
	"unicode"

	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	id_translations "github.com/go-playground/validator/v10/translations/id"

	res_global "github.com/be-sistem-informasi-konveksi/common/response"
)

type (
	xValidator struct {
		validator *validator.Validate
		trans     ut.Translator
	}
)

type Validator interface {
	Validate(d interface{}) *res_global.BaseFormatError
}

func NewValidator() Validator {
	validate := validator.New()

	// custom validation

	validate.RegisterValidation("ulid", validateULID)
	// validate.RegisterValidation("url_cloud_storage", validateGoogleStorageURL)

	// default translations
	trans := (&translator{}).Translator()
	id_translations.RegisterDefaultTranslations(validate, trans)

	// validate.RegisterTranslation("url_cloud_storage", trans, func(ut ut.Translator) error {
	// 	return ut.Add("url_cloud_storage", "{0} tidak berupa url cloud storage yang valid", true)
	// }, func(ut ut.Translator, fe validator.FieldError) string {
	// 	t, _ := ut.T("url_cloud_storage", camelToSnake(fe.Field()))
	// 	return t
	// })

	// custom translations
	validate.RegisterTranslation("ulid", trans, func(ut ut.Translator) error {
		return ut.Add("ulid", "{0} tidak berupa ulid yang valid", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("ulid", fe.Field())
		return t
	})

	validate.RegisterTranslation("required_if", trans, func(ut ut.Translator) error {
		return ut.Add("required_if", "{0} wajib diisi ketika {1}", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		param := strings.Split(fe.Param(), " ")
		t, _ := ut.T("required_if", fe.Field(), convToReadAble(param[0])+" "+strings.ToLower(param[1]))
		return t
	})

	validate.RegisterTranslation("e164", trans, func(ut ut.Translator) error {
		return ut.Add("e164", "{0} harus berformat e164", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("e164", fe.Field())
		return t
	})

	validate.RegisterTranslation("datetime", trans, func(ut ut.Translator) error {
		return ut.Add("datetime", "{0} harus berformat {1}", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		var format string
		switch fe.Param() {
		case "2006-01-02T15:04:05Z07:00":
			format = "RFC3999"
		case "2006-01-02":
			format = "Tahun-Bulan-Tanggal"
		default:
			format = fe.Param()
		}
		t, _ := ut.T("datetime", fe.Field(), format)
		return t
	})

	validate.RegisterTranslation("required_without", trans, func(ut ut.Translator) error {
		return ut.Add("required_without", "{0} wajib diisi jika {1} tidak diisi", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("required_without", fe.Field(), convToReadAble(fe.Param()))
		return t
	})

	validate.RegisterTranslation("oneof", trans, func(ut ut.Translator) error {
		return ut.Add("oneof", "{0} harus berupa salah satu dari [{1}]", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("oneof", fe.Field(), strings.ReplaceAll(fe.Param(), " ", ","))
		return t
	})

	validate.RegisterTranslation("required_with", trans, func(ut ut.Translator) error {
		return ut.Add("required_with", "{0} wajib diisi jika {1} diisi", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		data := strings.Split(fe.Param(), " ")
		length := len(data)
		var param string = convToReadAble(data[0])
		if length > 1 {
			param += " atau "
			for i := 1; i < length; i++ {
				dat := convToReadAble(data[i])
				if i+1 == length {
					param += dat
					continue
				}
				param += dat + " atau "
			}
		}
		t, _ := ut.T("required_with", fe.Field(), param)
		return t
	})

	validate.RegisterTranslation("excluded_with", trans, func(ut ut.Translator) error {
		return ut.Add("excluded_with", "jika {0} diisi maka {1} harus kosong", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("excluded_with", fe.Field(), convToReadAble(fe.Param()))
		return t
	})

	return &xValidator{validate, trans}
}

// custom validation

// func validateGoogleStorageURL(fl validator.FieldLevel) bool {
// Parse the URL

// 	inputUrl := fl.Field().String()
// 	u, err := url.Parse(inputUrl)
// 	if err != nil {
// 		return false
// 	}

// 	if u.Scheme != "https" {
// 		return false
// 	}

// 	if u.Hostname() != "storage.googleapis.com" {
// 		return false
// 	}
// 	return true
// }

func validateULID(fl validator.FieldLevel) bool {
	ulid := fl.Field().String()
	ulidPkg := NewUlidPkg()
	_, err := ulidPkg.ParseStrict(ulid)
	return err == nil
}

func convToReadAble(s string) string {
	var buf bytes.Buffer
	var prevIsUpper bool
	for _, r := range s {
		if unicode.IsNumber(r) || !unicode.IsLetter(r) {
			continue
		}
		if unicode.IsUpper(r) {
			if prevIsUpper {
				buf.WriteRune(unicode.ToLower(r))
			} else {
				if buf.Len() > 0 {
					buf.WriteRune(' ')
				}
				buf.WriteRune(unicode.ToLower(r))
				prevIsUpper = true
			}
		} else {
			buf.WriteRune(r)
			prevIsUpper = false
		}
	}

	return buf.String()
}

type str struct {
	value string
}

var strPool = sync.Pool{
	New: func() any {
		return new(str)
	},
}

func (x *xValidator) Validate(d interface{}) *res_global.BaseFormatError {

	errs := x.validator.Struct(d)

	field := strPool.Get().(*str)
	param := strPool.Get().(*str)
	message := strPool.Get().(*str)
	if errs != nil {
		errosMsg := make([]string, len(errs.(validator.ValidationErrors)))
		for i, err := range errs.(validator.ValidationErrors) {
			field.value = convToReadAble(err.Field())
			message.value = strings.ReplaceAll(err.Translate(x.trans), err.Field(), field.value)
			if err.Param() != "" && err.Param() == "2006" {
				message.value = strings.ReplaceAll(message.value, ".", "")
			}
			errosMsg[i] = message.value
			strPool.Put(field)
			strPool.Put(param)
			strPool.Put(message)
		}
		return res_global.ErrorRes(400, "Bad Request", errosMsg)

	}

	return nil
}
