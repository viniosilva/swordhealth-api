package service_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/viniosilva/swordhealth-api/internal/service"
	"github.com/viniosilva/swordhealth-api/mock"
)

func TestAuthServiceLogin(t *testing.T) {
	var cases = map[string]struct {
		inputBasicAuth   string
		mocking          func(userRepository *mock.MockUserRepository, cryptoService *mock.MockCryptoService)
		expectedUsername string
		expectedPassword string
		expectedErr      error
	}{
		"should return username and password": {
			inputBasicAuth:   "Basic dXNlcm5hbWU6MTEyMjMzNDQ1NQ==",
			expectedUsername: "username",
			expectedPassword: "1122334455",
		},
		"should throw invalid authorization error": {
			inputBasicAuth: "auth",
			expectedErr:    fmt.Errorf("invalid authorization"),
		},
		"should throw error on decode auth": {
			inputBasicAuth: "Basic ",
			expectedErr:    fmt.Errorf("invalid authorization"),
		},
	}
	for name, cs := range cases {
		t.Run(name, func(t *testing.T) {
			// given
			ctx := context.Background()
			authService := service.NewAuthService()

			// when
			username, password, err := authService.DecodeBasicAuth(ctx, cs.inputBasicAuth)

			// then
			assert.Equal(t, cs.expectedErr, err)
			assert.Equal(t, cs.expectedUsername, username)
			assert.Equal(t, cs.expectedPassword, password)
		})
	}
}

func BenchmarkAuthServiceLogin(b *testing.B) {
	// given
	ctx := context.Background()

	authService := service.NewAuthService()

	// when
	for i := 0; i < b.N; i++ {
		authService.DecodeBasicAuth(ctx, "Basic dXNlcm5hbWU6MTEyMjMzNDQ1NQ==")
	}
}
