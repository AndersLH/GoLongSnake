<style>
.container{
    height: fit-content;
    width: 97%;
    margin: 20px;
    padding: 10px;
    border: 1px solid black;
}

.grid{
    width: fit-content;
}

.gridCell{
    border: 1px solid black;
    width: 40px;
    height: 40px;
    /* border-radius: 10px; */
}
</style>
<template>
    <v-container>
        <div>{{this.snake.id}}</div>
        <div class="container">
            <v-container class="grid">
                <v-row v-for="column, row in grid" :key="row">
                    <v-col v-for="cell, i in column" :key="i" class="gridCell" :id="[row + '-' + i]">
                        <div>{{cell}}</div>
                    </v-col>
                </v-row>
            </v-container>

        </div>
    </v-container>
</template>

<script>
    export default {

        name: 'App',
        components: {},
        data: () => ({
            gridSize: {
                x: 0,
                y: 0,
            },
            playerColors: [],
            grid: [[]],
            colors: ["red", "blue", "green", "yellow", "orange", "purple", "pink", "saddlebrown", "lime", "grey"],
            snake: {
                username: "user01",
                // color: [Math.floor(Math.random() * 255), Math.floor(Math.random() * 255), Math.floor(Math.random() * 255)],
                id: null,
                direction: "",
                x: 0,
                y: 0,
                other: null,
            },

            connected: false,
            wscon: WebSocket.prototype,
            msg: {
                msgdata: null,
                masgtype: null,
            },
            lostws: false,
        }),
        mounted: function () {
            this.joinGame()
            
            // add an event listener for keypress
            window.addEventListener('keydown', this.handleKeyPress)
        },
        methods: {
            //Read keypress by user
            handleKeyPress: function (e) {
                const keyCode = String(e.key).toLowerCase();
                if (keyCode == "arrowup" || keyCode == "arrowdown" || keyCode == "arrowleft" || keyCode == "arrowright"){
                    this.keyPress(keyCode)
                }
            },

            //Initiate websocket connection with server
            joinGame(){
                const host = process.env.NODE_ENV === 'development' ? "localhost:8080" : window.location.host;
                const url = `ws://${host}/newplayer`

                this.wscon = new WebSocket( url );

                this.wscon.onclose = function() {
                    console.log("disconnected from ws!")
                }

                //Handle incoming message from server
                this.wscon.onmessage = function(e){
                    this.handleMessage(e.data)
                }.bind(this);

                this.wscon.onerror = function(e){
                    this.lostws = true
                    console.log("ERROR", e) 
                }
            },
            
            //=============================================================
            //=============================================================
            //                  Server -> Player
            //=============================================================
            //=============================================================
            
            //Handle messages from backend websocket connection
            handleMessage(msg){
                try{ msg = JSON.parse(msg) }catch(e){ return }

                switch(msg.msgtype){
                    case "move": this.playerMove(msg.msgdata); break;
                    case "initgrid": this.initGrid(msg.msgdata); break;
                    case "updategrid": this.updateGrid(msg.msgdata); break;
                    case "finished": this.finishedGame(msg.msgdata); break;
                    default:
                        console.log("default handlemessage vue")
                }
            },

            initGrid(size){
                this.gridSize = size
                this.snake.id = size.playerid
                for(let y = 0; y < this.gridSize.y; y++){
                    this.grid[y] = []
                    for(let x = 0; x < this.gridSize.x; x++){
                        this.grid[y][x] = "-"
                    }
                }
            },
            //CURRENTLY CLONE
            updateGrid(upGrid){

                this.snake.id = upGrid.playerid

                //Create updated grid for joining player
                for(let y = 0; y < upGrid.y; y++){
                    this.grid[y] = []
                    for(let x = 0; x < upGrid.x; x++){
                        if(upGrid.grid[y][x] != 0){
                            this.grid[y][x] = upGrid.grid[y][x]
                            //Add color
                            // document.getElementById(y+"-"+x).style.backgroundColor = "red" //this.colors[move.playerid % 10];
                        } else {
                            this.grid[y][x] = "-"
                        }
                    }
                }
            },

            //Move snake
            playerMove(move){
                //Reset tail
                // document.getElementById(this.snake.y+"-"+this.snake.x).style.backgroundColor = "white";
                // document.getElementById(this.snake.y+"-"+this.snake.x).innerHTML = "-"
                
                this.snake.x = move.x
                this.snake.y = move.y
                //Set snake to number
                this.grid[this.snake.y][this.snake.x] = move.playerid

                //Set player color for themselves
                if (move.playerid == this.snake.id){
                    document.getElementById(this.snake.y+"-"+this.snake.x).style.backgroundColor = "cyan";
                }

                //Set color of enemies
                if (move.playerid != this.snake.id){
                    //Set color of enemies, (currently loops which isnt great)
                    document.getElementById(move.y+"-"+move.x).style.backgroundColor = this.colors[move.playerid % 10];
                }
        
            },

            finishedGame(msg){
                console.log(msg.message)
            },

            //=============================================================
            //=============================================================
            //                  Player -> Server
            //=============================================================
            //=============================================================

            //When player hits an arrow key
            keyPress(dir){
                this.msg = {
                    msgtype: "move",
                    msgdata: dir
                }
                this.wscon.send(JSON.stringify(this.msg));
            }
        },
    }
</script>