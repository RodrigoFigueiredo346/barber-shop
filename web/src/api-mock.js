// Mock API usando localStorage — substitui o backend Go pra demo

function getStore(key, fallback = []) {
    return JSON.parse(localStorage.getItem('mock_' + key) || JSON.stringify(fallback))
}
function setStore(key, data) {
    localStorage.setItem('mock_' + key, JSON.stringify(data))
}
function nextId(key) {
    const items = getStore(key)
    return items.length ? Math.max(...items.map(i => i.id)) + 1 : 1
}

// Seed serviços padrão se não existirem
if (!localStorage.getItem('mock_services')) {
    setStore('services', [
        { id: 1, name: 'Corte padrão', duration: 30, price: 3000, active: true },
        { id: 2, name: 'Corte degradê', duration: 30, price: 4000, active: true },
        { id: 3, name: 'Barba', duration: 30, price: 2000, active: true },
        { id: 4, name: 'Sobrancelha', duration: 5, price: 500, active: true },
        { id: 5, name: 'Cabelo e barba', duration: 45, price: 6000, active: true },
        { id: 6, name: 'Degradê e sobrancelha', duration: 30, price: 4500, active: true },
    ])
}

// Seed schedules padrão
if (!localStorage.getItem('mock_schedules')) {
    const scheds = []
    // Ter a Sáb, 2 faixas: manhã e tarde
    for (let d = 2; d <= 6; d++) {
        scheds.push({ id: scheds.length + 1, day_of_week: d, slot: 1, start_time: '09:00', end_time: '12:00', active: true })
        scheds.push({ id: scheds.length + 1, day_of_week: d, slot: 2, start_time: '13:00', end_time: '18:00', active: true })
    }
    setStore('schedules', scheds)
}

if (!localStorage.getItem('mock_settings')) {
    setStore('settings', { reminder_minutes: 60 })
}

function generateSlots(date) {
    const d = new Date(date + 'T12:00:00')
    const dow = d.getDay()
    const scheds = getStore('schedules').filter(s => s.day_of_week === dow && s.active)
    const slots = []
    for (const sched of scheds) {
        let [h, m] = sched.start_time.split(':').map(Number)
        const [eh, em] = sched.end_time.split(':').map(Number)
        while (h < eh || (h === eh && m < em)) {
            slots.push(`${String(h).padStart(2, '0')}:${String(m).padStart(2, '0')}`)
            m += 30; if (m >= 60) { h++; m = 0 }
        }
    }
    // Filtra passados se hoje
    const now = new Date()
    const today = `${now.getFullYear()}-${String(now.getMonth() + 1).padStart(2, '0')}-${String(now.getDate()).padStart(2, '0')}`
    const nowMin = now.getHours() * 60 + now.getMinutes()
    const blocked = getStore('blocked_slots').filter(b => b.date === date).map(b => b.time)
    const appts = getStore('appointments').filter(a => a.date === date && a.status === 'scheduled').map(a => a.time?.slice(0, 5))
    return slots.filter(s => {
        if (blocked.includes(s)) return false
        if (appts.includes(s)) return false
        if (date === today) {
            const [sh, sm] = s.split(':').map(Number)
            if (sh * 60 + sm <= nowMin) return false
        }
        return true
    })
}

function delay(ms = 150) { return new Promise(r => setTimeout(r, ms)) }

