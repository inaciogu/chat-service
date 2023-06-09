package handlers

import (
	chatcompletion "chat-service/internal/useCases/chatCompletion"
	"encoding/json"
	"io"
	"net/http"
)

type WebChatGptHandler struct {
	CompletionUseCase chatcompletion.ChatCompletionUseCase
	Config            chatcompletion.ChatCompletionConfigInputDTO
	AuthToken         string
}

func NewWebChatGptHandler(completionUseCase chatcompletion.ChatCompletionUseCase, config chatcompletion.ChatCompletionConfigInputDTO, token string) *WebChatGptHandler {
	return &WebChatGptHandler{
		CompletionUseCase: completionUseCase,
		Config:            config,
		AuthToken:         token,
	}
}

func (h *WebChatGptHandler) Handle(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	if r.Header.Get("Authorization") != h.AuthToken {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	body, err := io.ReadAll(r.Body)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	if !json.Valid(body) {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	var input chatcompletion.ChatCompletionInputDTO

	err = json.Unmarshal(body, &input)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	input.Config = h.Config

	output, err := h.CompletionUseCase.Execute(r.Context(), input)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(output)
}
