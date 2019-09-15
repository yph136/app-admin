package utils

import (
	"crypto/rand"
	"encoding/hex"
)

// MustUUID 创建UUID
func MustUUID() string {
	uuid, err := NewUUID()
	if err != nil {
		return ""
	}
	return uuid
}

// NewUUID 创建UUID,参考: http://github.com/google/uuid
func NewUUID() (string, error) {
	var buf [16]byte
	_, err := rand.Read(buf[:])
	if err != nil {
		return "", err
	}

	buf[6] = (buf[6] & 0x0f) | 0x40
	buf[8] = (buf[8] & 0x3f) | 0x80

	dst := make([]byte, 36)
	hex.Encode(dst, buf[:4])
	dst[8] = '-'
	hex.Encode(dst[9:13], buf[4:6])
	dst[13] = '-'
	hex.Encode(dst[14:18], buf[6:8])
	dst[18] = '-'
	hex.Encode(dst[19:23], buf[8:10])
	dst[23] = '-'
	hex.Encode(dst[24:], buf[10:])

	return string(dst), nil
}
