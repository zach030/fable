package vector

import "context"

type Vector interface {
	Insert(ctx context.Context, value string, vector []float32) error
	CreateIndex(ctx context.Context, field string) error
	Search(ctx context.Context, field string, vector []float32) error
}
