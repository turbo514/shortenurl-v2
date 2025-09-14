package util

import (
	"crypto/ecdsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"time"
)

import "github.com/golang-jwt/jwt/v5"

func GenerateToken(identify string, key []byte) (string, error) {
	claims := jwt.RegisteredClaims{
		Issuer:    "api-gateway service",
		Subject:   identify,
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 999)), // TODO: 修改时间
	}

	token := jwt.NewWithClaims(jwt.SigningMethodES256, claims)

	privateKey, err := loadECDSAPrivateKey(key)
	if err != nil {
		return "", err
	}
	return token.SignedString(privateKey)

}

func loadECDSAPrivateKey(pemData []byte) (*ecdsa.PrivateKey, error) {
	// 解析 PEM 数据
	block := &pem.Block{}
	block, pemData = pem.Decode(pemData)
	if block == nil || block.Type != "EC PRIVATE KEY" {
		return nil, fmt.Errorf("failed to decode PEM block containing ECDSA private key")
	}

	// 解析 ECDSA 私钥
	return x509.ParseECPrivateKey(block.Bytes)
}
