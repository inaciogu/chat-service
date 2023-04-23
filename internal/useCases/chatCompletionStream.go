package usecases

import (
	"chat-service/internal/domain/entities"
	"chat-service/internal/domain/gateways"
	"context"
	"errors"
	"io"
	"strings"

	"github.com/sashabaranov/go-openai"
)

type ChatCompletionUseCase struct {
	ChatGateway  gateways.ChatGateway
	OpenAIClient openai.Client
	Stream       chan ChatCompletionOutputDTO
}

type ChatCompletionConfigInputDTO struct {
	Model                string
	ModelMaxTokens       int
	Temperature          float32
	TopP                 float32
	N                    int
	Stop                 []string
	MaxTokens            int
	PresencePenalty      float32
	FrequencyPenalty     float32
	InitialSystemMessage string
}

type ChatCompletionInputDTO struct {
	ChatID      string
	UserID      string
	UserMessage string
	Config      ChatCompletionConfigInputDTO
}

type ChatCompletionOutputDTO struct {
	ChatID  string
	UserID  string
	Content string
}

func NewChatCompletionUseCase(chatGateway gateways.ChatGateway, openAIClient openai.Client, stream chan ChatCompletionOutputDTO) *ChatCompletionUseCase {
	return &ChatCompletionUseCase{
		ChatGateway:  chatGateway,
		OpenAIClient: openAIClient,
		Stream:       stream,
	}
}

func (uc *ChatCompletionUseCase) Execute(ctx context.Context, input ChatCompletionInputDTO) (*ChatCompletionOutputDTO, error) {
	chat, err := uc.ChatGateway.FindChatById(ctx, input.ChatID)

	if err != nil {
		if err.Error() == "chat not found" {
			chat, err = createNewChat(input)
			if err != nil {
				return nil, errors.New("failed to create new chat: " + err.Error())
			}
			err = uc.ChatGateway.CreateChat(ctx, chat)
			if err != nil {
				return nil, errors.New("failed to persist new chat: " + err.Error())
			}
		} else {
			return nil, errors.New("failed to find chat: " + err.Error())
		}
	}

	userMessage, err := entities.NewMessage("user", input.UserMessage, chat.Config.Model)

	if err != nil {
		return nil, errors.New("failed to create user message: " + err.Error())
	}

	err = chat.AddMessage(userMessage)

	if err != nil {
		return nil, errors.New("failed to add user message to chat: " + err.Error())
	}

	messages := []openai.ChatCompletionMessage{}

	for _, message := range chat.Messages {
		messages = append(messages, openai.ChatCompletionMessage{
			Role:    message.Role,
			Content: message.Content,
		})
	}

	response, err := uc.OpenAIClient.CreateChatCompletionStream(ctx, openai.ChatCompletionRequest{
		Model:            chat.Config.Model.Name,
		Messages:         messages,
		MaxTokens:        chat.Config.MaxTokens,
		Temperature:      chat.Config.Temperature,
		TopP:             chat.Config.TopP,
		Stop:             chat.Config.Stop,
		PresencePenalty:  chat.Config.PresencePenalty,
		FrequencyPenalty: chat.Config.FrequencyPenalty,
		Stream:           true,
	})

	if err != nil {
		return nil, errors.New("failed to create chat completion stream: " + err.Error())
	}

	var fullResponse strings.Builder

	for {
		response, err := response.Recv()

		if errors.Is(err, io.EOF) {
			break
		}

		if err != nil {
			return nil, errors.New("failed to receive chat completion response: " + err.Error())
		}

		fullResponse.WriteString(response.Choices[0].Delta.Content)

		r := ChatCompletionOutputDTO{
			ChatID:  chat.ID,
			UserID:  input.UserID,
			Content: fullResponse.String(),
		}

		uc.Stream <- r
	}

	assistant, err := entities.NewMessage("assistant", fullResponse.String(), chat.Config.Model)

	if err != nil {
		return nil, errors.New("failed to create assistant message: " + err.Error())
	}

	err = chat.AddMessage(assistant)

	if err != nil {
		return nil, errors.New("failed to add assistant message to chat: " + err.Error())
	}

	err = uc.ChatGateway.SaveChat(ctx, chat)

	if err != nil {
		return nil, errors.New("failed to save chat: " + err.Error())
	}
	return &ChatCompletionOutputDTO{
		ChatID:  chat.ID,
		UserID:  input.UserID,
		Content: fullResponse.String(),
	}, nil
}

func createNewChat(input ChatCompletionInputDTO) (*entities.Chat, error) {
	model := entities.NewModel(input.Config.Model, input.Config.ModelMaxTokens)
	chatConfig := &entities.ChatConfig{
		Model:            model,
		Temperature:      input.Config.Temperature,
		TopP:             input.Config.TopP,
		N:                input.Config.N,
		Stop:             input.Config.Stop,
		MaxTokens:        input.Config.MaxTokens,
		PresencePenalty:  input.Config.PresencePenalty,
		FrequencyPenalty: input.Config.FrequencyPenalty,
	}

	initialSystemMessage, err := entities.NewMessage("system", input.Config.InitialSystemMessage, model)

	if err != nil {
		return nil, errors.New("failed to create initial system message: " + err.Error())
	}

	chat, err := entities.NewChat(input.UserID, initialSystemMessage, chatConfig)

	if err != nil {
		return nil, errors.New("failed to create chat: " + err.Error())
	}

	return chat, nil
}
