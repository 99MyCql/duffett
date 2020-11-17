package pkg

import (
	"crypto/md5"
	"fmt"
)

// Md5Encode Md5 字符串加密
func Md5Encode(data string) string {
	return fmt.Sprintf(
		"%x",
		md5.Sum([]byte(data)),
	)
}
