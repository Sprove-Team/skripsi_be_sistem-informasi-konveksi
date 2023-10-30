package helper

import (
	"bytes"
	"strings"
	"unicode"

	"github.com/go-playground/validator/v10"
	id_translations "github.com/go-playground/validator/v10/translations/id"

	"github.com/be-sistem-informasi-konveksi/common/response"
)

type (
	xValidator struct {
		validator *validator.Validate
	}
)

type Validator interface {
	Validate(d interface{}) []response.BaseFormatError
}

func NewValidator() Validator {
	validate := validator.New()
	validate.RegisterValidation("uuidv4_no_hyphens", validateUUIDv4WithoutHyphens)
	return &xValidator{validate}
}

func validateUUIDv4WithoutHyphens(fl validator.FieldLevel) bool {
	uuid := fl.Field().String()
	uuidGoogle := NewGoogleUUID()
	return uuidGoogle.CheckValidUUID(uuid)
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

// func wordToSnake(s string) string {
// 	ss := strings.Split(s, " ")
// 	wg := &sync.WaitGroup{}
// 	for i, d := range ss {
// 		wg.Add(1)
// 		go func(i int, d string) {
// 			defer wg.Done()
// 			ss[i] = camelToSnake(d)
// 		}(i, d)
// 	}
// 	wg.Wait()
// 	return strings.Join(ss, " ")
// }

func (x *xValidator) Validate(d interface{}) []response.BaseFormatError {
	trans := (&translator{}).Translator()

	id_translations.RegisterDefaultTranslations(x.validator, trans)

	validatorErrors := []response.BaseFormatError{}

	errs := x.validator.Struct(d)

	if errs != nil {
		for _, err := range errs.(validator.ValidationErrors) {
			var elem response.BaseFormatError
			elem.ValueInput = err.Value()
			// fmt.Println(err.Tag())
			if err.Tag() == "uuidv4_no_hyphens" {
				elem.ErrorMessage = camelToSnake(err.Field()) + " tidak berupa uuid versi 4"
				validatorErrors = append(validatorErrors, elem)
				continue
			}
			// fmt.Println()
			elem.ErrorMessage = strings.ReplaceAll(err.Translate(trans), err.Field(), camelToSnake(err.Field()))
			validatorErrors = append(validatorErrors, elem)

		}
	}
	return validatorErrors
}
