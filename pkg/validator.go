package pkg

import (
	"bytes"
	"net/url"
	"strings"
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
	Validate(d interface{}) []response.BaseFormatError
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

func (x *xValidator) Validate(d interface{}) []response.BaseFormatError {
	validatorErrors := []response.BaseFormatError{}

	errs := x.validator.Struct(d)

	if errs != nil {
		for _, err := range errs.(validator.ValidationErrors) {
			var elem response.BaseFormatError
			field := camelToSnake(err.Field())
			elem.FieldName = field
			elem.Message = strings.ReplaceAll(err.Translate(x.trans), err.Field(), field)
			validatorErrors = append(validatorErrors, elem)
		}
	}
	return validatorErrors
}
