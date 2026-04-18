package handlers

import (
	"crypto/subtle"
	"net/http"
	"strconv"

	"barber-app/internal/models"
	"barber-app/internal/repository"

	"github.com/go-chi/chi/v5"
)

type AdminHandler struct {
	scheduleRepo    *repository.ScheduleRepo
	blockedSlotRepo *repository.BlockedSlotRepo
	appointmentRepo *repository.AppointmentRepo
	settingsRepo    *repository.SettingsRepo
	adminUser       string
	adminPassword   string
}

func NewAdminHandler(
	sr *repository.ScheduleRepo,
	br *repository.BlockedSlotRepo,
	ar *repository.AppointmentRepo,
	str *repository.SettingsRepo,
	user, password string,
) *AdminHandler {
	return &AdminHandler{
		scheduleRepo:    sr,
		blockedSlotRepo: br,
		appointmentRepo: ar,
		settingsRepo:    str,
		adminUser:       user,
		adminPassword:   password,
	}
}

// BasicAuth middleware para rotas admin.
func (h *AdminHandler) BasicAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, pass, ok := r.BasicAuth()
		if !ok ||
			subtle.ConstantTimeCompare([]byte(user), []byte(h.adminUser)) != 1 ||
			subtle.ConstantTimeCompare([]byte(pass), []byte(h.adminPassword)) != 1 {
			w.Header().Set("WWW-Authenticate", `Basic realm="admin"`)
			jsonError(w, http.StatusUnauthorized, "Credenciais inválidas")
			return
		}
		next.ServeHTTP(w, r)
	})
}

// UpsertSchedule cria ou atualiza horário de um dia da semana.
func (h *AdminHandler) UpsertSchedule(w http.ResponseWriter, r *http.Request) {
	var s models.Schedule
	if err := decodeJSON(r, &s); err != nil {
		jsonError(w, http.StatusBadRequest, "JSON inválido")
		return
	}
	if err := h.scheduleRepo.Upsert(r.Context(), s); err != nil {
		jsonError(w, http.StatusInternalServerError, "Erro ao salvar horário")
		return
	}
	jsonResponse(w, http.StatusOK, map[string]string{"message": "Salvo"})
}

// GetSchedules retorna todos os horários configurados.
func (h *AdminHandler) GetSchedules(w http.ResponseWriter, r *http.Request) {
	schedules, err := h.scheduleRepo.GetAll(r.Context())
	if err != nil {
		jsonError(w, http.StatusInternalServerError, "Erro ao buscar horários")
		return
	}
	jsonResponse(w, http.StatusOK, schedules)
}

// BlockSlot bloqueia um horário específico.
func (h *AdminHandler) BlockSlot(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Date string `json:"date"`
		Time string `json:"time"`
	}
	if err := decodeJSON(r, &req); err != nil {
		jsonError(w, http.StatusBadRequest, "JSON inválido")
		return
	}
	if err := h.blockedSlotRepo.Block(r.Context(), req.Date, req.Time); err != nil {
		jsonError(w, http.StatusInternalServerError, "Erro ao bloquear")
		return
	}
	jsonResponse(w, http.StatusOK, map[string]string{"message": "Bloqueado"})
}

// UnblockSlot desbloqueia um horário.
func (h *AdminHandler) UnblockSlot(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Date string `json:"date"`
		Time string `json:"time"`
	}
	if err := decodeJSON(r, &req); err != nil {
		jsonError(w, http.StatusBadRequest, "JSON inválido")
		return
	}
	if err := h.blockedSlotRepo.Unblock(r.Context(), req.Date, req.Time); err != nil {
		jsonError(w, http.StatusInternalServerError, "Erro ao desbloquear")
		return
	}
	jsonResponse(w, http.StatusOK, map[string]string{"message": "Desbloqueado"})
}

// GetAppointmentsByDate retorna agendamentos de uma data.
func (h *AdminHandler) GetAppointmentsByDate(w http.ResponseWriter, r *http.Request) {
	date := r.URL.Query().Get("date")
	if date == "" {
		jsonError(w, http.StatusBadRequest, "Parâmetro 'date' é obrigatório")
		return
	}
	appointments, err := h.appointmentRepo.GetByDate(r.Context(), date)
	if err != nil {
		jsonError(w, http.StatusInternalServerError, "Erro ao buscar agendamentos")
		return
	}
	jsonResponse(w, http.StatusOK, appointments)
}

// GetBookedSlots retorna os horários agendados de uma data (pra tela de bloquear).
func (h *AdminHandler) GetBookedSlots(w http.ResponseWriter, r *http.Request) {
	date := r.URL.Query().Get("date")
	if date == "" {
		jsonError(w, http.StatusBadRequest, "Parâmetro 'date' é obrigatório")
		return
	}
	slots, err := h.appointmentRepo.GetBookedSlotsByDate(r.Context(), date)
	if err != nil {
		jsonError(w, http.StatusInternalServerError, "Erro ao buscar horários agendados")
		return
	}
	jsonResponse(w, http.StatusOK, slots)
}

// CancelAppointment cancela um agendamento (admin).
func (h *AdminHandler) CancelAppointment(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		jsonError(w, http.StatusBadRequest, "ID inválido")
		return
	}
	if err := h.appointmentRepo.AdminCancel(r.Context(), id); err != nil {
		jsonError(w, http.StatusInternalServerError, "Erro ao cancelar")
		return
	}
	jsonResponse(w, http.StatusOK, map[string]string{"message": "Cancelado"})
}

// GetSettings retorna as configurações.
func (h *AdminHandler) GetSettings(w http.ResponseWriter, r *http.Request) {
	s, err := h.settingsRepo.Get(r.Context())
	if err != nil {
		jsonError(w, http.StatusInternalServerError, "Erro ao buscar configurações")
		return
	}
	jsonResponse(w, http.StatusOK, s)
}

// UpdateSettings atualiza as configurações.
func (h *AdminHandler) UpdateSettings(w http.ResponseWriter, r *http.Request) {
	var s models.Settings
	if err := decodeJSON(r, &s); err != nil {
		jsonError(w, http.StatusBadRequest, "JSON inválido")
		return
	}
	if err := h.settingsRepo.Update(r.Context(), s); err != nil {
		jsonError(w, http.StatusInternalServerError, "Erro ao atualizar")
		return
	}
	jsonResponse(w, http.StatusOK, map[string]string{"message": "Atualizado"})
}
