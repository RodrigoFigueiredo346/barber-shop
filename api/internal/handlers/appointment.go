package handlers

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"barber-app/internal/repository"
	"barber-app/internal/services"

	"github.com/go-chi/chi/v5"
)

type AppointmentHandler struct {
	repo        *repository.AppointmentRepo
	slotService *services.SlotService
}

func NewAppointmentHandler(repo *repository.AppointmentRepo, slotService *services.SlotService) *AppointmentHandler {
	return &AppointmentHandler{repo: repo, slotService: slotService}
}

// GetAvailableSlots retorna horários disponíveis para uma data.
func (h *AppointmentHandler) GetAvailableSlots(w http.ResponseWriter, r *http.Request) {
	date := r.URL.Query().Get("date")
	if date == "" {
		jsonError(w, http.StatusBadRequest, "Parâmetro 'date' é obrigatório (YYYY-MM-DD)")
		return
	}
	slots, err := h.slotService.GetAvailableSlots(r.Context(), date)
	if err != nil {
		jsonError(w, http.StatusBadRequest, err.Error())
		return
	}
	jsonResponse(w, http.StatusOK, map[string][]string{"slots": slots})
}

// Create agenda um horário.
func (h *AppointmentHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req struct {
		ClientID int    `json:"client_id"`
		Date     string `json:"date"`
		Time     string `json:"time"`
	}
	if err := decodeJSON(r, &req); err != nil {
		jsonError(w, http.StatusBadRequest, "JSON inválido")
		return
	}

	// Verifica limite de 3 agendamentos ativos
	count, err := h.repo.CountActiveByClient(r.Context(), req.ClientID)
	if err != nil {
		jsonError(w, http.StatusInternalServerError, "Erro ao verificar agendamentos")
		return
	}
	if count >= 3 {
		jsonError(w, http.StatusForbidden, "Limite de 3 agendamentos ativos atingido")
		return
	}

	// Verifica se slot está disponível
	taken, _ := h.repo.IsSlotTaken(r.Context(), req.Date, req.Time+":00")
	if taken {
		jsonError(w, http.StatusConflict, "Horário já ocupado")
		return
	}

	appointment, err := h.repo.Create(r.Context(), req.ClientID, req.Date, req.Time+":00")
	if err != nil {
		jsonError(w, http.StatusInternalServerError, "Erro ao agendar")
		return
	}
	jsonResponse(w, http.StatusCreated, appointment)
}

// GetByClient retorna agendamentos de um cliente.
func (h *AppointmentHandler) GetByClient(w http.ResponseWriter, r *http.Request) {
	clientID, err := strconv.Atoi(chi.URLParam(r, "clientID"))
	if err != nil {
		jsonError(w, http.StatusBadRequest, "client_id inválido")
		return
	}
	appointments, err := h.repo.GetByClient(r.Context(), clientID)
	if err != nil {
		jsonError(w, http.StatusInternalServerError, "Erro ao buscar agendamentos")
		return
	}
	jsonResponse(w, http.StatusOK, appointments)
}

// Cancel cancela um agendamento do cliente (mínimo 30 min de antecedência).
func (h *AppointmentHandler) Cancel(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(chi.URLParam(r, "id"))
	var req struct {
		ClientID int `json:"client_id"`
	}
	if err := decodeJSON(r, &req); err != nil {
		jsonError(w, http.StatusBadRequest, "JSON inválido")
		return
	}

	// Busca o agendamento pra validar antecedência
	appt, err := h.repo.GetByID(r.Context(), id)
	if err != nil {
		jsonError(w, http.StatusNotFound, "Agendamento não encontrado")
		return
	}

	// Monta o datetime do agendamento
	apptTime, err := time.Parse("2006-01-02 15:04:05", fmt.Sprintf("%s %s", appt.Date, appt.Time))
	if err != nil {
		// Tenta sem segundos
		apptTime, err = time.Parse("2006-01-02 15:04", fmt.Sprintf("%s %s", appt.Date, appt.Time[:5]))
		if err != nil {
			jsonError(w, http.StatusInternalServerError, "Erro ao processar data do agendamento")
			return
		}
	}

	if time.Until(apptTime) < 30*time.Minute {
		jsonError(w, http.StatusForbidden, "Só é possível cancelar com pelo menos 30 minutos de antecedência")
		return
	}

	if err := h.repo.Cancel(r.Context(), id, req.ClientID); err != nil {
		jsonError(w, http.StatusInternalServerError, "Erro ao cancelar")
		return
	}
	jsonResponse(w, http.StatusOK, map[string]string{"message": "Cancelado"})
}
