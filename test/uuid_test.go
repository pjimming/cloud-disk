package test

import (
	"testing"

	uuid "github.com/satori/go.uuid"
)

func TestGenerateUUID(t *testing.T) {
	v4 := uuid.NewV4().String()
	println(v4)
}
