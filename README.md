## Запуск приложения и зависимостей

1. Склонировать репозиторий.

```shell
git clone https://github.com/frutonanny/wallet-service
```

2. Перейти в репозиторий.
3. Выполнить команду запуска приложения и его зависимостей.

```shell
docker-compose -f deployments/docker-compose.yaml --profile dev up --build --detach
```

## Простейший сценарий тестирования приложения

1. Пополняем кошелек пользователя с userID на некоторую сумму.

```
curl -X POST --location "http://localhost:8081/v1/add" \
    -H "Content-Type: application/json" \
    -d "{
          \"userID\": 1,
          \"cash\": 1000
        }"
```

В ответе ожидаем получить баланс пользователя с учетом пополнения.

```json
{
  "data": {
    "balance": 1000
  }
}
```

2. Проверим прямым запросом, что кошелек пользователя создался, и баланс был пополнен.

```
curl -X POST --location "http://localhost:8081/v1/getBalance" \
    -H "Content-Type: application/json" \
    -d "{
          \"userID\": 1
        }"
```

В ответе ожидаем получить ответ, аналогичный ответу из п. 1.

```json
{
  "data": {
    "balance": 1000
  }
}
```

3. Зарезервируем сумму некую у пользователя с userID для оплаты заказа с orderID услуги serviceID.

```
curl -X POST --location "http://localhost:8081/v1/reserve" \
    -H "Content-Type: application/json" \
    -d "{
          \"userID\": 1,
          \"serviceID\": 1,
          \"orderID\": 1,
          \"price\": 500
        }"
```

В ответе ожидаем получить актуальный баланс пользователя без учета зарезервированной суммы. Т.е. 1000 - 500 = 500ю

```json
{
  "data": {
    "balance": 500
  }
}
```

4. Проверим прямым запросом, что сумма была точно зарезервирована и не учитывается в балансе.

```
curl -X POST --location "http://localhost:8081/v1/getBalance" \
    -H "Content-Type: application/json" \
    -d "{
          \"userID\": 1
        }"
```

В ответе ожидаем получить ответ, аналогичный ответу из п. 3.

```json
{
  "data": {
    "balance": 500
  }
}
```

5. Списываем средства у пользователя userID, зарезервированные по заказу orderID.

В качестве суммы списания будем указывать сумму меньшую, чем зарезервировали в п. 3 по заказу orderID. Ожидаем, что
разница между резервом и списанием вернется на баланс пользователя.

```
curl -X POST --location "http://localhost:8081/v1/writeOff" \
    -H "Content-Type: application/json" \
    -d "{
          \"userID\": 1,
          \"serviceID\": 1,
          \"orderID\": 1,
          \"price\": 250
        }"
```

В ответе получим актуальный баланс пользователя. 500 + (500 - 250) = 750.

```json
{
  "data": {
    "balance": 750
  }
}
```

6. Запрашиваем историю транзакций для пользователя userID.

```
curl -X POST --location "http://localhost:8081/v1/getTransactions" \
    -H "Content-Type: application/json" \
    -d "{
          \"userID\": 1,
          \"limit\": 10,
          \"offset\": 0,
          \"sortBy\": \"amount\",
          \"direction\": \"desc\"
        }"
```

Ожидаем в ответе получить три транзакции:

- Зачисление (п. 1).
- Резервирование (п. 3).
- Списание (п. 5).

```json
{
  "data": {
    "transactions": [
      {
        "amount": 1000,
        "createdAt": "2022-11-06T13:05:16.985182Z",
        "description": "Зачисление средств"
      },
      {
        "amount": 500,
        "createdAt": "2022-11-06T13:05:29.546299Z",
        "description": "Резервирование средств по заказу 1"
      },
      {
        "amount": 250,
        "createdAt": "2022-11-06T13:05:49.73709Z",
        "description": "Списание средств по заказу 1"
      }
    ]
  }
}
```
