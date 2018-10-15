// main.go project main.go
package main

import (
	"flag"
	"log"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/sessions"
)

var (
	sessStore         = sessions.NewCookieStore([]byte("33446a9dcf9ea060a0a6532b166da32f304af0de"))
	configurationPath = flag.String("config", "config.json", "Путь к файлу конфигурации")
	config            = loadConfig(*configurationPath)
)

func init() {

	sessStore.Options = &sessions.Options{
		Domain:   "localhost",
		Path:     "/",
		MaxAge:   3600 * 8, // 8 hours
		HttpOnly: true,
	}

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
	}

	http.HandleFunc("/logout", serveHome)
	http.HandleFunc("/login", serveLogin(db))
	http.HandleFunc("/home", serveHome)
	http.HandleFunc("/webres/", serveWebRes)
	http.HandleFunc("/api/", serveApi(db))
	//http.HandleFunc("/ws", serveWs)
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
