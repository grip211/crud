package models

// тут мы будем описывать структуры для чтения

type Product struct {
	Id       int     `db:"id"`
	Model    string  `db:"model"`
	Company  string  `db:"company"`
	Quantity int     `db:"quantity"`
	Price    float32 `db:"price"`
}
