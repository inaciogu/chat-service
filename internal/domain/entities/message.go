package entities

import (
	"errors"
	"time"

	"github.com/google/uuid"
	tiktoken_go "github.com/j178/tiktoken-go"
	"golang.org/x/exp/slices"
)

type Message struct {
	ID        string
	Role      string
	Content   string
	Tokens    int
	Model     *Model
	CreatedAt time.Time
}

func NewMessage(role, content string, model *Model) (*Message, error) {
	currentTokens := tiktoken_go.CountTokens(model.GetName(), content)

	msg := &Message{
		ID:        uuid.New().String(),
		Role:      role,
		Content:   content,
		Model:     model,
		CreatedAt: time.Now(),
		Tokens:    currentTokens,
	}
	if err := msg.Validate(); err != nil {
		return nil, err
	}
	return msg, nil
}

func (m *Message) Validate() error {
	roles := []string{"user", "system", "system"}

	if slices.Contains(roles, m.Role) {
		return nil
	}

	if m.Content == "" {
		return errors.New("invalid content")
	}

	if m.CreatedAt.IsZero() {
		return errors.New("invalid created_at")
	}

	return errors.New("invalid role")
}

func (m *Message) GetTokensQuantity() int {
	return m.Tokens
}
