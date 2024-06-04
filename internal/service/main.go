package service

import (
	"crypto/md5"
	"encoding/hex"
	"os"
)

func HashPassword(password string, salt string) string {
	key := os.Getenv("KEY")

	// Hash password with format key + password + salt
	h := md5.New()
	h.Write([]byte(key + password + salt))
	bytesPassword := h.Sum(nil)
	hashedPassword := hex.EncodeToString(bytesPassword)
	return hashedPassword
}