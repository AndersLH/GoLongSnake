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
}
</style>
<template>
    <v-container>
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
    import {getGridSize} from '@/http/http';
    export default {

        name: 'App',
        components: {},
        data: () => ({
            gridSize: {
                x: 0,
                y: 0,
            }, 
            grid: [[]],
            snake: {
                username: "user01",
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
            }
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
            getSize: async function(){
                //
                getGridSize()
                    .then((response)=>{
                        this.gridSize = response;
                    })
                    .catch(()=>{console.log("Something went wrong")})
                    .finally(()=>{
                        //Initialize board
                        for(let y = 0; y < this.gridSize.y; y++){
                            this.grid[y] = []
                            for(let x = 0; x < this.gridSize.x; x++){
                                this.grid[y][x] = "0"
                            }
                        }
                    })
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
                    console.log("ERROR", e) 
                }

                // this.ws.send(JSON.stringify({message: this.newMessage}));
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
                    default:
                        // this.$refs.game.handleMessage(msg)
                        console.log("default handlemessage vue")
                }
            },

            initGrid(size){
                this.gridSize = size
                for(let y = 0; y < this.gridSize.y; y++){
                    this.grid[y] = []
                    for(let x = 0; x < this.gridSize.x; x++){
                        this.grid[y][x] = "-"
                    }
                }
            },
            //CURRENTLY CLONE
            updateGrid(size){
                this.gridSize = size
                for(let y = 0; y < this.gridSize.y; y++){
                    this.grid[y] = []
                    for(let x = 0; x < this.gridSize.x; x++){
                        this.grid[y][x] = "-"
                    }
                }
            },

            //Move snake
            playerMove(move){
                //Previous cell
                // document.getElementById(this.snake.y+"-"+this.snake.x).style.backgroundColor = "white";
                
                this.snake.x = move.x
                this.snake.y = move.y
                //Set snake to S
                this.grid[this.snake.y][this.snake.x] = move.playerid
                //Set current to red
                if(move.playerid %2 == 0){
                    document.getElementById(this.snake.y+"-"+this.snake.x).style.backgroundColor = "blue";
                } else {
                    document.getElementById(this.snake.y+"-"+this.snake.x).style.backgroundColor = "red";
                }
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