package claim

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"time"
)

var secret []byte

const TokenExpireDuration = 5 * time.Minute

type AuthClaim struct {
	UserName string `json:"userName"`
	jwt.StandardClaims
}

func GenToken(userName string) (string, error) {
	c := AuthClaim{
		UserName: userName,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(TokenExpireDuration).Unix(),
			Issuer:    "CC",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, c)
	return token.SignedString(secret)
}

func ParseToken(tokenStr string) (*AuthClaim, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &AuthClaim{}, func(tk *jwt.Token) (interface{}, error) {
		return secret, nil
	})
	if err != nil {
		return nil, err
	}
	if claim, ok := token.Claims.(*AuthClaim); ok && token.Valid {
		return claim, nil
	}
	return nil, errors.New("Invalid token ")
}
