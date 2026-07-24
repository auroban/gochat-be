package security

type PasswordHasher interface {
	Hash(rawPassword string) (string, error)
}
