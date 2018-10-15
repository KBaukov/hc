package db

import (
	"log"

	"github.com/jmoiron/sqlx"
)

const (
	vacanciesQuery = `SELECT id, name, salary, experience, city 
					  FROM vacancy
					  ORDER BY name`
)

// database структура подключения к базе данных
type database struct {
	conn *sqlx.DB
}

// dbService представляет интерфейс взаимодействия с базой данных
type dbService interface {
	auth(string, string) bool
}

func (db database) auth(login string, password string) bool {
	log.Println("------------------")

	return true
}

// newDB открывает соединение с базой данных и создаёт основную структуру сервиса
func newDB(connectionString string) (database, error) {
	dbConn, err := sqlx.Open("mysql", connectionString)
	log.Println(connectionString)
	return database{conn: dbConn}, err
}
