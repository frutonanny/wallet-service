package get_balance

type Request struct {
	UserID int64 // Идентификатор пользователя.
}

type Response struct {
	Balance int64 // Текущий баланс пользователя в копейках.
}
