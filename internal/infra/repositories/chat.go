package repositories

import (
	"chat-service/internal/domain/entities"
	"chat-service/internal/infra/db"
	"context"
	"database/sql"
	"errors"
	"time"
)

type ChatRepositoryMySQL struct {
	DB      *sql.DB
	Queries *db.Queries
}

func NewChatRepositoryMySQL(dbt *sql.DB) *ChatRepositoryMySQL {
	return &ChatRepositoryMySQL{
		DB:      dbt,
		Queries: db.New(dbt),
	}
}

func (r *ChatRepositoryMySQL) CreateChat(ctx context.Context, chat *entities.Chat) error {
	err := r.Queries.CreateChat(
		ctx,
		db.CreateChatParams{
			ID:               chat.ID,
			UserID:           chat.UserID,
			InitialMessageID: chat.InitialSystemMessage.ID,
			Status:           chat.Status,
			TokenUsage:       int32(chat.TokenUsage),
			Model:            chat.Config.Model.Name,
			Temperature:      float64(chat.Config.Temperature),
			TopP:             float64(chat.Config.TopP),
			N:                int32(chat.Config.N),
			Stop:             chat.Config.Stop[0],
			MaxTokens:        int32(chat.Config.MaxTokens),
			PresencePenalty:  float64(chat.Config.PresencePenalty),
			FrequencyPenalty: float64(chat.Config.FrequencyPenalty),
			CreatedAt:        time.Now(),
			UpdatedAt:        time.Now(),
		},
	)

	if err != nil {
		return err
	}

	err = r.Queries.AddMessage(
		ctx,
		db.AddMessageParams{
			ID:        chat.InitialSystemMessage.ID,
			ChatID:    chat.ID,
			Role:      chat.InitialSystemMessage.Role,
			Content:   chat.InitialSystemMessage.Content,
			Tokens:    int32(chat.InitialSystemMessage.Tokens),
			CreatedAt: chat.InitialSystemMessage.CreatedAt,
		},
	)

	if err != nil {
		return err
	}

	return nil
}

func (r *ChatRepositoryMySQL) FindChatById(ctx context.Context, chatId string) (*entities.Chat, error) {
	chat := &entities.Chat{}

	response, err := r.Queries.FindChatById(ctx, chatId)

	if err != nil {
		return nil, errors.New("chat not found")
	}

	chat.ID = response.ID
	chat.UserID = response.UserID
	chat.Status = response.Status
	chat.TokenUsage = int(response.TokenUsage)
	chat.Config = &entities.ChatConfig{
		Model: &entities.Model{
			Name:      response.Model,
			MaxTokens: int(response.ModelMaxTokens),
		},
		Temperature:      float32(response.Temperature),
		TopP:             float32(response.TopP),
		N:                int(response.N),
		Stop:             []string{response.Stop},
		MaxTokens:        int(response.MaxTokens),
		PresencePenalty:  float32(response.PresencePenalty),
		FrequencyPenalty: float32(response.FrequencyPenalty),
	}

	messages, err := r.Queries.FindMessagesByChatId(ctx, chatId)

	if err != nil {
		return nil, err
	}

	for _, message := range messages {
		chat.Messages = append(chat.Messages, &entities.Message{
			ID:        message.ID,
			Role:      message.Role,
			Content:   message.Content,
			Tokens:    int(message.Tokens),
			Model:     &entities.Model{Name: message.Model},
			CreatedAt: message.CreatedAt,
		})
	}

	erasedMessages, err := r.Queries.FindErasedMessagesByChatId(ctx, chatId)

	if err != nil {
		return nil, err
	}

	for _, erasedMessage := range erasedMessages {
		chat.ErasedMessages = append(chat.ErasedMessages, &entities.Message{
			ID:        erasedMessage.ID,
			Role:      erasedMessage.Role,
			Content:   erasedMessage.Content,
			Tokens:    int(erasedMessage.Tokens),
			Model:     &entities.Model{Name: erasedMessage.Model},
			CreatedAt: erasedMessage.CreatedAt,
		})
	}

	return chat, nil
}

func (r *ChatRepositoryMySQL) SaveChat(ctx context.Context, chat *entities.Chat) error {
	params := db.SaveChatParams{
		ID:               chat.ID,
		UserID:           chat.UserID,
		Status:           chat.Status,
		TokenUsage:       int32(chat.TokenUsage),
		Model:            chat.Config.Model.Name,
		ModelMaxTokens:   int32(chat.Config.Model.MaxTokens),
		Temperature:      float64(chat.Config.Temperature),
		TopP:             float64(chat.Config.TopP),
		N:                int32(chat.Config.N),
		Stop:             chat.Config.Stop[0],
		MaxTokens:        int32(chat.Config.MaxTokens),
		PresencePenalty:  float64(chat.Config.PresencePenalty),
		FrequencyPenalty: float64(chat.Config.FrequencyPenalty),
		UpdatedAt:        time.Now(),
	}

	err := r.Queries.SaveChat(
		ctx,
		params,
	)
	if err != nil {
		return err
	}

	err = r.Queries.DeleteChatMessages(ctx, chat.ID)
	if err != nil {
		return err
	}

	err = r.Queries.DeleteErasedChatMessages(ctx, chat.ID)
	if err != nil {
		return err
	}

	i := 0
	for _, message := range chat.Messages {
		err = r.Queries.AddMessage(
			ctx,
			db.AddMessageParams{
				ID:        message.ID,
				ChatID:    chat.ID,
				Content:   message.Content,
				Role:      message.Role,
				Tokens:    int32(message.Tokens),
				Model:     chat.Config.Model.Name,
				CreatedAt: message.CreatedAt,
				OrderMsg:  int32(i),
				Erased:    false,
			},
		)
		if err != nil {
			return err
		}
		i++
	}

	i = 0
	for _, message := range chat.ErasedMessages {
		err = r.Queries.AddMessage(
			ctx,
			db.AddMessageParams{
				ID:        message.ID,
				ChatID:    chat.ID,
				Content:   message.Content,
				Role:      message.Role,
				Tokens:    int32(message.Tokens),
				Model:     chat.Config.Model.Name,
				CreatedAt: message.CreatedAt,
				OrderMsg:  int32(i),
				Erased:    true,
			},
		)
		if err != nil {
			return err
		}
		i++
	}
	return nil
}
