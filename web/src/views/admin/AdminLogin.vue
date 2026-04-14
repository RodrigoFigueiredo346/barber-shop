<template>
  <div class="container" style="display:flex;flex-direction:column;justify-content:center;min-height:100vh">
    <h1 class="title">🔒 Admin</h1>
    <input class="input" v-model="user" placeholder="Usuário" />
    <input class="input" v-model="password" placeholder="Senha" type="password" />
    <p v-if="error" class="error">{{ error }}</p>
    <button class="btn btn-primary" @click="handleLogin" :disabled="loading">
      {{ loading ? 'Verificando...' : 'Entrar' }}
    </button>
    <div class="spacer"></div>
    <button class="btn btn-secondary" @click="$router.push('/')">Voltar</button>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { api } from '../../api'

const user = ref('')
const password = ref('')
const error = ref('')
const loading = ref(false)
const router = useRouter()

async function handleLogin() {
  error.value = ''
  if (!user.value || !password.value) {
    error.value = 'Preencha usuário e senha'
    return
  }
  loading.value = true
  try {
    await api.adminLogin(user.value, password.value)
    localStorage.setItem('barber_admin', user.value + ':' + password.value)
    router.push('/admin')
  } catch (e) {
    error.value = 'Credenciais inválidas'
  } finally {
    loading.value = false
  }
}
</script>
