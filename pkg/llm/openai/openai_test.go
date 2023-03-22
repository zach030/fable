package openai

import (
	"context"
	"testing"
)

var (
	openaiClient *OpenAI
	ak           = ""
	proxyURL     = ""
)

func TestOpenAI_Embedding(t *testing.T) {
	openaiClient = NewOpenAI(ak, proxyURL)
	res, err := openaiClient.Embedding(context.Background(), []string{"The dawn of Dark Mode"})
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("res=%+v", res)
}
