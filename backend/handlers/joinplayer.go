package handlers

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/websocket"
)

// Upgrades request to websocket
var upgrader = websocket.Upgrader{
	ReadBufferSize:  2048,
	WriteBufferSize: 2048,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

type Player struct {
	// The actual websocket connection.
	conn     *websocket.Conn
	Username string `json:"username"`
	Id       int    `json:"id"`
	X        int    `json:"x"`
	Y        int    `json:"y"`
}

type ClientMsg struct {
	MsgType string      `json:"msgtype"`
	MsgData interface{} `json:"msgdata"`
}

var playerID = 0
var playerList []Player

// Upgrade HTTP connection to websocket
func ServeWs(w http.ResponseWriter, r *http.Request) {
	//Upgrades connection
	wscon, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	//Make new player
	player := newPlayer(wscon)
	playerList = append(playerList, *player)
	playerID += 1

	fmt.Println("New Client joined the hub!")
	fmt.Println(player)

	//Single player
	_ = wscon.WriteJSON(
		ClientMsg{
			MsgType: "initplayer",
			MsgData: player},
	)
	//All players
	for _, p := range playerList {
		_ = p.conn.WriteJSON(
			ClientMsg{
				MsgType: "hola",
				MsgData: "Just joined: " + strconv.Itoa(player.Id)},
		)
	}

}

func newPlayer(conn *websocket.Conn) *Player {
	return &Player{
		conn: conn,
		Id:   playerID,
	}
}
