package randutils

import (
	"crypto/rand"
	"encoding/base64"
)

func GenerateId() (string, bool) {
	// use crypto/rand to generate a random id

	data := make([]byte, 16)
	if size, err := rand.Read(data); err != nil || size != len(data) {
		return "", false
	}

	return base64.StdEncoding.EncodeToString(data), true
}
