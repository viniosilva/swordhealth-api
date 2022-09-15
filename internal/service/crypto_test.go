package service_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/viniosilva/swordhealth-api/internal/exception"
	"github.com/viniosilva/swordhealth-api/internal/model"
	"github.com/viniosilva/swordhealth-api/internal/service"
)

func TestCryptoServiceHash(t *testing.T) {
	var cases = map[string]struct {
		injectKey    string
		inputValue   string
		expectedHash string
	}{
		"should return hash": {
			injectKey:    "key",
			inputValue:   "S3cR31",
			expectedHash: "c70a5040e8f1bad417435911e93d030ac8894dd7af3fc613d0af7a59dd50ccc0",
		},
	}
	for name, cs := range cases {
		t.Run(name, func(t *testing.T) {
			// given
			cryptoService := service.NewCryptoService(cs.injectKey, "", 0)

			// when
			hash := cryptoService.Hash(cs.inputValue)

			// then
			assert.Equal(t, cs.expectedHash, hash)
		})
	}
}

func BenchmarkCryptoServiceHash(b *testing.B) {
	// given
	cryptoService := service.NewCryptoService("key", "", 0)

	// when
	for i := 0; i < b.N; i++ {
		cryptoService.Hash("value")
	}
}

func TestCryptoServiceEncryptJwt(t *testing.T) {
	var cases = map[string]struct {
		injectKey       string
		injectExpiresIn int64
		inputSub        int
		inputClaims     map[string]interface{}
		expectedError   error
	}{
		"should return accessToken": {
			injectKey:       "key",
			injectExpiresIn: 900000,
			inputSub:        1,
			inputClaims: map[string]interface{}{
				"username": "username",
				"role":     model.UserRoleTechnician,
			},
		},
	}
	for name, cs := range cases {
		t.Run(name, func(t *testing.T) {
			// given
			ctx := context.Background()
			cryptoService := service.NewCryptoService("", cs.injectKey, cs.injectExpiresIn)

			// when
			_, err := cryptoService.EncryptJwt(ctx, cs.inputSub, cs.inputClaims)

			// then
			assert.Equal(t, cs.expectedError, err)
		})
	}
}

func BenchmarkCryptoServiceEncryptJwt(b *testing.B) {
	// given
	ctx := context.Background()
	cryptoService := service.NewCryptoService("", "key", 900000)

	// when
	for i := 0; i < b.N; i++ {
		cryptoService.EncryptJwt(ctx, 1, map[string]interface{}{})
	}
}

func TestCryptoServiceDecryptJwt(t *testing.T) {
	var cases = map[string]struct {
		injectKey       string
		injectExpiresIn int64
		inputSub        int
		inputClaims     map[string]interface{}
		expectedClaims  map[string]interface{}
		expectedError   error
	}{
		"should return claim": {
			injectKey:       "key",
			injectExpiresIn: 900000,
			inputSub:        1,
			inputClaims: map[string]interface{}{
				"username": "username",
				"role":     model.UserRoleTechnician,
			},
			expectedClaims: map[string]interface{}{
				"exp":      float64(0),
				"iat":      float64(0),
				"sub":      float64(1),
				"username": "username",
				"role":     string(model.UserRoleTechnician),
			},
		},
		"should throw expired jwt exception": {
			injectKey:       "key",
			injectExpiresIn: 0,
			inputSub:        1,
			inputClaims: map[string]interface{}{
				"username": "username",
				"role":     model.UserRoleTechnician,
			},
			expectedError: &exception.ExpiredTokenException{Message: "token is expired"},
		},
	}
	for name, cs := range cases {
		t.Run(name, func(t *testing.T) {
			// given
			ctx := context.Background()
			cryptoService := service.NewCryptoService("", cs.injectKey, cs.injectExpiresIn)
			accessToken, _ := cryptoService.EncryptJwt(ctx, cs.inputSub, cs.inputClaims)

			// when
			claims, err := cryptoService.DecryptJwt(ctx, accessToken)
			if err == nil {
				claims["exp"] = float64(0)
				claims["iat"] = float64(0)
			}

			// then
			assert.Equal(t, cs.expectedClaims, claims)
			assert.Equal(t, cs.expectedError, err)
		})
	}
}

func BenchmarkCryptoServiceDecryptJwt(b *testing.B) {
	// given
	ctx := context.Background()
	cryptoService := service.NewCryptoService("", "key", 900000)
	accessToken, _ := cryptoService.EncryptJwt(ctx, 1, map[string]interface{}{})

	// when
	for i := 0; i < b.N; i++ {
		cryptoService.DecryptJwt(ctx, accessToken)
	}
}
