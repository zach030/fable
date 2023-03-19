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

func NewOpenAI(apikey string) *OpenAI {
	client := openai.NewClient(apikey)
	return &OpenAI{
		client:      client,
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
	a.client.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
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
	return "", nil
}

func (a *OpenAI) PreparePrompt(ctx context.Context, input []string, embeddings []string) (string, error) {
	return "", nil
}
