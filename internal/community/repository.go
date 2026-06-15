package community

import (
	"github.com/Abiggj/pawwp/internal/database"
	"github.com/google/uuid"
)

type Repository interface {
	CreateChannel(channel *Channel) error
	FindChannelByID(id string) (*Channel, error)
	ListActiveChannels() ([]Channel, error)
	UpdateChannel(channel *Channel) error

	CreatePost(post *ChannelPost) error
	FindPostByID(id string) (*ChannelPost, error)
	ListPostsByChannelID(channelID uuid.UUID) ([]ChannelPost, error)
	ListPostsBySubscribedChannels(accountID uuid.UUID) ([]ChannelPost, error)
	UpdatePost(post *ChannelPost) error

	Subscribe(sub *Subscription) error
	Unsubscribe(accountID, channelID uuid.UUID) error
	GetSubscriptions(accountID uuid.UUID) ([]Subscription, error)
}

type repository struct{}

func NewRepository() Repository {
	return &repository{}
}

func (r *repository) CreateChannel(channel *Channel) error {
	return database.DB.Create(channel).Error
}

func (r *repository) FindChannelByID(id string) (*Channel, error) {
	var c Channel
	err := database.DB.Where("id = ?", id).First(&c).Error
	return &c, err
}

func (r *repository) ListActiveChannels() ([]Channel, error) {
	var channels []Channel
	err := database.DB.Where("is_active = ?", true).Find(&channels).Error
	return channels, err
}

func (r *repository) UpdateChannel(channel *Channel) error {
	return database.DB.Save(channel).Error
}

func (r *repository) CreatePost(post *ChannelPost) error {
	return database.DB.Create(post).Error
}

func (r *repository) FindPostByID(id string) (*ChannelPost, error) {
	var p ChannelPost
	err := database.DB.Where("id = ?", id).First(&p).Error
	return &p, err
}

func (r *repository) ListPostsByChannelID(channelID uuid.UUID) ([]ChannelPost, error) {
	var posts []ChannelPost
	err := database.DB.Where("channel_id = ?", channelID).Order("created_at desc").Find(&posts).Error
	return posts, err
}

func (r *repository) ListPostsBySubscribedChannels(accountID uuid.UUID) ([]ChannelPost, error) {
	var posts []ChannelPost
	// JOIN posts with subscriptions
	err := database.DB.Table("channel_posts").
		Joins("JOIN subscriptions ON subscriptions.channel_id = channel_posts.channel_id").
		Where("subscriptions.account_id = ?", accountID).
		Order("channel_posts.created_at desc").
		Find(&posts).Error
	return posts, err
}

func (r *repository) UpdatePost(post *ChannelPost) error {
	return database.DB.Save(post).Error
}

func (r *repository) Subscribe(sub *Subscription) error {
	return database.DB.Create(sub).Error
}

func (r *repository) Unsubscribe(accountID, channelID uuid.UUID) error {
	return database.DB.Delete(&Subscription{}, "account_id = ? AND channel_id = ?", accountID, channelID).Error
}

func (r *repository) GetSubscriptions(accountID uuid.UUID) ([]Subscription, error) {
	var subs []Subscription
	err := database.DB.Where("account_id = ?", accountID).Find(&subs).Error
	return subs, err
}
