package models

import "database/sql"

// тут мы будем описывать структуры для чтения

type Product struct {
	ID       int      `db:"id" json:"id"`
	Model    string   `db:"model" json:"model"`
	Company  string   `db:"company" json:"company"`
	Quantity int      `db:"quantity" json:"quantity"`
	Price    float32  `db:"price" json:"price"`
	Features Features `db:"features" json:"features"`
}

type Features struct {
	CPU     sql.NullInt32 `db:"cpu" json:"cpu"`
	Memory  sql.NullInt32 `db:"memory" json:"memory"`
	Display sql.NullInt32 `db:"display" json:"display"`
	Camera  sql.NullInt32 `db:"camera" json:"camera"`
}
