package community

import (
	"errors"

	"github.com/Abiggj/pawwp/internal/account"
	"github.com/google/uuid"
)

type Service interface {
	CreateChannel(req CreateChannelRequest, accountID uuid.UUID) (*Channel, error)
	ListChannels() ([]Channel, error)
	GetChannel(id string) (*Channel, error)
	UpdateChannel(id string, req UpdateChannelRequest, accountID uuid.UUID) (*Channel, error)

	CreatePost(channelID string, req CreatePostRequest, authorID uuid.UUID) (*ChannelPost, error)
	ListPosts(channelID string) ([]ChannelPost, error)
	GetSubscribedFeed(accountID uuid.UUID) ([]ChannelPost, error)
	UpdatePostStatus(postID string, status string, accountID uuid.UUID) (*ChannelPost, error)

	Subscribe(channelID string, accountID uuid.UUID) error
	Unsubscribe(channelID string, accountID uuid.UUID) error
	GetSubscribedChannels(accountID uuid.UUID) ([]Channel, error)

	AddChannelAdmin(channelID string, adminAccountID string, role string, requesterID uuid.UUID) error
	ListChannelAdmins(channelID string) ([]ChannelAdmin, error)

	CreateTownhallQuestion(req CreateQuestionRequest, authorID uuid.UUID) (*TownhallQuestion, error)
	ListTownhallQuestions(category string) ([]TownhallQuestion, error)
	CreateTownhallAnswer(questionID string, content string, authorID uuid.UUID) (*TownhallAnswer, error)
	VoteTownhallQuestion(questionID string) error
}

type service struct {
	repo        Repository
	accountRepo account.Repository
}

func NewService(repo Repository, accountRepo account.Repository) Service {
	return &service{repo: repo, accountRepo: accountRepo}
}

type CreateChannelRequest struct {
	Name         string `json:"name"`
	Description  string `json:"description"`
	CoverageArea string `json:"coverage_area"`
}

type UpdateChannelRequest struct {
	Name         string `json:"name"`
	Description  string `json:"description"`
	CoverageArea string `json:"coverage_area"`
	IsActive     *bool  `json:"is_active"`
}

type CreatePostRequest struct {
	Type     string `json:"type"` // sos, adoption, donation
	Title    string `json:"title"`
	Content  string `json:"content"`
	MediaURL string `json:"media_url"`
	Status   string `json:"status"` // optional, e.g., urgent
}

type CreateQuestionRequest struct {
	Title    string `json:"title"`
	Content  string `json:"content"`
	Category string `json:"category"`
	IsUrgent bool   `json:"is_urgent"`
}

func (s *service) CreateChannel(req CreateChannelRequest, accountID uuid.UUID) (*Channel, error) {
	acc, err := s.accountRepo.FindByID(accountID.String())
	if err != nil {
		return nil, errors.New("account not found")
	}

	if acc.Type != "shelter" {
		return nil, errors.New("unauthorized: only shelter accounts can create channels")
	}

	channel := &Channel{
		AccountID:    accountID,
		Name:         req.Name,
		Description:  req.Description,
		CoverageArea: req.CoverageArea,
	}

	err = s.repo.CreateChannel(channel)
	if err == nil {
		_ = s.repo.AddChannelAdmin(&ChannelAdmin{
			ChannelID: channel.ID,
			AccountID: accountID,
			Role:      "owner",
		})
	}
	return channel, err
}

func (s *service) ListChannels() ([]Channel, error) {
	return s.repo.ListActiveChannels()
}

func (s *service) GetChannel(id string) (*Channel, error) {
	return s.repo.FindChannelByID(id)
}

func (s *service) UpdateChannel(id string, req UpdateChannelRequest, accountID uuid.UUID) (*Channel, error) {
	channel, err := s.repo.FindChannelByID(id)
	if err != nil {
		return nil, errors.New("channel not found")
	}

	if channel.AccountID != accountID {
		return nil, errors.New("unauthorized")
	}

	channel.Name = req.Name
	channel.Description = req.Description
	channel.CoverageArea = req.CoverageArea
	if req.IsActive != nil {
		channel.IsActive = *req.IsActive
	}

	err = s.repo.UpdateChannel(channel)
	return channel, err
}

func (s *service) CreatePost(channelID string, req CreatePostRequest, authorID uuid.UUID) (*ChannelPost, error) {
	cID, err := uuid.Parse(channelID)
	if err != nil {
		return nil, errors.New("invalid channel id")
	}

	channel, err := s.repo.FindChannelByID(channelID)
	if err != nil {
		return nil, errors.New("channel not found")
	}

	isAdmin, _ := s.repo.IsChannelAdmin(cID, authorID)
	if channel.AccountID != authorID && !isAdmin {
		return nil, errors.New("unauthorized: author is not a designated admin or owner of this channel")
	}

	if req.Type != "sos" && req.Type != "adoption" && req.Type != "donation" {
		return nil, errors.New("invalid post type")
	}

	post := &ChannelPost{
		ChannelID: cID,
		AuthorID:  authorID,
		Type:      req.Type,
		Title:     req.Title,
		Content:   req.Content,
		MediaURL:  req.MediaURL,
		Status:    req.Status,
	}

	err = s.repo.CreatePost(post)
	return post, err
}

