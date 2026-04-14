import { createRouter, createWebHistory } from 'vue-router'
import { useAuthStore } from '../stores/auth'

const routes = [
    { path: '/', name: 'welcome', component: () => import('../views/Welcome.vue') },
    { path: '/login', name: 'login', component: () => import('../views/Login.vue') },
    { path: '/register', name: 'register', component: () => import('../views/Register.vue') },
    { path: '/home', name: 'home', component: () => import('../views/Home.vue'), meta: { auth: true } },
    { path: '/agendar', name: 'book', component: () => import('../views/Book.vue'), meta: { auth: true } },
    { path: '/meus-horarios', name: 'myAppointments', component: () => import('../views/MyAppointments.vue'), meta: { auth: true } },
    { path: '/admin/login', name: 'adminLogin', component: () => import('../views/admin/AdminLogin.vue') },
    { path: '/admin', name: 'admin', component: () => import('../views/admin/AdminPanel.vue'), meta: { admin: true } },
]

const router = createRouter({
    history: createWebHistory(),
    routes
})

router.beforeEach((to) => {
    const auth = useAuthStore()
    if (to.meta.auth && !auth.isLoggedIn) return { name: 'welcome' }
    if (to.meta.admin && !localStorage.getItem('barber_admin')) return { name: 'adminLogin' }
    if ((to.name === 'welcome' || to.name === 'login') && auth.isLoggedIn) return { name: 'home' }
})

export default router
