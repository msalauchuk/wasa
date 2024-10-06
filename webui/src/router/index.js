import {
    createRouter,
    createWebHashHistory
} from 'vue-router'
import Feed from '../views/Feed.vue'
import Login from '../views/Login.vue'
import Search from '../views/Search.vue'
import ProfileView from '../views/Profile.vue'
import PageNotFoundView from '../views/PageNotFound.vue'
import SettingsView from '../views/Settings.vue'

const router = createRouter({
    history: createWebHashHistory(
        import.meta.env.BASE_URL),
    routes: [{
            path: '/',
            redirect: '/login'
        },
        {
            path: '/login',
            component: Login
        },
        {
            path: '/home',
            component: Feed
        },
        {
            path: '/search',
            component: Search
        },
        {
            path: '/users/:id',
            component: ProfileView

        },
        {
            path: '/users/:id/settings',
            component: SettingsView

        },
        {
            path: "/:catchAll(.*)",
            component: PageNotFoundView
        },
    ]
})

export default router
