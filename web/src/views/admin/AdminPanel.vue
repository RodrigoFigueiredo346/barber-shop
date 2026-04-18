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
      <p style="font-size:0.8rem;color:#aaa;margin-bottom:1rem">Use até 4 faixas por dia para configurar intervalos</p>
      <div v-for="(day, i) in dayNames" :key="i" class="card" style="margin-bottom:0.75rem">
        <div style="font-weight:600;margin-bottom:0.5rem">{{ day }}</div>
        <div v-for="s in 2" :key="s" style="display:flex;align-items:center;gap:0.5rem;margin-bottom:0.25rem">
          <input type="checkbox" v-model="formSchedules[schedKey(i,s)].active" />
          <select class="input" style="width:auto;margin:0;flex:1;padding:0.5rem" v-model="formSchedules[schedKey(i,s)].start_time" :disabled="!formSchedules[schedKey(i,s)].active">
            <option v-for="t in timeOptions" :key="'s'+i+s+t" :value="t">{{ t }}</option>
          </select>
          <span style="font-size:0.85rem">às</span>
          <select class="input" style="width:auto;margin:0;flex:1;padding:0.5rem" v-model="formSchedules[schedKey(i,s)].end_time" :disabled="!formSchedules[schedKey(i,s)].active">
            <option v-for="t in timeOptions" :key="'e'+i+s+t" :value="t">{{ t }}</option>
          </select>
        </div>
      </div>
      <div class="spacer"></div>
      <button class="btn btn-primary" @click="saveSchedules" :disabled="savingSchedules">
        {{ savingSchedules ? 'Salvando...' : 'Salvar horários' }}
      </button>
    </div>

    <!-- Agendamentos -->
    <div v-if="tab === 'appointments'">
      <p class="subtitle">Agendamentos do dia</p>
      <input class="input" type="date" v-model="apptDate" @change="loadAppointments" />
      <div v-if="loadingAppts" style="color:#aaa">Carregando...</div>
      <div v-for="a in appointments" :key="a.id" class="card card-compact" style="display:flex;justify-content:space-between;align-items:center">
        <div>
          <span style="font-weight:600">{{ a.time?.slice(0,5) }}</span>
          <span style="margin-left:0.5rem">{{ a.client_name }}</span>
          <span v-if="a.service_name" style="color:#4ade80;margin-left:0.5rem;font-size:0.8rem">{{ a.service_name }}</span>
        </div>
        <button @click="cancelAppt(a.id)" style="background:#e94560;color:#fff;border:none;padding:0.3rem 0.6rem;border-radius:6px;cursor:pointer;font-size:0.75rem">✕</button>
      </div>
      <p v-if="!loadingAppts && apptDate && !appointments.length" style="color:#aaa">Nenhum agendamento</p>
    </div>

    <!-- Bloquear horários -->
    <div v-if="tab === 'block'">
      <p class="subtitle">Bloquear/Desbloquear horários</p>
      <input class="input" type="date" v-model="blockDate" @change="loadBlockSlots" />
      <div style="display:flex;gap:1rem;margin-bottom:0.75rem;font-size:0.8rem;color:#aaa" v-if="blockDate">
        <span>🟢 Agendado</span>
        <span>🔴 Bloqueado</span>
        <span>⚪ Livre</span>
      </div>
      <div v-if="blockDate" class="slot-grid">
        <div
          v-for="slot in allDaySlots"
          :key="slot"
          class="slot"
          :class="{ 'slot-blocked': blockedSet.has(slot), 'slot-booked': bookedSet.has(slot) }"
          @click="!bookedSet.has(slot) && toggleBlock(slot)"
          :style="bookedSet.has(slot) ? 'cursor:default;opacity:0.9' : ''"
        >
          {{ slot }}
          <span v-if="bookedSet.has(slot)" style="font-size:0.7rem;display:block;color:#4ade80">✅</span>
          <span v-else-if="blockedSet.has(slot)" style="font-size:0.7rem;display:block">🔒</span>
          <span v-else style="font-size:0.7rem;display:block;color:#aaa">—</span>
        </div>
      </div>
    </div>

    <!-- Serviços -->
    <div v-if="tab === 'services'">
      <p class="subtitle">Gerenciar serviços</p>
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
          <button class="btn btn-primary" @click="saveService" style="flex:1">{{ editingService ? 'Atualizar' : 'Adicionar' }}</button>
          <button v-if="editingService" class="btn btn-secondary" @click="cancelEditSvc" style="flex:1">Cancelar</button>
        </div>
      </div>
      <div v-for="s in adminServices" :key="s.id" class="card card-compact" style="display:flex;justify-content:space-between;align-items:center">
        <div>
          <span style="font-weight:600">{{ s.name }}</span>
          <span style="font-size:0.85rem;color:#aaa;margin-left:0.5rem">{{ s.duration }}min</span>
          <span style="color:#4ade80;margin-left:0.5rem">R$ {{ (s.price / 100).toFixed(2).replace('.', ',') }}</span>
          <span v-if="!s.active" style="color:#e94560;margin-left:0.25rem">(inativo)</span>
        </div>
        <div style="display:flex;gap:0.5rem">
          <button @click="editService(s)" style="background:#0f3460;color:#eee;border:none;padding:0.3rem 0.5rem;border-radius:6px;cursor:pointer;font-size:0.8rem">✏️</button>
          <button @click="deleteService(s.id)" style="background:#e94560;color:#fff;border:none;padding:0.3rem 0.5rem;border-radius:6px;cursor:pointer;font-size:0.8rem">🗑️</button>
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
    </div>

    <!-- Dialog de sucesso -->
    <Teleport to="body">
      <div v-if="successMsg" class="dialog-overlay" @click="successMsg = ''">
        <div class="dialog-box" @click.stop>
          <div style="font-size:2rem;margin-bottom:0.5rem">✅</div>
          <p style="font-weight:600">{{ successMsg }}</p>
          <button class="btn btn-primary" @click="successMsg = ''" style="margin-top:1rem">OK</button>
        </div>
      </div>
    </Teleport>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted, computed } from 'vue'
