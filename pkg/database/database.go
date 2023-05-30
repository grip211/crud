package database

import builder "github.com/doug-martin/goqu/v9"

// тут мы пишем интефейс, для получения доступа к пулу коннектов базы данных
// его мы и будем использовать, а не частный случай MySQL

type Pool interface {
	Builder() *builder.Database
}
