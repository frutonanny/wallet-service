package repositories

import "errors"

// Ошибки, о которых необходимо сообщить другому сервису / пользователю.
var (
	ErrRepoNotEnoughAmount = errors.New("not enough cash")
	ErrRepoWalletNotFound  = errors.New("wallet not found")
)
