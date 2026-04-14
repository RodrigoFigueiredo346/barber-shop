<template>
  <div class="container" style="display:flex;flex-direction:column;justify-content:center;min-height:100vh">
    <h1 class="title">Cadastrar</h1>
    <input class="input" v-model="phone" placeholder="DDD + Número (ex: 11999998888)" maxlength="11" inputmode="numeric" />
    <input v-if="showName" class="input" v-model="name" placeholder="Seu nome" />
    <p v-if="error" class="error">{{ error }}</p>
    <button class="btn btn-primary" @click="handleRegister" :disabled="loading">
      {{ loading ? 'Aguarde...' : (showName ? 'Finalizar cadastro' : 'Continuar') }}
    </button>
    <div class="spacer"></div>
    <button class="btn btn-secondary" @click="$router.push('/login')">Já tenho cadastro</button>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '../stores/auth'
import { api } from '../api'

const phone = ref('')
const name = ref('')
const showName = ref(false)
const error = ref('')
const loading = ref(false)
const router = useRouter()
const auth = useAuthStore()

async function handleRegister() {
  error.value = ''
  const cleaned = phone.value.replace(/\D/g, '')
  if (cleaned.length < 10 || cleaned.length > 11) {
    error.value = 'Informe DDD + número (10 ou 11 dígitos)'
    return
  }

  loading.value = true
  try {
    if (!showName.value) {
      const { exists } = await api.checkPhone(cleaned)
      if (exists) {
        error.value = 'Telefone já cadastrado. Faça login.'
        return
      }
      showName.value = true
      return
    }

    if (!name.value.trim()) {
      error.value = 'Informe seu nome'
      return
    }

    const client = await api.register(name.value.trim(), cleaned)
    auth.setClient(client)
    router.push('/home')
  } catch (e) {
    error.value = e.message
  } finally {
    loading.value = false
  }
}
</script>
