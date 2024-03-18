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
            
            <v-card style="width:fit-content">
                <v-icon
                    size="small"
                    @click="update('u')"
                    >
                    mdi-arrow-up
                </v-icon>
                <v-icon
                    size="small"
                    @click="update('d')"
                    >
                    mdi-arrow-down
                </v-icon>
                <v-icon
                    size="small"
                    @click="update('l')"
                    >
                    mdi-arrow-left
                </v-icon>
                <v-icon
                    size="small"
                    @click="update('r')"
                    >
                    mdi-arrow-right
                </v-icon>
            </v-card>
            <v-container class="grid">
                <v-row v-for="column, row in grid" :key="row">
                    <v-col v-for="cell, i in column" :key="i" class="gridCell">
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
            gridSize: 8, 
            grid: [[]],
            snake: {
                name: "user01",
                direction: "",
                x: 0,
                y: 0,
            },
        }),

        methods: {
            update: async function(dir) {
                this.snake.direction = dir
                changeName(this.snake)
                    .then((response)=>{
                        this.snake = response;
                    })
                    .catch(()=>{console.log("Something went wrong")})
                    .finally()
            },
        },
        created(){
            //Request board size from server
            for(let x = 0; x < this.gridSize; x++){
                this.grid[x] = []
                for(let y = 0; y < this.gridSize; y++){
                    this.grid[x][y] = "0"
                }
            }
        },
    }
</script>