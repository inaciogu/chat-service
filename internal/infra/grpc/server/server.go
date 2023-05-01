package server

import (
	chatcompletionstream "chat-service/internal/useCases/chatCompletionStream"
	"net"

	"chat-service/internal/infra/grpc/pb"
	service "chat-service/internal/infra/grpc/service"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type GRPCServer struct {
	ChatCompletionStreamUseCase chatcompletionstream.ChatCompletionUseCase
	ChatConfigStream            chatcompletionstream.ChatCompletionStreamConfigInputDTO
	ChatService                 service.ChatService
	Port                        string
	AuthToken                   string
	StreamChannel               chan chatcompletionstream.ChatCompletionStreamOutputDTO
}

func NewGRPCServer(
	chatCompletionStreamUseCase chatcompletionstream.ChatCompletionUseCase,
	chatConfigStream chatcompletionstream.ChatCompletionStreamConfigInputDTO,
	streamChannel chan chatcompletionstream.ChatCompletionStreamOutputDTO,
	port string,
	authToken string,
) *GRPCServer {
	chatService := service.NewChatService(chatCompletionStreamUseCase, chatConfigStream, streamChannel)

	return &GRPCServer{
		ChatCompletionStreamUseCase: chatCompletionStreamUseCase,
		ChatConfigStream:            chatConfigStream,
		Port:                        port,
		AuthToken:                   authToken,
		StreamChannel:               streamChannel,
		ChatService:                 *chatService,
	}
}

func (s *GRPCServer) AuthInterceptor(srv interface{}, streamServer grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	ctx := streamServer.Context()

	meta, ok := metadata.FromIncomingContext(ctx)

	if !ok {
		return status.Error(codes.Unauthenticated, "missing context metadata")
	}

	token := meta.Get("authorization")

	if len(token) == 0 {
		return status.Error(codes.Unauthenticated, "missing authorization token")
	}

	if token[0] != s.AuthToken {
		return status.Error(codes.Unauthenticated, "invalid authorization token")
	}
	return handler(srv, streamServer)
}

func (s *GRPCServer) Start() {
	options := []grpc.ServerOption{
		grpc.StreamInterceptor(s.AuthInterceptor),
	}

	grpcServer := grpc.NewServer(options...)

	pb.RegisterChatServiceServer(grpcServer, &s.ChatService)

	listener, err := net.Listen("tcp", ":"+s.Port)

	if err != nil {
		panic(err.Error())
	}

	if err := grpcServer.Serve(listener); err != nil {
		panic(err.Error())
	}
}
