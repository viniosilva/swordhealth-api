package service

import (
	"context"
	"encoding/base64"
	"fmt"
	"strings"
)

//go:generate mockgen -destination=../../mock/auth_service_mock.go -package=mock . AuthService
type AuthService interface {
	DecodeBasicAuth(ctx context.Context, authorization string) (string, string, error)
}

type authService struct{}

func NewAuthService() AuthService {
	return &authService{}
}

func (impl *authService) DecodeBasicAuth(ctx context.Context, authorization string) (string, string, error) {
	splitedBasicAuth := strings.Split(authorization, " ")
	if strings.ToLower(splitedBasicAuth[0]) != "basic" || len(splitedBasicAuth) != 2 {
		return "", "", fmt.Errorf("invalid authorization")
	}

	basicAuth := splitedBasicAuth[1]
	decodedBasicAuth, _ := base64.StdEncoding.DecodeString(basicAuth)

	auth := strings.Split(string(decodedBasicAuth), ":")
	if len(auth) != 2 {
		return "", "", fmt.Errorf("invalid authorization")
	}

	return auth[0], auth[1], nil
}
