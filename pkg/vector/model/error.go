package model

import "errors"

const (
	OpenAIEmbeddingDimensions = 1536
)

var (
	ErrInconsistentLength = errors.New("inconsistent length between value and vector")
)
