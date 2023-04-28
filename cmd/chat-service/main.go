package main

import (
	"chat-service/configs"
	"chat-service/internal/infra/repositories"
	"chat-service/internal/infra/web/handlers"
	"chat-service/internal/infra/web/server"
	chatcompletion "chat-service/internal/useCases/chatCompletion"
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/sashabaranov/go-openai"
)

func main() {
	configs := configs.LoadConfig(".")

	connection, err := sql.Open(configs.DBDriver, fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true&multiStatements=true", configs.DBUser, configs.DBPassword, configs.DBHost, configs.DBPort, configs.DBName))

	if err != nil {
		panic(err)
	}

	defer connection.Close()

	repo := repositories.NewChatRepositoryMySQL(connection)
	client := openai.NewClient(configs.OpenAIApiKey)

	chatConfig := chatcompletion.ChatCompletionConfigInputDTO{
		Model:                configs.Model,
		ModelMaxTokens:       configs.ModelMaxTokens,
		Temperature:          float32(configs.Temperature),
		TopP:                 float32(configs.TopP),
		N:                    configs.N,
		Stop:                 configs.Stop,
		MaxTokens:            configs.MaxTokens,
		InitialSystemMessage: configs.InitialChatMessage,
	}

	/* chatConfigStream := chatcompletionstream.ChatCompletionStreamConfigInputDTO{
		Model:                configs.Model,
		ModelMaxTokens:       configs.ModelMaxTokens,
		Temperature:          float32(configs.Temperature),
		TopP:                 float32(configs.TopP),
		N:                    configs.N,
		Stop:                 configs.Stop,
		MaxTokens:            configs.ModelMaxTokens,
		InitialSystemMessage: configs.InitialChatMessage,
	} */

	chatCompletion := chatcompletion.NewChatCompletionUseCase(repo, client)

	/* streamChannel := make(chan chatcompletionstream.ChatCompletionStreamOutputDTO) */
	/* chatCompletionStream := chatcompletionstream.NewChatCompletionStreamUseCase(repo, client, streamChannel) */

	webServer := server.NewWebServer(":" + configs.WebServerPort)

	webServerChatHandler := handlers.NewWebChatGptHandler(*chatCompletion, chatConfig, configs.AuthToken)

	webServer.AddHandler("/chat", webServerChatHandler.Handle)

	fmt.Println("Server is running on port " + configs.WebServerPort)
	webServer.Start()
}
