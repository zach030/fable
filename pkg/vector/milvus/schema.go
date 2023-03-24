package milvus

import (
	"encoding/json"
	"strconv"

	"github.com/milvus-io/milvus-sdk-go/v2/entity"
	"github.com/zach030/fable/pkg/vector/model"
)

const (
	filedContentKey    = "content_key"
	filedContent       = "content"
	filedContentVector = "content_vector"
	filedMetadata      = "metadata"
)

var (
	defaultCollection   = `fable`
	defaultValueColumn  = "content"
	defaultVectorColumn = "content_vector"

	fableDefaultCollection = &entity.Schema{
		CollectionName: defaultCollection,
		Description:    "this is the example collection for fableh",
		AutoID:         false,
		Fields: []*entity.Field{
			{
				Name:       "id",
				DataType:   entity.FieldTypeInt64, // int64 only for now
				PrimaryKey: true,
				AutoID:     true,
			},
			{
				Name:     "content_vector",
				DataType: entity.FieldTypeFloatVector,
				TypeParams: map[string]string{
					entity.TypeParamDim: strconv.Itoa(model.OpenAIEmbeddingDimensions),
				},
			},
			{
				Name:     "content",
				DataType: entity.FieldTypeVarChar,
				TypeParams: map[string]string{
					entity.TypeParamMaxLength: "800",
				},
			},
			{
				Name:     "content_key",
				DataType: entity.FieldTypeVarChar,
				TypeParams: map[string]string{
					entity.TypeParamMaxLength: "256",
				},
			},
			{
				Name:     "metadata",
				DataType: entity.FieldTypeVarChar,
				TypeParams: map[string]string{
					entity.TypeParamMaxLength: "800",
				},
			},
		},
	}
)

func newContentKey(key string) *entity.ColumnVarChar {
	return entity.NewColumnVarChar(filedContentKey, []string{key})
}

func newContent(val string) *entity.ColumnVarChar {
	return entity.NewColumnVarChar(filedContent, []string{val})
}

func newMetadata(metadata model.ChunkMetadata) *entity.ColumnVarChar {
	var str string
	buf, err := json.Marshal(metadata)
	if err == nil {
		str = string(buf)
	}
	return entity.NewColumnVarChar(filedMetadata, []string{str})
}

func newContentVector(vec []float32) *entity.ColumnFloatVector {
	vecs := make([][]float32, 0)
	return entity.NewColumnFloatVector(filedContentVector, model.OpenAIEmbeddingDimensions, append(vecs, vec))
}
