package internal_errors

import "errors"

var (
	ErrNotEnoughCash  = errors.New("not enough cash")
	ErrWalletNotFound = errors.New("wallet not found")
)
