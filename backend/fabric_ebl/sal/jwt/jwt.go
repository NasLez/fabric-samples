package jwt

import (
	"context"
	"fabric_ebl/biz_error"
	"github.com/bytedance/gopkg/util/logger"
	"github.com/golang-jwt/jwt/v5"
)

var secret = []byte("qSqbqPdB/sLcJxexrr9OnjpzKHsidoHg4vGmQdmhevY=")

func GenerateToken(ctx context.Context, userID int64) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID,
	})

	tokenString, err := token.SignedString(secret)
	if err != nil {
		logger.CtxErrorf(ctx, "GenerateToken failed, err = %v", err)
		return "", err
	}
	return tokenString, err
}

func ValidateToken(ctx context.Context, tokenString string) (map[string]any, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return secret, nil
	})
	if err != nil {
		logger.CtxErrorf(ctx, "ValidateToken failed, err = %v", err)
		return nil, err
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, biz_error.TokenError
}
