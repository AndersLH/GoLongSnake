import axios from 'axios';

let api;

function startApi(){
    const host = process.env.NODE_ENV === 'development' ? "localhost:8080" : window.location.host;
    api = axios.create({
        withCredentials: true,
        baseURL: `http://${host}/`
    });

    api.interceptors.response.use((res) => {
        return res.data;
    });
}

startApi();

//APIs
// export const joinGame = (payload) => {return api.post('/joingame', payload)}
// export const moveSnake = (payload) => {return api.post('/movesnake', payload)}
export const getGridSize = (payload) => {return api.get('/getgridsize', {payload})}

