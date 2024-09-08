package biz

import (
	"crypto/md5"
	"encoding/hex"
	"io"
)

func GetMd5(str string) string {
	hash := md5.New()
	_, _ = io.WriteString(hash, str)
	return hex.EncodeToString(hash.Sum(nil))
}
