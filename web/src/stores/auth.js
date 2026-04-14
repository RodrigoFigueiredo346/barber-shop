import { defineStore } from 'pinia'
import { ref, computed } from 'vue'

export const useAuthStore = defineStore('auth', () => {
    const client = ref(JSON.parse(localStorage.getItem('barber_client') || 'null'))

    const isLoggedIn = computed(() => !!client.value)

    function setClient(c) {
        client.value = c
        localStorage.setItem('barber_client', JSON.stringify(c))
    }

    function logout() {
        client.value = null
        localStorage.removeItem('barber_client')
    }

    return { client, isLoggedIn, setClient, logout }
})