import { useRouter } from 'vue-router'
import { api } from '../../api'

const router = useRouter()
const tab = ref('appointments')
const successMsg = ref('')
const tabs = [
  { id: 'appointments', label: '📋 Agendamentos' },
  { id: 'schedules', label: '📅 Horários' },
  { id: 'block', label: '🔒 Bloquear' },
  { id: 'services', label: '💈 Serviços' },
  { id: 'settings', label: '⚙️ Config' },
]
const dayNames = ['Dom', 'Seg', 'Ter', 'Qua', 'Qui', 'Sex', 'Sáb']

function showSuccess(msg) {
  successMsg.value = msg
}

// Schedules - até 4 faixas por dia
const schedules = reactive({})
const formSchedules = reactive({})

function schedKey(day, slot) { return `${day}_${slot}` }

// Inicializa formulário vazio
for (let d = 0; d < 7; d++) {
  for (let s = 1; s <= 2; s++) {
    formSchedules[schedKey(d, s)] = { day_of_week: d, slot: s, start_time: '09:00', end_time: '18:00', active: false }
  }
}

const timeOptions = []
for (let h = 0; h < 24; h++) {
  for (let m = 0; m < 60; m += 30) {
    timeOptions.push(`${String(h).padStart(2, '0')}:${String(m).padStart(2, '0')}`)
  }
}

const savingSchedules = ref(false)

