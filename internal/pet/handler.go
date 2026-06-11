package pet

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/Abiggj/pawwp/internal/middleware"
	"github.com/Abiggj/pawwp/internal/response"
	"github.com/google/uuid"
)

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) getAccountID(r *http.Request) (uuid.UUID, error) {
	accountIDRaw := r.Context().Value(middleware.AccountIDKey)
	if accountIDRaw == nil {
		return uuid.Nil, errors.New("unauthorized")
	}
	return uuid.Parse(accountIDRaw.(string))
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var req CreatePetRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid request")
		return
	}

	accountID, err := h.getAccountID(r)
	if err != nil {
		response.Error(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	pet, err := h.service.Create(req, accountID)
	if err != nil {
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	response.JSON(w, http.StatusCreated, pet)
}

func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		response.Error(w, http.StatusBadRequest, "Missing ID")
		return
	}

	pet, err := h.service.GetByID(id)
	if err != nil {
		response.Error(w, http.StatusNotFound, "Pet not found")
		return
	}

	response.JSON(w, http.StatusOK, pet)
}

func (h *Handler) ListByAccount(w http.ResponseWriter, r *http.Request) {
	accountIDStr := r.PathValue("account_id")
	accountID, err := uuid.Parse(accountIDStr)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid account ID")
		return
	}

	pets, err := h.service.GetByAccountID(accountID)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	response.JSON(w, http.StatusOK, pets)
}

func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var req UpdatePetRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid request")
		return
	}

	accountID, err := h.getAccountID(r)
	if err != nil {
		response.Error(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	pet, err := h.service.Update(id, req, accountID)
	if err != nil {
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	response.JSON(w, http.StatusOK, pet)
}

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	accountID, err := h.getAccountID(r)
	if err != nil {
		response.Error(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	err = h.service.Delete(id, accountID)
	if err != nil {
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	response.JSON(w, http.StatusOK, map[string]string{"message": "Pet deleted successfully"})
}
