package handlers

import (
	"dsproject/backend/structs"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/websocket"
)

// Upgrades request to websocket
var upgrader = websocket.Upgrader{
	ReadBufferSize:  2048,
	WriteBufferSize: 2048,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

const GRIDX int = 16
const GRIDY int = 10

var gameOn bool = false //Start game, not implemented
var gameFinished bool = false

// Incrementing player id
var playerID = 1

// List of players
var playerList []*structs.Player

// Upgrade HTTP connection to websocket
func NewPlayer(w http.ResponseWriter, r *http.Request) {
	//Upgrade to websocket
	wscon, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	//Make new player
	player := structs.Player{
		Conn: wscon,
		Id:   playerID,
	}
	playerList = append(playerList, &player)
	playerID += 1

	//Set grid size
	gridSize := structs.GridSize{
		X:        GRIDX,
		Y:        GRIDY,
		PlayerId: player.Id,
	}
	//On first player join, initiate grid with size
	if len(playerList) == 1 {

		fmt.Println(getGrid())
		createGrid()
		fmt.Println(getGrid())

		startGrid := structs.ClientMsg{
			MsgType: "initgrid",
			MsgData: gridSize,
		}
		//Send grid size to player
		err := player.Conn.WriteJSON(startGrid)
		if err != nil {
			return // Maybe do a retry / drop connection here
		}
	} else {
		//If not first player, set board content for new player

		updatedGrid := structs.UpdatedStruct{
			X:        GRIDX,
			Y:        GRIDY,
			Grid:     getGrid(),
			PlayerId: player.Id,
		}

		startGrid := structs.ClientMsg{
			MsgType: "updategrid",
			MsgData: updatedGrid,
		}
		//Send grid size to player
		err := player.Conn.WriteJSON(startGrid)
		if err != nil {
			return // Maybe do a retry or drop connection here
		}
	}

	fmt.Println("New Client " + strconv.Itoa(player.Id) + " joined the hub!")

	player.Done = false
	SocketListener(&player)
}

// Listen for messages from players
func SocketListener(player *structs.Player) {
	for !player.Done {
		time.Sleep(250 * time.Millisecond)
		msg := structs.ClientMsg{}
		err := player.Conn.ReadJSON(&msg)
		if err != nil || gameFinished {
			fmt.Println("Error joinplayer.go userid:", player.Id)
			fmt.Println(err)
			//Remove player after losing connection
			player.Conn.Close()
			for i, p := range playerList {
				if p.Id == player.Id {
					playerList = append(playerList[:i], playerList[i+1:]...)
				}
			}
			break
		}
		messageHandler(player, &msg)
	}
}

// Handle messages from players
func messageHandler(player *structs.Player, msg *structs.ClientMsg) {

	//Handle messages from players
	switch msg.MsgType {
	case "move":
		SnakeMove(player, msg)
		break
	default:
		fmt.Println("default messagehandler")
		break
	}
}

func setGameFinished() {
	gameFinished = true
}
