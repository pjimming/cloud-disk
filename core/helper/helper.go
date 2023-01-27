package helper

import (
	"bytes"
	"cloud-disk/core/define"
	"context"
	"crypto/md5"
	"crypto/tls"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"net/smtp"
	"net/url"
	"path"
	"strconv"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/jordan-wright/email"
	uuid "github.com/satori/go.uuid"
	"github.com/tencentyun/cos-go-sdk-v5"
)

// Md5
// md5加密
func Md5(s string) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(s)))
}

// GenerateToken
// 生成Token
func GenerateToken(id int, identity, name string, seconds int) (string, error) {
	// id, identity, name
	uc := define.UserClaim{
		Id:       id,
		Identity: identity,
		Name:     name,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Second * time.Duration(seconds)).Unix(),
		},
	}
	// 利用jwt生成一个token
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
	e.To = []string{mail}
	e.Subject = "Jimmy Cloud-Disk验证码发送测试"
	e.HTML = []byte("您的验证码为：<h1>" + code + "</h1>")
	err := e.SendWithTLS("smtp.126.com:465", smtp.PlainAuth("", "panjm2001@126.com", define.MailPassword, "smtp.126.com"),
		&tls.Config{InsecureSkipVerify: true, ServerName: "smtp.126.com"})
	if err != nil {
		return err
	}
	return nil
}

// RandCode
// 生成随机验证码
func RandCode() string {
	s := "1234567890"
	code := ""
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < define.CodeLength; i++ {
		code += string(s[rand.Intn(len(s))])
	}
	return code
}

// UUID
// 生成UUID
func UUID() string {
	return uuid.NewV4().String()
}

// CosUpload
// 上传文件到腾讯云，返回文件路径和error
func CosUpload(r *http.Request) (string, error) {
	u, _ := url.Parse(define.CosBucket)
	b := &cos.BaseURL{BucketURL: u}
	client := cos.NewClient(b, &http.Client{
		Transport: &cos.AuthorizationTransport{
			SecretID:  define.TencentSecretId,
			SecretKey: define.TencentSecretKey,
		},
	})

	file, fileHeader, err := r.FormFile("file")
	if err != nil {
		panic(err)
	}

	key := "test/" + UUID() + path.Ext(fileHeader.Filename) // 文件路径

	_, err = client.Object.Put(
		context.Background(), key, file, nil,
	)
	if err != nil {
		panic(err)
	}
	return define.CosBucket + "/" + key, err
}

// AnalyzeToken
// 解析 Token
func AnalyzeToken(token string) (*define.UserClaim, error) {
	uc := new(define.UserClaim)
	claims, err := jwt.ParseWithClaims(token, uc, func(token *jwt.Token) (interface{}, error) {
		return []byte(define.JwtKey), nil
	})
	if err != nil {
		return nil, err
	}
	// claims里面没有值
	if !claims.Valid {
		return uc, errors.New("token is invalid")
	}
	return uc, err
}

// InitPartUpload
// 分片上传初始化
func InitPartUpload(ext string) (string, string, error) {
	u, _ := url.Parse(define.CosBucket)
	b := &cos.BaseURL{BucketURL: u}
	client := cos.NewClient(b, &http.Client{
		Transport: &cos.AuthorizationTransport{
			SecretID:  define.TencentSecretId,
			SecretKey: define.TencentSecretKey,
		},
	})

	name := "test/" + UUID() + ext
	v, _, err := client.Object.InitiateMultipartUpload(context.Background(), name, nil)
	if err != nil {
		return "", "", err
	}
	return name, v.UploadID, nil
}

// CosPartUpload
// 上传分片
func CosPartUpload(r *http.Request) (string, error) {
	u, _ := url.Parse(define.CosBucket)
	b := &cos.BaseURL{BucketURL: u}
	client := cos.NewClient(b, &http.Client{
		Transport: &cos.AuthorizationTransport{
			SecretID:  define.TencentSecretId,
			SecretKey: define.TencentSecretKey,
		},
	})

	key := r.PostForm.Get("key")
	uploadId := r.PostForm.Get("upload_id")
	partNumber, err := strconv.Atoi(r.PostForm.Get("part_number"))
	if err != nil {
		return "", err
	}
	f, _, err := r.FormFile("file")
	if err != nil {
		return "", err
	}
	buf := bytes.NewBuffer(nil)
	io.Copy(buf, f)

	// opt 可选
	resp, err := client.Object.UploadPart(
		context.Background(), key, uploadId, partNumber, bytes.NewReader(buf.Bytes()), nil,
	)
	if err != nil {
		return "", err
	}
	return strings.Trim(resp.Header.Get("ETag"), "\""), nil
}

// 完成分片上传
func FinishPartUpload(key, uploadId string, co []cos.Object) error {
	u, _ := url.Parse(define.CosBucket)
	b := &cos.BaseURL{BucketURL: u}
	client := cos.NewClient(b, &http.Client{
		Transport: &cos.AuthorizationTransport{
			SecretID:  define.TencentSecretId,
			SecretKey: define.TencentSecretKey,
		},
	})

	opt := &cos.CompleteMultipartUploadOptions{}
	opt.Parts = append(opt.Parts, co...)
	_, _, err := client.Object.CompleteMultipartUpload(
		context.Background(), key, uploadId, opt,
	)
	return err
}
