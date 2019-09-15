package utils

import (
	"crypto/md5"
	"encoding/hex"
)

// generate 32 bit MD5
func GenerateID(text string) string {
	ctx := md5.New()
	ctx.Write([]byte(text))
	return hex.EncodeToString(ctx.Sum(nil))
}
