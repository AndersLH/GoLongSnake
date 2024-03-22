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

            <v-btn color="green-lighten-3" @click="update">Change name</v-btn>
            <v-text-field
                    v-model="this.snake.name"
            />
            
                
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
    import {changeName, initBoard} from '@/http/http';
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
        }),
        mounted: function () {
            
            this.joinGame()
            this.getSize()

            
            // add an event listener for keypress
            window.addEventListener('keydown', this.handleKeyPress)
        },

        methods: {
            //Read keypress by user
            handleKeyPress: function (e) {
                const keyCode = String(e.key).toLowerCase();
                if (keyCode == "arrowup" || keyCode == "arrowdown" || keyCode == "arrowleft" || keyCode == "arrowright"){
                    this.update(keyCode)
                }
            },
            //Send update request to backend
            update: async function(dir) {
                this.snake.direction = dir
                changeName(this.snake)
                    //HTTP response
                    .then((response)=>{
                        //Set previous to white
                        document.getElementById(this.snake.y+"-"+this.snake.x).style.backgroundColor = "white";

                        //Set snake to S
                        this.snake.x = response.x;
                        this.snake.y = response.y;
                        this.grid[this.snake.y][this.snake.x] = "S"
                        //Set current to red
                        document.getElementById(this.snake.y+"-"+this.snake.x).style.backgroundColor = "red";
                    })
                    .catch(()=>{console.log("Something went wrong")})
                    .finally()
            },
            getSize: async function(){
                //
                initBoard()
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
            joinGame(){

                const host = process.env.NODE_ENV === 'development' ? "localhost:8080" : window.location.host;
                const url = `ws://${host}/ws`

                this.wscon = new WebSocket( url );
                //Console log for ws connection
                // this.ws.addEventListener('open', (event) => { console.log("Connected :D") });

                this.wscon.onclose = function() {
                    console.log("disconnected from ws!")
                }

                this.wscon.onmessage = function(e){
                    console.log(e.data)
                    this.handleMessage(e.data)
                }.bind(this);

                this.wscon.onerror = function(e){
                    console.log("ERROR", e)
                }
            },
            handleMessage(msg){
                try{ msg = JSON.parse(msg) }catch(e){ return }

                switch(msg.msgtype){
                    case "initplayer": this.snake = msg.msgdata; break;
                    case "move": this.updateOther(msg.msgdata); break;
                    case "hola": console.log("hype"); break;
                    default:
                        // this.$refs.game.handleMessage(msg)
                        console.log("def")
                }
            },
            updateOther(data){
                if (this.snake.id != data.other) {
                    console.log("other")
                }
            }
        },
        // created(){
        //     //Request board size 
        //     // this.joinGame()
        //     this.getSize()        
        // },
    }
</script>