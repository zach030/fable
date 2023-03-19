package openai

import (
	"context"
	"testing"
)

var (
	openaiClient *OpenAI
	ak           = ""
)

func TestOpenAI_Embedding(t *testing.T) {
	openaiClient = NewOpenAI(ak)
	res, err := openaiClient.Embedding(context.Background(), []string{"what's the meaning of life"})
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("res=%+v", res)
}
