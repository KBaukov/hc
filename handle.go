package main

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/gob"
	"encoding/hex"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/sessions"
	"github.com/gorilla/websocket"
)

var (
	sessStore = sessions.NewCookieStore([]byte("33446a9dcf9ea060a0a6532b166da32f304af0de"))
)

func init() {
	gob.Register(User{})
	gob.Register(websocket.Conn{})

	sessStore.Options = &sessions.Options{
		Domain:   "localhost",
		Path:     "/",
		MaxAge:   3600 * 8, // 8 hours
		HttpOnly: true,
	}
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

		if r.Method == "POST" {

			login := r.PostFormValue("username")
			pass := r.PostFormValue("password")

			pass, err := hashPass(pass)
			if err != nil {
				log.Printf("Ошибка хеширования", err)
			}

			//log.Println("pass: ", pass)

			users, err := db.auth(login, pass)
			if err != nil {
				http.Error(w, "Ошибка обработки запроса", http.StatusInternalServerError)
				log.Printf("Ошибка авторизации (логин: %v): %v", login, err)
			}

			if len(users) != 1 {
				http.Redirect(w, r, "/login", 403)
				return
			}

			u := *users[0]
			log.Println("Login user: ", u)
			createSession(w, r, u, "user")

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

		//##############
		if r.URL.Path == "/api/devices" {
			devices, err := db.getDevices()
			apiDataResponse(w, devices, err)
		}
		if r.URL.Path == "/api/device/edit" {
			id := r.PostFormValue("id")
			intId, err := strconv.Atoi(id)
			if err != nil {
				log.Println("err:", err.Error())
			}

			devType := r.PostFormValue("type")
			devName := r.PostFormValue("name")
			ip := r.PostFormValue("ip")
			actFlag := r.PostFormValue("active_flag")
			descr := r.PostFormValue("description")

			_, err = db.editDevice(intId, devType, devName, ip, actFlag, descr)
			apiDataResponse(w, []int{}, err)
		}
		if r.URL.Path == "/api/device/delete" {
			id := r.PostFormValue("id")
			intId, err := strconv.Atoi(id)
			if err != nil {
				log.Println("err:", err.Error())
			}

			_, err = db.delDevice(intId)
			apiDataResponse(w, []int{}, err)
		}
		//###################
		if r.URL.Path == "/api/users" {
			users, err := db.getUsers()
			apiDataResponse(w, users, err)
		}
		if r.URL.Path == "/api/user/edit" {
			id := r.PostFormValue("id")
			intId, err := strconv.Atoi(id)
			if err != nil {
				log.Println("err:", err.Error())
			}

			login := r.PostFormValue("login")
			pass := r.PostFormValue("pass")
			pass, err = hashPass(pass)
			if err != nil {
				log.Printf("Ошибка хеширования", err)
			}
			usrType := r.PostFormValue("user_type")
			actFlag := r.PostFormValue("active_flag")
			lastVs := r.PostFormValue("last_visit")
			lastV, err := time.Parse("2006-01-02T00:00:00Z", lastVs)
			if err != nil {
				log.Println("date forma validation error:", err.Error())
			}

			_, err = db.editUser(intId, login, pass, usrType, actFlag, lastV)
			apiDataResponse(w, []int{}, err)
		}
		if r.URL.Path == "/api/user/delete" {
			id := r.PostFormValue("id")
			intId, err := strconv.Atoi(id)
			if err != nil {
				log.Println("err:", err.Error())
			}

			_, err = db.delUser(intId)
			apiDataResponse(w, []int{}, err)
		}
		//##################
		if r.URL.Path == "/api/maps" {
			maps, err := db.getMaps()
			apiDataResponse(w, maps, err)
		}
		if r.URL.Path == "/api/maps/edit" {
			ids := r.PostFormValue("id")
			id, err := strconv.Atoi(ids)
			if err != nil {
				log.Println("err:", err.Error())
			}
			title := r.PostFormValue("title")
			pict := r.PostFormValue("pict")

			ws := r.PostFormValue("w")
			wi, err := strconv.Atoi(ws)
			if err != nil {
				log.Println("err:", err.Error())
			}
			hs := r.PostFormValue("h")
			h, err := strconv.Atoi(hs)
			if err != nil {
				log.Println("err:", err.Error())
			}
			descr := r.PostFormValue("description")

			if err != nil {
				log.Println("date forma validation error:", err.Error())
			}

			_, err = db.editMap(id, title, pict, wi, h, descr)
			apiDataResponse(w, []int{}, err)
		}
		if r.URL.Path == "/api/maps/delete" {
			ids := r.PostFormValue("id")
			id, err := strconv.Atoi(ids)
			if err != nil {
				log.Println("err:", err.Error())
			}

			_, err = db.delMap(id)
			apiDataResponse(w, []int{}, err)
		}
		//##############################################
		if r.URL.Path == "/api/sensors" {
			mapIds := r.PostFormValue("map_id")
			mapId, err := strconv.Atoi(mapIds)
			if err != nil {
				log.Println("err:", err.Error())
			} //////
			sens, err := db.getMapSensors(mapId)
			apiDataResponse(w, sens, err)
		}
		if r.URL.Path == "/api/sensors/edit" {
			ids := r.PostFormValue("id")
			id, err := strconv.Atoi(ids)
			if err != nil {
				log.Println("err:", err.Error())
			} //////

			mapIds := r.PostFormValue("map_id")
			mapId, err := strconv.Atoi(mapIds)
			if err != nil {
				log.Println("err:", err.Error())
			} /////

			devIds := r.PostFormValue("dev_id")
			devId, err := strconv.Atoi(devIds)
			if err != nil {
				log.Println("err:", err.Error())
			} /////

			xks := r.PostFormValue("xk")
			xk, err := strconv.ParseFloat(xks, 64)
			if err != nil {
				log.Println("err:", err.Error())
			} /////

			yks := r.PostFormValue("yk")
			yk, err := strconv.ParseFloat(yks, 64)
			if err != nil {
				log.Println("err:", err.Error())
			} /////

			sensorType := r.PostFormValue("type")
			pict := r.PostFormValue("pict")
			descr := r.PostFormValue("description")

			_, err = db.editMapSensor(id, mapId, devId, sensorType, xk, yk, pict, descr)
			apiDataResponse(w, []int{}, err)
		}
		if r.URL.Path == "/api/sensors/delete" {
			ids := r.PostFormValue("id")
			id, err := strconv.Atoi(ids)
			if err != nil {
				log.Println("err:", err.Error())
			}

			_, err = db.delMapSensor(id)
			apiDataResponse(w, []int{}, err)
		}
		if r.URL.Path == "/api/sensors/lastid" {
			sens, err := db.getLastId("map_sensors")
			apiDataResponse(w, sens, err)
		}

		//##############################################
		if r.URL.Path == "/api/kotel/getvalues" {
			data, err := db.getKotelData()
			//data := KotelData{0, 30.56, 45.12, 2.21, 11, 30.0, 45.0, 2.25, 11, 25.0, time.Now()}
			apiDataResponse(w, data, err)
		}
		if r.URL.Path == "/api/kotel/setdest" {

			tp, err := strconv.ParseFloat(r.PostFormValue("tp"), 64)
			to, err := strconv.ParseFloat(r.PostFormValue("to"), 64)
			pr, err := strconv.ParseFloat(r.PostFormValue("pr"), 64)
			kw, err := strconv.Atoi(r.PostFormValue("kw"))
			desttp, err := strconv.ParseFloat(r.PostFormValue("desttp"), 64)
			destto, err := strconv.ParseFloat(r.PostFormValue("destto"), 64)
			destpr, err := strconv.ParseFloat(r.PostFormValue("destpr"), 64)
			destkw, err := strconv.Atoi(r.PostFormValue("destkw"))
			destc, err := strconv.ParseFloat(r.PostFormValue("desttc"), 64)
			if err != nil {

			}

			kd, err := db.getKotelData()

			if tp == 0 {
				tp = kd.TP
			}
			if pr == 0 {
				pr = kd.PR
			}
			if kw == 0 {
				kw = kd.KW
			}
			if desttp == 0 {
				desttp = kd.DESTTP
			}
			if destto == 0 {
				destto = kd.DESTTO
			}
			if destkw == 0 {
				destkw = kd.DESTKW
			}
			if destpr == 0 {
				destpr = kd.DESTPR
			}
			if destc == 0 {
				destc = kd.DESTС
			}

			err = db.updtKotelData(tp, to, pr, kw, desttp, destto, destpr, destkw, destc)

			apiDataResponse(w, []int{}, err)

		}
		if r.URL.Path == "/api/kotel/pressbutt" {
			var (
				msg     string
				err     error
				kotelId string
			)

			kotelId, err = db.getKotelID()
			if err != nil || kotelId == "" {
				err = errors.New("Котел не найден")
			}

			butt := r.PostFormValue("button")

			ws := wsConnections[kotelId]
			if ws == nil {
				err = errors.New("Сессия не активна")
			} else {
				msg = "{\"action\":\"pessButton\", \"butt\":\"" + butt + "\"}"
				log.Printf("Sending message to %s: %s", kotelId, msg)
				err = ws.WriteMessage(1, []byte(butt))
				if err != nil {
					log.Println("Sending message error:", err)
				}
			}

			apiDataResponse(w, msg, err)
		}

		return
	}
}

