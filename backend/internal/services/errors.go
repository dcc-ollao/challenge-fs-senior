package services

import "errors"

var (
	ErrForbidden                = errors.New("forbidden")
	ErrInvalidCredentials       = errors.New("invalid credentials")
	ErrEmailAlreadyExists       = errors.New("email already exists")
	ErrNotFound                 = errors.New("not found")
	ErrEmailAndPasswordRequired = errors.New("email and password are required")
	ErrCannotDeleteOwnUser      = errors.New("cannot delete own user")
	ErrInvalidToken             = errors.New("invalid token")
	ErrJWTSecretNotSet          = errors.New("JWT_SECRET is not set")
	ErrBadRequest               = errors.New("bad request")
)
