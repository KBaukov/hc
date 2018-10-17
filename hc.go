// main.go project main.go
package main

import (
	"flag"
	"log"
	"net"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/websocket"
)

type Listener chan net.Conn

var (
	configurationPath = flag.String("config", "config.json", "Путь к файлу конфигурации")
	config            = loadConfig(*configurationPath)
	wsConnections     = make(map[string]*websocket.Conn)
)

func init() {

}

func main() {
	flag.Parse()

	if config.LoggerPath != "" {
		// Логер только добавляет данные
		logFile, err := os.OpenFile(config.LoggerPath, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
		if err != nil {
			log.Printf("Ошибка открытия файла лога: %v", err)
		} else {
			defer logFile.Close()
			log.SetOutput(logFile)
		}
	}

	db, err := newDB(config.DbConnectionString)
	if err != nil {
		log.Printf("Не удалось подключиться к базе данных: %v", err)
	} else {
		_, err = db.conn.Exec("SET AUTOCOMMIT=1;")
		if err != nil {
			log.Printf("Не удалось установить настройки базы данных: %v", err)
		}
	}

	//go server.Serve(l)

	http.HandleFunc("/logout", serveHome)
	http.HandleFunc("/login", serveLogin(db))
	http.HandleFunc("/home", serveHome)
	http.HandleFunc("/webres/", serveWebRes)
	http.HandleFunc("/api/", serveApi(db))
	//http.HandleFunc("/ws", STOMPserver)
	http.HandleFunc("/ws", serveWs)
	http.HandleFunc("/homews", serveHomeWs)
	//log.Fatal(http.ListenAndServe(*addr, nil))

	listenString := config.Server.Address + ":" + config.Server.Port
	log.Print("Запуск сервера: ", listenString)

	if config.Server.TLS {
		err = http.ListenAndServeTLS(listenString, config.Server.CertificatePath, config.Server.KeyPath, nil)
	} else {
		err = http.ListenAndServe(listenString, nil)
	}
	if err != nil {
		log.Printf("Ошибка веб-сервера: %v", err)
	}

}

//func STOMPserver(w http.ResponseWriter, req *http.Request) {
//	deviceId := req.Header.Get("deviceid")
//	log.Println("incoming request from: ", deviceId)
//	// Create a net.Conn from the websocket
//	ws, err := upgraderr.Upgrade(w, req, nil)
//	if err != nil {
//		log.Println("error getting websocket conn:", err)
//		return
//	}

//	//conn := ws.UnderlyingConn()

//	// Send this connection to our STOMP server
//	l <- ws.UnderlyingConn()
//}

//func (l Listener) Accept() (c net.Conn, err error) {

//	log.Println("Accept:")
//	return <-l, nil
//}

//func (l Listener) Close() (err error) {
//	log.Println("Close:")
//	return nil
//}

//func (l Listener) Addr() (a net.Addr) {
//	log.Println("Addr:")
//	return nil
//}
