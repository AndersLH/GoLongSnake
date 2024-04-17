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

var gameOn bool = false
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
	player.Leader = false

	//Set grid size
	gridSize := structs.GridSize{
		X:      GRIDX,
		Y:      GRIDY,
		Player: &player,
	}
	//On first player join, initiate grid with size
	if len(playerList) == 1 {
		createGrid()

		//set first player leader
		player.Leader = true

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
			X:      GRIDX,
			Y:      GRIDY,
			Grid:   getGrid(),
			Player: &player,
		}

		startGrid := structs.ClientMsg{
			MsgType: "updategrid",
			MsgData: updatedGrid,
		}
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
		//Only allow leader to start game
		for !gameOn {
			if player.Leader {
				break
			}
			msg := structs.ClientMsg{}
			player.Conn.ReadJSON(&msg)
		}
		//Annonuce game started for players
		if !player.Leader {
			start := structs.ClientMsg{
				MsgType: "announcestart",
				MsgData: player,
			}

			err := player.Conn.WriteJSON(start)
			if err != nil {
				return // Maybe do a retry or drop connection here
			}
		}

		time.Sleep(250 * time.Millisecond)
		msg := structs.ClientMsg{}
		err := player.Conn.ReadJSON(&msg)

		//Temporary game finished:
		if gameFinished {

			createGrid() //reset grid
			//Set grid size
			for _, p := range playerList {
				gridSize := structs.GridSize{
					X:      GRIDX,
					Y:      GRIDY,
					Player: p,
				}
				startGrid := structs.ClientMsg{
					MsgType: "initgrid",
					MsgData: gridSize,
				}
				//Send grid size to player
				err := p.Conn.WriteJSON(startGrid)
				if err != nil {
					return // Maybe do a retry / drop connection here
				}
			}
			player.X = 0
			player.Y = 0
			gameOn = false
			gameFinished = false
			// return
		}

		if err != nil {
			fmt.Println("Error joinplayer.go userid:", player.Id)
			fmt.Println(err)

			//Remove player after losing connection
			player.Conn.Close()
			for i, p := range playerList {
				if p.Id == player.Id {
					playerList = append(playerList[:i], playerList[i+1:]...)
				}
			}

			//Assign new leader if leader leaves
			if len(playerList) > 0 && player.Leader {
				playerList[0].Leader = true
				newLeader := structs.ClientMsg{
					MsgType: "newleader",
					MsgData: player,
				}

				err := playerList[0].Conn.WriteJSON(newLeader)
				if err != nil {
					return // Maybe do a retry or drop connection here
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
	case "start":
		startGame()
	case "move":
		if !gameOn {
			return
		}
		SnakeMove(player, msg)
	default:
		fmt.Println("default messagehandler")
	}
}

func startGame() {
	gameOn = true
}

func setGameFinished() {
	gameFinished = true
}
