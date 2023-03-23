package model

type InsertRequest struct {
	Collection string
	ContentKey string
	Content    []string
	Embeddings [][]float32
	Metadata   []ChunkMetadata
}

func NewInsertRequest(collection, key string, content []string, embedding [][]float32) *InsertRequest {
	if len(content) != len(embedding) {
		return nil
	}
	return &InsertRequest{
		Collection: collection,
		ContentKey: key,
		Content:    content,
		Embeddings: embedding,
	}
}
