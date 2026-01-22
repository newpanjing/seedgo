package shared

import (
	"errors"
	"seedgo/internal/global"
	"seedgo/internal/model"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type MyCustomClaims struct {
	UserID   model.ID `json:"userId"`
	Username string   `json:"username"`
	TenantID model.ID `json:"tenantId"`
	Super    bool     `json:"super"`
	jwt.RegisteredClaims
}

func GenerateToken(userID model.ID, username string, tenantID model.ID, super bool) (string, error) {
	// 默认2小时过期
	expire := global.Config.JWT.Expire
	if expire == 0 {
		expire = 7200
	}

	claims := MyCustomClaims{
		UserID:   userID,
		Username: username,
		TenantID: tenantID,
		Super:    super,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(expire) * time.Second)),
			Issuer:    "smart_butler",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(global.Config.JWT.Secret))
}

func ParseToken(tokenString string) (*MyCustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &MyCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(global.Config.JWT.Secret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*MyCustomClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}
