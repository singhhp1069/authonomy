package utils

import (
	"identitysphere-api/models"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/spf13/viper"
)

var secret = viper.GetString("service.jwt_encryption_key")
var encryptionKey = []byte(secret)
var validity = 24 * time.Hour

type CustomClaims struct {
	AppDID         string                              `json:"app_id"`
	CredentialJWTs models.IssueOAuthCredentialResponse `json:"credential_jwts"`
	jwt.StandardClaims
}

func CreateAccessToken(appDID string, credentialJWTs models.IssueOAuthCredentialResponse) (string, error) {
	expirationTime := time.Now().Add(validity)
	claims := CustomClaims{
		AppDID:         appDID,
		CredentialJWTs: credentialJWTs,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(encryptionKey)

	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func ValidateAccessToken(tokenString string) (*CustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return encryptionKey, nil
	})

	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		return claims, nil
	} else {
		return nil, err
	}
}
