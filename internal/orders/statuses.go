package orders

const (
	StatusReserved   = "reserved"    // Деньги зарезервированы для оплаты заказа.
	StatusWrittenOff = "written_off" // Заказ оплачен.
	StatusCancelled  = "cancelled"   // Заказ отменен.
)

func IsOrderReserved(status string) bool {
	return status == StatusReserved
}
