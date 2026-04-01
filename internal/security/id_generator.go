package security

import (
	"github.com/google/uuid"

	"github.com/mikkkkkkka/what-i-know-api/internal/service"
)

type UUIDGenerator struct{}

func NewUUIDGenerator() *UUIDGenerator {
	return &UUIDGenerator{}
}

func (g *UUIDGenerator) Generate() (string, error) {
	return uuid.NewString(), nil
}

var _ service.IDGenerator = (*UUIDGenerator)(nil)
