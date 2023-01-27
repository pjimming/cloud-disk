package test

import (
	"bytes"
	"cloud-disk/core/define"
	"context"
	"fmt"
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

// 分片上传初始化
func TestInitPartUpload(t *testing.T) {
	u, _ := url.Parse("https://jimmy-cloud-disk-1304996341.cos.ap-shanghai.myqcloud.com/")
	b := &cos.BaseURL{BucketURL: u}
	client := cos.NewClient(b, &http.Client{
		Transport: &cos.AuthorizationTransport{
			SecretID:  define.TencentSecretId,
			SecretKey: define.TencentSecretKey,
		},
	})

	name := "test/exampleobject.mp4"
	v, _, err := client.Object.InitiateMultipartUpload(context.Background(), name, nil)
	if err != nil {
		t.Fatal(err)
	}
	UploadID := v.UploadID // 167479887249200a3bef7f87ba4e2e2f7e5179cdbf1b1c66a8271497695900229c4ed9f0b8
	fmt.Println(UploadID)
}

// 上传分片
func TestPartUpload(t *testing.T) {
	u, _ := url.Parse("https://jimmy-cloud-disk-1304996341.cos.ap-shanghai.myqcloud.com/")
	b := &cos.BaseURL{BucketURL: u}
	client := cos.NewClient(b, &http.Client{
		Transport: &cos.AuthorizationTransport{
			SecretID:  define.TencentSecretId,
			SecretKey: define.TencentSecretKey,
		},
	})

	key := "test/exampleobject.mp4"
	// f, err := os.ReadFile("vedio/0.chunk") // md5: bd1a6e47fb1a96b8b0e1ffed99da2147
	f, err := os.ReadFile("vedio/1.chunk") // md5: f56d3293d59e678f1bef4344214a67ab
	if err != nil {
		t.Fatal(err)
	}
	UploadID := "167479887249200a3bef7f87ba4e2e2f7e5179cdbf1b1c66a8271497695900229c4ed9f0b8"
	// opt 可选
	resp, err := client.Object.UploadPart(
		context.Background(), key, UploadID, 2, bytes.NewReader(f), nil,
	)
	if err != nil {
		t.Fatal(err)
	}
	PartETag := resp.Header.Get("ETag")
	fmt.Println(PartETag)
}

// 完成分片上传
func TestFinishPartUpload(t *testing.T) {
	u, _ := url.Parse("https://jimmy-cloud-disk-1304996341.cos.ap-shanghai.myqcloud.com/")
	b := &cos.BaseURL{BucketURL: u}
	client := cos.NewClient(b, &http.Client{
		Transport: &cos.AuthorizationTransport{
			SecretID:  define.TencentSecretId,
			SecretKey: define.TencentSecretKey,
		},
	})

	key := "test/exampleobject.mp4"
	UploadID := "167479887249200a3bef7f87ba4e2e2f7e5179cdbf1b1c66a8271497695900229c4ed9f0b8"

	opt := &cos.CompleteMultipartUploadOptions{}
	opt.Parts = append(opt.Parts,
		cos.Object{PartNumber: 1, ETag: "bd1a6e47fb1a96b8b0e1ffed99da2147"},
		cos.Object{PartNumber: 2, ETag: "f56d3293d59e678f1bef4344214a67ab"},
	)
	_, _, err := client.Object.CompleteMultipartUpload(
		context.Background(), key, UploadID, opt,
	)
	if err != nil {
		t.Fatal(err)
	}
}
