package errs

import "errors"

var (
	ErrInvalidInput     = errors.New("invalid input")
	ErrUnauthorized     = errors.New("unauthorized")
	ErrForbidden        = errors.New("forbidden")
	ErrNotFound         = errors.New("not found")
	ErrAlreadyExists    = errors.New("already exists")
	ErrAlreadyReserved  = errors.New("item already reserved")
	ErrWishlistMismatch = errors.New("item does not belong to wishlist")
)
