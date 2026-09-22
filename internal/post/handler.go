package post

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
	var req CreatePostRequest
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

	post, err := h.service.CreatePost(req, accountID)
	if err != nil {
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	response.JSON(w, http.StatusCreated, post)
}

func (h *Handler) GetFeed(w http.ResponseWriter, r *http.Request) {
	posts, err := h.service.Get48hFeed()
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	response.JSON(w, http.StatusOK, posts)
}

func (h *Handler) GetShowcase(w http.ResponseWriter, r *http.Request) {
	petIDStr := r.PathValue("pet_id")
	petID, err := uuid.Parse(petIDStr)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid pet ID")
		return
	}

	posts, err := h.service.GetPetShowcase(petID)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	response.JSON(w, http.StatusOK, posts)
}

func (h *Handler) ToggleShowcase(w http.ResponseWriter, r *http.Request) {
	postID := r.PathValue("post_id")
	accountID, err := h.getAccountID(r)
	if err != nil {
		response.Error(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	post, err := h.service.ToggleShowcase(postID, accountID)
	if err != nil {
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	response.JSON(w, http.StatusOK, post)
}

func (h *Handler) GetArchive(w http.ResponseWriter, r *http.Request) {
	petIDStr := r.PathValue("pet_id")
	petID, err := uuid.Parse(petIDStr)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid pet ID")
		return
	}

	accountID, err := h.getAccountID(r)
	if err != nil {
		response.Error(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	posts, err := h.service.GetPetArchive(petID, accountID)
	if err != nil {
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	response.JSON(w, http.StatusOK, posts)
}

func (h *Handler) ToggleBoop(w http.ResponseWriter, r *http.Request) {
	postID := r.PathValue("post_id")
	accountID, err := h.getAccountID(r)
	if err != nil {
		response.Error(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	err = h.service.ToggleBoop(postID, accountID)
	if err != nil {
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	response.JSON(w, http.StatusOK, map[string]string{"message": "Boop toggled"})
}

func (h *Handler) AddWoof(w http.ResponseWriter, r *http.Request) {
	postID := r.PathValue("post_id")
	var req struct {
		Content string `json:"content"`
	}
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil || req.Content == "" {
		response.Error(w, http.StatusBadRequest, "Invalid request")
		return
	}

	accountID, err := h.getAccountID(r)
	if err != nil {
		response.Error(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	woof, err := h.service.AddWoof(postID, req.Content, accountID)
	if err != nil {
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	response.JSON(w, http.StatusCreated, woof)
}

func (h *Handler) GetWoofs(w http.ResponseWriter, r *http.Request) {
	postID := r.PathValue("post_id")
	woofs, err := h.service.GetPostWoofs(postID)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	response.JSON(w, http.StatusOK, woofs)
}

func (h *Handler) GetBoops(w http.ResponseWriter, r *http.Request) {
	postID := r.PathValue("post_id")
	boops, err := h.service.GetPostBoops(postID)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	response.JSON(w, http.StatusOK, boops)
}

type RecordViewRequest struct {
	DwellSeconds int    `json:"dwell_seconds"`
	ViewerKey    string `json:"viewer_key"`
}

func (h *Handler) RecordView(w http.ResponseWriter, r *http.Request) {
	postID := r.PathValue("post_id")
	var req RecordViewRequest
	_ = json.NewDecoder(r.Body).Decode(&req)

	if req.DwellSeconds == 0 {
		req.DwellSeconds = 2
	}

	var accountID *uuid.UUID
	accID, err := h.getAccountID(r)
	if err == nil {
		accountID = &accID
	}

	viewsCount, isNew, err := h.service.RecordView(postID, accountID, req.ViewerKey, req.DwellSeconds)
	if err != nil {
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	response.JSON(w, http.StatusOK, map[string]interface{}{
		"views_count":       viewsCount,
		"is_authentic_view": isNew,
		"threshold_met":     req.DwellSeconds >= 2,
	})
}
