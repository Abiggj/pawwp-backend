package account

import (
	"encoding/json"
	"net/http"

	"github.com/Abiggj/pawwp/internal/middleware"
	"github.com/Abiggj/pawwp/internal/response"
)

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if req.Email == "" || req.Password == "" || req.Type == "" {
		response.Error(w, http.StatusBadRequest, "Missing required fields")
		return
	}

	account, err := h.service.Register(req)
	if err != nil {
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	response.JSON(w, http.StatusCreated, account)
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid request")
		return
	}

	authResp, err := h.service.Login(req)
	if err != nil {
		response.Error(w, http.StatusUnauthorized, err.Error())
		return
	}

	response.JSON(w, http.StatusOK, authResp)
}

func (h *Handler) Refresh(w http.ResponseWriter, r *http.Request) {
	var req struct {
		RefreshToken string `json:"refresh_token"`
	}
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil || req.RefreshToken == "" {
		response.Error(w, http.StatusBadRequest, "Invalid request")
		return
	}

	authResp, err := h.service.Refresh(req.RefreshToken)
	if err != nil {
		response.Error(w, http.StatusUnauthorized, err.Error())
		return
	}

	response.JSON(w, http.StatusOK, authResp)
}

func (h *Handler) Me(w http.ResponseWriter, r *http.Request) {
	accountID := r.Context().Value(middleware.AccountIDKey).(string)
	account, err := h.service.GetByID(accountID)
	if err != nil {
		response.Error(w, http.StatusNotFound, "Account not found")
		return
	}
	response.JSON(w, http.StatusOK, account)
}

func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	accountID := r.Context().Value(middleware.AccountIDKey).(string)
	var req UpdateRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid request")
		return
	}

	account, err := h.service.Update(accountID, req)
	if err != nil {
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	response.JSON(w, http.StatusOK, account)
}

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	accountID := r.Context().Value(middleware.AccountIDKey).(string)
	err := h.service.Delete(accountID)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	response.JSON(w, http.StatusOK, map[string]string{"message": "Account deleted successfully"})
}
