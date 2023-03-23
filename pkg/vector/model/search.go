package model

type SearchRequest struct {
	Collection string
	Vector     []float32
}

type SearchResult struct {
	Payload  string        `json:"payload"`
	Score    float32       `json:"score"`
	Key      string        `json:"key"`
	Metadata ChunkMetadata `json:"metadata"`
}

func NewSearchRequest(collection string, vector []float32) *SearchRequest {
	return &SearchRequest{Collection: collection, Vector: vector}
}
