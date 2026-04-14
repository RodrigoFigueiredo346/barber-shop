<template>
  <div class="container" style="padding-top:2rem">
    <div style="display:flex;justify-content:space-between;align-items:center;margin-bottom:1.5rem">
      <h1 class="title" style="margin:0">⚙️ Painel Admin</h1>
      <button @click="handleLogout" style="background:none;border:1px solid #555;color:#aaa;padding:0.4rem 0.8rem;border-radius:6px;cursor:pointer;font-size:0.8rem">Sair</button>
    </div>

    <!-- Tabs -->
    <div style="display:flex;gap:0.5rem;margin-bottom:1.5rem;flex-wrap:wrap">
      <button v-for="t in tabs" :key="t.id" class="tab-btn" :class="{ active: tab === t.id }" @click="tab = t.id">
        {{ t.label }}
      </button>
    </div>

    <!-- Horários de funcionamento -->
    <div v-if="tab === 'schedules'">
      <p class="subtitle">Horários por dia da semana</p>
      <div v-for="(day, i) in dayNames" :key="i" class="card" style="display:flex;align-items:center;gap:0.5rem;flex-wrap:wrap">
        <label style="width:50px">
          <input type="checkbox" v-model="formSchedules[i].active" />
          {{ day }}
        </label>
        <select class="input" style="width:auto;margin:0;flex:1" v-model="formSchedules[i].start_time" :disabled="!formSchedules[i].active">
          <option v-for="t in timeOptions" :key="'s'+i+t" :value="t">{{ t }}</option>
        </select>
        <span>às</span>
        <select class="input" style="width:auto;margin:0;flex:1" v-model="formSchedules[i].end_time" :disabled="!formSchedules[i].active">
          <option v-for="t in timeOptions" :key="'e'+i+t" :value="t">{{ t }}</option>
        </select>
      </div>
      <div class="spacer"></div>
      <button class="btn btn-primary" @click="saveSchedules" :disabled="savingSchedules">
        {{ savingSchedules ? 'Salvando...' : 'Salvar horários' }}
      </button>
      <p v-if="schedulesMsg" style="color:#4ade80;text-align:center;margin-top:0.5rem">{{ schedulesMsg }}</p>
    </div>

    <!-- Agendamentos -->
    <div v-if="tab === 'appointments'">
      <p class="subtitle">Agendamentos do dia</p>
      <input class="input" type="date" v-model="apptDate" @change="loadAppointments" />
      <div v-if="loadingAppts" style="color:#aaa">Carregando...</div>
      <div v-for="a in appointments" :key="a.id" class="card" style="display:flex;justify-content:space-between;align-items:center">
        <div>
          <div style="font-weight:600">{{ a.time?.slice(0,5) }} - {{ a.client_name }}</div>
        </div>
        <button @click="cancelAppt(a.id)" style="background:#e94560;color:#fff;border:none;padding:0.4rem 0.8rem;border-radius:6px;cursor:pointer;font-size:0.8rem">Cancelar</button>
      </div>
      <p v-if="!loadingAppts && apptDate && !appointments.length" style="color:#aaa">Nenhum agendamento</p>
    </div>

    <!-- Bloquear horários -->
    <div v-if="tab === 'block'">
      <p class="subtitle">Bloquear/Desbloquear horários</p>
      <input class="input" type="date" v-model="blockDate" @change="loadBlockSlots" />
      <div v-if="blockDate" class="slot-grid">
        <div
          v-for="slot in allDaySlots"
          :key="slot"
          class="slot"
          :class="{ selected: blockedSet.has(slot) }"
          @click="toggleBlock(slot)"
        >
          {{ slot }}
          <span style="font-size:0.7rem;display:block">{{ blockedSet.has(slot) ? '🔒' : '✓' }}</span>
        </div>
      </div>
    </div>

    <!-- Serviços -->
    <div v-if="tab === 'services'">
      <p class="subtitle">Gerenciar serviços</p>

      <!-- Formulário novo serviço -->
      <div class="card" style="margin-bottom:1rem">
        <p style="font-weight:600;margin-bottom:0.5rem">{{ editingService ? 'Editar serviço' : 'Novo serviço' }}</p>
        <input class="input" v-model="svcForm.name" placeholder="Nome do serviço" />
        <div style="display:flex;gap:0.5rem">
          <div style="flex:1">
            <label style="color:#aaa;font-size:0.8rem">Duração (min)</label>
            <select class="input" v-model.number="svcForm.duration">
              <option v-for="d in durationOptions" :key="d" :value="d">{{ d }} min</option>
            </select>
          </div>
          <div style="flex:1">
            <label style="color:#aaa;font-size:0.8rem">Preço (R$)</label>
            <input class="input" v-model="svcForm.priceDisplay" placeholder="30,00" @blur="parsePriceInput" />
          </div>
        </div>
        <div style="display:flex;gap:0.5rem">
          <button class="btn btn-primary" @click="saveService" style="flex:1">
            {{ editingService ? 'Atualizar' : 'Adicionar' }}
          </button>
          <button v-if="editingService" class="btn btn-secondary" @click="cancelEdit" style="flex:1">Cancelar</button>
        </div>
      </div>

      <!-- Lista de serviços -->
      <div v-for="s in adminServices" :key="s.id" class="card" style="display:flex;justify-content:space-between;align-items:center">
        <div>
          <div style="font-weight:600">{{ s.name }}</div>
          <div style="font-size:0.85rem;color:#aaa">
            {{ s.duration }} min - <span style="color:#4ade80">R$ {{ (s.price / 100).toFixed(2).replace('.', ',') }}</span>
            <span v-if="!s.active" style="color:#e94560"> (inativo)</span>
          </div>
        </div>
        <div style="display:flex;gap:0.5rem">
          <button @click="editService(s)" style="background:#0f3460;color:#eee;border:none;padding:0.4rem 0.6rem;border-radius:6px;cursor:pointer;font-size:0.8rem">✏️</button>
          <button @click="deleteService(s.id)" style="background:#e94560;color:#fff;border:none;padding:0.4rem 0.6rem;border-radius:6px;cursor:pointer;font-size:0.8rem">🗑️</button>
        </div>
      </div>

      <p v-if="!adminServices.length" style="color:#aaa;text-align:center">Nenhum serviço cadastrado</p>
    </div>

    <!-- Configurações -->
    <div v-if="tab === 'settings'">
      <p class="subtitle">Configurações</p>
      <label style="color:#aaa;font-size:0.875rem">Lembrete (minutos antes)</label>
      <input class="input" type="number" v-model.number="reminderMinutes" min="5" max="1440" />
      <button class="btn btn-primary" @click="saveSettings">Salvar</button>
      <p v-if="settingsMsg" style="color:#4ade80;text-align:center;margin-top:0.5rem">{{ settingsMsg }}</p>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted, computed } from 'vue'
