package llm

import "context"

type Llm interface {
	Embedding(ctx context.Context, input []string) ([]float32, error)
	Completion(ctx context.Context, prompt string) ([]string, error)
	PreparePrompt(lan int, context []string, question string) string
}
