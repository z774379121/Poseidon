package common

import (
	"crypto/md5"
	"fmt"
)

func GetMD5(srcData string) string {

	return fmt.Sprintf("%x", md5.Sum([]byte(srcData)))

}
