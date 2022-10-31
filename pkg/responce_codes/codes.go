package responce_codes

const (
	// WalletNotFound - кошелек пользователя не найден. Передан незнакомый идентификатор пользователя
	WalletNotFound = "wallet_not_found"

	// NotEnoughAmount - у пользователя недостаточно средств для списания требуемой суммы.
	NotEnoughAmount = "not_enough_cash"

	// InternalError - внутренняя ошибка.
	InternalError = "internal_error"
)
