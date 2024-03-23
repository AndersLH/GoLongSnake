package handlers

import (
	"dsproject/backend/structs"
	"fmt"
	"net/http"
)

// Move snake from player direction
func SnakeMove(player *structs.Player, msg *structs.ClientMsg) {
	ogX := player.X
	ogY := player.Y
	switch msg.MsgData {
	case "arrowup":
		player.Y -= 1
	case "arrowdown":
		player.Y += 1
	case "arrowright":
		player.X += 1
	case "arrowleft":
		player.X -= 1
	default:
		return
	}

	//Teleport on borders
	switch {
	case player.X >= GRIDX:
		player.X = 0
	case player.X < 0:
		player.X = GRIDX - 1
	}

	switch {
	case player.Y >= GRIDY:
		player.Y = 0
	case player.Y < 0:
		player.Y = GRIDY - 1
	}

	//Collision detection
	for _, p := range playerList {
		fmt.Println(len(playerList), p.X, p.Y, player.X, player.Y)
		if p.X == player.X && p.Y == player.Y && p.Id != player.Id {
			player.X = ogX
			player.Y = ogY
			return
		}
	}

	//New player move
	newMove := structs.PlayerMove{
		X:        player.X,
		Y:        player.Y,
		PlayerId: player.Id,
	}
	moveMsg := structs.ClientMsg{
		MsgType: "move",
		MsgData: newMove,
	}
	// //Send grid size to all players
	for _, p := range playerList {
		//Send grid size to player
		err := p.Conn.WriteJSON(moveMsg)
		if err != nil {
			//Deal with error
			continue
		}
	}
}

func cors(w http.ResponseWriter, r *http.Request) bool {
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, PUT, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:8081")

	return r.Method == http.MethodOptions // If method is Options, return true to allow logic in handler return early - prevents double work.
}
