import {
    createRouter,
    createWebHashHistory
} from 'vue-router'
import HomeView from '../views/Feed.vue'
import LoginView from '../views/Login.vue'
import SearchView from '../views/Search.vue'
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
            component: LoginView
        },
        {
            path: '/home',
            component: HomeView
        },
        {
            path: '/search',
            component: SearchView
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
