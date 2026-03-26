package security

import (
	"crypto/rand"
	"encoding/base64"

	"what-i-know-api/internal/usecase"
)

const defaultTokenBytes = 32

type RandomTokenGenerator struct {
	size int
}

func NewRandomTokenGenerator(size int) *RandomTokenGenerator {
	if size <= 0 {
		size = defaultTokenBytes
	}

	return &RandomTokenGenerator{size: size}
}

func (g *RandomTokenGenerator) Generate() (string, error) {
	tokenBytes := make([]byte, g.size)
	if _, err := rand.Read(tokenBytes); err != nil {
		return "", err
	}

	return base64.RawURLEncoding.EncodeToString(tokenBytes), nil
}

var _ usecase.TokenGenerator = (*RandomTokenGenerator)(nil)
