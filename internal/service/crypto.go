package service

import (
	"crypto/hmac"
	"crypto/sha256"
	"fmt"
)

//go:generate mockgen -destination=../../mock/crypto_service_mock.go -package=mock . CryptoService
type CryptoService interface {
	Hash(value string) string
}

type cryptoService struct {
	key string
}

func NewCryptoService(key string) CryptoService {
	return &cryptoService{
		key: key,
	}
}

func (impl *cryptoService) Hash(value string) string {
	h := hmac.New(sha256.New, []byte(impl.key))
	h.Write([]byte(value))

	return fmt.Sprintf("%x", h.Sum(nil))
}