async function loadSchedules() {
  try {
    const data = await api.getSchedules()
    if (data) {
      data.forEach(s => {
        const slot = s.slot || 1
        const key = schedKey(s.day_of_week, slot)
        schedules[key] = s
        formSchedules[key] = {
          day_of_week: s.day_of_week,
          slot,
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
  try {
    for (const key of Object.keys(formSchedules)) {
      const s = formSchedules[key]
      await api.upsertSchedule(s)
      schedules[key] = { ...s }
    }
    showSuccess('Horários salvos!')
  } catch (e) {
    alert('Erro ao salvar: ' + e.message)
  } finally {
    savingSchedules.value = false
  }
}

// Appointments
const apptDate = ref(new Date().toISOString().split('T')[0])
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
    showSuccess('Agendamento cancelado')
  } catch (e) { alert(e.message) }
}

// Block slots
const blockDate = ref('')
const blockedSlots = ref([])
const bookedSlots = ref([])
const allDaySlots = ref([])
const blockedSet = computed(() => new Set(blockedSlots.value.map(s => { const t = s.time || s; return t.length > 5 ? t.slice(0,5) : t })))
const bookedSet = computed(() => new Set(bookedSlots.value.map(s => s.length > 5 ? s.slice(0,5) : s)))

async function loadBlockSlots() {
  if (!blockDate.value) return
  try {
    const d = new Date(blockDate.value + 'T12:00:00')
    const dayOfWeek = d.getDay()
    const slots = []
    for (let s = 1; s <= 2; s++) {
      const sched = formSchedules[schedKey(dayOfWeek, s)]
      if (sched && sched.active) {
        let [sh, sm] = sched.start_time.split(':').map(Number)
        const [eh, em] = sched.end_time.split(':').map(Number)
        while (sh < eh || (sh === eh && sm < em)) {
          const t = `${String(sh).padStart(2,'0')}:${String(sm).padStart(2,'0')}`
          if (!slots.includes(t)) slots.push(t)
          sm += 30
          if (sm >= 60) { sh++; sm = 0 }
        }
      }
    }
    slots.sort()
    allDaySlots.value = slots

    const available = await api.getSlots(blockDate.value)
    const availableSet = new Set(available.slots || [])
    const booked = await api.getBookedSlots(blockDate.value)
    bookedSlots.value = booked || []
    const bookedNorm = new Set((booked || []).map(s => s.length > 5 ? s.slice(0,5) : s))
    const blocked = slots.filter(s => !availableSet.has(s) && !bookedNorm.has(s)).map(s => ({ time: s }))
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
      await api.updateService(editingService.value.id, { name: svcForm.name, duration: svcForm.duration, price: svcForm.price, active: svcForm.active })
    } else {
      await api.createService({ name: svcForm.name, duration: svcForm.duration, price: svcForm.price })
    }
    resetSvcForm()
    await loadServices()
    showSuccess(editingService.value ? 'Serviço atualizado!' : 'Serviço criado!')
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

function cancelEditSvc() { editingService.value = null; resetSvcForm() }

function resetSvcForm() {
  editingService.value = null
  svcForm.name = ''; svcForm.duration = 30; svcForm.price = 0; svcForm.priceDisplay = ''; svcForm.active = true
}

async function deleteService(id) {
  if (!confirm('Remover este serviço?')) return
  try { await api.deleteService(id); await loadServices(); showSuccess('Serviço removido') } catch (e) { alert(e.message) }
}

// Settings
const reminderMinutes = ref(60)

async function loadSettings() {
  try { const data = await api.getSettings(); reminderMinutes.value = data.reminder_minutes } catch (e) { console.error(e) }
}

async function saveSettings() {
  try { await api.updateSettings({ reminder_minutes: reminderMinutes.value }); showSuccess('Configurações salvas!') } catch (e) { alert(e.message) }
}

function handleLogout() { localStorage.removeItem('barber_admin'); router.push('/admin/login') }

onMounted(() => {
  loadSchedules()
  loadSettings()
  loadServices()
  loadAppointments()
})
</script>

<style scoped>
.tab-btn { background:#16213e;border:1px solid #333;color:#aaa;padding:0.5rem 0.75rem;border-radius:8px;cursor:pointer;font-size:0.85rem }
.tab-btn.active { background:#e94560;color:#fff;border-color:#e94560 }
.slot-blocked { background:#5c1a2a!important;border-color:#e94560!important }
.slot-booked { background:#1a3a2a!important;border-color:#4ade80!important }
.card-compact { padding:0.6rem 0.75rem }
.dialog-overlay { position:fixed;top:0;left:0;right:0;bottom:0;background:rgba(0,0,0,0.6);display:flex;align-items:center;justify-content:center;z-index:999 }
.dialog-box { background:#16213e;border-radius:16px;padding:2rem;text-align:center;max-width:320px;width:90% }
</style>
