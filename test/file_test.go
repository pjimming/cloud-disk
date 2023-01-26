package test

import (
	"crypto/md5"
	"fmt"
	"io/ioutil"
	"math"
	"os"
	"strconv"
	"testing"
)

// 分片大小
const chunkSize = 10 * 1024 * 1024 // 10MB

// 文件分片
func TestChunkFile(t *testing.T) {
	fileInfo, err := os.Stat("vedio/test.mp4")
	if err != nil {
		t.Fatal(err)
	}
	// 分片个数
	chunkNum := math.Ceil(float64(fileInfo.Size()) / chunkSize)
	myFile, err := os.OpenFile("vedio/test.mp4", os.O_RDONLY, 0666)
	if err != nil {
		t.Fatal(err)
	}
	b := make([]byte, chunkSize)
	for i := 0; i < int(chunkNum); i++ {
		// 指定文件起始位置
		myFile.Seek(int64(i*chunkSize), 0)
		if chunkSize > fileInfo.Size()-int64(i*chunkSize) {
			b = make([]byte, fileInfo.Size()-int64(i*chunkSize))
		}
		myFile.Read(b)

		f, err := os.OpenFile("vedio/"+strconv.Itoa(i)+".chunk", os.O_WRONLY|os.O_CREATE, os.ModePerm)
		if err != nil {
			t.Fatal(err)
		}
		f.Write(b)
		f.Close()
	}
	myFile.Close()
}

// 合并分片
func TestMerge(t *testing.T) {
	myFile, err := os.OpenFile("vedio/merge.mp4", os.O_CREATE|os.O_WRONLY|os.O_APPEND, os.ModePerm)
	if err != nil {
		t.Fatal(err)
	}
	fileInfo, err := os.Stat("vedio/test.mp4")
	if err != nil {
		t.Fatal(err)
	}
	// 分片个数
	chunkNum := math.Ceil(float64(fileInfo.Size()) / chunkSize)
	for i := 0; i < int(chunkNum); i++ {
		f, err := os.OpenFile("vedio/"+strconv.Itoa(i)+".chunk", os.O_RDONLY, os.ModePerm)
		if err != nil {
			t.Fatal(err)
		}
		b, err := ioutil.ReadAll(f)
		if err != nil {
			t.Fatal(err)
		}
		myFile.Write(b)
		f.Close()
	}
	myFile.Close()
}

// 一致性校验(md5)
func TestCheck(t *testing.T) {
	// 获取第一个文件
	f1, err := os.OpenFile("vedio/test.mp4", os.O_RDONLY, 0666)
	if err != nil {
		t.Fatal(err)
	}
	b1, err := ioutil.ReadAll(f1)
	if err != nil {
		t.Fatal(err)
	}
	// 获取第二个文件
	f2, err := os.OpenFile("vedio/merge.mp4", os.O_RDONLY, 0666)
	if err != nil {
		t.Fatal(err)
	}
	b2, err := ioutil.ReadAll(f2)
	if err != nil {
		t.Fatal(err)
	}
	// md5校验
	s1 := fmt.Sprintf("%x", md5.Sum(b1))
	s2 := fmt.Sprintf("%x", md5.Sum(b2))
	fmt.Println(s1)
	fmt.Println(s2)
	fmt.Println(s1 == s2)

	f1.Close()
	f2.Close()
}
