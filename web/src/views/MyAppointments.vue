<template>
  <div class="container" style="padding-top:2rem">
    <h1 class="title">Meus Agendamentos</h1>

    <div v-if="loading" style="text-align:center;color:#aaa">Carregando...</div>

    <div v-for="a in appointments" :key="a.id" class="card">
      <div style="display:flex;justify-content:space-between;align-items:center">
        <div>
          <div style="font-weight:600">{{ formatDate(a.date) }} às {{ a.time?.slice(0,5) }}</div>
          <div v-if="a.service_name" style="font-size:0.85rem;color:#4ade80">{{ a.service_name }}</div>
          <div style="font-size:0.875rem;color:#aaa">📅 Agendado</div>
        </div>
        <button
          @click="handleCancel(a.id)"
          style="background:#e94560;color:#fff;border:none;padding:0.5rem 1rem;border-radius:6px;cursor:pointer;font-size:0.875rem"
        >
          Cancelar
        </button>
      </div>
    </div>

    <p v-if="!loading && !appointments.length" style="text-align:center;color:#aaa">
      Nenhum agendamento encontrado
    </p>

    <div class="spacer"></div>
    <button class="btn btn-secondary" @click="$router.push('/home')">Voltar</button>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useAuthStore } from '../stores/auth'
import { api } from '../api'

const auth = useAuthStore()
const appointments = ref([])
const loading = ref(true)

onMounted(loadAppointments)

async function loadAppointments() {
  loading.value = true
  try {
    const data = await api.getMyAppointments(auth.client.id)
    appointments.value = data || []
  } catch (e) {
    console.error(e)
  } finally {
    loading.value = false
  }
}

async function handleCancel(id) {
  if (!confirm('Cancelar este agendamento?')) return
  try {
    await api.cancelAppointment(id, auth.client.id)
    await loadAppointments()
  } catch (e) {
    alert(e.message)
  }
}

function formatDate(d) {
  if (!d) return ''
  const [y, m, day] = d.split('T')[0].split('-')
  return `${day}/${m}/${y}`
}
</script>