import { useRouter } from 'vue-router'
import { api } from '../../api'

const router = useRouter()
const tab = ref('schedules')
const tabs = [
  { id: 'schedules', label: '📅 Horários' },
  { id: 'appointments', label: '📋 Agendamentos' },
  { id: 'block', label: '🔒 Bloquear' },
  { id: 'services', label: '💈 Serviços' },
  { id: 'settings', label: '⚙️ Config' },
]
const dayNames = ['Dom', 'Seg', 'Ter', 'Qua', 'Qui', 'Sex', 'Sáb']

// Schedules
const schedules = reactive({})
const formSchedules = reactive(
  Array.from({ length: 7 }, (_, i) => ({ day_of_week: i, start_time: '09:00', end_time: '18:00', active: false }))
)
const savingSchedules = ref(false)
const schedulesMsg = ref('')

// Gera opções de meia em meia hora (00:00 até 23:30)
const timeOptions = []
for (let h = 0; h < 24; h++) {
  for (let m = 0; m < 60; m += 30) {
    timeOptions.push(`${String(h).padStart(2, '0')}:${String(m).padStart(2, '0')}`)
  }
}

async function loadSchedules() {
  try {
    const data = await api.getSchedules()
    if (data) {
      data.forEach(s => {
        schedules[s.day_of_week] = s
        formSchedules[s.day_of_week] = {
          day_of_week: s.day_of_week,
          start_time: s.start_time?.slice(0, 5) || '09:00',
          end_time: s.end_time?.slice(0, 5) || '18:00',
          active: s.active
        }
      })
    }
  } catch (e) { console.error(e) }
}

async function saveSchedules() {
  savingSchedules.value = true
  schedulesMsg.value = ''
  try {
    for (const s of formSchedules) {
      await api.upsertSchedule(s)
      schedules[s.day_of_week] = { ...s }
    }
    schedulesMsg.value = 'Horários salvos!'
    setTimeout(() => schedulesMsg.value = '', 2000)
  } catch (e) {
    alert('Erro ao salvar: ' + e.message)
  } finally {
    savingSchedules.value = false
  }
}

// Appointments
const apptDate = ref('')
const appointments = ref([])
const loadingAppts = ref(false)

async function loadAppointments() {
  if (!apptDate.value) return
  loadingAppts.value = true
  try {
    const data = await api.getAdminAppointments(apptDate.value)
    appointments.value = data || []
  } catch (e) { console.error(e) }
  finally { loadingAppts.value = false }
}

async function cancelAppt(id) {
  if (!confirm('Cancelar este agendamento?')) return
  try {
    await api.adminCancelAppointment(id)
    await loadAppointments()
  } catch (e) { alert(e.message) }
}

// Block slots
const blockDate = ref('')
const blockedSlots = ref([])
const allDaySlots = ref([])
const blockedSet = computed(() => new Set(blockedSlots.value.map(s => s.time?.slice(0,5))))

