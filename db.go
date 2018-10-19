package main

import (
	"errors"
	"log"
	"time"

	"github.com/jmoiron/sqlx"
)

const (
	authQuery = `SELECT * FROM hc.users WHERE users.login = ? AND users.pass = ?`

	getUsersQuery = `SELECT * FROM hc.users ORDER BY id`
	addUserQuery  = `INSERT INTO hc.users (login, pass, user_type, active_flag, last_visit, id) VALUES (?,?,?,?,?,?)`
	updUserQuery  = `UPDATE hc.users SET login=?, pass=?, user_type=?, active_flag=?, last_visit=? WHERE id=?`
	delUserQuery  = `DELETE FROM hc.users WHERE id=?`
	//lastUsrIdQuery = `SELECT max(id) FROM hc.users`

	getDevicesQuery = `SELECT * FROM hc.devices ORDER BY id`
	addDeviceQuery  = `INSERT INTO hc.devices (type, name, ip, active_flag, description, id) VALUES (?,?,?,?,?,?)`
	updDeviceQuery  = `UPDATE hc.devices SET type = ?, name = ?, ip = ?, active_flag = ?, description = ? WHERE id=?`
	delDeviceQuery  = `DELETE FROM hc.devices WHERE id=?`

	kotelDevIdQuery    = `SELECT name FROM hc.devices WHERE type="KotelController"`
	getKotelDataQuery  = `SELECT DEVICE_ID , TP, 'TO', PR, KW, DEST_TP, DEST_TO, DEST_PR, DEST_KW, DEST_TC FROM hc.kotel`
	updKotelDevIdQuery = `UPDATE hc.kotel SET device_id = ? WHERE 1`
	updKotelDataQuery  = `UPDATE hc.kotel SET tp= ?, to= ?, pr= ?, kw= ?, desttp= ?, desto= ?, destpr= ?, destkw= ?, destc= ? WHERE 1`

	getMapsQuery = `SELECT * FROM hc.maps ORDER BY id`
	addMapQuery  = `INSERT INTO hc.maps (title, pict, w, h, description, id) VALUES (?,?,?,?,?,?)`
	updMapQuery  = `UPDATE hc.maps SET title=?, pict=?, w=?, h=?, description=?  WHERE id=?`
	delMapQuery  = `DELETE FROM hc.maps WHERE id=?`
	//lastMapIdQuery = `SELECT max(id) as id FROM hc.maps`

	getMapSensorsQuery = `SELECT * FROM hc.map_sensors WHERE map_id= ? ORDER BY id`
	addMapSensorQuery  = `INSERT INTO hc.map_sensors (map_id, device_id, type, xk, yk, pict, description, id) VALUES (?,?,?,?,?,?,?,?)`
	updMapSensorQuery  = `UPDATE hc.map_sensors SET map_id=?, device_id=?, type=?, xk=?, yk=?, pict=?, description=? WHERE id=?`
	delMapSensorQuery  = `DELETE FROM hc.map_sensors WHERE id=?`
	//lastMapSensorIdQuery = `SELECT max(id) as id FROM hc.map_sensors`
)

// database структура подключения к базе данных
type database struct {
	conn *sqlx.DB
}

// dbService представляет интерфейс взаимодействия с базой данных
type dbService interface {
	auth(string, string) ([]*User, error)

	getLastId(table string) (int, error)

	getUsers() ([]User, error)
	editUser(id int, login string, pass string, userType string, actFlag string, lastV time.Time) (bool, error)
	delUser(id int) (bool, error)

	getDevices() ([]Device, error)
	editDevice(id int, devType string, devName string, ip string, actFlag string, descr string) (bool, error)
	delDevice(id int) (bool, error)

	getKotelID() (string, error)
	getKotelData() (KotelData, error)
	updtKotelData(tp float64, to float64, pr float64, kw int, desttp float64, desto float64, destpr float64, destkw int, destc float64) error

	getMaps() ([]Map, error)
	editMap(id int, title string, pict string, w int, h int, descr string) (bool, error)
	delMap(id int) (bool, error)

	getMapSensors(mapId int) ([]MapSensor, error)
	editMapSensor(id int, mapId int, devId int, sensorType string, xk float64, yk float64, pict string, descr string) (bool, error)
	delMapSensor(id int) (bool, error)
}

