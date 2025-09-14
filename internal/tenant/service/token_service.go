package service

import (
	"context"
	"crypto/ecdsa"
	"github.com/google/uuid"
	"github.com/turbo514/shortenurl-v2/shared/util"
)

type ITokenService interface {
	GenerateToken(ctx context.Context, userID, tenantID uuid.UUID) (string, error)
}

type TokenService struct {
	privateKey *ecdsa.PrivateKey
	publicKey  *ecdsa.PublicKey
}

func NewTokenService(privateKeyPath, publicKeyPath string) *TokenService {
	//block, _ := pem.Decode([]byte(privateKeyPath))
	//x509.ParseECPrivateKey(block.Bytes)
	return &TokenService{}
}

func (t TokenService) GenerateToken(ctx context.Context, userID, tenantID uuid.UUID) (string, error) {
	token, err := util.GenerateToken(map[string]any{
		"user_id":   userID,
		"tenant_id": tenantID,
	}, []byte{}, "Tenant Service")
	if err != nil {
		return "", err
	}
	return token, nil
}
