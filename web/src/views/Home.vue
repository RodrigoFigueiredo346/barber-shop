<template>
  <div class="container" style="padding-top:2rem">
    <h1 class="title">✂️ Olá, {{ auth.client?.name }}</h1>

    <button class="btn btn-primary" @click="goToBook" style="margin-bottom:1rem">Agendar horário</button>
    <button class="btn btn-secondary" @click="$router.push('/meus-horarios')" style="margin-bottom:1rem">Meus agendamentos</button>
    <button class="btn btn-secondary" @click="handleLogout" style="opacity:0.6;margin-bottom:2rem">Sair</button>

    <!-- Seleção de serviço -->
    <p class="subtitle">Escolha o serviço</p>

    <div v-if="loading" style="text-align:center;color:#aaa">Carregando serviços...</div>

    <div v-if="!loading && !services.length" style="text-align:center;color:#aaa">
      Nenhum serviço disponível no momento
    </div>

    <div
      v-for="s in services"
      :key="s.id"
      class="service-card"
      :class="{ selected: selectedService?.id === s.id }"
      @click="selectService(s)"
    >
      <div class="service-radio">
        <div class="radio-outer">
          <div v-if="selectedService?.id === s.id" class="radio-inner"></div>
        </div>
      </div>
      <div class="service-info">
        <div class="service-name">{{ s.name }}</div>
        <div class="service-details">
          <span>{{ s.duration }} min</span>
          <span class="service-price">R$ {{ (s.price / 100).toFixed(2).replace('.', ',') }}</span>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '../stores/auth'
import { api } from '../api'

const router = useRouter()
const auth = useAuthStore()
const services = ref([])
const selectedService = ref(JSON.parse(localStorage.getItem('barber_service') || 'null'))
const loading = ref(true)

onMounted(async () => {
  try {
    const data = await api.getServices()
    services.value = data || []
  } catch (e) {
    console.error(e)
  } finally {
    loading.value = false
  }
})

function selectService(s) {
  selectedService.value = s
  localStorage.setItem('barber_service', JSON.stringify(s))
}

function goToBook() {
  if (!selectedService.value) {
    alert('Selecione um serviço primeiro')
    return
  }
  router.push('/agendar')
}

function handleLogout() {
  auth.logout()
  localStorage.removeItem('barber_service')
  router.push('/')
}
</script>

<style scoped>
.service-card {
  display: flex;
  align-items: center;
  gap: 1rem;
  background: #16213e;
  border: 2px solid transparent;
  border-radius: 12px;
  padding: 1rem;
  margin-bottom: 0.75rem;
  cursor: pointer;
  transition: all 0.2s;
}

.service-card:hover {
  border-color: #e94560;
}

.service-card.selected {
  border-color: #e94560;
  background: #1a2744;
}

.service-radio {
  flex-shrink: 0;
}

.radio-outer {
  width: 22px;
  height: 22px;
  border-radius: 50%;
  border: 2px solid #555;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: border-color 0.2s;
}

.service-card.selected .radio-outer {
  border-color: #e94560;
}

.radio-inner {
  width: 12px;
  height: 12px;
  border-radius: 50%;
  background: #e94560;
}

.service-info {
  flex: 1;
}

.service-name {
  font-weight: 600;
  font-size: 1rem;
  margin-bottom: 0.25rem;
}

.service-details {
  display: flex;
  gap: 0.75rem;
  font-size: 0.85rem;
  color: #aaa;
}

.service-price {
  color: #4ade80;
  font-weight: 600;
}
</style>