// newDB открывает соединение с базой данных
func newDB(connectionString string) (database, error) {
	dbConn, err := sqlx.Open("mysql", connectionString)
	log.Println(connectionString)
	return database{conn: dbConn}, err
}

//#################################################################
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
		u := User{uid, login, pass, active, userType, lastV}
		users = append(users, &u)
	}

	return users, err
}

func (db database) getLastId(table string) (int, error) {
	var lastId int
	stmt, err := db.conn.Prepare("SELECT max(id) as id FROM hc." + table)
	if err != nil {
		return -1, err
	}
	defer stmt.Close()

	rows, err := stmt.Query()
	for rows.Next() {
		err = rows.Scan(&lastId)
		if err != nil {
			return -1, err
		}
	}

	return lastId, err
}

// ############## Users ############################
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
		u := User{uid, login, pass, active, userType, lastV}
		users = append(users, u)
	}

	return users, err
}

func (db database) editUser(id int, login string, pass string, userType string, actFlag string, lastV time.Time) (bool, error) {

	var lastId int

	execQuery := updUserQuery

	lastId, err := db.getLastId("users")
	if err != nil {
		return false, err
	}

	if id > lastId {
		execQuery = addUserQuery
	}

	stmt, err := db.conn.Prepare(execQuery)
	if err != nil {
		return false, err
	}

	_, err = stmt.Exec(login, pass, userType, actFlag, lastV, id)
	if err != nil {
		return false, err
	}

	return true, err
}

func (db database) delUser(id int) (bool, error) {

	stmt, err := db.conn.Prepare(delUserQuery)
	if err != nil {
		return false, err
	}
	defer stmt.Close()

	_, err = stmt.Exec(id)
	if err != nil {
		return false, err
	}

	return true, err
}

//##################### Devices ########################################

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

func (db database) editDevice(id int, devType string, devName string, ip string, actFlag string, descr string) (bool, error) {

	var (
		lastId int
		kId    string
	)
	if devType == "KotelController" {
		kId, _ = db.getKotelID()
		if kId != devName {
			return false, errors.New("Устройство с типом KotelController уже существует. Такое устройсто может быть только одно.")
		} else {
			stmt, err := db.conn.Prepare(updKotelDevIdQuery)
			if err != nil {
				return false, err
			}
			_, err = stmt.Exec(devName)
			if err != nil {
				return false, err
			}
		}

	}

	execQuery := updDeviceQuery

	lastId, err := db.getLastId("device")
	if err != nil {
		return false, err
	}

	if id > lastId {
		execQuery = addDeviceQuery
	}

	stmt, err := db.conn.Prepare(execQuery)
	if err != nil {
		return false, err
	}

	_, err = stmt.Exec(devType, devName, ip, actFlag, descr, id)
	if err != nil {
		return false, err
	}

	return true, err
}

func (db database) delDevice(id int) (bool, error) {

	stmt, err := db.conn.Prepare(delDeviceQuery)
	if err != nil {
		return false, err
	}
	defer stmt.Close()

	_, err = stmt.Exec(id)
	if err != nil {
		return false, err
	}

	return true, err
}

//################# Kotel #####################
func (db database) getKotelID() (string, error) {
	var kId string
	stmt, err := db.conn.Prepare(kotelDevIdQuery)
	if err != nil {
		return kId, err
	}
	defer stmt.Close()

	rows, err := stmt.Query()

	for rows.Next() {
		err = rows.Scan(&kId)
		if err != nil {
			return kId, err
		}
	}

	return kId, err
}

