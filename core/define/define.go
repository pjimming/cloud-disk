package define

import (
	"os"

	"github.com/golang-jwt/jwt/v4"
)

type UserClaim struct {
	Id       int
	Identity string
	Name     string
	jwt.StandardClaims
}

// JwtKey 签名
var JwtKey = "cloud-disk-key"

// 邮件发送密码
var MailPassword = os.Getenv("MailPassWord_For_panjm2001@126.com")

// 验证码长度
var CodeLength = 6

// 验证码过期时间(s)
var CodeExpire = 300

// 腾讯云密钥ID
var TencentSecretId = os.Getenv("TencentSecretID")

// 腾讯云密钥KEY
var TencentSecretKey = os.Getenv("TencentSecretKey")

// 腾讯云存储桶路径
var CosBucket = "https://jimmy-cloud-disk-1304996341.cos.ap-shanghai.myqcloud.com/"

// 分页默认参数
var PageSize = 20
