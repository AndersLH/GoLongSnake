import { createWebHistory, createRouter } from "vue-router";
import SnakeGame from "@/views/SnakeGame.vue";
import NotFound from "@/views/NotFound.vue"

const routes = [
    {
        path: "/",
        alias: ['/index.html'],
        name: "home",
        component: SnakeGame 
    },
    {
        path: "/:pathMatch(.*)*",
        name: "NotFound",
        component: NotFound
    }
];

const router = createRouter({
  history: createWebHistory(),
  routes,
});

export default router;
