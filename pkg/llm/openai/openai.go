package openai

import (
	"context"

	"github.com/sashabaranov/go-openai"
)

type OpenAI struct {
	client      *openai.Client
	apiKey      string
	gptModel    string
	temperature float32
}

func NewOpenAI(apikey, baseURL string) *OpenAI {
	config := openai.DefaultConfig(apikey)
	config.BaseURL = baseURL
	return &OpenAI{
		client:      openai.NewClientWithConfig(config),
		apiKey:      apikey,
		gptModel:    openai.GPT3Dot5Turbo,
		temperature: 0.9,
	}
}

func (a *OpenAI) Embedding(ctx context.Context, input []string) ([]float32, error) {
	rsp, err := a.client.CreateEmbeddings(ctx, openai.EmbeddingRequest{
		Input: input,
		Model: openai.AdaEmbeddingV2,
	})
	if err != nil {
		return nil, err
	}
	embeddings := make([]float32, 0)
	for _, embedding := range rsp.Data {
		embeddings = append(embeddings, embedding.Embedding...)
	}
	return embeddings, nil
}

func (a *OpenAI) Completion(ctx context.Context, prompt string) (string, error) {
	_, err := a.client.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
		Model:            a.gptModel,
		Messages:         nil,
		MaxTokens:        0,
		Temperature:      a.temperature,
		TopP:             0,
		N:                0,
		Stream:           false,
		Stop:             nil,
		PresencePenalty:  0,
		FrequencyPenalty: 0,
		LogitBias:        nil,
		User:             "",
	})
	if err != nil {
		return "", err
	}
	return "", nil
}

func (a *OpenAI) PreparePrompt(ctx context.Context, input []string, embeddings []string) (string, error) {
	return "", nil
}
