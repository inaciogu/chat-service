package entities

import (
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/tiktoken-go/tokenizer"
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
	enc, err := tokenizer.Get(tokenizer.Cl100kBase)

	if err != nil {
		return nil, err
	}
	tokenIds, _, _ := enc.Encode(content)

	currentTokens := len(tokenIds)

	fmt.Println("currentTokens", currentTokens)

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
	roles := []string{"user", "system", "assistant"}

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
