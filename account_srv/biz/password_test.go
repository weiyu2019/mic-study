package biz

import (
	"github.com/anaskhan96/go-password-encoder"
	"testing"
)

func TestGetMd5(t *testing.T) {
	salt, encodedPwd := password.Encode("happy", nil)
	println(salt)
	println(encodedPwd)
	res := password.Verify("happy", salt, encodedPwd, nil)
	println(res)
}
