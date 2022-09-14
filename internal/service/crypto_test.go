package service_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
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
			cryptoService := service.NewCryptoService(cs.injectKey)

			// when
			hash := cryptoService.Hash(cs.inputValue)

			// then
			assert.Equal(t, cs.expectedHash, hash)
		})
	}
}

func BenchmarkCryptoServiceHash(b *testing.B) {
	// given
	cryptoService := service.NewCryptoService("key")

	// when
	for i := 0; i < b.N; i++ {
		cryptoService.Hash("value")
	}
}
