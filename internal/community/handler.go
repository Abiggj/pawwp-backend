package community

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

func (h *Handler) CreateChannel(w http.ResponseWriter, r *http.Request) {
	var req CreateChannelRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	accountID, err := h.getAccountID(r)
	if err != nil {
		response.Error(w, http.StatusUnauthorized, err.Error())
		return
	}

	channel, err := h.service.CreateChannel(req, accountID)
	if err != nil {
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	response.JSON(w, http.StatusCreated, channel)
}

func (h *Handler) ListChannels(w http.ResponseWriter, r *http.Request) {
	channels, err := h.service.ListChannels()
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	response.JSON(w, http.StatusOK, channels)
}

func (h *Handler) GetChannel(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	channel, err := h.service.GetChannel(id)
	if err != nil {
		response.Error(w, http.StatusNotFound, "Channel not found")
		return
	}
	response.JSON(w, http.StatusOK, channel)
}

func (h *Handler) UpdateChannel(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var req UpdateChannelRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid request")
		return
	}

	accountID, err := h.getAccountID(r)
	if err != nil {
		response.Error(w, http.StatusUnauthorized, err.Error())
		return
	}

	channel, err := h.service.UpdateChannel(id, req, accountID)
	if err != nil {
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	response.JSON(w, http.StatusOK, channel)
}

func (h *Handler) CreatePost(w http.ResponseWriter, r *http.Request) {
	channelID := r.PathValue("id")
	var req CreatePostRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	authorID, err := h.getAccountID(r)
	if err != nil {
		response.Error(w, http.StatusUnauthorized, err.Error())
		return
	}

	post, err := h.service.CreatePost(channelID, req, authorID)
	if err != nil {
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	response.JSON(w, http.StatusCreated, post)
}

func (h *Handler) ListPosts(w http.ResponseWriter, r *http.Request) {
	channelID := r.PathValue("id")
	posts, err := h.service.ListPosts(channelID)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	response.JSON(w, http.StatusOK, posts)
}

func (h *Handler) GetFeed(w http.ResponseWriter, r *http.Request) {
	accountID, err := h.getAccountID(r)
	if err != nil {
		response.Error(w, http.StatusUnauthorized, err.Error())
		return
	}

	posts, err := h.service.GetSubscribedFeed(accountID)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	response.JSON(w, http.StatusOK, posts)
}

func (h *Handler) UpdatePostStatus(w http.ResponseWriter, r *http.Request) {
	postID := r.PathValue("post_id")
	var req struct {
		Status string `json:"status"`
	}
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil || req.Status == "" {
		response.Error(w, http.StatusBadRequest, "Invalid request")
		return
	}

	accountID, err := h.getAccountID(r)
	if err != nil {
		response.Error(w, http.StatusUnauthorized, err.Error())
		return
	}

	post, err := h.service.UpdatePostStatus(postID, req.Status, accountID)
	if err != nil {
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	response.JSON(w, http.StatusOK, post)
}

func (h *Handler) Subscribe(w http.ResponseWriter, r *http.Request) {
	channelID := r.PathValue("id")
	accountID, err := h.getAccountID(r)
	if err != nil {
		response.Error(w, http.StatusUnauthorized, err.Error())
		return
	}

	err = h.service.Subscribe(channelID, accountID)
	if err != nil {
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	response.JSON(w, http.StatusOK, map[string]string{"message": "Subscribed successfully"})
}

func (h *Handler) Unsubscribe(w http.ResponseWriter, r *http.Request) {
	channelID := r.PathValue("id")
	accountID, err := h.getAccountID(r)
	if err != nil {
		response.Error(w, http.StatusUnauthorized, err.Error())
		return
	}

	err = h.service.Unsubscribe(channelID, accountID)
	if err != nil {
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	response.JSON(w, http.StatusOK, map[string]string{"message": "Unsubscribed successfully"})
}

func (h *Handler) GetMySubscriptions(w http.ResponseWriter, r *http.Request) {
	accountID, err := h.getAccountID(r)
	if err != nil {
		response.Error(w, http.StatusUnauthorized, err.Error())
		return
	}

	channels, err := h.service.GetSubscribedChannels(accountID)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	response.JSON(w, http.StatusOK, channels)
}

func (h *Handler) ListChannelAdmins(w http.ResponseWriter, r *http.Request) {
	channelID := r.PathValue("id")
	admins, err := h.service.ListChannelAdmins(channelID)
	if err != nil {
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	response.JSON(w, http.StatusOK, admins)
}

func (h *Handler) AddChannelAdmin(w http.ResponseWriter, r *http.Request) {
	channelID := r.PathValue("id")
	var req struct {
		AccountID string `json:"account_id"`
		Role      string `json:"role"`
	}
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil || req.AccountID == "" {
		response.Error(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	accountID, err := h.getAccountID(r)
	if err != nil {
		response.Error(w, http.StatusUnauthorized, err.Error())
		return
	}

	err = h.service.AddChannelAdmin(channelID, req.AccountID, req.Role, accountID)
	if err != nil {
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	response.JSON(w, http.StatusOK, map[string]string{"message": "Admin added successfully"})
}

func (h *Handler) ListQuestions(w http.ResponseWriter, r *http.Request) {
	category := r.URL.Query().Get("category")
	questions, err := h.service.ListTownhallQuestions(category)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	response.JSON(w, http.StatusOK, questions)
}

func (h *Handler) CreateQuestion(w http.ResponseWriter, r *http.Request) {
	var req CreateQuestionRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil || req.Title == "" || req.Content == "" {
		response.Error(w, http.StatusBadRequest, "Title and content are required")
		return
	}

	accountID, err := h.getAccountID(r)
	if err != nil {
		response.Error(w, http.StatusUnauthorized, err.Error())
		return
	}

	question, err := h.service.CreateTownhallQuestion(req, accountID)
	if err != nil {
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	response.JSON(w, http.StatusCreated, question)
}

func (h *Handler) CreateAnswer(w http.ResponseWriter, r *http.Request) {
	questionID := r.PathValue("id")
	var req struct {
		Content string `json:"content"`
	}
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil || req.Content == "" {
		response.Error(w, http.StatusBadRequest, "Content is required")
		return
	}

	accountID, err := h.getAccountID(r)
	if err != nil {
		response.Error(w, http.StatusUnauthorized, err.Error())
		return
	}

	answer, err := h.service.CreateTownhallAnswer(questionID, req.Content, accountID)
	if err != nil {
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	response.JSON(w, http.StatusCreated, answer)
}

func (h *Handler) VoteQuestion(w http.ResponseWriter, r *http.Request) {
	questionID := r.PathValue("id")
	err := h.service.VoteTownhallQuestion(questionID)
	if err != nil {
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	response.JSON(w, http.StatusOK, map[string]string{"message": "Voted successfully"})
}

