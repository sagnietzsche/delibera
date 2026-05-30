package internal

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/joho/godotenv"
	openai "github.com/sashabaranov/go-openai"
)

type LLMProvider interface {
	Prompt(ctx context.Context, userMessage string) (string, error)
}

type OpenAICompatibleProvider struct {
	client      *openai.Client
	model       string
	maxTokens   int
	temperature float32
	topP        float32
}

func NewInferenceProviderFromEnv() (*OpenAICompatibleProvider, error) {
	_ = godotenv.Load()

	apiKey := os.Getenv(InferenceAPIKeyEnv)
	if apiKey == "" {
		return nil, fmt.Errorf("%s is required", InferenceAPIKeyEnv)
	}

	config := openai.DefaultConfig(apiKey)
	config.BaseURL = InferenceBaseURL

	return &OpenAICompatibleProvider{
		client:      openai.NewClientWithConfig(config),
		model:       Model,
		maxTokens:   512,
		temperature: 1.0,
		topP:        1.0,
	}, nil
}

func (p *OpenAICompatibleProvider) Prompt(ctx context.Context, userMessage string) (string, error) {
	userMessage = strings.TrimSpace(userMessage)
	if userMessage == "" {
		return "", errors.New("user message is required")
	}

	stream, err := p.client.CreateChatCompletionStream(
		ctx,
		openai.ChatCompletionRequest{
			Model: p.model,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleUser,
					Content: userMessage,
				},
			},
			MaxTokens:   p.maxTokens,
			Temperature: p.temperature,
			TopP:        p.topP,
		},
	)
	if err != nil {
		return "", err
	}
	defer stream.Close()

	var response strings.Builder
	for {
		resp, err := stream.Recv()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return "", err
		}
		if len(resp.Choices) == 0 {
			continue
		}
		response.WriteString(resp.Choices[0].Delta.Content)
	}

	return response.String(), nil
}

func CallInference(ctx context.Context, userMessage string) (string, error) {
	provider, err := NewInferenceProviderFromEnv()
	if err != nil {
		return "", err
	}

	return provider.Prompt(ctx, userMessage)
}
