package models

type Product struct {
	Id      int     `db:"id"`
	Model   string  `db:"model"`
	Company string  `db:"company"`
	Price   float32 `db:"price"`
}
