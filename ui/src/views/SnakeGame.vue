<style>
.container{
    height: 77vh;
    width: 97%;
    margin: 20px;
    padding: 10px;
    border: 1px solid black;
}
</style>
<template>
    <v-card>
        <div class="container">

            <v-btn color="green-lighten-3" @click="update">Change name</v-btn>
            <v-text-field
                    v-model="this.snake.name"
            />
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


        </div>
    </v-card>
</template>

<script>
    import {changeName} from '@/http/http';
    export default {

        name: 'App',
        components: {},
        data: () => ({
            snake: {
                name: "user01",
                direction: "",
                x: 1,
                y: 22,
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
        mounted: async function () {
            //Stuff that happens once 
        }
    }
</script>