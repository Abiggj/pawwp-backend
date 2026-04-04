package pet

import (
	"encoding/json"
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

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {

	var req CreatePetRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid request")
		return
	}

	accountIDRaw := r.Context().Value(middleware.AccountIDKey)
	if accountIDRaw == nil {
		response.Error(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	accountIDStr := accountIDRaw.(string)
	accountID, err := uuid.Parse(accountIDStr)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "Invalid account ID")
		return
	}

	pet, err := h.service.Create(req, accountID)
	if err != nil {
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	response.JSON(w, http.StatusCreated, pet)
}
