package service

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"fmt"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v4"
	log "github.com/sirupsen/logrus"
	"github.com/viniosilva/swordhealth-api/internal/exception"
)

//go:generate mockgen -destination=../../mock/crypto_service_mock.go -package=mock . CryptoService
type CryptoService interface {
	Hash(value string) string
	EncryptJwt(ctx context.Context, sub interface{}, claims map[string]interface{}) (string, error)
	DecryptJwt(ctx context.Context, accessToken string) (map[string]interface{}, error)
}

type cryptoService struct {
	hashKey   string
	jwtKey    string
	expiresAt int64
}

func NewCryptoService(hashKey, jwtKey string, expiresAt int64) CryptoService {
	return &cryptoService{
		hashKey:   hashKey,
		jwtKey:    jwtKey,
		expiresAt: expiresAt,
	}
}

func (impl *cryptoService) Hash(value string) string {
	h := hmac.New(sha256.New, []byte(impl.hashKey))
	h.Write([]byte(value))

	return fmt.Sprintf("%x", h.Sum(nil))
}

func (impl *cryptoService) EncryptJwt(ctx context.Context, sub interface{}, claims map[string]interface{}) (string, error) {
	jwtClaims := jwt.MapClaims{
		"iat": jwt.NewNumericDate(time.Now()),
		"exp": jwt.NewNumericDate(time.Now().Add(time.Millisecond * time.Duration(impl.expiresAt))),
		"sub": sub,
	}
	for k, v := range claims {
		jwtClaims[k] = v
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwtClaims)

	accessToken, err := token.SignedString([]byte(impl.jwtKey))
	if err != nil {
		log.WithContext(ctx).WithFields(log.Fields{
			"trace": "internal.service.crypto.encryptjwt",
		}).Error(err.Error())
		return "", err
	}

	return accessToken, nil
}

func (impl *cryptoService) DecryptJwt(ctx context.Context, accessToken string) (map[string]interface{}, error) {
	token, err := jwt.Parse(accessToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(impl.jwtKey), nil
	})
	if err != nil {
		if e, ok := err.(*jwt.ValidationError); ok && e.Error() == "Token is expired" {
			return nil, &exception.ExpiredTokenException{
				Message: strings.ToLower(err.Error()),
			}
		}

		log.WithContext(ctx).WithFields(log.Fields{
			"trace": "internal.service.crypto.decryptjwt",
		}).Error(err.Error())

		return nil, err
	}

	return token.Claims.(jwt.MapClaims), nil
}
