package gateways

import (
	"chat-service/internal/domain/entities"
	"context"
)

type ChatGateway interface {
	CreateChat(ctx context.Context, chat *entities.Chat) error
	FindChatById(ctx context.Context, id string) (*entities.Chat, error)
	SaveChat(ctx context.Context, chat *entities.Chat) error
}
