package utils

import tg "github.com/pkoukk/tiktoken-go"

const (
	OpenAITokenLimit = 4096
)

func TokensNum(input string) int {
	encoder, err := tg.GetEncoding("cl100k_base")
	if err != nil {
		return 0
	}
	return len(encoder.Encode(input, nil, nil))
}
