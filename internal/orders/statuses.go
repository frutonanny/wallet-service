package orders

const (
	StatusReserved   = "reserved"    // Деньги зарезервированы по заказу.
	StatusWrittenOff = "written_off" // Деньги по заказу списаны.
	StatusCancelled  = "cancelled"   // Заказ отменен.
)

func IsOrderReserved(status string) bool {
	return status == StatusReserved
}
