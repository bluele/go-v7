package v7

import (
	"errors"
)

var (
	errNegativeInt   = errors.New("go-v7: unexpected value for Uint64")
	errUndefinedType = errors.New("go-v7: undefined type")
)
