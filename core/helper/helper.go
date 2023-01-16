package helper

import (
	"cloud-disk/core/define"
	"crypto/md5"
	"crypto/tls"
	"fmt"
	"net/smtp"

	"github.com/golang-jwt/jwt/v4"
	"github.com/jordan-wright/email"
)

// Md5
// md5加密
func Md5(s string) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(s)))
}

// GenerateToken
// 生成Token
func GenerateToken(id int, identity, name string) (string, error) {
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

// MailCodeSend
// 邮箱验证码发送
func MailCodeSend(mail, code string) error {
	e := email.NewEmail()
	e.From = "Jimmy Cloud-Disk <panjm2001@126.com>"
	e.To = []string{"2556323541@qq.com"}
	e.Subject = "Jimmy Cloud-Disk验证码发送测试"
	e.HTML = []byte("您的验证码为：<h1>" + code + "</h1>")
	err := e.SendWithTLS("smtp.126.com:465", smtp.PlainAuth("", "panjm2001@126.com", define.MailPassword, "smtp.126.com"),
		&tls.Config{InsecureSkipVerify: true, ServerName: "smtp.126.com"})
	if err != nil {
		return err
	}
	return nil
}
