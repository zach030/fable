package model

type SearchResult struct {
	Payload string  `json:"payload"`
	Score   float32 `json:"score"`
}
