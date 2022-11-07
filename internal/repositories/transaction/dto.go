package transaction

import "time"

const (
	Asc    Direction = "asc"
	Desc   Direction = "desc"
	Amount SortBy    = "amount"
	Date   SortBy    = "created_at"
)

type SortBy string

type Direction string

type Transaction struct {
	Type      string
	Payload   []byte
	Amount    int64
	CreatedAt time.Time
}
