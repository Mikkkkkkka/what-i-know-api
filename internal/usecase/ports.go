package usecase

type PasswordHasher interface {
	Hash(password string) (string, error)
	Compare(hashedPassword, password string) error
}

type IDGenerator interface {
	Generate() (string, error)
}
