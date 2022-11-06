package transactions

import (
	"encoding/json"
	"fmt"
)

const (
	typeEnrollment = "enrollment"
)

type payload struct {
	OrderID int64 `json:"order_id"`
}

type addPayload struct {
	Type string `json:"type"`
}

func EnrollmentPayload() (json.RawMessage, error) {
	d := addPayload{
		Type: typeEnrollment,
	}

	b, err := json.Marshal(d)
	if err != nil {
		return nil, fmt.Errorf("marshal payload: %v", err)
	}

	return b, nil
}

func ReservationPayload(orderID int64) (json.RawMessage, error) {
	return commonPayload(orderID)
}

func WriteOffPayload(orderID int64) (json.RawMessage, error) {
	return commonPayload(orderID)
}

func CancelPayload(orderID int64) (json.RawMessage, error) {
	return commonPayload(orderID)
}

func commonPayload(orderID int64) (json.RawMessage, error) {
	d := payload{
		OrderID: orderID,
	}

	b, err := json.Marshal(d)
	if err != nil {
		return nil, fmt.Errorf("marshal payload: %v", err)
	}

	return b, nil
}

// GetOrderID вытаскивает номер заказа из переданного
func GetOrderID(raw json.RawMessage) (int64, error) {
	p := payload{}

	if err := json.Unmarshal(raw, &p); err != nil {
		return 0, fmt.Errorf("unmarshal payload: %v", err)
	}

	return p.OrderID, nil
}
