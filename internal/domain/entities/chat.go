package entities

import (
	"errors"

	"github.com/google/uuid"
	"golang.org/x/exp/slices"
)

type ChatConfig struct {
	Model            *Model
	Temperature      float32
	TopP             float32
	N                int
	Stop             []string
	MaxTokens        int
	PresencePenalty  float32
	FrequencyPenalty float32
}

type Chat struct {
	ID                   string
	UserID               string
	Status               string
	TokenUsage           int
	Config               *ChatConfig
	InitialSystemMessage *Message
	Messages             []*Message
	ErasedMessages       []*Message
}

func NewChat(userID string, initialSystemMessage *Message, config *ChatConfig) (*Chat, error) {
	chat := &Chat{
		ID:                   uuid.New().String(),
		UserID:               userID,
		InitialSystemMessage: initialSystemMessage,
		Config:               config,
		Status:               "active",
		TokenUsage:           0,
	}
	if err := chat.Validate(); err != nil {
		return nil, err
	}
	chat.AddMessage(initialSystemMessage)

	return chat, nil
}

func (c *Chat) Validate() error {
	status := []string{"active", "ended"}

	if c.UserID == "" {
		return errors.New("invalid user id")
	}

	if !slices.Contains(status, c.Status) {
		return errors.New("invalid status")
	}

	if c.Config.Temperature < 0 || c.Config.Temperature > 2 {
		return errors.New("invalid temperature")
	}
	return nil
}

func (c *Chat) AddMessage(msg *Message) error {
	if c.Status == "ended" {
		return errors.New("chat is ended")
	}
	for {
		if c.Config.Model.GetMaxTokens() >= msg.GetTokensQuantity()+c.TokenUsage {
			c.Messages = append(c.Messages, msg)
			c.RefreshTokenUsage()
			break
		}
		c.ErasedMessages = append(c.Messages, c.Messages[0])
		c.Messages = c.Messages[1:]
		c.RefreshTokenUsage()
	}
	return nil
}

func (c *Chat) GetMessages() []*Message {
	return c.Messages
}

func (c *Chat) CountMessages() int {
	return len(c.Messages)
}

func (c *Chat) End() {
	c.Status = "ended"
}

func (c *Chat) RefreshTokenUsage() {
	c.TokenUsage = 0
	for i := range c.Messages {
		c.TokenUsage += c.Messages[i].GetTokensQuantity()
	}
}
