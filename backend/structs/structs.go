package structs

import "github.com/gorilla/websocket"

type Player struct {
	// The actual websocket connection.
	Conn     *websocket.Conn
	Username string `json:"username"`
	Id       int    `json:"id"`
	X        int    `json:"x"`
	Y        int    `json:"y"`
	Done     bool
	// Red string `json:"color"`
}

type PlayerMove struct {
	X        int `json:"x"`
	Y        int `json:"y"`
	PlayerId int `json:"playerid"`
}

type ClientMsg struct {
	MsgType string      `json:"msgtype"`
	MsgData interface{} `json:"msgdata"`
}

// Grid size
type GridSize struct {
	X        int `json:"x"`
	Y        int `json:"y"`
	PlayerId int `json:"playerid"`
}

type UpdatedStruct struct {
	X        int     `json:"x"`
	Y        int     `json:"y"`
	Grid     [][]int `json:"grid"`
	PlayerId int     `json:"playerid"`
}

type FinishedMsg struct {
	Message string `json:"message"`
}
