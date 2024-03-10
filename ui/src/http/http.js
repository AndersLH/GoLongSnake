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


// Authentication
// export const getExample = (params) => {return api.get('/somepath/', {params})} // Remember the curly braces for URL param-unfolding
export const changeName = (payload) => {return api.post('/start', payload)}

