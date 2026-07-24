package domainerror

import "errors"

var (
	ErrProfileAlreadyExists = errors.New("profile already exists")
	ErrNoProfileFound       = errors.New("no profile found")
)
