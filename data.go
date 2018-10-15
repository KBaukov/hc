package main

import (
	"time"
)

type User struct {
	ID          int       `json:"id"`
	LOGIN       string    `json:"login"`
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

//type ApiError struct {
//	ERROR_CODE   int    `json:"errorCode"`
//	EROR_MESSAGE string `json:"errorMEssage"`
//}

type ApiResp struct {
	SUCCESS bool        `json:"success"`
	DATA    interface{} `json:"data"`
	MSG     string      `json:"msg"`
}
