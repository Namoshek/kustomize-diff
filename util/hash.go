package util

import (
	"crypto/md5"
	"encoding/hex"
)

func CalculateMD5AsString(text string) string {
	hash := md5.Sum([]byte(text))

	return hex.EncodeToString(hash[:])
}
