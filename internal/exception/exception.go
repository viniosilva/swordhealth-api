package exception

import (
	"strings"

	"github.com/go-playground/validator/v10"
)

func FormatBindingErrors(err error) string {
	if validationErrors, ok := err.(validator.ValidationErrors); !ok {
		return "invalid payload"
	} else {
		errors := []string{}
		for _, e := range validationErrors {
			errors = append(errors, e.Error())
		}

		return strings.Join(errors, "; ")
	}
}
