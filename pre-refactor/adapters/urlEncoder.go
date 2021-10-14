package adapters

import (
	"math/rand"
	"strings"
	"time"
)

const (
	alphabet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	length   = uint64(len(alphabet))
)

// EncodeLink creates a hash based on the number passed in.
func EncodeLink(number uint64) string {
	var encodedBuilder strings.Builder
	encodedBuilder.Grow(11)

	for ; number > 0; number = number / length {
		encodedBuilder.WriteByte(alphabet[(number % length)])
	}

	return encodedBuilder.String()
}

// CreateURL creates a hash which is used for mapping to a URL.
func CreateURL() string {
	hash := EncodeLink(uint64(time.Now().UTC().Unix() * rand.Int63()))
	return hash
}
