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
