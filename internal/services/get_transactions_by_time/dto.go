package get_transactions_by_time

import (
	"fmt"
	"time"

	"github.com/frutonanny/wallet-service/internal/repositories/transaction"
	"github.com/frutonanny/wallet-service/internal/transactions"
)

type Transaction struct {
	Description string
	Amount      int64
	CreatedAt   time.Time
}

// adaptTxs преобразует список транзакций, полученный из базы, в список транзакций, который отдает метод.
func adaptTxs(txs []transaction.Transaction) ([]Transaction, error) {
	result := make([]Transaction, 0, len(txs))

	for i := range txs {
		tx, err := adaptTx(txs[i])
		if err != nil {
			return nil, fmt.Errorf("adapt tx: %v", err)
		}

		result = append(result, tx)
	}

	return result, nil
}

func adaptTx(tx transaction.Transaction) (Transaction, error) {
	desc, err := getTxDescription(tx.Type, tx.Payload)
	if err != nil {
		return Transaction{}, fmt.Errorf("get tx description: %v", err)
	}

	return Transaction{
		Description: desc,
		Amount:      tx.Amount,
		CreatedAt:   tx.CreatedAt,
	}, nil
}

// getTxDescription() - создает описание транзакции в зависимости от полученного типа.
func getTxDescription(txType string, payload []byte) (string, error) {
	if txType == transactions.TypeAdd {
		return "Зачисление средств", nil
	}

	orderID, err := transactions.GetOrderID(payload)
	if err != nil {
		return "", fmt.Errorf("get order id: %v", err)
	}

	switch txType {
	case transactions.TypeReserve:
		return fmt.Sprintf("Резервирование средств по заказу %d", orderID), nil
	case transactions.TypeWriteOff:
		return fmt.Sprintf("Списание средств по заказу %d", orderID), nil
	case transactions.TypeCancel:
		return fmt.Sprintf("Отмена резервирования средств по заказу %d", orderID), nil
	default:
		// Сознательно не возвращаем ошибку, чтобы не блокировать показ всех остальных транзакций,
		// если такое случится.
		return "Неизвестный тип транзакции", nil
	}
}
