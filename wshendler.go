// Copyright 2015 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bufio"
	//"flag"
	"io"
	"log"
	"net/http"
	//"os"
	"strings"
	//"os/exec"
	//"encoding/hex"
	"time"

	"github.com/gorilla/websocket"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Maximum message size allowed from peer.
	maxMessageSize = 8192

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Time to wait before force close on connection.
	closeGracePeriod = 10 * time.Second
)

func pumpStdin(ws *websocket.Conn, w io.Writer) {
	defer ws.Close()
	ws.SetReadLimit(maxMessageSize)
	ws.SetReadDeadline(time.Now().Add(pongWait))
	ws.SetPongHandler(func(string) error { ws.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		_, message, err := ws.ReadMessage()
		if err != nil {
			break
		}
		message = append(message, '\n')
		if _, err := w.Write(message); err != nil {
			break
		}
	}
}

func pumpStdout(ws *websocket.Conn, r io.Reader, done chan struct{}) {
	defer func() {}()
	s := bufio.NewScanner(r)
	for s.Scan() {
		ws.SetWriteDeadline(time.Now().Add(writeWait))
		if err := ws.WriteMessage(websocket.TextMessage, s.Bytes()); err != nil {
			ws.Close()
			break
		}
	}
	if s.Err() != nil {
		log.Println("scan:", s.Err())
	}
	close(done)

	ws.SetWriteDeadline(time.Now().Add(writeWait))
	ws.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
	time.Sleep(closeGracePeriod)
	ws.Close()
}

func ping(ws *websocket.Conn, done chan struct{}) {
	ticker := time.NewTicker(pingPeriod)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			if err := ws.WriteControl(websocket.PingMessage, []byte{}, time.Now().Add(writeWait)); err != nil {
				log.Println("ping:", err)
			}
		case <-done:
			return
		}
	}
}

func internalError(ws *websocket.Conn, msg string, err error) {
	log.Println(msg, err)
	ws.WriteMessage(websocket.TextMessage, []byte("Internal server error."))
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		origin := r.Header.Get("Origin")
		log.Println("Origin:", origin)

		return true
	},
}

func serveWs(w http.ResponseWriter, r *http.Request) {
	var ws *websocket.Conn
	var err error
	deviceId := r.Header.Get("deviceid")
	log.Println("incoming request from: ", deviceId)

	ws, err = upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("upgrade:", err)
		return
	}
	//defer ws.Close()

	log.Println("Ws Connection: ", ws)
	go wsProcessor(ws)

}

func wsProcessor(c *websocket.Conn) {
	defer c.Close()
	for {
		mt, message, err := c.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			break
		}
		log.Printf("recv: %s", message)
		msg := string(message)

		if strings.Contains(msg, "{action:connect") {
			if !sendMsg(c, mt, "{action:connect,success:true}") {
				break
			}
		}

		if strings.Contains(msg, "Important") {

			if !sendMsg(c, mt, "Ok. I'm...") {
				break
			}
		}

	}

}

func sendMsg(c *websocket.Conn, n int, m string) bool {
	err := c.WriteMessage(n, []byte(m))
	if err != nil {
		log.Println("write:", err)
		return false
	}
	return true
}

func serveHomeWs(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/homews" {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	http.ServeFile(w, r, "./webres/html/ws.html")
}
