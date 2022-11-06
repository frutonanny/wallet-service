package pkg

// Виды услуг: map[идентификатор услуги]наименование услуги
var services = map[int64]string{
	1: "More views",
	2: "Highlighting",
	3: "XL-advert",
}

// ExistServiceID - проверяет, что получили существующий идентификатор услуги
func ExistServiceID(ID int64) bool {
	_, ok := services[ID]
	return ok
}
