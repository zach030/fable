package utils

import (
	"fmt"
	"log"
)

const (
	contentLength    = 800
	singleTokenLimit = 400
)

type SplitResult struct {
	Content       string
	ContentLength int
	ContentTokens int
}

func TextSplit(input string) ([]SplitResult, error) {
	length := len(input)
	batch := length/contentLength + 1
	var result []SplitResult
	for i := 0; i < batch; i++ {
		start, end := i*contentLength, (i+1)*contentLength
		if end > len(input) {
			end = len(input)
		}
		str := input[start:end]
		tokens := TokensNum(str)
		fmt.Printf("str length=%v tokens=%v\n", len(str), tokens)
		if tokens > singleTokenLimit {
			log.Fatal("exceed token limit")
		}
		result = append(result, SplitResult{Content: str, ContentLength: len(str), ContentTokens: tokens})
	}
	return result, nil
}
