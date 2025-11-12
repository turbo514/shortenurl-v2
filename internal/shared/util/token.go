package util

import (
	"crypto/ecdsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

// FIXME: 以后再来改造,目标是类型安全和拓展性
// FIXME: 改用ES256

var hs256key = []byte("123456")

type GeneralClaims struct {
	//Type string
	M map[string]interface{}
	jwt.RegisteredClaims
}

//type UserClaims struct {
//	UserID   uuid.UUID `json:"user_id"`
//	TenantID uuid.UUID `json:"tenant_id"`
//	jwt.RegisteredClaims
//}
//
//func NewUserClaims(userID uuid.UUID, tenantID uuid.UUID) *UserClaims {
//	return &UserClaims{
//		UserID:   userID,
//		TenantID: tenantID,
//	}
//}

func GenerateToken(data map[string]any, key []byte, issuer string) (string, error) {
	claims := GeneralClaims{
		M: data,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    issuer,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24 * 7)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	//privateKey, err := loadECDSAPrivateKey(key)
	//if err != nil {
	//	return "", err
	//}

	return token.SignedString(hs256key)
}

func ParseToken(tokenString string, key []byte) (*GeneralClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &GeneralClaims{}, func(t *jwt.Token) (any, error) {
		return hs256key, nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("非法Token")
	}

	claims, ok := token.Claims.(*GeneralClaims)
	if !ok {
		return nil, errors.New("Token解析失败")
	}

	return claims, nil
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
