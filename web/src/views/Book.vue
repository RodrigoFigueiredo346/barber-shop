<template>
  <div class="container" style="padding-top:2rem">
    <h1 class="title">Agendar</h1>

    <!-- Calendário -->
    <div class="card">
      <div style="display:flex;justify-content:space-between;align-items:center;margin-bottom:1rem">
        <button class="nav-btn" @click="prevMonth">&lt;</button>
        <span style="font-weight:600;text-transform:capitalize">{{ monthLabel }}</span>
        <button class="nav-btn" @click="nextMonth">&gt;</button>
      </div>

      <div class="weekdays">
        <span v-for="d in weekdays" :key="d">{{ d }}</span>
      </div>

      <div class="calendar-grid">
        <div v-for="blank in firstDayOffset" :key="'b'+blank" class="day-cell empty"></div>
        <div
          v-for="day in daysInMonth"
          :key="day"
          class="day-cell"
          :class="{
            selected: selectedDay === day,
            today: isToday(day),
            past: isPast(day)
          }"
          @click="!isPast(day) && selectDay(day)"
        >
          {{ day }}
        </div>
      </div>
    </div>

    <!-- Horários -->
    <div v-if="selectedDay">
      <p class="subtitle">Horários em {{ selectedDay }}/{{ String(currentMonth+1).padStart(2,'0') }}</p>

      <div v-if="loadingSlots" style="text-align:center;color:#aaa">Carregando...</div>

      <div v-if="slots.length" class="slot-grid">
        <div
          v-for="slot in slots"
          :key="slot"
          class="slot"
          :class="{ selected: selectedSlot === slot }"
          @click="selectedSlot = slot"
        >
          {{ slot }}
        </div>
      </div>

      <p v-if="!loadingSlots && !slots.length" style="text-align:center;color:#aaa">
        Nenhum horário disponível
      </p>
    </div>

    <p v-if="error" class="error">{{ error }}</p>
    <p v-if="success" style="color:#4ade80;text-align:center">{{ success }}</p>

    <button
      v-if="selectedSlot"
      class="btn btn-primary"
      @click="handleBook"
      :disabled="booking"
      style="margin-top:1rem"
    >
      {{ booking ? 'Agendando...' : `Confirmar ${selectedSlot}` }}
    </button>

    <div class="spacer"></div>
    <button class="btn btn-secondary" @click="$router.push('/home')">Voltar</button>
  </div>
</template>

<script setup>
import { ref, computed } from 'vue'
import { useAuthStore } from '../stores/auth'
import { api } from '../api'

const auth = useAuthStore()
const now = new Date()
const currentYear = ref(now.getFullYear())
const currentMonth = ref(now.getMonth())
const selectedDay = ref(null)
const selectedSlot = ref('')
const slots = ref([])
const loadingSlots = ref(false)
const booking = ref(false)
const error = ref('')
const success = ref('')

const weekdays = ['Dom', 'Seg', 'Ter', 'Qua', 'Qui', 'Sex', 'Sáb']

const monthLabel = computed(() => {
  const d = new Date(currentYear.value, currentMonth.value)
  return d.toLocaleDateString('pt-BR', { month: 'long', year: 'numeric' })
})

const daysInMonth = computed(() => new Date(currentYear.value, currentMonth.value + 1, 0).getDate())
const firstDayOffset = computed(() => new Date(currentYear.value, currentMonth.value, 1).getDay())

function isToday(day) {
  return day === now.getDate() && currentMonth.value === now.getMonth() && currentYear.value === now.getFullYear()
}

function isPast(day) {
  const d = new Date(currentYear.value, currentMonth.value, day)
  const today = new Date(now.getFullYear(), now.getMonth(), now.getDate())
  return d < today
}

function prevMonth() {
  if (currentMonth.value === 0) {
    currentMonth.value = 11
    currentYear.value--
  } else {
    currentMonth.value--
  }
  selectedDay.value = null
  slots.value = []
  selectedSlot.value = ''
}

function nextMonth() {
  if (currentMonth.value === 11) {
    currentMonth.value = 0
    currentYear.value++
  } else {
    currentMonth.value++
  }
  selectedDay.value = null
  slots.value = []
  selectedSlot.value = ''
}

async function selectDay(day) {
  selectedDay.value = day
  selectedSlot.value = ''
  error.value = ''
  success.value = ''
  loadingSlots.value = true
  const date = `${currentYear.value}-${String(currentMonth.value+1).padStart(2,'0')}-${String(day).padStart(2,'0')}`
  try {
    const data = await api.getSlots(date)
    slots.value = data.slots || []
  } catch (e) {
    error.value = e.message
  } finally {
    loadingSlots.value = false
  }
}

async function handleBook() {
  error.value = ''
  success.value = ''
  booking.value = true
  const date = `${currentYear.value}-${String(currentMonth.value+1).padStart(2,'0')}-${String(selectedDay.value).padStart(2,'0')}`
  try {
    const service = JSON.parse(localStorage.getItem('barber_service') || 'null')
    await api.createAppointment(auth.client.id, service?.id || null, date, selectedSlot.value)
    success.value = `Agendado para ${selectedDay.value}/${String(currentMonth.value+1).padStart(2,'0')} às ${selectedSlot.value}`
    selectedSlot.value = ''
    await selectDay(selectedDay.value)
  } catch (e) {
    error.value = e.message
  } finally {
    booking.value = false
  }
}
</script>

<style scoped>
.weekdays {
  display: grid;
  grid-template-columns: repeat(7, 1fr);
  text-align: center;
  font-size: 0.75rem;
  color: #aaa;
  margin-bottom: 0.5rem;
}
.calendar-grid {
  display: grid;
  grid-template-columns: repeat(7, 1fr);
  gap: 4px;
}
.day-cell {
  aspect-ratio: 1;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 8px;
  cursor: pointer;
  font-size: 0.9rem;
  transition: all 0.15s;
}
.day-cell:not(.empty):not(.past):hover { background: #0f3460; }
.day-cell.selected { background: #e94560; font-weight: 700; }
.day-cell.today { border: 1px solid #e94560; }
.day-cell.past { color: #555; cursor: default; }
.day-cell.empty { cursor: default; }
.nav-btn {
  background: none;
  border: 1px solid #555;
  color: #eee;
  width: 36px;
  height: 36px;
  border-radius: 8px;
  cursor: pointer;
  font-size: 1rem;
}
.nav-btn:hover { border-color: #e94560; }
</style>
