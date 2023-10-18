package helper

import (
	"strings"

	"github.com/go-playground/validator/v10"
	id_translations "github.com/go-playground/validator/v10/translations/id"
)

type (
	ErrorResponse struct {
		ValueInput   interface{} `json:"value_input"`
		ErrorMessage string      `json:"error_message"`
	}

	xValidator struct {
		validator *validator.Validate
	}
)

type Validator interface {
	Validate(d interface{}) []ErrorResponse
}

func NewValidator() Validator {
	validate := validator.New()
	return &xValidator{validate}
}

func (x *xValidator) Validate(d interface{}) []ErrorResponse {
	trans := (&translator{}).Translator()

	id_translations.RegisterDefaultTranslations(x.validator, trans)

	validatorErrors := []ErrorResponse{}

	errs := x.validator.Struct(d)

	if errs != nil {
		for _, err := range errs.(validator.ValidationErrors) {
			var elem ErrorResponse
			elem.ValueInput = err.Value()
			elem.ErrorMessage = strings.ToLower(err.Translate(trans))
			validatorErrors = append(validatorErrors, elem)

		}
	}
	return validatorErrors
}
