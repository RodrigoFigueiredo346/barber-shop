<template>
  <div class="container" style="padding-top:2rem">
    <h1 class="title">✂️ Olá, {{ auth.client?.name }}</h1>

    <button class="btn btn-primary" @click="goToBook" style="margin-bottom:1rem">Agendar horário</button>
    <button class="btn btn-secondary" @click="$router.push('/meus-horarios')" style="margin-bottom:1rem">Meus agendamentos</button>
    <button class="btn btn-secondary" @click="handleLogout" style="opacity:0.6;margin-bottom:2rem">Sair</button>

    <!-- Seleção de serviços (múltiplos) -->
    <p class="subtitle">Escolha os serviços</p>

    <div v-if="loading" style="text-align:center;color:#aaa">Carregando serviços...</div>

    <div v-if="!loading && !services.length" style="text-align:center;color:#aaa">
      Nenhum serviço disponível no momento
    </div>

    <div
      v-for="s in services"
      :key="s.id"
      class="service-card"
      :class="{ selected: isSelected(s.id) }"
      @click="toggleService(s)"
    >
      <div class="service-check">
        <div class="check-outer">
          <span v-if="isSelected(s.id)" class="check-inner">✓</span>
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

    <div v-if="selectedServices.length" style="margin-top:1rem;padding:0.75rem;background:#16213e;border-radius:8px;font-size:0.85rem">
      <div style="color:#aaa;margin-bottom:0.25rem">Selecionados: {{ selectedServices.length }} serviço(s)</div>
      <div style="color:#4ade80;font-weight:600">Total: R$ {{ totalPrice }} · {{ totalDuration }} min</div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '../stores/auth'
import { api } from '../api'

const router = useRouter()
const auth = useAuthStore()
const services = ref([])
const selectedServices = ref(JSON.parse(localStorage.getItem('barber_services') || '[]'))
const loading = ref(true)

onMounted(async () => {
  try {
    const data = await api.getServices()
    services.value = data || []
  } catch (e) { console.error(e) }
  finally { loading.value = false }
})

function isSelected(id) {
  return selectedServices.value.some(s => s.id === id)
}

function toggleService(s) {
  const idx = selectedServices.value.findIndex(x => x.id === s.id)
  if (idx >= 0) {
    selectedServices.value.splice(idx, 1)
  } else {
    selectedServices.value.push(s)
  }
  localStorage.setItem('barber_services', JSON.stringify(selectedServices.value))
}

const totalPrice = computed(() => {
  const cents = selectedServices.value.reduce((sum, s) => sum + s.price, 0)
  return (cents / 100).toFixed(2).replace('.', ',')
})

const totalDuration = computed(() => selectedServices.value.reduce((sum, s) => sum + s.duration, 0))

function goToBook() {
  if (!selectedServices.value.length) {
    alert('Selecione pelo menos um serviço')
    return
  }
  router.push('/agendar')
}

function handleLogout() {
  auth.logout()
  localStorage.removeItem('barber_services')
  router.push('/')
}
</script>

<style scoped>
.service-card {
  display:flex;align-items:center;gap:1rem;background:#16213e;border:2px solid transparent;
  border-radius:12px;padding:1rem;margin-bottom:0.75rem;cursor:pointer;transition:all 0.2s;
}
.service-card:hover { border-color:#e94560 }
.service-card.selected { border-color:#e94560;background:#1a2744 }
.service-check { flex-shrink:0 }
.check-outer {
  width:22px;height:22px;border-radius:4px;border:2px solid #555;
  display:flex;align-items:center;justify-content:center;transition:all 0.2s;
}
.service-card.selected .check-outer { border-color:#e94560;background:#e94560 }
.check-inner { color:#fff;font-size:0.8rem;font-weight:700 }
.service-info { flex:1 }
.service-name { font-weight:600;font-size:1rem;margin-bottom:0.25rem }
.service-details { display:flex;gap:0.75rem;font-size:0.85rem;color:#aaa }
.service-price { color:#4ade80;font-weight:600 }
</style>
