package milvus

import "github.com/milvus-io/milvus-sdk-go/v2/entity"

var (
	defaultCollection      = `fable`
	defaultValueColumn     = "content"
	defaultVectorColumn    = "content_vector"
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
					entity.TypeParamDim: "768",
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
