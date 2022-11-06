package repositories

import "errors"

// Ошибки, о которых необходимо сообщить другому сервису / пользователю.
var (
	ErrRepoNotEnoughCash         = errors.New("not enough cash")
	ErrRepoWalletNotFound        = errors.New("wallet not found")
	ErrRepoOrderNotFound         = errors.New("order not found")
	ErrRepoNotEnoughReservedCash = errors.New("not enough reserved cash")
)
