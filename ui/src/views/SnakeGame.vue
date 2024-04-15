<style>
body{
    overflow: hidden;
}

.container{
    height: fit-content;
    width: 97%;
    margin: 20px;
    padding: 10px;
}

.grid{
    width: fit-content;
    border: 10px solid black;
    border-radius: 10px;
    /* box-shadow: 0px 0px 10px rgba(0, 0, 0, 0.8); */
    animation: rainbow 20s linear infinite alternate;
}

.playertitle{
    text-align: center;
    animation: wiggle 5s linear infinite, flashRainbow 10s infinite;
    margin: 0px;
    font-weight: bold;
    font-size: 18px;
    text-shadow: 2px 2px 4px rgba(0, 0, 0, 0.5);
    overflow:hidden;
}

.gridCell{
    border: 1px solid rgba(237, 237, 237, 100);
    width: 45px;
    height: 45px;
    border-radius: 5px;
}

@keyframes wiggle {
    0% {
        transform: rotate(0deg);
    }
    25% {
        transform: rotate(30deg);
    }
    75%{
        transform: rotate(-30deg);
    }
    100%{
        transform: rotate(0deg);
    }
}

@keyframes flashRainbow {
    0% {
        color: violet;
    }
    15% {
        color: orange;
    }
    30% {
        color: yellow;
    }
    45% {
        color: green;
    }
    60% {
        color: blue;
    }
    85% {
        color: indigo;
    }
    100% {
        color: violet;
    }
}

@keyframes rainbow {
    0% {
        box-shadow: 0px 0px 100px 50px violet;
    }
    15% {
        box-shadow: 0px 0px 100px 50px orange;
    }
    30% {
        box-shadow: 0px 0px 100px 50px yellow;
    }
    45% {
        box-shadow: 0px 0px 100px 50px green;
    }
    60% {
        box-shadow: 0px 0px 100px 50px blue;
    }
    85% {
        box-shadow: 0px 0px 100px 50px indigo;
    }
    100% {
        box-shadow: 0px 0px 100px 50px violet;
    }
}

.snakeHeadRight{
    border-top-right-radius: 30px;
    border-bottom-right-radius: 30px;
    /* border-right: 2px solid black; */
    transition: all 0.25s ease;
}

.snakeHeadLeft{
    border-top-left-radius: 30px;
    border-bottom-left-radius: 30px;
    /* border-left: 2px solid black; */
    transition: all 0.25s ease;
}
.snakeHeadDown{
    border-bottom-left-radius: 30px;
    border-bottom-right-radius: 30px;
    /* border-bottom: 2px solid black; */
    transition: all 0.25s ease;
}
.snakeHeadUp{
    border-top-right-radius: 30px;
    border-top-left-radius: 30px;
    /* border-top: 2px solid black; */
    transition: all 0.25s ease;

}
</style>
<template>
    <v-container>
        <div class="playertitle">Player {{this.snake.id}}</div>
        <div>{{this.finishedmsg}}</div>
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
            canMove: true,
            interval: null,
            finishedmsg: "",
        
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
                const timeint = 275;
                switch (keyCode){
                    case "arrowup": clearInterval(this.interval); this.interval = setInterval(this.keyUp, timeint); break; 
                    case "arrowdown": clearInterval(this.interval); this.interval = setInterval(this.keyDown, timeint); break; 
                    case "arrowright": clearInterval(this.interval); this.interval = setInterval(this.keyRight, timeint); break; 
                    case "arrowleft": clearInterval(this.interval); this.interval = setInterval(this.keyLeft, timeint); break; 
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
                        this.grid[y][x] = " "
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
                        } else {
                            this.grid[y][x] = " "
                        }
                    }
                }
            },

            //Move snake
            playerMove(move){
                
                //Set snake to number
                this.grid[move.y][move.x] = "";//move.playerid

                document.getElementById(move.oldy+"-"+move.oldx).classList.remove("snakeHeadUp");
                document.getElementById(move.oldy+"-"+move.oldx).classList.remove("snakeHeadDown");
                document.getElementById(move.oldy+"-"+move.oldx).classList.remove("snakeHeadRight");
                document.getElementById(move.oldy+"-"+move.oldx).classList.remove("snakeHeadLeft");

                switch (move.dir){
                    case "u": document.getElementById(move.y+"-"+move.x).classList.add("snakeHeadUp"); break;
                    case "d": document.getElementById(move.y+"-"+move.x).classList.add("snakeHeadDown"); break;
                    case "r": document.getElementById(move.y+"-"+move.x).classList.add("snakeHeadRight"); break;
                    case "l": document.getElementById(move.y+"-"+move.x).classList.add("snakeHeadLeft"); break;

                }

                //Set player color for themselves
                if (move.playerid == this.snake.id){
                    //Color of snake
                    document.getElementById(move.y+"-"+move.x).style.backgroundColor = "cyan";
                }

                //Set color of enemies
                if (move.playerid != this.snake.id){
                    //Set color of enemy snakes, (currently loops which isnt great)
                    document.getElementById(move.y+"-"+move.x).style.backgroundColor = this.colors[move.playerid % 10];
                }

                
        
            },

            finishedGame(msg){
                this.finishedmsg = msg.message
                console.log(msg.message)
            },

            //=============================================================
            //=============================================================
            //                  Player -> Server
            //=============================================================
            //=============================================================

            //When player hits an arrow key
            keyUp(){
                this.snake.direction = "u"
                this.msg = {
                    msgtype: "move",
                    msgdata: "arrowup"
                }
                
                this.wscon.send(JSON.stringify(this.msg));
            },
            keyDown(){
                this.snake.direction = "d"
                this.msg = {
                    msgtype: "move",
                    msgdata: "arrowdown"
                }
                
                this.wscon.send(JSON.stringify(this.msg));
            },
            keyRight(){
                this.snake.direction = "r"
                this.msg = {
                    msgtype: "move",
                    msgdata: "arrowright"
                }
                
                this.wscon.send(JSON.stringify(this.msg));
            },
            keyLeft(){
                this.snake.direction = "l"
                this.msg = {
                    msgtype: "move",
                    msgdata: "arrowleft"
                }
                
                this.wscon.send(JSON.stringify(this.msg));
            },
        },
    }
</script>