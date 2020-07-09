package common

import (
	"gin_vue_project/model"
	"github.com/dgrijalva/jwt-go"
	"time"
)

var jwtkey = []byte("key")

type Claims struct {
	UserId uint
	jwt.StandardClaims
}

func GetToken(user model.User) (string, error) {
	expiredTime := time.Now().Add(7 * 24 * time.Hour)
	claims := &Claims{
		UserId: user.ID,
		StandardClaims: jwt.StandardClaims{
			// 到期时间
			ExpiresAt: expiredTime.Unix(),
			// 颁发时间
			IssuedAt: time.Now().Unix(),
			// 颁发者
			Issuer: "jizhihaosama",
			// token 名称
			Subject: "user-token",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtkey)

	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func ParseToken(tokenString string) (*jwt.Token, *Claims, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtkey, nil
	})

	return token, claims, err

}
