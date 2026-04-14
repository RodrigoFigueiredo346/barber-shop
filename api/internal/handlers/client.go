package handlers

import (
	"net/http"
	"regexp"

	"barber-app/internal/repository"
)

type ClientHandler struct {
	repo *repository.ClientRepo
}

func NewClientHandler(repo *repository.ClientRepo) *ClientHandler {
	return &ClientHandler{repo: repo}
}

var phoneRegex = regexp.MustCompile(`^\d{10,11}$`)

// Login verifica se o telefone existe e retorna o cliente.
func (h *ClientHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Phone string `json:"phone"`
	}
	if err := decodeJSON(r, &req); err != nil {
		jsonError(w, http.StatusBadRequest, "JSON inválido")
		return
	}
	if !phoneRegex.MatchString(req.Phone) {
		jsonError(w, http.StatusBadRequest, "Telefone inválido. Use DDD + número (10 ou 11 dígitos)")
		return
	}

	client, err := h.repo.FindByPhone(r.Context(), req.Phone)
	if err != nil {
		jsonError(w, http.StatusNotFound, "Cliente não encontrado")
		return
	}
	jsonResponse(w, http.StatusOK, client)
}

// Register cria um novo cliente.
func (h *ClientHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Name  string `json:"name"`
		Phone string `json:"phone"`
	}
	if err := decodeJSON(r, &req); err != nil {
		jsonError(w, http.StatusBadRequest, "JSON inválido")
		return
	}
	if !phoneRegex.MatchString(req.Phone) {
		jsonError(w, http.StatusBadRequest, "Telefone inválido. Use DDD + número (10 ou 11 dígitos)")
		return
	}
	if req.Name == "" {
		jsonError(w, http.StatusBadRequest, "Nome é obrigatório")
		return
	}

	// Verifica se já existe
	existing, _ := h.repo.FindByPhone(r.Context(), req.Phone)
	if existing != nil {
		jsonError(w, http.StatusConflict, "Telefone já cadastrado")
		return
	}

	client, err := h.repo.Create(r.Context(), req.Name, req.Phone)
	if err != nil {
		jsonError(w, http.StatusInternalServerError, "Erro ao cadastrar")
		return
	}
	jsonResponse(w, http.StatusCreated, client)
}

// CheckPhone verifica se um telefone já existe (usado no fluxo de cadastro).
func (h *ClientHandler) CheckPhone(w http.ResponseWriter, r *http.Request) {
	phone := r.URL.Query().Get("phone")
	if !phoneRegex.MatchString(phone) {
		jsonError(w, http.StatusBadRequest, "Telefone inválido")
		return
	}
	client, _ := h.repo.FindByPhone(r.Context(), phone)
	jsonResponse(w, http.StatusOK, map[string]bool{"exists": client != nil})
}
