package utils

import (
	"crypto/rand"
	"io"
	"time"
)

var table = [...]byte{'1', '2', '3', '4', '5', '6', '7', '8', '9', '0'}

//EncodeToString generates random code
func EncodeToString(max int) string {
	b := make([]byte, max)
	n, err := io.ReadAtLeast(rand.Reader, b, max)
	if n != max {
		panic(err)
	}
	for i := 0; i < len(b); i++ {
		b[i] = table[int(b[i])%len(table)]
	}
	return string(b)
}

// EncodingUTF8GetString converts byte array to string
func EncodingUTF8GetString(byteArray []byte) string {
	return string(byteArray)
}

// EncodingUTF8GetBytes converts string to byte array
func EncodingUTF8GetBytes(str string) []byte {
	return []byte(str)
}

// RefBool return pointer to boolean value
func RefBool(value bool) *bool {
	b := value
	return &b
}

// RefTime return pointer to boolean value
func RefTime(value time.Time) *time.Time {
	b := value
	return &b
}
