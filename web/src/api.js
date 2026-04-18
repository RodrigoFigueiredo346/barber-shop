const BASE = '/api'

async function request(url, options = {}) {
    const res = await fetch(BASE + url, {
        headers: { 'Content-Type': 'application/json', ...options.headers },
        ...options
    })
    const data = await res.json()
    if (!res.ok) throw new Error(data.error || 'Erro desconhecido')
    return data
}

function adminHeaders() {
    const creds = localStorage.getItem('barber_admin')
    if (!creds) return {}
    return { Authorization: 'Basic ' + btoa(creds) }
}

export const api = {
    // Cliente
    login: (phone) => request('/clients/login', { method: 'POST', body: JSON.stringify({ phone }) }),
    register: (name, phone) => request('/clients/register', { method: 'POST', body: JSON.stringify({ name, phone }) }),
    checkPhone: (phone) => request(`/clients/check?phone=${phone}`),

    // Slots e agendamentos
    getSlots: (date) => request(`/slots?date=${date}`),
    createAppointment: (clientId, serviceId, date, time) =>
        request('/appointments', { method: 'POST', body: JSON.stringify({ client_id: clientId, service_id: serviceId, date, time }) }),
    getMyAppointments: (clientId) => request(`/appointments/client/${clientId}`),
    cancelAppointment: (id, clientId) =>
        request(`/appointments/${id}/cancel`, { method: 'PUT', body: JSON.stringify({ client_id: clientId }) }),

    // Admin
    adminLogin: (user, password) => {
        const creds = user + ':' + password
        return request('/admin/settings', { headers: { Authorization: 'Basic ' + btoa(creds) } })
    },
    getSchedules: () => request('/admin/schedules', { headers: adminHeaders() }),
    upsertSchedule: (schedule) =>
        request('/admin/schedules', { method: 'PUT', body: JSON.stringify(schedule), headers: adminHeaders() }),
    blockSlot: (date, time) =>
        request('/admin/blocked-slots', { method: 'POST', body: JSON.stringify({ date, time }), headers: adminHeaders() }),
    unblockSlot: (date, time) =>
        request('/admin/blocked-slots', { method: 'DELETE', body: JSON.stringify({ date, time }), headers: adminHeaders() }),
    getAdminAppointments: (date) =>
        request(`/admin/appointments?date=${date}`, { headers: adminHeaders() }),
    getBookedSlots: (date) =>
        request(`/admin/booked-slots?date=${date}`, { headers: adminHeaders() }),
    adminCancelAppointment: (id) =>
        request(`/admin/appointments/${id}`, { method: 'DELETE', headers: adminHeaders() }),
    getSettings: () => request('/admin/settings', { headers: adminHeaders() }),
    updateSettings: (settings) =>
        request('/admin/settings', { method: 'PUT', body: JSON.stringify(settings), headers: adminHeaders() }),

    // Serviços
    getServices: () => request('/services'),
    getAdminServices: () => request('/admin/services', { headers: adminHeaders() }),
    createService: (service) =>
        request('/admin/services', { method: 'POST', body: JSON.stringify(service), headers: adminHeaders() }),
    updateService: (id, service) =>
        request(`/admin/services/${id}`, { method: 'PUT', body: JSON.stringify(service), headers: adminHeaders() }),
    deleteService: (id) =>
        request(`/admin/services/${id}`, { method: 'DELETE', headers: adminHeaders() }),
}
