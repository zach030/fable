package openai

import (
	"context"
	"fmt"

	"github.com/sashabaranov/go-openai"
)

const (
	ch = 1
	en = 2
)

const (
	enPromptTemplate = `You are an AI assistant for the open source library LangChain. The documentation is located at https://langchain.readthedocs.io.
You are given the following extracted parts of a long document and a question. Provide a conversational answer with a hyperlink to the documentation.
You should only use hyperlinks that are explicitly listed as a source in the context. Do NOT make up a hyperlink that is not listed.
If the question includes a request for code, provide a code block directly from the documentation.
If you don't know the answer, just say "Hmm, I'm not sure." Don't try to make up an answer.
If the question is not about LangChain, politely inform them that you are tuned to only answer questions about LangChain.
Question: %s
=========
%v
=========
Answer in Markdown:`

	chPromptTemplate = `你是一位助理，接下来我会给你从长篇文章中抽取的一部分和输入的提问，如果你不知道答案，请回答：“我不太了解这个问题”，不要编造答案，在回答时避免和原文内容完全一致，尽量用你自己的词语，保证回答精确，有帮助，无攻击性。
我提供的资料：%v,
用户的问题: %s`
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
	res, err := a.client.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
		Model: a.gptModel,
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleSystem,
				Content: "You are a helpful assistant that accurately answers queries. Use the text provided to form your answer, but avoid copying word-for-word from the essays. Try to use your own words when possible. Keep your answer under 5 sentences. Be accurate, helpful, concise, and clear.",
			},
			{
				Role:    openai.ChatMessageRoleUser,
				Content: prompt,
			},
		},
		Temperature: a.temperature,
		Stream:      true,
	})
	if err != nil {
		return "", err
	}
	fmt.Println(res)
	return "", nil
}

func (a *OpenAI) PreparePrompt(lan int, context []string, question string) string {
	switch lan {
	case en:
		return fmt.Sprintf(enPromptTemplate, question, context)
	case ch:
		return fmt.Sprintf(chPromptTemplate, question, context)
	}
	return ""
}
