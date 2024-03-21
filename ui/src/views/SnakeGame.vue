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
    import {changeName} from '@/http/http';
    export default {

        name: 'App',
        components: {},
        data: () => ({
            gridSize: {
                x: 16,
                y: 8,
            }, 
            grid: [[]],
            snake: {
                name: "user01",
                direction: "",
                x: 0,
                y: 0,
            },
            snakecss: "#ffffff",
            backgroundcss: "#000000",
        }),
        mounted: function () {
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
                        this.snake = response;
                        this.grid[this.snake.y][this.snake.x] = "S"
                        //Set current to red
                        document.getElementById(this.snake.y+"-"+this.snake.x).style.backgroundColor = "red";
                    })
                    .catch(()=>{console.log("Something went wrong")})
                    .finally()
            },
        },
        created(){
            //Request board size from server
            for(let y = 0; y < this.gridSize.y; y++){
                this.grid[y] = []
                for(let x = 0; x < this.gridSize.x; x++){
                    this.grid[y][x] = "0"
                }
            }
        },
    }
</script>