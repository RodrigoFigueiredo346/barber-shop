<template>
  <div class="container" style="display:flex;flex-direction:column;justify-content:center;min-height:100vh">
    <h1 class="title">Entrar</h1>
    <input class="input" v-model="phone" placeholder="DDD + Número (ex: 11999998888)" maxlength="11" inputmode="numeric" />
    <p v-if="error" class="error">{{ error }}</p>
    <button class="btn btn-primary" @click="handleLogin" :disabled="loading">
      {{ loading ? 'Verificando...' : 'Entrar' }}
    </button>
    <div class="spacer"></div>
    <button class="btn btn-secondary" @click="$router.push('/register')">Não tenho cadastro</button>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '../stores/auth'
import { api } from '../api'

const phone = ref('')
const error = ref('')
const loading = ref(false)
const router = useRouter()
const auth = useAuthStore()

async function handleLogin() {
  error.value = ''
  const cleaned = phone.value.replace(/\D/g, '')
  if (cleaned.length < 10 || cleaned.length > 11) {
    error.value = 'Informe DDD + número (10 ou 11 dígitos)'
    return
  }
  loading.value = true
  try {
    const client = await api.login(cleaned)
    auth.setClient(client)
    router.push('/home')
  } catch (e) {
    error.value = 'Telefone não encontrado. Faça seu cadastro.'
  } finally {
    loading.value = false
  }
}
</script>
