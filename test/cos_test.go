package test

import (
	"bytes"
	"cloud-disk/core/define"
	"context"
	"net/http"
	"net/url"
	"os"
	"testing"

	"github.com/tencentyun/cos-go-sdk-v5"
)

func TestFileUploadByFilePath(t *testing.T) {
	u, _ := url.Parse("https://jimmy-cloud-disk-1304996341.cos.ap-shanghai.myqcloud.com/")
	b := &cos.BaseURL{BucketURL: u}
	client := cos.NewClient(b, &http.Client{
		Transport: &cos.AuthorizationTransport{
			// 通过环境变量获取密钥
			// 设置了环境变量一定要重新启动！！！
			SecretID:  define.TencentSecretId,
			SecretKey: define.TencentSecretKey,
		},
	})

	key := "test/filepath1.png"

	_, _, err := client.Object.Upload(
		context.Background(), key, "./img/favicon.png", nil,
	)
	if err != nil {
		panic(err)
	}
}

func TestFileUploadByReader(t *testing.T) {
	u, _ := url.Parse("https://jimmy-cloud-disk-1304996341.cos.ap-shanghai.myqcloud.com/")
	b := &cos.BaseURL{BucketURL: u}
	client := cos.NewClient(b, &http.Client{
		Transport: &cos.AuthorizationTransport{
			// 通过环境变量获取密钥
			// 设置了环境变量一定要重新启动！！！
			SecretID:  define.TencentSecretId,
			SecretKey: define.TencentSecretKey,
		},
	})

	key := "test/reader.png"
	f, err := os.ReadFile("./img/favicon.png")
	if err != nil {
		t.Fatal(err)
		return
	}
	_, err = client.Object.Put(
		context.Background(), key, bytes.NewReader(f), nil,
	)
	if err != nil {
		panic(err)
	}
}