func (s *service) ListPosts(channelID string) ([]ChannelPost, error) {
	cID, err := uuid.Parse(channelID)
	if err != nil {
		return nil, errors.New("invalid channel id")
	}
	return s.repo.ListPostsByChannelID(cID)
}

func (s *service) GetSubscribedFeed(accountID uuid.UUID) ([]ChannelPost, error) {
	return s.repo.ListPostsBySubscribedChannels(accountID)
}

func (s *service) UpdatePostStatus(postID string, status string, accountID uuid.UUID) (*ChannelPost, error) {
	post, err := s.repo.FindPostByID(postID)
	if err != nil {
		return nil, errors.New("post not found")
	}

	channel, err := s.repo.FindChannelByID(post.ChannelID.String())
	if err != nil {
		return nil, errors.New("associated channel not found")
	}

	if post.AuthorID != accountID && channel.AccountID != accountID {
		return nil, errors.New("unauthorized")
	}

	post.Status = status
	err = s.repo.UpdatePost(post)
	return post, err
}

func (s *service) Subscribe(channelID string, accountID uuid.UUID) error {
	cID, err := uuid.Parse(channelID)
	if err != nil {
		return errors.New("invalid channel id")
	}

	sub := &Subscription{
		AccountID: accountID,
		ChannelID: cID,
	}
	return s.repo.Subscribe(sub)
}

func (s *service) Unsubscribe(channelID string, accountID uuid.UUID) error {
	cID, err := uuid.Parse(channelID)
	if err != nil {
		return errors.New("invalid channel id")
	}
	return s.repo.Unsubscribe(accountID, cID)
}

func (s *service) GetSubscribedChannels(accountID uuid.UUID) ([]Channel, error) {
	subs, err := s.repo.GetSubscriptions(accountID)
	if err != nil {
		return nil, err
	}

	var channels []Channel
	for _, sub := range subs {
		c, err := s.repo.FindChannelByID(sub.ChannelID.String())
		if err == nil {
			channels = append(channels, *c)
		}
	}
	return channels, nil
}

func (s *service) AddChannelAdmin(channelID string, adminAccountID string, role string, requesterID uuid.UUID) error {
	cID, err := uuid.Parse(channelID)
	if err != nil {
		return errors.New("invalid channel id")
	}
	adminUUID, err := uuid.Parse(adminAccountID)
	if err != nil {
		return errors.New("invalid admin account id")
	}

	channel, err := s.repo.FindChannelByID(channelID)
	if err != nil {
		return errors.New("channel not found")
	}

	if channel.AccountID != requesterID {
		return errors.New("unauthorized: only channel owner can add admins")
	}

	if role == "" {
		role = "admin"
	}

	return s.repo.AddChannelAdmin(&ChannelAdmin{
		ChannelID: cID,
		AccountID: adminUUID,
		Role:      role,
	})
}

func (s *service) ListChannelAdmins(channelID string) ([]ChannelAdmin, error) {
	cID, err := uuid.Parse(channelID)
	if err != nil {
		return nil, errors.New("invalid channel id")
	}
	return s.repo.ListChannelAdmins(cID)
}

func (s *service) CreateTownhallQuestion(req CreateQuestionRequest, authorID uuid.UUID) (*TownhallQuestion, error) {
	acc, err := s.accountRepo.FindByID(authorID.String())
	if err != nil {
		return nil, errors.New("account not found")
	}

	authorName := acc.Email
	if acc.Bio != "" {
		authorName = acc.Bio
	}

	category := req.Category
	if category == "" {
		category = "general"
	}

	q := &TownhallQuestion{
		AuthorID:   authorID,
		AuthorName: authorName,
		AuthorType: acc.Type,
		Title:      req.Title,
		Content:    req.Content,
		Category:   category,
		IsUrgent:   req.IsUrgent,
	}
	err = s.repo.CreateQuestion(q)
	return q, err
}

func (s *service) ListTownhallQuestions(category string) ([]TownhallQuestion, error) {
	return s.repo.ListQuestions(category)
}

func (s *service) CreateTownhallAnswer(questionID string, content string, authorID uuid.UUID) (*TownhallAnswer, error) {
	qID, err := uuid.Parse(questionID)
	if err != nil {
		return nil, errors.New("invalid question id")
	}

	acc, err := s.accountRepo.FindByID(authorID.String())
	if err != nil {
		return nil, errors.New("account not found")
	}

	authorName := acc.Email
	if acc.Bio != "" {
		authorName = acc.Bio
	}

	isVerified := acc.Type == "shelter" || acc.IsVerified

	a := &TownhallAnswer{
		QuestionID: qID,
		AuthorID:   authorID,
		AuthorName: authorName,
		AuthorType: acc.Type,
		Content:    content,
		IsVerified: isVerified,
	}
	err = s.repo.CreateAnswer(a)
	return a, err
}

func (s *service) VoteTownhallQuestion(questionID string) error {
	return s.repo.VoteQuestion(questionID)
}

