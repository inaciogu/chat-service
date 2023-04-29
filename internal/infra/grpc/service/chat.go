package service

import (
	"chat-service/internal/infra/grpc/pb"
	chatcompletionstream "chat-service/internal/useCases/chatCompletionStream"
)

type ChatService struct {
	pb.UnimplementedChatServiceServer
	ChatCompletionStreamUseCase chatcompletionstream.ChatCompletionUseCase
	ChatConfigStream            chatcompletionstream.ChatCompletionStreamConfigInputDTO
	StreamChannel               chan chatcompletionstream.ChatCompletionStreamOutputDTO
}

func NewChatService(
	chatCompletionStreamUseCase chatcompletionstream.ChatCompletionUseCase,
	chatConfigStream chatcompletionstream.ChatCompletionStreamConfigInputDTO,
	streamChannel chan chatcompletionstream.ChatCompletionStreamOutputDTO,
) *ChatService {
	return &ChatService{
		ChatCompletionStreamUseCase: chatCompletionStreamUseCase,
		ChatConfigStream:            chatConfigStream,
		StreamChannel:               streamChannel,
	}
}

func (s *ChatService) ChatStream(req *pb.ChatRequest, stream pb.ChatService_ChatStreamServer) error {
	chatConfig := chatcompletionstream.ChatCompletionStreamConfigInputDTO{
		Model:                s.ChatConfigStream.Model,
		ModelMaxTokens:       s.ChatConfigStream.ModelMaxTokens,
		Temperature:          s.ChatConfigStream.Temperature,
		TopP:                 s.ChatConfigStream.TopP,
		N:                    s.ChatConfigStream.N,
		Stop:                 s.ChatConfigStream.Stop,
		MaxTokens:            s.ChatConfigStream.MaxTokens,
		InitialSystemMessage: s.ChatConfigStream.InitialSystemMessage,
	}

	input := chatcompletionstream.ChatCompletionStreamInputDTO{
		UserID:      req.GetUserId(),
		UserMessage: req.GetUserMessage(),
		ChatID:      req.GetChatId(),
		Config:      chatConfig,
	}

	ctx := stream.Context()

	consume := func() {
		for msg := range s.StreamChannel {
			stream.Send(&pb.ChatResponse{
				ChatId:  msg.ChatID,
				UserId:  msg.UserID,
				Content: msg.Content,
			})
		}
	}

	go consume()

	_, err := s.ChatCompletionStreamUseCase.Execute(ctx, input)

	if err != nil {
		return err
	}
	return nil
}