func (db database) getKotelData() (KotelData, error) {
	var kData KotelData

	stmt, err := db.conn.Prepare(getKotelDataQuery)
	if err != nil {
		return kData, err
	}
	defer stmt.Close()

	rows, err := stmt.Query()
	if err != nil {
		return kData, err
	}

	for rows.Next() {
		var (
			deviceId string
			tp       float64
			to       float64
			pr       float64
			kw       int
			desttp   float64
			destto   float64
			destpr   float64
			destkw   int
			destc    float64
		)
		//SELECT DEVICE_ID , TP, TO, PR, KW, DEST_TP, DEST_TO, DEST_PR, DEST_KW, DEST_TC FROM hc.kotel
		err = rows.Scan(&deviceId, &tp, &to, &pr, &kw, &desttp, &destto, &destpr, &destkw, &destc)
		if err != nil {
			return kData, err
		}

		kData = KotelData{deviceId, tp, to, pr, kw, desttp, destto, destpr, destkw, destc}
	}

	return kData, err
}

func (db database) updtKotelData(tp float64, to float64, pr float64, kw int, desttp float64, desto float64, destpr float64, destkw int, destc float64) error {

	stmt, err := db.conn.Prepare(updKotelDevIdQuery)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(tp, to, pr, kw, desttp, desto, destpr, destkw, destc)
	if err != nil {
		return err
	}

	return err
}

//################## Maps #########################
func (db database) getMaps() ([]Map, error) {
	maps := make([]Map, 0)
	stmt, err := db.conn.Prepare(getMapsQuery)
	if err != nil {
		return maps, err
	}
	defer stmt.Close()

	rows, err := stmt.Query()

	for rows.Next() {
		var (
			id    int
			title string
			pict  string
			w     int
			h     int
			descr string
		)

		err = rows.Scan(&id, &title, &pict, &w, &h, &descr)
		if err != nil {
			return maps, err
		}
		m := Map{id, title, pict, w, h, descr}
		maps = append(maps, m)
	}

	return maps, err
}

func (db database) editMap(id int, title string, pict string, w int, h int, descr string) (bool, error) {

	var lastId int

	execQuery := updMapQuery

	lastId, err := db.getLastId("maps")
	if err != nil {
		return false, err
	}

	if id > lastId {
		execQuery = addMapQuery
	}

	stmt, err := db.conn.Prepare(execQuery)
	if err != nil {
		return false, err
	}

	_, err = stmt.Exec(title, pict, w, h, descr, id)
	if err != nil {
		return false, err
	}

	return true, err
}

func (db database) delMap(id int) (bool, error) {

	stmt, err := db.conn.Prepare(delMapQuery)
	if err != nil {
		return false, err
	}
	defer stmt.Close()

	_, err = stmt.Exec(id)
	if err != nil {
		return false, err
	}

	return true, err
}

//###################### MapSensors #####################
func (db database) getMapSensors(mapId int) ([]MapSensor, error) {
	mapSens := make([]MapSensor, 0)
	stmt, err := db.conn.Prepare(getMapSensorsQuery)
	if err != nil {
		return mapSens, err
	}
	defer stmt.Close()

	rows, err := stmt.Query(mapId)

	for rows.Next() {
		var (
			id       int
			mapId    int
			devId    int
			sensType string
			xk       float32
			yk       float32
			pict     string
			descr    string
		)

		err = rows.Scan(&id, &mapId, &devId, &sensType, &xk, &yk, &pict, &descr)
		if err != nil {
			return mapSens, err
		}
		s := MapSensor{id, mapId, devId, sensType, xk, yk, pict, descr}
		mapSens = append(mapSens, s)
	}

	return mapSens, err
}

func (db database) editMapSensor(id int, mapId int, devId int, sensorType string, xk float64, yk float64, pict string, descr string) (bool, error) {

	var lastId int

	execQuery := updMapSensorQuery

	lastId, err := db.getLastId("map_sensors")
	if err != nil {
		return false, err
	}

	if id > lastId {
		execQuery = addMapSensorQuery
	}

	stmt, err := db.conn.Prepare(execQuery)
	if err != nil {
		return false, err
	}

	_, err = stmt.Exec(mapId, devId, sensorType, xk, yk, pict, descr, id)
	if err != nil {
		return false, err
	}

	return true, err
}

func (db database) delMapSensor(id int) (bool, error) {

	stmt, err := db.conn.Prepare(delMapSensorQuery)
	if err != nil {
		return false, err
	}
	defer stmt.Close()

	_, err = stmt.Exec(id)
	if err != nil {
		return false, err
	}

	return true, err
}
