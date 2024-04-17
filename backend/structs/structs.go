package structs

import (
	"sync"

	"github.com/gorilla/websocket"
)

type Player struct {
	// The actual websocket connection.
	Conn     *websocket.Conn
	Username string `json:"username"`
	Id       int    `json:"id"`
	X        int    `json:"x"`
	Y        int    `json:"y"`
	Dir      string `json:"dir"`
	Done     bool
	Lockmx   sync.Mutex
	Leader   bool `json:"leader"`
	// Red string `json:"color"`
}

type PlayerMove struct {
	X        int    `json:"x"`
	Y        int    `json:"y"`
	OldX     int    `json:"oldx"`
	OldY     int    `json:"oldy"`
	Dir      string `json:"dir"`
	PlayerId int    `json:"playerid"`
}

type ClientMsg struct {
	MsgType string      `json:"msgtype"`
	MsgData interface{} `json:"msgdata"`
}

// Grid size
type GridSize struct {
	X      int     `json:"x"`
	Y      int     `json:"y"`
	Player *Player `json:"player"`
}

type UpdatedStruct struct {
	X      int     `json:"x"`
	Y      int     `json:"y"`
	Grid   [][]int `json:"grid"`
	Player *Player `json:"player"`
}

type FinishedMsg struct {
	Message string `json:"message"`
}
