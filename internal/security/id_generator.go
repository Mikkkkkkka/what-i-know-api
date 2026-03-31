package security

import (
	"github.com/google/uuid"

	"github.com/mikkkkkkka/what-i-know-api/internal/usecase"
)

type UUIDGenerator struct{}

func NewUUIDGenerator() *UUIDGenerator {
	return &UUIDGenerator{}
}

func (g *UUIDGenerator) Generate() (string, error) {
	return uuid.NewString(), nil
}

var _ usecase.IDGenerator = (*UUIDGenerator)(nil)
