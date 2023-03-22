package vector

import (
	"context"

	"github.com/zach030/fable/pkg/vector/model"
)

type Vector interface {
	Insert(ctx context.Context, collection string, value []string, vector [][]float32) error
	Search(ctx context.Context, collection, field string, vector []float32) ([]model.SearchResult, error)
}
