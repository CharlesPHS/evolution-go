package message_model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Message struct {
	Id              string    `json:"id" gorm:"type:uuid;primaryKey"`
	InstanceId      string    `json:"instanceId" gorm:"index"`
	MessageID       string    `json:"message_id" gorm:"unique"`
	Timestamp       time.Time `json:"timestamp" gorm:"index"`
	Status          string    `json:"status"`
	Source          string    `json:"source" gorm:"index"`
	Chat            string    `json:"chat" gorm:"index"`
	FromMe          bool      `json:"fromMe"`
	Content         string    `json:"content" gorm:"type:text"`
	MessageType     string    `json:"messageType"`
	QuotedMessageID string    `json:"quotedMessageId"`
}

func (m *Message) BeforeCreate(tx *gorm.DB) (err error) {
	m.Id = uuid.New().String()
	return
}
