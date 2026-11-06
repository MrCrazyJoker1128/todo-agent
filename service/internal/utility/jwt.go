package utility

import (
	"context"
	"errors"
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	UserId string `json:"userId"`
	jwt.RegisteredClaims
}

func jwtKey(ctx context.Context) string {
	if key := g.Cfg().MustGet(ctx, "jwt.key", "todo-agent").String(); key != "" {
		return key
	}
	return "todo-agent"
}

func jwtExpireSeconds(ctx context.Context) int64 {
	if expire := g.Cfg().MustGet(ctx, "jwt.expire", 86400).Int64(); expire > 0 {
		return expire
	}
	return 86400
}

func GenerateToken(ctx context.Context, userId string) (string, error) {
	key := jwtKey(ctx)
	expire := jwtExpireSeconds(ctx)

	claims := Claims{
		UserId: userId,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(expire) * time.Second)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(key))
}

func ParseToken(ctx context.Context, tokenString string) (*Claims, error) {
	key := jwtKey(ctx)

	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (any, error) {
		return []byte(key), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}
