package exception

import (
	"strings"

	"github.com/go-playground/validator/v10"
)

type error interface {
	Error() string
}

type ForeignKeyConstraintException struct {
	Message string
}

func (impl *ForeignKeyConstraintException) Error() string {
	return impl.Message
}

type NotFoundException struct {
	Message string
}

func (impl *NotFoundException) Error() string {
	return impl.Message
}

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
