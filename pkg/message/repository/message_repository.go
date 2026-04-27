package message_repository

import (
	message_model "github.com/EvolutionAPI/evolution-go/pkg/message/model"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type MessageRepository interface {
	InsertMessage(message message_model.Message) error
	GetMessageByID(messageID string) (*message_model.Message, error)
	DeleteAllMessages() (int64, error)
	GetLatestMessageID(source string) (string, string, error)
	GetChatMessages(instanceId, chat string, limit, offset int) ([]message_model.Message, error)
	ListChats(instanceId string, limit int) ([]message_model.Message, error)
}

type messageRepository struct {
	db *gorm.DB
}

func (m *messageRepository) InsertMessage(message message_model.Message) error {
	return m.db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "message_id"}},
		DoUpdates: clause.AssignmentColumns([]string{"timestamp", "status", "source", "chat", "from_me", "content", "message_type", "quoted_message_id"}),
	}).Create(&message).Error
}

func (m *messageRepository) GetMessageByID(messageID string) (*message_model.Message, error) {
	var message message_model.Message
	err := m.db.Where("message_id = ?", messageID).First(&message).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &message, nil
}

func (m *messageRepository) DeleteAllMessages() (int64, error) {
	result := m.db.Exec("DELETE FROM messages")
	return result.RowsAffected, result.Error
}

func (m *messageRepository) GetLatestMessageID(source string) (string, string, error) {
	var message message_model.Message
	err := m.db.Where("source = ?", source).Order("timestamp DESC").First(&message).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return "", "", nil
		}
		return "", "", err
	}
	return message.MessageID, message.Timestamp.Format("2006-01-02 15:04:05"), nil
}

// GetChatMessages retorna mensagens de um chat específico, ordenadas da mais recente para a mais antiga.
func (m *messageRepository) GetChatMessages(instanceId, chat string, limit, offset int) ([]message_model.Message, error) {
	var messages []message_model.Message
	query := m.db.Where("instance_id = ? AND chat = ?", instanceId, chat).
		Order("timestamp DESC")
	if limit > 0 {
		query = query.Limit(limit)
	}
	if offset > 0 {
		query = query.Offset(offset)
	}
	err := query.Find(&messages).Error
	return messages, err
}

// ListChats retorna a última mensagem de cada chat da instância (preview de conversas).
func (m *messageRepository) ListChats(instanceId string, limit int) ([]message_model.Message, error) {
	var messages []message_model.Message

	subQuery := m.db.Model(&message_model.Message{}).
		Select("MAX(timestamp) as max_ts, chat").
		Where("instance_id = ?", instanceId).
		Group("chat")

	query := m.db.Model(&message_model.Message{}).
		Joins("JOIN (?) as latest ON messages.chat = latest.chat AND messages.timestamp = latest.max_ts", subQuery).
		Where("messages.instance_id = ?", instanceId).
		Order("messages.timestamp DESC")

	if limit > 0 {
		query = query.Limit(limit)
	}

	err := query.Find(&messages).Error
	return messages, err
}

func NewMessageRepository(db *gorm.DB) MessageRepository {
	return &messageRepository{db: db}
}
