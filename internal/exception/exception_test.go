package exception_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
	"github.com/viniosilva/swordhealth-api/internal/dto"
	"github.com/viniosilva/swordhealth-api/internal/exception"
)

func TestErrorFormatBindingErrors(t *testing.T) {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("enum", func(fl validator.FieldLevel) bool { return true })
	}

	var cases = map[string]struct {
		inputError     error
		expectedErrors string
	}{
		`should return "invalid payload"`: {
			inputError:     fmt.Errorf("error"),
			expectedErrors: "invalid payload",
		},
		"should return list errors": {
			inputError: binding.Validator.ValidateStruct(&dto.CreateUserDto{
				Username: "",
				Email:    "email",
				Password: "0",
			}),
			expectedErrors: strings.Join([]string{
				"Key: 'CreateUserDto.Username' Error:Field validation for 'Username' failed on the 'required' tag",
				"Key: 'CreateUserDto.Email' Error:Field validation for 'Email' failed on the 'email' tag",
				"Key: 'CreateUserDto.Password' Error:Field validation for 'Password' failed on the 'min' tag",
			}, "; "),
		},
	}
	for name, cs := range cases {
		t.Run(name, func(t *testing.T) {
			// when
			errors := exception.FormatBindingErrors(cs.inputError)

			// then
			assert.Equal(t, errors, cs.expectedErrors)
		})
	}
}

func BenchmarkErrorFormatBindingErrors(b *testing.B) {
	// given
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("enum", func(fl validator.FieldLevel) bool { return true })
	}

	inputErrors := binding.Validator.ValidateStruct(&dto.CreateUserDto{
		Username: "",
		Email:    "email",
		Password: "0",
	})

	// when
	for i := 0; i < b.N; i++ {
		exception.FormatBindingErrors(inputErrors)
	}
}
