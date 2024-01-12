package pkg

import (
	"bytes"
	"fmt"
	"net/url"
	"strings"
	"sync"
	"unicode"

	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	id_translations "github.com/go-playground/validator/v10/translations/id"

	"github.com/be-sistem-informasi-konveksi/common/response"
)

type (
	xValidator struct {
		validator *validator.Validate
		trans     ut.Translator
	}
)

type Validator interface {
	Validate(d interface{}) *response.BaseFormatError
}

func NewValidator() Validator {
	validate := validator.New()

	// custom validation

	validate.RegisterValidation("ulid", validateULID)
	validate.RegisterValidation("url_cloud_storage", validateGoogleStorageURL)

	// default translations
	trans := (&translator{}).Translator()
	id_translations.RegisterDefaultTranslations(validate, trans)

	// custom translations

	validate.RegisterTranslation("ulid", trans, func(ut ut.Translator) error {
		return ut.Add("ulid", "{0} tidak berupa ulid yang valid", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("ulid", camelToSnake(fe.Field()))
		return t
	})

	validate.RegisterTranslation("url_cloud_storage", trans, func(ut ut.Translator) error {
		return ut.Add("url_cloud_storage", "{0} tidak berupa url cloud storage yang valid", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("url_cloud_storage", camelToSnake(fe.Field()))
		return t
	})

	validate.RegisterTranslation("required_if", trans, func(ut ut.Translator) error {
		return ut.Add("required_if", "{0} wajib diisi ketika {1}", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		param := strings.Split(fe.Param(), " ")
		t, _ := ut.T("required_if", fe.Field(), camelToSnake(param[0])+" "+strings.ToLower(param[1]))
		return t
	})

	validate.RegisterTranslation("e164", trans, func(ut ut.Translator) error {
		return ut.Add("e164", "{0} harus berformat e164", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("e164", camelToSnake(fe.Field()))
		return t
	})

	validate.RegisterTranslation("datetime", trans, func(ut ut.Translator) error {
		return ut.Add("datetime", "{0} harus berformat {1}", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("datetime", camelToSnake(fe.Field()), fe.Param())
		return t
	})

	return &xValidator{validate, trans}
}

// custom validation

func validateGoogleStorageURL(fl validator.FieldLevel) bool {
	// Parse the URL

	inputUrl := fl.Field().String()
	u, err := url.Parse(inputUrl)
	if err != nil {
		return false
	}

	if u.Scheme != "https" {
		return false
	}

	if u.Hostname() != "storage.googleapis.com" {
		return false
	}
	return true
}

func validateULID(fl validator.FieldLevel) bool {
	ulid := fl.Field().String()
	ulidPkg := NewUlidPkg()
	_, err := ulidPkg.ParseStrict(ulid)
	return err == nil
}

func camelToSnake(s string) string {
	var buf bytes.Buffer
	var prevIsUpper bool
	for _, r := range s {
		if unicode.IsUpper(r) {
			if prevIsUpper {
				buf.WriteRune(unicode.ToLower(r))
			} else {
				if buf.Len() > 0 {
					buf.WriteRune('_')
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

func (x *xValidator) Validate(d interface{}) *response.BaseFormatError {
	validatorErrors := map[string][]string{} // Use a map for validation errors

	errs := x.validator.Struct(d)

	field := strPool.Get().(*str)
	fullField := strPool.Get().(*str)
	message := strPool.Get().(*str)
	namespace := make([]string, 0, 2)
	if errs != nil {
		for _, err := range errs.(validator.ValidationErrors) {
			// strings.Split(err.StructField())
			field.value = camelToSnake(err.Field())
			namespace = strings.Split(err.Namespace(), ".")
			if len(namespace) > 2 {
				// Add prefix "detail_invoice." for nested fields
				fullField.value = fmt.Sprintf("%s.%s", camelToSnake(namespace[1]), field.value)
				message.value = strings.ReplaceAll(err.Translate(x.trans), err.Field(), fullField.value)
				validatorErrors[fullField.value] = append(validatorErrors[fullField.value], message.value)
			} else {
				message.value = strings.ReplaceAll(err.Translate(x.trans), err.Field(), field.value)
				validatorErrors[field.value] = append(validatorErrors[field.value], message.value)
			}
			// Reset and put the objects back to the pool
			namespace = namespace[:0] // Reset slice without allocating new memory
			strPool.Put(field)
			strPool.Put(fullField)
			strPool.Put(message)
		}
	}

	if len(validatorErrors) > 0 {
		// You can customize the error code and status based on your requirements
		return response.ErrorRes(400, "Bad Request", validatorErrors)
	}

	return nil // No validation errors
}
