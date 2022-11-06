package transactions

import "time"

const (
	TypeAdd      = "incoming_transfer"
	TypeReserve  = "reservation"
	TypeWriteOff = "write_off"
	TypeCancel   = "cancel"
)

type Transaction struct {
	Description string
	Amount      int64
	CreatedAt   time.Time
}
