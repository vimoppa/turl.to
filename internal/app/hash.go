package app

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
)

// GenerateHash generates a 10 length hash string from  key.
func GenerateHash(key string) string {
	fmt.Println(key)
	h := md5.New()
	h.Write([]byte(key))
	return hex.EncodeToString(h.Sum(nil))[:10]
}