func apiDataResponse(w http.ResponseWriter, data interface{}, err error) {
	errMsg := ""
	succes := true

	if err != nil {
		//http.Error(w, "Ошибка обработки запроса", http.StatusInternalServerError)
		log.Printf("Ошибка: %v", err)
		errMsg = err.Error()
		succes = false
	}

	dataResp := ApiResp{SUCCESS: succes, DATA: data, MSG: errMsg}

	json, err := json.Marshal(dataResp)
	if err != nil {
		//http.Error(w, "Ошибка формирования ответа", http.StatusInternalServerError)
		log.Printf("Ошибка маршалинга: %v", err)
		return
	}
	w.Header().Set("Content-type", "application/json; charset=utf-8")
	_, err = w.Write(json)
	if err != nil {
		log.Printf("Ошибка записи результата запроса: %v", err)
	}
}

//########################## helpers ############################

func hashPass(p string) (string, error) {
	h := sha256.New()
	_, err := h.Write([]byte(p))
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(h.Sum(nil)), nil
}

func createSession(w http.ResponseWriter, r *http.Request, o interface{}, key string) {

	session, err := sessStore.Get(r, "session-name")
	if err != nil {
		log.Printf("Error getting session: %v", err)
	}

	session.Values[key] = o //User{"Pogi", "Points", ""}
	session.Save(r, w)

	log.Println("Session initiated")
}

func getSession(w http.ResponseWriter, r *http.Request) *sessions.Session {
	session, err := sessStore.Get(r, "session-name")
	if err != nil {
		log.Printf("Error getting session: %v", err)
		session, err = sessStore.New(r, "session-name")
	}
	return session
}
