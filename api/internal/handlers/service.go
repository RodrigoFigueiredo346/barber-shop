package handlers

import (
	"net/http"
	"strconv"

	"barber-app/internal/models"
	"barber-app/internal/repository"

	"github.com/go-chi/chi/v5"
)

type ServiceHandler struct {
	repo *repository.ServiceRepo
}

func NewServiceHandler(repo *repository.ServiceRepo) *ServiceHandler {
	return &ServiceHandler{repo: repo}
}

// GetActive retorna serviços ativos (rota pública).
func (h *ServiceHandler) GetActive(w http.ResponseWriter, r *http.Request) {
	services, err := h.repo.GetActive(r.Context())
	if err != nil {
		jsonError(w, http.StatusInternalServerError, "Erro ao buscar serviços")
		return
	}
	jsonResponse(w, http.StatusOK, services)
}

// GetAll retorna todos os serviços (admin).
func (h *ServiceHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	services, err := h.repo.GetAll(r.Context())
	if err != nil {
		jsonError(w, http.StatusInternalServerError, "Erro ao buscar serviços")
		return
	}
	jsonResponse(w, http.StatusOK, services)
}

// Create cria um serviço (admin).
func (h *ServiceHandler) Create(w http.ResponseWriter, r *http.Request) {
	var s models.Service
	if err := decodeJSON(r, &s); err != nil {
		jsonError(w, http.StatusBadRequest, "JSON inválido")
		return
	}
	if s.Name == "" {
		jsonError(w, http.StatusBadRequest, "Nome é obrigatório")
		return
	}
	s.Active = true
	created, err := h.repo.Create(r.Context(), s)
	if err != nil {
		jsonError(w, http.StatusInternalServerError, "Erro ao criar serviço")
		return
	}
	jsonResponse(w, http.StatusCreated, created)
}

// Update atualiza um serviço (admin).
func (h *ServiceHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		jsonError(w, http.StatusBadRequest, "ID inválido")
		return
	}
	var s models.Service
	if err := decodeJSON(r, &s); err != nil {
		jsonError(w, http.StatusBadRequest, "JSON inválido")
		return
	}
	s.ID = id
	if err := h.repo.Update(r.Context(), s); err != nil {
		jsonError(w, http.StatusInternalServerError, "Erro ao atualizar")
		return
	}
	jsonResponse(w, http.StatusOK, map[string]string{"message": "Atualizado"})
}

// Delete remove um serviço (admin).
func (h *ServiceHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		jsonError(w, http.StatusBadRequest, "ID inválido")
		return
	}
	if err := h.repo.Delete(r.Context(), id); err != nil {
		jsonError(w, http.StatusInternalServerError, "Erro ao remover")
		return
	}
	jsonResponse(w, http.StatusOK, map[string]string{"message": "Removido"})
}