export const api = {
    // Cliente
    async login(phone) {
        await delay()
        const clients = getStore('clients')
        const c = clients.find(c => c.phone === phone)
        if (!c) throw new Error('Cliente não encontrado')
        return c
    },

    async register(name, phone) {
        await delay()
        const clients = getStore('clients')
        if (clients.find(c => c.phone === phone)) throw new Error('Telefone já cadastrado')
        const c = { id: nextId('clients'), name, phone, created_at: new Date().toISOString() }
        clients.push(c)
        setStore('clients', clients)
        return c
    },

    async checkPhone(phone) {
        await delay()
        const clients = getStore('clients')
        return { exists: !!clients.find(c => c.phone === phone) }
    },

    // Slots
    async getSlots(date) {
        await delay()
        return { slots: generateSlots(date) }
    },

    // Agendamentos
    async createAppointment(clientId, serviceIds, date, time) {
        await delay()
        const appts = getStore('appointments')
        const active = appts.filter(a => a.client_id === clientId && a.status === 'scheduled')
        if (active.length >= 3) throw new Error('Limite de 3 agendamentos ativos atingido')
        if (appts.find(a => a.date === date && a.time?.slice(0, 5) === time && a.status === 'scheduled'))
            throw new Error('Horário já ocupado')
        const services = getStore('services')
        const names = (serviceIds || []).map(id => services.find(s => s.id === id)?.name).filter(Boolean)
        const a = {
            id: nextId('appointments'), client_id: clientId, date, time: time + ':00',
            status: 'scheduled', created_at: new Date().toISOString(),
            service_ids: serviceIds || [], service_name: names.join(', '), service_names: names
        }
        appts.push(a)
        setStore('appointments', appts)
        return a
    },

    async getMyAppointments(clientId) {
        await delay()
        return getStore('appointments').filter(a => a.client_id === clientId && a.status === 'scheduled')
    },

    async cancelAppointment(id, clientId) {
        await delay()
        const appts = getStore('appointments')
        const idx = appts.findIndex(a => a.id === id && a.client_id === clientId)
        if (idx < 0) throw new Error('Agendamento não encontrado')
        const a = appts[idx]
        const apptDate = new Date(`${a.date}T${a.time?.slice(0, 5)}:00`)
        if (apptDate - new Date() < 30 * 60 * 1000) throw new Error('Só é possível cancelar com pelo menos 30 minutos de antecedência')
        appts.splice(idx, 1)
        setStore('appointments', appts)
        return { message: 'Cancelado' }
    },

    // Serviços
    async getServices() {
        await delay()
        return getStore('services').filter(s => s.active)
    },

    // Admin
    async adminLogin(user, password) {
        await delay()
        if (user === 'admin' && password === 'admin123') return { reminder_minutes: 60 }
        throw new Error('Credenciais inválidas')
    },

    async getSchedules() {
        await delay()
        return getStore('schedules')
    },

    async upsertSchedule(schedule) {
        await delay()
        const scheds = getStore('schedules')
        const idx = scheds.findIndex(s => s.day_of_week === schedule.day_of_week && s.slot === schedule.slot)
        if (idx >= 0) { scheds[idx] = { ...scheds[idx], ...schedule } }
        else { schedule.id = nextId('schedules'); scheds.push(schedule) }
        setStore('schedules', scheds)
        return { message: 'Salvo' }
    },

    async blockSlot(date, time) {
        await delay()
        const slots = getStore('blocked_slots')
        if (!slots.find(s => s.date === date && s.time === time)) {
            slots.push({ id: nextId('blocked_slots'), date, time })
            setStore('blocked_slots', slots)
        }
        return { message: 'Bloqueado' }
    },

    async unblockSlot(date, time) {
        await delay()
        const slots = getStore('blocked_slots').filter(s => !(s.date === date && s.time === time))
        setStore('blocked_slots', slots)
        return { message: 'Desbloqueado' }
    },

    async getAdminAppointments(date) {
        await delay()
        const appts = getStore('appointments').filter(a => a.date === date && a.status === 'scheduled')
        const clients = getStore('clients')
        return appts.map(a => ({ ...a, client_name: clients.find(c => c.id === a.client_id)?.name || 'Desconhecido' }))
    },

    async getBookedSlots(date) {
        await delay()
        return getStore('appointments').filter(a => a.date === date && a.status === 'scheduled').map(a => a.time)
    },

    async adminCancelAppointment(id) {
        await delay()
        const appts = getStore('appointments').filter(a => a.id !== id)
        setStore('appointments', appts)
        return { message: 'Cancelado' }
    },

    async getSettings() {
        await delay()
        return getStore('settings', { reminder_minutes: 60 })
    },

    async updateSettings(settings) {
        await delay()
        setStore('settings', settings)
        return { message: 'Atualizado' }
    },

    async getAdminServices() {
        await delay()
        return getStore('services')
    },

    async createService(service) {
        await delay()
        const services = getStore('services')
        service.id = nextId('services')
        service.active = true
        services.push(service)
        setStore('services', services)
        return service
    },

    async updateService(id, service) {
        await delay()
        const services = getStore('services')
        const idx = services.findIndex(s => s.id === id)
        if (idx >= 0) services[idx] = { ...services[idx], ...service, id }
        setStore('services', services)
        return { message: 'Atualizado' }
    },

    async deleteService(id) {
        await delay()
        setStore('services', getStore('services').filter(s => s.id !== id))
        return { message: 'Removido' }
    },
}
