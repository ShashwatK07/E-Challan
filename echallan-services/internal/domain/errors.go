package domain

import "errors"

var (
	ErrInvalidTenant          = errors.New("invalid tenant id")
	ErrCitizenNotFound        = errors.New("citizen not found")
	ErrUnauthorized          = errors.New("unauthorized role")
	ErrBusinessServiceMissing = errors.New("business service is missing")
)
