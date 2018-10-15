package main

import (
	"log"
	"time"

	"github.com/jmoiron/sqlx"
)

const (
	authQuery       = `SELECT * FROM hc.users WHERE users.login = ? AND users.pass = ?`
	getUsersQuery   = `SELECT * FROM hc.users ORDER BY id`
	getDevicesQuery = `SELECT * FROM hc.devices ORDER BY id`
)

// database структура подключения к базе данных
type database struct {
	conn *sqlx.DB
}

// dbService представляет интерфейс взаимодействия с базой данных
type dbService interface {
	auth(string, string) ([]*User, error)
	getUsers() ([]User, error)
	getDevices() ([]Device, error)
}

func (db database) auth(login string, password string) ([]*User, error) {

	users := make([]*User, 0)
	//err := db.conn.Select(&users, authQuery, login, password)

	stmt, err := db.conn.Prepare(authQuery)
	if err != nil {
		return users, err
	}
	defer stmt.Close()

	rows, err := stmt.Query(login, password)

	for rows.Next() {
		var uid int
		var login string
		var pass string
		var active string
		var userType string
		var lastV time.Time
		err = rows.Scan(&uid, &login, &pass, &active, &userType, &lastV)
		if err != nil {
			return users, err
		}
		u := User{uid, login, active, userType, lastV}
		users = append(users, &u)
	}

	return users, err
}

//func (db database) vacancies() ([]*vacancy, error) {
//	vacancies := make([]*vacancy, 0)
//	err := db.conn.Select(&vacancies, vacanciesQuery)
//	return vacancies, err
//}

// newDB открывает соединение с базой данных и создаёт основную структуру сервиса
func newDB(connectionString string) (database, error) {
	dbConn, err := sqlx.Open("mysql", connectionString)
	log.Println(connectionString)
	return database{conn: dbConn}, err
}

func (db database) getUsers() ([]User, error) {
	users := make([]User, 0)
	stmt, err := db.conn.Prepare(getUsersQuery)
	if err != nil {
		return users, err
	}
	defer stmt.Close()

	rows, err := stmt.Query()

	for rows.Next() {
		var (
			uid      int
			login    string
			pass     string
			active   string
			userType string
			lastV    time.Time
		)
		err = rows.Scan(&uid, &login, &pass, &active, &userType, &lastV)
		if err != nil {
			return users, err
		}
		u := User{uid, login, active, userType, lastV}
		users = append(users, u)
	}

	return users, err
}

func (db database) getDevices() ([]Device, error) {
	devices := make([]Device, 0)
	stmt, err := db.conn.Prepare(getDevicesQuery)
	if err != nil {
		return devices, err
	}
	defer stmt.Close()

	rows, err := stmt.Query()

	for rows.Next() {
		var (
			did    int
			typ    string
			name   string
			ip     string
			active string
			descr  string
		)

		err = rows.Scan(&did, &typ, &name, &ip, &active, &descr)
		if err != nil {
			return devices, err
		}
		d := Device{did, typ, name, ip, active, descr}
		devices = append(devices, d)
	}

	return devices, err
}
