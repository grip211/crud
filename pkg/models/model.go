package models

import "database/sql"

// тут мы будем описывать структуры для чтения

type Product struct {
	ID       int      `db:"id"`
	Model    string   `db:"model"`
	Company  string   `db:"company"`
	Quantity int      `db:"quantity"`
	Price    float32  `db:"price"`
	Features Features `db:"features"`
}

type Features struct {
	CPU     sql.NullInt32 `db:"cpu"`
	Memory  sql.NullInt32 `db:"memory"`
	Display sql.NullInt32 `db:"display"`
	Camera  sql.NullInt32 `db:"camera"`
}
