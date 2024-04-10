package handlers

import (
	"dsproject/backend/structs"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

var grid [][]int

// Move snake from player direction
func SnakeMove(player *structs.Player, msg *structs.ClientMsg) {
	ogX := player.X
	ogY := player.Y
	switch msg.MsgData {
	case "arrowup":
		player.Y -= 1
		player.Dir = "u"
	case "arrowdown":
		player.Y += 1
		player.Dir = "d"
	case "arrowright":
		player.X += 1
		player.Dir = "r"
	case "arrowleft":
		player.X -= 1
		player.Dir = "l"
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
		OldX:     ogX,
		OldY:     ogY,
		PlayerId: player.Id,
		Dir:      player.Dir,
	}
	grid[player.Y][player.X] = player.Id
	moveMsg := structs.ClientMsg{
		MsgType: "move",
		MsgData: newMove,
	}

	//Print board
	for _, row := range getGrid() {
		fmt.Println(strings.Trim(fmt.Sprint(row), "[]"))
	}
	fmt.Println("-")

	//Send grid size to all players
	for _, p := range playerList {
		//Send grid size to player
		p.Lockmx.Lock()
		err := p.Conn.WriteJSON(moveMsg)
		p.Lockmx.Unlock()
		if err != nil {
			//Deal with error
			continue
		}
	}

	// ======= FINISH ======= //

	//Check for finish
	isFinished := true
	for x := range grid {
		for y := range grid[x] {
			if grid[x][y] == 0 {
				isFinished = false
			}
		}
	}

	//Send finish message to all
	if isFinished {

		//Make score slice
		countScore := make(map[int]int)
		for _, p := range playerList {
			countScore[p.Id] = 0
		}

		for x := range grid {
			for y := range grid[x] {
				countScore[grid[x][y]]++
			}
		}

		//Does not deal with ties
		var winner int
		maxScore := 0
		for p, s := range countScore {
			if s > maxScore {
				maxScore = s
				winner = p
			}
		}

		finishStruct := structs.FinishedMsg{
			Message: "Game over, the winner is " + strconv.Itoa(winner) + " with a score of " + strconv.Itoa(maxScore),
		}

		finishedMsg := structs.ClientMsg{
			MsgType: "finished",
			MsgData: finishStruct,
		}

		setGameFinished()

		for _, p := range playerList {
			//Send grid size to player
			err := p.Conn.WriteJSON(finishedMsg)
			if err != nil {
				//Deal with error
				continue
			}
		}
	}
}

// Create an empty grid
func createGrid() {
	emptyGrid := make([][]int, GRIDY)
	for x := range emptyGrid {
		emptyGrid[x] = make([]int, GRIDX)
		for y := range emptyGrid[x] {
			emptyGrid[x][y] = 0
		}
	}
	grid = emptyGrid
}

func getGrid() [][]int {
	return grid

}

func cors(w http.ResponseWriter, r *http.Request) bool {
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, PUT, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:8081")

	return r.Method == http.MethodOptions // If method is Options, return true to allow logic in handler return early - prevents double work.
}
