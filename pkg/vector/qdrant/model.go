package qdrant

type CommonResponse struct {
	Result interface{} `json:"result"`
	Status interface{} `json:"status"`
	Time   float64     `json:"time"`
}

type UpsertPointRequest struct {
	Points []Point `json:"points"`
}

type Point struct {
	ID      string      `json:"id"`
	Payload interface{} `json:"payload"`
	Vector  []float32   `json:"vector"`
}

type PointSearchRequest struct {
	Params      map[string]interface{} `json:"params"`
	Vector      []float32              `json:"vector"`
	Limit       int                    `json:"limit"`
	WithPayload bool                   `json:"with_payload"`
	WithVector  bool                   `json:"with_vector"`
}

type SearchResult struct {
	ID      string      `json:"id"`
	Version int         `json:"version"`
	Score   float64     `json:"score"`
	Payload interface{} `json:"payload"`
	Vector  []float32   `json:"vector,omitempty"`
}
