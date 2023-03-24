package model

type SearchRequest struct {
	Collection string
	ContentKey string
	Vector     []float32
}

type SearchResult struct {
	Payload  string        `json:"payload"`
	Score    float32       `json:"score"`
	Key      string        `json:"key"`
	Metadata ChunkMetadata `json:"metadata"`
}

func NewSearchRequest(collection, contentKey string, vector []float32) *SearchRequest {
	return &SearchRequest{Collection: collection, ContentKey: contentKey, Vector: vector}
}
