package test

import (
	"cloud-disk/core/define"
	"crypto/tls"
	"net/smtp"
	"testing"

	"github.com/jordan-wright/email"
)

func TestSendMail(t *testing.T) {
	e := email.NewEmail()
	e.From = "Jimmy Cloud-Disk <panjm2001@126.com>"
	e.To = []string{"2556323541@qq.com"}
	e.Subject = "Jimmy Cloud-Disk验证码发送测试"
	e.HTML = []byte("您的验证码为：<h1>123456</h1>")
	err := e.SendWithTLS("smtp.126.com:465", smtp.PlainAuth("", "panjm2001@126.com", define.MailPassword, "smtp.126.com"),
		&tls.Config{InsecureSkipVerify: true, ServerName: "smtp.126.com"})
	if err != nil {
		t.Fatal(err)
	}
}