async function loadBlockSlots() {
  if (!blockDate.value) return
  try {
    // Gera todos os slots do dia baseado no schedule
    const d = new Date(blockDate.value + 'T12:00:00')
    const dayOfWeek = d.getDay()
    const sched = schedules[dayOfWeek]
    const slots = []
    if (sched && sched.active) {
      let [sh, sm] = sched.start_time.split(':').map(Number)
      const [eh, em] = sched.end_time.split(':').map(Number)
      while (sh < eh || (sh === eh && sm < em)) {
        slots.push(`${String(sh).padStart(2,'0')}:${String(sm).padStart(2,'0')}`)
        sm += 30
        if (sm >= 60) { sh++; sm = 0 }
      }
    }
    allDaySlots.value = slots

    // Busca bloqueados via API (usa a rota de slots pra comparar)
    const available = await api.getSlots(blockDate.value)
    const availableSet = new Set(available.slots || [])
    // Bloqueados = todos os slots do dia que NÃO estão disponíveis e NÃO estão agendados
    // Simplificação: busca os blocked_slots diretamente
    // Como não temos rota pública pra blocked_slots, vamos inferir
    const blocked = slots.filter(s => !availableSet.has(s)).map(s => ({ time: s + ':00' }))
    blockedSlots.value = blocked
  } catch (e) { console.error(e) }
}

async function toggleBlock(slot) {
  try {
    if (blockedSet.value.has(slot)) {
      await api.unblockSlot(blockDate.value, slot)
    } else {
      await api.blockSlot(blockDate.value, slot)
    }
    await loadBlockSlots()
  } catch (e) { alert(e.message) }
}

// Services
const adminServices = ref([])
const editingService = ref(null)
const svcForm = reactive({ name: '', duration: 30, price: 0, priceDisplay: '', active: true })
const durationOptions = [5, 10, 15, 20, 30, 45, 60, 90, 120]

async function loadServices() {
  try {
    const data = await api.getAdminServices()
    adminServices.value = data || []
  } catch (e) { console.error(e) }
}

function parsePriceInput() {
  const cleaned = svcForm.priceDisplay.replace(/[^\d,\.]/g, '').replace(',', '.')
  const val = parseFloat(cleaned)
  svcForm.price = isNaN(val) ? 0 : Math.round(val * 100)
  svcForm.priceDisplay = isNaN(val) ? '' : val.toFixed(2).replace('.', ',')
}

async function saveService() {
  parsePriceInput()
  if (!svcForm.name.trim()) { alert('Nome é obrigatório'); return }
  try {
    if (editingService.value) {
      await api.updateService(editingService.value.id, {
        name: svcForm.name, duration: svcForm.duration, price: svcForm.price, active: svcForm.active
      })
    } else {
      await api.createService({ name: svcForm.name, duration: svcForm.duration, price: svcForm.price })
    }
    resetSvcForm()
    await loadServices()
  } catch (e) { alert(e.message) }
}

function editService(s) {
  editingService.value = s
  svcForm.name = s.name
  svcForm.duration = s.duration
  svcForm.price = s.price
  svcForm.priceDisplay = (s.price / 100).toFixed(2).replace('.', ',')
  svcForm.active = s.active
}

function cancelEdit() {
  editingService.value = null
  resetSvcForm()
}

function resetSvcForm() {
  editingService.value = null
  svcForm.name = ''
  svcForm.duration = 30
  svcForm.price = 0
  svcForm.priceDisplay = ''
  svcForm.active = true
}

async function deleteService(id) {
  if (!confirm('Remover este serviço?')) return
  try {
    await api.deleteService(id)
    await loadServices()
  } catch (e) { alert(e.message) }
}

// Settings
const reminderMinutes = ref(60)
const settingsMsg = ref('')

async function loadSettings() {
  try {
    const data = await api.getSettings()
    reminderMinutes.value = data.reminder_minutes
  } catch (e) { console.error(e) }
}

async function saveSettings() {
  settingsMsg.value = ''
  try {
    await api.updateSettings({ reminder_minutes: reminderMinutes.value })
    settingsMsg.value = 'Salvo!'
    setTimeout(() => settingsMsg.value = '', 2000)
  } catch (e) { alert(e.message) }
}

function handleLogout() {
  localStorage.removeItem('barber_admin')
  router.push('/admin/login')
}

onMounted(() => {
  loadSchedules()
  loadSettings()
  loadServices()
})
</script>

<style scoped>
.tab-btn {
  background: #16213e;
  border: 1px solid #333;
  color: #aaa;
  padding: 0.5rem 0.75rem;
  border-radius: 8px;
  cursor: pointer;
  font-size: 0.85rem;
}
.tab-btn.active {
  background: #e94560;
  color: #fff;
  border-color: #e94560;
}
</style>
