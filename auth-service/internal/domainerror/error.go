package domainerror

import "errors"

var (
	ErrUserAlreadyExists = errors.New("user already exists")
	ErrNoUserFound       = errors.New("no user found")
)
