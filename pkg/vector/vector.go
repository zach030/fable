package vector

import (
	"context"

	"github.com/zach030/fable/pkg/vector/model"
)

type Vector interface {
	Insert(ctx context.Context, req *model.InsertRequest) error
	Search(ctx context.Context, req *model.SearchRequest) ([]model.SearchResult, error)
}
