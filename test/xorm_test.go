package test

import (
	"bytes"
	"cloud-disk/core/models"
	"encoding/json"
	"fmt"
	"testing"

	_ "github.com/go-sql-driver/mysql"

	"xorm.io/xorm"
)

func TestXormTest(t *testing.T) {
	// 使用Engine引擎，建立mysql连接。
	engine, err := xorm.NewEngine("mysql", "root:123456@tcp(127.0.0.1:3306)/cloud_disk?charset=utf8mb4&parseTime=True&loc=Local")
	if err != nil {
		t.Fatal(err)
	}
	// 建立UserBasic的数组切片，初始个数为0。
	data := make([]*models.UserBasic, 0)
	/*
		xorm.engine.Find(*): 查询多条数据使用Find方法，
		Find方法的第一个参数为slice的指针或Map指针，即为查询后返回的结果，
		第二个参数可选，为查询的条件struct的指针。
	*/
	err = engine.Find(&data)
	if err != nil {
		t.Fatal(err)
	}
	/*
		json.Marshal(): 将数据编码成json字符串，
		结构体必须是大写字母开头的成员才会被JSON处理到，小写字母开头的成员不会有影响。
	*/
	b, err := json.Marshal(data)
	if err != nil {
		t.Fatal(err)
	}
	dst := new(bytes.Buffer)
	/*
		json.Indent(要复制到的dst，src，前缀prefix，indent string) error
	*/
	err = json.Indent(dst, b, "", " ")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(dst.String())
}
