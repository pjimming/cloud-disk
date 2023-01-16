package helper

import (
	"cloud-disk/core/define"
	"crypto/md5"
	"fmt"

	"github.com/golang-jwt/jwt/v4"
)

func Md5(s string) string { // md5加密
	return fmt.Sprintf("%x", md5.Sum([]byte(s)))
}

func GenerateToken(id int, identity, name string) (string, error) { // 生成Token
	// id, identity, name
	uc := define.UserClaim{
		Id:       id,
		Identity: identity,
		Name:     name,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, uc)
	tokenString, err := token.SignedString([]byte(define.JwtKey))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}
