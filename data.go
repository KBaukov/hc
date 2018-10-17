package main

import (
	"time"
)

type User struct {
	ID          int       `json:"id"`
	LOGIN       string    `json:"login"`
	PASS        string    `json:"pass"`
	ACTIVE_FLAG string    `json:"active_flag"`
	USER_TYPE   string    `json:"user_type"`
	LAST_VISIT  time.Time `json:"last_visit"`
}

type Device struct {
	ID          int    `json:"id"`
	TYPE        string `json:"type"`
	NAME        string `json:"name"`
	IP          string `json:"ip"`
	ACTIVE_FLAG string `json:"active_flag"`
	DESCRIPTION string `json:"description"`
}

type ApiError struct {
	ERROR_CODE   int    `json:"errorCode"`
	EROR_MESSAGE string `json:"errorMEssage"`
}

type ApiResp struct {
	SUCCESS bool        `json:"success"`
	DATA    interface{} `json:"data"`
	MSG     string      `json:"msg"`
}

type Map struct {
	ID          int    `json:"id"`
	TITLE       string `json:"title"`
	PICT        string `json:"pict"`
	W           int    `json:"w"`
	H           int    `json:"h"`
	DESCRIPTION string `json:"description"`
}

type MapSensor struct {
	ID          int     `json:"id"`
	MAP_ID      int     `json:"map_id"`
	DEVICE_ID   int     `json:"device_id"`
	TYPE        string  `json:"type"`
	XK          float32 `json:"xk"`
	YK          float32 `json:"yk"`
	PICT        string  `json:"pict"`
	DESCRIPTION string  `json:"description"`
}

type KotelData struct {
	DEVICE_ID int       `json:"device_id"`
	TP        float64   `json:"tp"`
	TO        float64   `json:"to"`
	PR        float32   `json:"pr"`
	KW        int       `json:"kw"`
	DESTTP    float64   `json:"desttp"`
	DESTTO    float64   `json:"destto"`
	DESTPR    float64   `json:"destpr"`
	DESTKW    int       `json:"destkw"`
	DESTС     float64   `json:"destс"`
	DATE      time.Time `json:"date"`
}
