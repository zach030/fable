package milvus

import (
	"strconv"

	"github.com/milvus-io/milvus-sdk-go/v2/entity"
	"github.com/zach030/fable/pkg/vector/model"
)

var (
	defaultCollection      = `medium_articles`
	defaultValueColumn     = "title"
	defaultVectorColumn    = "title_vector"
	fableDefaultCollection = &entity.Schema{
		CollectionName: defaultCollection,
		Description:    "this is the example collection for fableh",
		AutoID:         false,
		Fields: []*entity.Field{
			{
				Name:       "ID",
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
					entity.TypeParamMaxLength: "255",
				},
			},
		},
	}
)
