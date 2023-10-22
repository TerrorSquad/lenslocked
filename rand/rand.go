package rand

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
)

func Bytes(n int) ([]byte, error) {
	b := make([]byte, n)
	nRead, err := rand.Read(b)
	if err != nil {
		return nil, fmt.Errorf("bytes: %w", err)
	}
	if nRead < n {
		return nil, fmt.Errorf("bytes: didn't read enough bytes: %w", err)
	}
	return b, nil
}

// String returns a random string using the crypto/rand package
// n is the number of bytes being used to generate the random string
func String(n int) (string, error) {
	const characters = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz-"
	bytes, err := Bytes(n)
	if err != nil {
		return "", fmt.Errorf("string: %w", err)
	}
	return base64.URLEncoding.EncodeToString(bytes), nil
}

const SessionTokenBytes = 32

func SessionToken() (string, error) {
	return String(SessionTokenBytes)
}
