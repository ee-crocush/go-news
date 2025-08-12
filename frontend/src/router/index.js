import { createRouter, createWebHistory } from 'vue-router'
import News from '../views/News.vue'
import NewsDetail from '../views/NewsDetail.vue'

const routes = [
    {
        path: '/',
        name: 'News',
        component: News
    },
    {
        path: '/news/:id',
        name: 'NewsDetail',
        component: NewsDetail,
        props: true
    }
]

const router = createRouter({
    history: createWebHistory(import.meta.env.BASE_URL),
    routes
})

export default router