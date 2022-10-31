package add

type Request struct {
	UserID int64 // Идентификатор пользователя.
	Cash   int64 // Сумма в копейках.
}

type Response struct {
	Balance int64 // Текущий баланс пользователя в копейках после пополнения.
}
