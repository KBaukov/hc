package main

import (
	"encoding/base64"
	"encoding/gob"
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/sessions"
)

func init() {
	gob.Register(User{})

}

func serveHome(w http.ResponseWriter, r *http.Request) {

	log.Println("Path: " + r.URL.Path)

	if r.URL.Path == "/logout" {
		session := getSession(w, r)
		log.Println("++++Session: ", session)
		delete(session.Values, "user")
		session.Save(r, w)
		log.Println("----Session: ", session)
		now := []byte(time.Now().String())
		sha := base64.URLEncoding.EncodeToString(now)
		http.Redirect(w, r, "/home?dc="+sha, 301)
		//http.ServeFile(w, r, "./webres/html/login.html")
	}

	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if r.URL.Path != "/home" {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	} else {
		session := getSession(w, r)
		log.Println("Session: ", session)
		if session.Values["user"] == nil {
			log.Println("Redirect: /login")
			http.Redirect(w, r, "/login", 301)
		} else {
			http.ServeFile(w, r, "./webres/html/main.html")
		}
	}

}

func serveLogin(db dbService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		if r.Method != "GET" && r.Method != "POST" {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		//		sess := getSession(w, r)

		//		if sess.Values["user"] != nil {
		//			log.Println("Login user: ", sess.Values["user"])
		//			//http.Redirect(w, r, "/home", 301)
		//		}

		if r.Method == "POST" {

			login := r.PostFormValue("username")
			pass := r.PostFormValue("password")

			users, err := db.auth(login, pass)
			if err != nil {
				http.Error(w, "Ошибка обработки запроса", http.StatusInternalServerError)
				log.Printf("Ошибка авторизации (логин: %v): %v", login, err)
				//http.Error(w, "Ошибка доступа", http.StatusNonAuthoritativeInfo)
				//return
			}

			if len(users) != 1 {
				http.Redirect(w, r, "/login", 401)
				return
			}

			u := *users[0]
			log.Println("Login user: ", u)
			createSession(w, r, u)

			http.Redirect(w, r, "/home", 301)

		}
		if r.Method == "GET" {
			http.ServeFile(w, r, "./webres/html/login.html")
		}

	}

}

func serveLogout(w http.ResponseWriter, r *http.Request) {
	session := getSession(w, r)
	log.Println("++++Session: ") //, session)
	session.Values["user"] = nil
	session.Save(r, w)
	log.Println("Session: ", session)

	http.Redirect(w, r, "/login", 301)

}

func serveWebRes(w http.ResponseWriter, r *http.Request) {
	if strings.Contains(r.URL.Path, "webres") {
		http.ServeFile(w, r, "."+r.URL.Path)
	}
}

func serveApi(db dbService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		log.Printf("incoming request in: %v", r.URL.Path)
		if strings.Contains(r.URL.Path, "/getvalues") {
			json := "{success:true,tp:30.56,to:45.12,kw:11,pr:2.23,desttp:60.00,destto:45.00,desttc:35.00,destkw:11}"

			b := []byte(json)
			w.Header().Set("Content-type", "application/json; charset=utf-8")
			_, err := w.Write(b)
			if err != nil {
				log.Printf("Ошибка записи результата запроса: %v", err)
			}
			return
		}
		if strings.Contains(r.URL.Path, "/devices") {
			devices, err := db.getDevices()
			if err != nil {
				http.Error(w, "Ошибка обработки запроса", http.StatusInternalServerError)
				log.Printf("Ошибка: %v", err)
			}
			data := ApiResp{SUCCESS: true, DATA: devices, MSG: ""}

			jd, err := json.Marshal(data)
			if err != nil {
				http.Error(w, "Ошибка формирования ответа", http.StatusInternalServerError)
				log.Printf("Ошибка маршалинга: %v", err)
				return
			}
			w.Header().Set("Content-type", "application/json; charset=utf-8")
			_, err = w.Write(jd)
			if err != nil {
				log.Printf("Ошибка записи результата запроса: %v", err)
			}
			return
		}
		if strings.Contains(r.URL.Path, "/users") {
			users, err := db.getUsers()
			if err != nil {
				http.Error(w, "Ошибка обработки запроса", http.StatusInternalServerError)
				log.Printf("Ошибка: %v", err)
			}
			data := ApiResp{SUCCESS: true, DATA: users, MSG: ""}

			ju, err := json.Marshal(data)
			if err != nil {
				http.Error(w, "Ошибка формирования ответа", http.StatusInternalServerError)
				log.Printf("Ошибка маршалинга: %v", err)
				return
			}
			w.Header().Set("Content-type", "application/json; charset=utf-8")
			_, err = w.Write(ju)
			if err != nil {
				log.Printf("Ошибка записи результата запроса: %v", err)
			}
			return
		}
	}
}

func createSession(w http.ResponseWriter, r *http.Request, u User) {

	//gob.Register(User{})

	session, err := sessStore.Get(r, "session-name")
	if err != nil {
		log.Printf("Error getting session: %v", err)
	}

	session.Values["user"] = u //User{"Pogi", "Points", ""}
	session.Save(r, w)

	log.Println("Session initiated")
}

func getSession(w http.ResponseWriter, r *http.Request) *sessions.Session {
	session, err := sessStore.Get(r, "session-name")
	if err != nil {
		log.Printf("Error getting session: %v", err)
		session, err = sessStore.New(r, "session-name")
	}
	//log.Println(session.Values["user"])
	return session
}
